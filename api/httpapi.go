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
    "go4api/types"                                                                                                                                  
    "go4api/utils"
    "go4api/assertion"
    "go4api/protocal/http"
    "reflect"
    "strings"
    "encoding/json"
    gjson "github.com/tidwall/gjson"
    "strconv"
)

type TestMessage struct {  
    FieldName interface{}
    AssertionKey  interface{}
    ExpValue interface{}
    ActualValue  interface{}
}


func HttpApi(wg *sync.WaitGroup, resultsChan chan types.TcRunResults, options map[string]string, pStart string, baseUrl string, 
        tc []interface{}, resultsDir string) {
    //
    defer wg.Done()
    //
    tcName := tc[0].(string)
    parentTestCase := tc[2].(string)
    tcJson := tc[3].(string)
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
    // fmt.Println("bodyText: ", reqPayload)
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
    var actualStatusCode int
    var actualHeader http.Header
    var actualBody []byte
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
    testResult, TestMessages := Compare(tcName, actualStatusCode, actualHeader, actualBody, expStatus, expHeader, expBody)
    //
    end_time := time.Now()
    end := end_time.String()
    // fmt.Println(tcName + " end: ", end)

    // (4). here to generate the outputs file if the Json has "outputs" field
    WriteOutputsDataToFile(testResult, tcJson, tcName, tc, actualBody)

    // (5). print to console
    resultPrintString := ""
    csvFileBase := ""
    TestMessages_len := 0
    // Note: if csvFile does not exist, the filepath.Base(csvFile) = ".", need to remove
    if filepath.Base(csvFile) == "." {
        csvFileBase = ""
    } else {
        csvFileBase = filepath.Base(csvFile)
    }

    if len(TestMessages) > 300 {
        TestMessages_len = 300
    } else {
        TestMessages_len = len(TestMessages)
    }


    resultPrintString1 := tcName + "," + strconv.Itoa(actualStatusCode) + "," + filepath.Base(jsonFile) + "," + csvFileBase + "," + rowCsv
    resultPrintString = resultPrintString1 + "," + testResult + "," + TestMessages[0:TestMessages_len]
    //
    fmt.Println(resultPrintString)


    // (6). write the channel to executor for scheduler and log
    // here can refactor to struct => done
    tcRunRes := types.TcRunResults {
        TcName : tcName,
        ParentTestCase : parentTestCase,
        TestResult : testResult,
        ActualStatusCode : strconv.Itoa(actualStatusCode),
        JsonFile_Base : filepath.Base(jsonFile),
        CsvFileBase : csvFileBase,
        RowCsv : rowCsv,
        Start : start,
        End : end,
        TestMessages : TestMessages,
        Start_time_UnixNano : start_time.UnixNano(),
        End_time_UnixNano : end_time.UnixNano(),
        Duration_UnixNano : end_time.UnixNano() - start_time.UnixNano(),
    }

    resultsChan <- tcRunRes

}

func Compare(tcName string, actualStatusCode int, actualHeader http.Header, actualBody []byte, 
        expStatusJson gjson.Result, expHeaderJson gjson.Result, expBodyJson gjson.Result) (string, string) {

    // the map for mapping the string and the related funciton to call
    fmt.Println("Compare: ", tcName)
    var testResults []bool

    var TestMessages []TestMessage

    // status
    expStatus := expStatusJson.Map()
    for assertionKey, expValue := range expStatus {
        // call the assertion function
        testResult := assertion.CallAssertion(assertionKey, actualStatusCode, expValue.Value())
        // fmt.Println("--> expStatus", assertionKey, actualStatusCode, expValue, reflect.TypeOf(actualStatusCode), reflect.TypeOf(expValue), testResult)
        if testResult == false {
            msg := TestMessage {
                    FieldName: "status",
                    AssertionKey:  assertionKey,
                    ExpValue: expValue.Value(),
                    ActualValue: actualStatusCode,
                }
            TestMessages = append(TestMessages, msg)
        }
        testResults = append(testResults, testResult)
    }

    // header
    // http.Header => map[string][]string
    expHeader := expHeaderJson.Map()
    for key, _ := range expHeader {
        expHeader_sub := expHeaderJson.Get(key).Map()
        //
        for assertionKey, expValue := range expHeader_sub {
            // as the http.Header has structure, so that here need to assert if the expValue in []string
            actualValue := strings.Join(actualHeader[key], ",")
            // call the assertion function
            testResult := assertion.CallAssertion(assertionKey, actualValue, expValue.Value())
            // fmt.Println("-> expHeader_sub", key, assertionKey, actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue.Value()), testResult)
            if testResult == false {
                msg := TestMessage {
                    FieldName: key,
                    AssertionKey:  assertionKey,
                    ExpValue: expValue.Value(),
                    ActualValue: actualValue,
                }
                TestMessages = append(TestMessages, msg)
            }
            testResults = append(testResults, testResult)
        } 
    }

    // body
    expBody := expBodyJson.Map()
    for key, value := range expBody {
        // Note, the below statement does not work, if the key starts with $, such as $.#, maybe bug for gjson???
        // expBody_sub := expBodyJson.Get(key).Map()
        // However, need to use value directly
        expBody_sub := value.Map()
        for assertionKey, expValue := range expBody_sub {
            // if path, then value - value, otherwise, key - value
            actualValue := GetActualValueBasedOnExpKeyAndActualBody(key, actualBody)
            // check the value gotten
            if actualValue == nil {
                msg := TestMessage {
                    FieldName: key,
                    AssertionKey:  assertionKey,
                    ExpValue: expValue.Value(),
                    ActualValue: actualValue,
                }
                TestMessages = append(TestMessages, msg)

                testResults = append(testResults, false)
            } else {
                // call the assertion function
                testResult := assertion.CallAssertion(assertionKey, actualValue, expValue.Value())
                // fmt.Println("-> expBody_sub", key, assertionKey, actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue), testResult)
                if testResult == false {
                    msg := TestMessage {
                        FieldName: key,
                        AssertionKey:  assertionKey,
                        ExpValue: expValue.Value(),
                        ActualValue: actualValue,
                    }
                    TestMessages = append(TestMessages, msg)
                }
                testResults = append(testResults, testResult)
            } 
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

    
    TestMessagesJson, _ := json.Marshal(TestMessages)
    TestMessagesJsonStr := string(TestMessagesJson)
    
    return finalResults, TestMessagesJsonStr

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

func GetActualValueBasedOnExpKeyAndActualBody(key string, actualBody []byte) interface{} {
    var actualValue interface{}
    // if key starts with "$.", it represents the path, for xml, json
    // if key == "text", it is plain text, represents its valu is the whole returned body
    //
    // parse it based on the json by default, need add logic for xml, and other format

    if key[0:2] == "$." {
        value := gjson.Get(string(actualBody), key[2:])
        actualValue = value.Value()
    } else {
        value := gjson.Get(string(actualBody), key)
        actualValue = value.Value()
    }

    // fmt.Println("actualValue: ", actualValue)
    return actualValue
}


func WriteOutputsDataToFile(testResult string, tcJson interface{}, tcName string, tc []interface{}, actualBody []byte) {
    var expOutputs []gjson.Result
    if testResult == "Success" {
        expOutputs = utils.GetExpectedOutputsFieldsForTC(tcJson, tcName)
        
        if len(expOutputs) > 0 {
            // get the actual value from actual body based on the fields in json outputs
            var keyStrList []string
            var valueStrList []interface{}
            //
            for _, item := range expOutputs {
                // item is {}
                for key, value := range item.Map() {
                    // for csv header
                    fmt.Println("key, value: ", key, value)
                    keyStrList = append(keyStrList, key)
                    //
                    // actualValue := GetActualValueBasedOnExpKeyAndActualBody(value.(gjson.Result).Value(), actualBody)
                    actualValue := "aa"
                    valueStrList = append(valueStrList, actualValue)
                }
            }
            // get the full path of outputsfile
            jsonFileName := strings.TrimRight(filepath.Base(tc[4].(string)), ".json")
            outputsFile := filepath.Join(filepath.Dir(tc[4].(string)), jsonFileName + "_outputs.csv")
            // write csv header
            utils.GenerateFileBasedOnVarOverride(strings.Join(keyStrList, ",") + "\n", outputsFile)

            valueStr := "aa,b,c\n"
            utils.GenerateFileBasedOnVarAppend(valueStr, outputsFile)
        }
    }
}

