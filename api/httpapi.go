/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package api

import (
    "fmt"
    "time"
    "sync"
    "os"
    "bytes"
    "mime/multipart"
    "path/filepath"
    "io"
    "net/http"                                                                                                                                            
    "go4api/utils"
    "go4api/assertion"
    "go4api/logger"
    "go4api/protocal/http"
    "reflect"
    "encoding/json"
    "strings"
    simplejson "github.com/bitly/go-simplejson"
    // "strconv"
)

func HttpApi(wg *sync.WaitGroup, resultsChan chan []interface{}, options map[string]string, pStart string, baseUrl string, 
        tc []interface{}, resultsDir string) {
    //
    defer wg.Done()
    //
    tcName := tc[0].(string)
    tcJson := tc[3].(*simplejson.Json)
    jsonFile := tc[4].(string)
    csvFile := tc[5].(string)
    rowCsv := tc[6].(string)
    //
    start := string(time.Now().String())
    // fmt.Println(tcName, " start: ", start)
    //
    apiPath, apiMethod := utils.GetRequestForTC(tcJson, tcName)
    // apiPath := "/api/operation/soldtos?pageIndex=1&pageSize=12"
    url := ""
    if strings.HasPrefix(strings.ToLower(apiPath), "http") {
        url = apiPath
    } else {
        url = baseUrl + apiPath
    }
    
    //
    // the map for mapping the string and the related funciton to call
    funcs := map[string]interface{} {
        "GET": protocal.HttpGet,
        "POST": protocal.HttpPost,
        "POSTMultipart": protocal.HttpPostMultipart,
    }
    // request payload(body)
    var reqPayload map[string]interface{}
    reqPayload = utils.GetRequestPayloadForTC(tcJson, tcName)
    //
    var bodyText *strings.Reader // init body
    bodyMultipart := &bytes.Buffer{}
    var boundary string
    apiMethodSelector := apiMethod
    mv := reflect.ValueOf(reqPayload)
    // Note, has 3 conditions: text (json), form, or multipart file upload
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        // case 1: multipart upload
        if k.Interface().(string) == "filename" {
            // Note, hardcode the name = excel here, potential bug
            bodyMultipart, boundary, _ = PrepMultipart(options["testhome"] + "/testresource/" + v.Interface().(string), "excel")
            apiMethodSelector = "POSTMultipart"
            break
        }
        // case 2: normal json
        if k.Interface().(string) == "text" {
            bodyText = PrepPostPayload(reqPayload)
            break
        }
    }

    // request headers
    var reqHeaders map[string]interface{}
    reqHeaders = utils.GetRequestHeadersForTC(tcJson, tcName)
    // set the boundary to headers, if multipart
    if boundary != "" {
        reqHeaders["Content-Type"] = boundary
    }
    // fmt.Println(tcName + " boundary: ", boundary)
   

    // < !! ----------- !! >

    // Actual response
    var actualStatusCode string
    var actualHeader http.Header
    var actualBody string
    // protocalChan := make(chan interface{}, 50)
    if apiMethodSelector == "POSTMultipart" {
        actualStatusCode, actualHeader, actualBody = protocal.CallHttpMethod(funcs, apiMethodSelector, url, apiMethod, reqHeaders, bodyMultipart)    
    } else {
        actualStatusCode, actualHeader, actualBody = protocal.CallHttpMethod(funcs, apiMethodSelector, url, apiMethod, reqHeaders, bodyText)
        }
    //
    // Expected response
    expStatus, expHeader, expBody := utils.GetExpectedResponseForTC(tcJson, tcName)
    // fmt.Println(actualStatusCode, actualHeader, actualBody)

    // compare
    testResult := Compare(tcName, actualStatusCode, actualHeader, actualBody, expStatus, expHeader, expBody)
    end := string(time.Now().String())
    // fmt.Println(tcName + " end: ", end)

    //
    resultPrintString := ""
    csvFileBase := ""
    resultReportString := ""
    // Note: if csvFile does not exist, the filepath.Base(csvFile) = ".", need to remove
    if filepath.Base(csvFile) == "." {
        csvFileBase = ""
    } else {
        csvFileBase = filepath.Base(csvFile)
    }
    if actualStatusCode == "200" {
        resultPrintString1 := tcName + "," + actualStatusCode + "," + filepath.Base(jsonFile) + "," + csvFileBase + "," + rowCsv
        resultPrintString = resultPrintString1 + "," + testResult + "," 

        resultReportString1 := tcName + "," + actualStatusCode + "," + filepath.Base(jsonFile) + "," + csvFileBase + "," + rowCsv
        resultReportString = resultReportString1 + "," + start + "," + end + "," + testResult + "," 
    } else {
        resultPrintString1 := tcName + "," + actualStatusCode + "," + filepath.Base(jsonFile) + "," + csvFileBase + "," + rowCsv
        resultPrintString = resultPrintString1 + "," + testResult + "," + actualBody

        resultReportString1 := tcName + "," + actualStatusCode + "," + filepath.Base(jsonFile) + "," + csvFileBase + "," + rowCsv 
        resultReportString = resultReportString1 + "," + start + "," + end + "," + testResult + "," + "actualBody" 
    }
    fmt.Println(resultPrintString)

    // put the execution log into results
    logger.WriteExecutionResults(resultReportString, pStart, resultsDir)


    // write the channel to executor for scheduler
    var resultsChanArray []interface{}
    resultsChanArray = append(resultsChanArray, tcName)
    resultsChanArray = append(resultsChanArray, testResult)
    resultsChan <- resultsChanArray

}

func Compare(tcName string, actualStatusCode string, actualHeader http.Header, actualBody string, 
        expStatus map[string]interface{}, expHeader map[string]interface{}, expBody map[string]interface{}) string {

    // the map for mapping the string and the related funciton to call
    funcs := map[string]interface{} {
        "Equals": assertion.Equals,
        "Contains": assertion.Contains,
        // "LargerThan": assertion.LargerThan,
    }

    var testResults []bool
    // expStatusCode := "200"
    // status
    for key, _ := range expStatus {
        // fmt.Println("expStatus", key, expStatus[key])
        // call the assertion function
        testResult, _ := assertion.CallAssertion(funcs, key, expStatus[key], actualStatusCode)
        // fmt.Println(tcName, "expStatus", testResult[0])
        testResults = append(testResults, testResult[0].Interface().(bool))
    }

    // header
    for key, _ := range expHeader {
        // fmt.Println("expHeader", key, actualHeader[key])
        bytesExpHeader, _ := json.Marshal(expHeader)
        res, _ := simplejson.NewJson(bytesExpHeader)

        expHeader_2, _ := res.Get(key).Map()

        for comp_key, _ := range expHeader_2 {
            // fmt.Println("expHeader_2", comp_key, expHeader_2[comp_key], actualHeader[key][0])
            // call the assertion function
            testResult, _ := assertion.CallAssertion(funcs, comp_key, expHeader_2[comp_key], actualHeader[key][0])
            // fmt.Println(tcName, "expHeader", testResult[0])
            testResults = append(testResults, testResult[0].Interface().(bool))
        } 
    }

    // body
    for key, _ := range expBody {
        // fmt.Println("expHeader", key, actualHeader[key])

        bytesExpBody, _ := json.Marshal(expBody)
        res, _ := simplejson.NewJson(bytesExpBody)
        //
        // bytesActualBody, _ := json.Marshal(actualBody)
        actualRes, _ := simplejson.NewJson([]byte(actualBody))
        //
        expBody_2, _ := res.Get(key).Map()

        for comp_key, _ := range expBody_2 {
            // Note: here does not know the type of the value, but use Int first for demo, may be a bug
            actualValue, _ := actualRes.Get(key).Int()
            // fmt.Println("expBody_2", comp_key, expBody_2[comp_key], actualValue)
            // call the assertion function
            testResult, _ := assertion.CallAssertion(funcs, comp_key, expBody_2[comp_key], actualValue)
            // fmt.Println(tcName, "expBody", testResult[0])
            testResults = append(testResults, testResult[0].Interface().(bool))
        }
    }

    // default finalResults
    finalResults := "Success"

    for key := range testResults {
        if testResults[key] == false {
            finalResults = "Fail"
            // fmt.Println(tcName + " results", testResults, "final results: ", finalResults)
            break
        }
    }
    
    return finalResults

}  


func PrepMultipart(path string, name string) (*bytes.Buffer, string, error) {
    fp, err := os.Open(path) 
    if err != nil {
        panic(err)
    }
    defer fp.Close()

    body := &bytes.Buffer{} // init body
    writer := multipart.NewWriter(body) // multipart
    
    // prepare the reader instances to encode
    params := map[string]io.Reader{
        name:  fp, // it is file
        // "other": strings.NewReader("hello world!"),
    }
    //
    for key, r := range params {
        var fw io.Writer
        if x, ok := r.(io.Closer); ok {
            defer x.Close()
        }
        // Add an file
        if x, ok := r.(*os.File); ok {
            if fw, err = writer.CreateFormFile(key, x.Name()); err != nil {
                return nil, "", err
            }
        } else {
            // Add other fields
            if fw, err = writer.CreateFormField(key); err != nil {
                return nil, "", err
            }
        }
        if _, err = io.Copy(fw, r); err != nil {
            return nil, "", err
        }
    }
    //
    err = writer.Close()
    if err != nil {
        return nil, "", err
    }
    // do not forget this
    boundary := writer.FormDataContentType()
    // fmt.Println("boundary", boundary)
    // ==> i.e. multipart/form-data; boundary=37b1e9deba0159aaf429d7183a9de344c532e50299532f7b4f7bdbbca435

    return body, boundary, nil

}


func PrepPostPayload(reqPayload map[string]interface{}) *strings.Reader {
    mv := reflect.ValueOf(reqPayload)
    var body *strings.Reader

    // Note, has 3 conditions: text (json), form, or multipart file upload
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        // fmt.Println("reqPayload", k, v)
        if k.Interface().(string) == "text" {
            body = strings.NewReader(v.Interface().(string))
            break
        }
    }

    return body
}
