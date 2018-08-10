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
    "net/url"                                                                                                                                       
    "go4api/utils"
    "go4api/assertion"
    "go4api/logger"
    "go4api/protocal/http"
    "reflect"
    "encoding/json"
    "strings"
    simplejson "github.com/bitly/go-simplejson"
    "strconv"
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
    start_time := time.Now()
    start := start_time.String()
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
        "POSTForm": protocal.HttpPostForm,
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
        // case 3: if Post, and the key does not have filename, text, then it would be PostForm
        bodyText = PrepPostFormPayload(reqPayload)
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

    // (1). Actual response
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
    // (2). Expected response
    expStatus, expHeader, expBody := utils.GetExpectedResponseForTC(tcJson, tcName)
    // fmt.Println(actualStatusCode, actualHeader, actualBody)

    // (3). compare
    testResult, testMessages := Compare(tcName, actualStatusCode, actualHeader, actualBody, expStatus, expHeader, expBody)
    //
    end_time := time.Now()
    end := end_time.String()
    // fmt.Println(tcName + " end: ", end)

    // (4). here to generate the outputs file if the Json has "outputs" field
    // WriteOutputsDataToFile(testResult, tcJson, tcName, tc, actualBody)

    // (5). print to console
    resultPrintString := ""
    csvFileBase := ""
    resultReportString := ""
    // Note: if csvFile does not exist, the filepath.Base(csvFile) = ".", need to remove
    if filepath.Base(csvFile) == "." {
        csvFileBase = ""
    } else {
        csvFileBase = filepath.Base(csvFile)
    }
    resultPrintString1 := tcName + "," + actualStatusCode + "," + filepath.Base(jsonFile) + "," + csvFileBase + "," + rowCsv
    resultPrintString = resultPrintString1 + "," + testResult + "," + strings.Join(testMessages, ":")

    resultReportString1 := tcName + "," + actualStatusCode + "," + filepath.Base(jsonFile) + "," + csvFileBase + "," + rowCsv 
    resultReportString2 := "," + start + "," + end + "," + strconv.FormatInt(start_time.UnixNano(), 10) + "," + strconv.FormatInt(end_time.UnixNano(), 10)
    resultReportString3 := "," + strconv.FormatInt(end_time.UnixNano() - start_time.UnixNano(), 10) + "," +  testResult + "," + `message`
    resultReportString =  resultReportString1 + resultReportString2 + resultReportString3
    //
    fmt.Println(resultPrintString)

    // (6). put the execution log into results
    logger.WriteExecutionResults(resultReportString, pStart, resultsDir)

    // write the channel to executor for scheduler
    var resultsChanArray []interface{}
    resultsChanArray = append(resultsChanArray, tcName)
    resultsChanArray = append(resultsChanArray, testResult)
    resultsChanArray = append(resultsChanArray, start)
    resultsChanArray = append(resultsChanArray, end)
    resultsChanArray = append(resultsChanArray, testMessages)
    resultsChanArray = append(resultsChanArray, start_time.UnixNano())
    resultsChanArray = append(resultsChanArray, end_time.UnixNano())
    resultsChanArray = append(resultsChanArray, end_time.UnixNano() - start_time.UnixNano())
    resultsChan <- resultsChanArray

}

func Compare(tcName string, actualStatusCode string, actualHeader http.Header, actualBody string, 
        expStatus map[string]interface{}, expHeader map[string]interface{}, expBody map[string]interface{}) (string, []string) {

    // the map for mapping the string and the related funciton to call
    // fmt.Println("Compare: ", tcName)
    var testResults []bool
    var testMessages []string
    // status
    for key, _ := range expStatus {
        // fmt.Println("expStatus", key, expStatus[key])
        // call the assertion function
        testResult := assertion.CallAssertion(key, actualStatusCode, expStatus[key])
        // fmt.Println("expStatus", key, actualStatusCode, expStatus[key], reflect.TypeOf(actualStatusCode), reflect.TypeOf(expStatus[key]), testResult)
        if testResult == false {
            testMessages = append(testMessages, "Status")
            testMessages = append(testMessages, actualStatusCode)
            // testMessages = append(testMessages, expStatus[key].(string))
        }
        testResults = append(testResults, testResult)
    }

    // header
    for key, _ := range expHeader {
        bytesExpHeader, _ := json.Marshal(expHeader)
        expRes, _ := simplejson.NewJson(bytesExpHeader)
        expHeader_2, _ := expRes.Get(key).Map()
        //
        for assertionKey, _ := range expHeader_2 {
            // http.Header has structure map[key:[] ...]
            actualValue := actualHeader[key][0]
            // call the assertion function
            testResult := assertion.CallAssertion(assertionKey, actualValue, expHeader_2[assertionKey])
            // fmt.Println("expHeader_2", key, assertionKey, actualValue, actualValue, reflect.TypeOf(actualValue), reflect.TypeOf(expHeader_2[assertionKey]), testResult)
            if testResult == false {
            testMessages = append(testMessages, "Header")
            // testMessages = append(testMessages, actualHeader)
            // testMessages = append(testMessages, expStatus[key].(string))
        }
            testResults = append(testResults, testResult)
        } 
    }

    // body
   for key, _ := range expBody {
        bytesExpBody, _ := json.Marshal(expBody)
        expRes, _ := simplejson.NewJson(bytesExpBody)
        expBody_2, _ := expRes.Get(key).Map()
        //
        for assertionKey, _ := range expBody_2 {
            // if path, then value - value, otherwise, key - value
            actualValue := GetActualValueBasedOnExpKeyAndActualBody(key, actualBody)
            // call the assertion function
            testResult := assertion.CallAssertion(assertionKey, actualValue, expBody_2[assertionKey])
            // fmt.Println("expBody_2", key, assertionKey, actualValue, expBody_2[assertionKey], reflect.TypeOf(actualValue), reflect.TypeOf(expBody_2[assertionKey]), testResult)
            if testResult == false {
            testMessages = append(testMessages, "Body")
            testMessages = append(testMessages, actualBody)
            // testMessages = append(testMessages, expStatus[key].(string))
        }
            testResults = append(testResults, testResult)
        }
    }

    // default finalResults
    finalResults := "Success"

    for key := range testResults {
        if testResults[key] == false {
            finalResults = "Fail"
            // fmt.Println(tcName + " results", testResults, "final results: ", finalResults, testMessages)
            break
        }
    }
    
    return finalResults, testMessages

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

func PrepPostFormPayload(reqPayload map[string]interface{}) *strings.Reader {
    mv := reflect.ValueOf(reqPayload)
    var body *strings.Reader

    // Note, has 3 conditions: text (json), form, or multipart file upload
    data := url.Values{}
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        data.Set(k.Interface().(string), v.Interface().(string))
        }
    body = strings.NewReader(data.Encode())

    return body
}

func GetActualValueBasedOnExpKeyAndActualBody(key string, actualBody string) interface{} {
    var actualValue interface{}
    // if key starts with "$.", it represents the path, for xml, json
    // if key == "text", it is plain text, represents its valu is the whole returned body
    //
    // parse it based on the json by default, need add logic for xml, and other format
    actualRes, _ := simplejson.NewJson([]byte(actualBody))
    //
    jsonKeyList := strings.Split(key, ".")
    lastItem := jsonKeyList[(len(jsonKeyList) - 1):(len(jsonKeyList))][0]
    // fmt.Println("lastItem: ", lastItem)
    //
    // for convinence, if need to covert all the value to type *simplejson.Json???
    if jsonKeyList[0] == "$" {
        switch lastItem {
            case "Count()": {
                if len(jsonKeyList) > 2 {
                    var jsonValue *simplejson.Json
                    for _, jsonKey := range jsonKeyList[0:(len(jsonKeyList) - 1)] {
                        jsonValue = actualRes.Get(jsonKey)
                    }
                    jsonValueList, _ := jsonValue.Array()
                    actualValue = len(jsonValueList)
                } else {
                    // here to deal with the special case, like returned body is []
                    jsonValueList, _ := actualRes.Array()
                    actualValue = len(jsonValueList)
                }
                    
            }
            default: {
                for _, jsonKey := range jsonKeyList {
                    actualValue = actualRes.Get(jsonKey)
                } 
            }
        }
    } else {
        for _, jsonKey := range jsonKeyList {
            actualValue = actualRes.Get(jsonKey)
        } 
    }

    // fmt.Println("actualValue: ", actualValue)
    return actualValue
}


func WriteOutputsDataToFile(testResult string, tcJson *simplejson.Json, tcName string, tc []interface{}, actualBody string) {
    var expOutputs []interface{}
    if testResult == "Success" {
        expOutputs = utils.GetExpectedOutputsFieldsForTC(tcJson, tcName)
        
        if len(expOutputs) > 0 {
            // get the actual value from actual body based on the fields in json outputs
            var keyStr string
            var valueStr string
            for _, item := range expOutputs {
                for key, value := range item.(map[string]interface{}) {
                    // for header
                    keyStr = keyStr + "," + key
                    //
                    actualValue := GetActualValueBasedOnExpKeyAndActualBody(key, actualBody)
                    // the actualValue can be int or *simplejson.Json
                    // here need to improve to provide smart way to handle with the type - string, int, float, bool
                    str := actualValue.(*simplejson.Json).MustString()
                    // if len(str) > 0 
                    fmt.Println("expOutputs -1 : ", str, reflect.TypeOf(str), reflect.TypeOf(actualValue))
                    fmt.Println("expOutputs: ", expOutputs, key, value, actualValue, actualValue.(*simplejson.Json), tc[0], tc[4])
                }
                
            }
            jsonFileName := strings.TrimRight(filepath.Base(tc[4].(string)), ".json")
            outputsFile := filepath.Join(filepath.Dir(tc[4].(string)), jsonFileName + "_outputs.csv")
            //
            utils.GenerateFileBasedOnVarOverride(keyStr, outputsFile)
            utils.GenerateFileBasedOnVarAppend(valueStr, outputsFile)
        }
    }
}

