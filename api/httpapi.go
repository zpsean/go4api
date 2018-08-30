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
    "go4api/testcase"                                                                                                                               
    "go4api/utils"
    "go4api/assertion"
    "go4api/protocal/http"
    "strings"
    "encoding/json"
    gjson "github.com/tidwall/gjson"
)

type TestMessage struct {  
    FieldName interface{}
    AssertionKey  interface{}
    ExpValue interface{}
    ActualValue  interface{}
}


func HttpApi(wg *sync.WaitGroup, resultsExeChan chan testcase.TestCaseExecutionInfo, options map[string]string, pStart string, baseUrl string, 
        tcData testcase.TestCaseDataInfo, resultsDir string) {
    //
    defer wg.Done()
    //
    tcName := tcData.TcName()
    //
    start_time := time.Now()
    start := start_time.String()
    //
    apiMethod := tcData.TestCase.ReqMethod()
    urlStr := tcData.TestCase.UrlEncode(baseUrl)
    //
    // request payload(body)
    var reqPayload map[string]interface{}
    reqPayload = tcData.TestCase.ReqPayload()
    //
    var bodyText *strings.Reader // init body
    bodyMultipart := &bytes.Buffer{}
    var boundary string
    //
    apiMethodSelector := apiMethod
    // Note, has 3 conditions: text (json), form, or multipart file upload
    for key, value := range reqPayload {
        // case 1: multipart upload
        if key == "filename" {
            bodyMultipart, boundary, _ = PrepMultipart(options["testhome"] + "/testresource/" + value.(string), "excel")
            apiMethodSelector = "POSTMultipart"
            break
        }
        // case 2: normal json
        if key == "text" {
            bodyText = PrepPostPayload(reqPayload)
            break
        }
        // case 3: if Post, and the key does not have filename, text, then it would be PostForm
        bodyText = PrepPostFormPayload(reqPayload)
    }
    //
    var reqHeaders map[string]interface{}
    reqHeaders = tcData.TestCase.ReqHeaders()
    // set the boundary to headers, if multipart
    if boundary != "" {
        reqHeaders["Content-Type"] = boundary
    }
    // fmt.Println(tcName + " boundary: ", boundary)
   

    // < !! ----------- !! >
    // the map for mapping the string and the related funciton to call
    funcs := map[string]interface{} {
        "GET": protocal.HttpGet,
        "POST": protocal.HttpPost,
        "POSTForm": protocal.HttpPostForm,
        "POSTMultipart": protocal.HttpPostMultipart,
        "PUT": protocal.HttpPut,
    }

    // (1). Actual response
    var actualStatusCode int
    var actualHeader http.Header
    var actualBody []byte
    // fmt.Println("----- to start call the http ------")
    if apiMethodSelector == "POSTMultipart" {
        actualStatusCode, actualHeader, actualBody = protocal.CallHttpMethod(funcs, apiMethodSelector, urlStr, apiMethod, reqHeaders, bodyMultipart)    
    } else {
        actualStatusCode, actualHeader, actualBody = protocal.CallHttpMethod(funcs, apiMethodSelector, urlStr, apiMethod, reqHeaders, bodyText)
        }
    //
    // (2). Expected response
    expStatus := tcData.TestCase.RespStatus()
    expHeader := tcData.TestCase.RespHeaders()
    expBody := tcData.TestCase.RespBody()
    // fmt.Println(actualStatusCode, actualHeader, actualBody)

    // (3). compare
    testResult, TestMessages := Compare(tcName, actualStatusCode, actualHeader, actualBody, expStatus, expHeader, expBody)
    //
    end_time := time.Now()
    end := end_time.String()
    // fmt.Println(tcName + " end: ", end)

    // (4). here to generate the outputs file if the Json has "outputs" field
    WriteOutputsDataToFile(testResult, tcData, actualBody)


    // get the TestCaseExecutionInfo
    tcExecution := testcase.TestCaseExecutionInfo {
        TestCaseDataInfo: tcData,
        TestResult: testResult,
        ActualStatusCode: actualStatusCode,
        StartTime: start,
        EndTime: end,
        TestMessages: TestMessages,
        StartTimeUnixNano: start_time.UnixNano(),
        EndTimeUnixNano: end_time.UnixNano(),
        DurationUnixNano: end_time.UnixNano() - start_time.UnixNano(),
    }


    // (5). print to console
    tcReportResults := tcExecution.TcConsoleResults()
    // repJson, _ := json.Marshal(tcReportResults)

    if tcReportResults.TestResult == "Fail" {
        length := len(string(actualBody))
        out_len := 0
        if length > 300 {
            out_len = 300
        } else {
            out_len = length
        }

        fmt.Printf("\n%-40s%-3s%-30s%-10s%-30s%-30s%-4s%d\n", tcReportResults.TcName, tcReportResults.Priority, tcReportResults.ParentTestCase, 
            tcReportResults.TestResult, tcReportResults.JsonFilePath, tcReportResults.CsvFile, tcReportResults.CsvRow,
            tcReportResults.ActualStatusCode)

        if tcReportResults.MutationInfo != nil {
            fmt.Println(tcReportResults.MutationInfo)
        }
        
        // fmt.Println(tcReportResults.MutationInfo)
        fmt.Println(tcReportResults.TestMessages)
        fmt.Println(string(actualBody)[0:out_len], "...")
    } else {
        // fmt.Println(string(repJson))

        // type TcConsoleResults struct { 
        //     TcName string 
        //     Priority string
        //     ParentTestCase string
        //     JsonFilePath string
        //     CsvFile string
        //     CsvRow string
        //     TestResult string  // Ready, Running, Success, Fail, ParentReady, ParentRunning, ParentFailed
        //     ActualStatusCode int
        //     TestMessages string
        // }

        fmt.Printf("\n%-40s%-3s%-30s%-10s%-30s%-30s%-4s%d\n", tcReportResults.TcName, tcReportResults.Priority, tcReportResults.ParentTestCase, 
            tcReportResults.TestResult, tcReportResults.JsonFilePath, tcReportResults.CsvFile, tcReportResults.CsvRow,
            tcReportResults.ActualStatusCode)

        if tcReportResults.MutationInfo != nil {
            fmt.Println(tcReportResults.MutationInfo)
        }
    }
    
    // (6). write the channel to executor for scheduler and log
    // here can refactor to struct => done
    
    resultsExeChan <- tcExecution

}

func Compare(tcName string, actualStatusCode int, actualHeader http.Header, actualBody []byte, 
        expStatus map[string]interface{}, expHeader map[string]interface{}, expBody map[string]interface{}) (string, string) {

    // the map for mapping the string and the related funciton to call
    // fmt.Println("Compare: ", tcName)
    var testResults []bool

    var TestMessages []TestMessage

    // status
    for assertionKey, expValue := range expStatus {
        // call the assertion function
        testResult := assertion.CallAssertion(assertionKey, actualStatusCode, expValue)
        // fmt.Println("--> expStatus", assertionKey, actualStatusCode, expValue, reflect.TypeOf(actualStatusCode), reflect.TypeOf(expValue), testResult)
        if testResult == false {
            msg := TestMessage {
                    FieldName: "status",
                    AssertionKey:  assertionKey,
                    ExpValue: expValue,
                    ActualValue: actualStatusCode,
                }
            TestMessages = append(TestMessages, msg)
        }
        testResults = append(testResults, testResult)
    }

    // header
    // http.Header => map[string][]string
    for key, value := range expHeader {
        expHeader_sub := value.(map[string]interface{})
        //
        for assertionKey, expValue := range expHeader_sub {
            // as the http.Header has structure, so that here need to assert if the expValue in []string
            actualValue := strings.Join(actualHeader[key], ",")
            // call the assertion function
            testResult := assertion.CallAssertion(assertionKey, actualValue, expValue)
            // fmt.Println("-> expHeader_sub", key, assertionKey, actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue.Value()), testResult)
            if testResult == false {
                msg := TestMessage {
                    FieldName: key,
                    AssertionKey:  assertionKey,
                    ExpValue: expValue,
                    ActualValue: actualValue,
                }
                TestMessages = append(TestMessages, msg)
            }
            testResults = append(testResults, testResult)
        } 
    }

    // body
    for key, value := range expBody {
        // Note, the below statement does not work, if the key starts with $, such as $.#, maybe bug for gjson???
        // expBody_sub := expBodyJson.Get(key).Map()
        // However, need to use value directly
        expBody_sub := value.(map[string]interface{})
        for assertionKey, expValue := range expBody_sub {
            // if path, then value - value, otherwise, key - value
            actualValue := GetActualValueBasedOnExpKeyAndActualBody(key, actualBody)
            // check the value gotten
            if actualValue == nil {
                msg := TestMessage {
                    FieldName: key,
                    AssertionKey:  assertionKey,
                    ExpValue: expValue,
                    ActualValue: actualValue,
                }
                TestMessages = append(TestMessages, msg)

                testResults = append(testResults, false)
            } else {
                // call the assertion function
                testResult := assertion.CallAssertion(assertionKey, actualValue, expValue)
                // fmt.Println("-> expBody_sub", key, assertionKey, actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue), testResult)
                if testResult == false {
                    msg := TestMessage {
                        FieldName: key,
                        AssertionKey:  assertionKey,
                        ExpValue: expValue,
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
    var body *strings.Reader

    for key, value := range reqPayload {
        if key == "text" {
            repJson, _ := json.Marshal(value)
            body = strings.NewReader(string(repJson))
            break
        }
    }

    return body
}

func PrepPostFormPayload(reqPayload map[string]interface{}) *strings.Reader {
    var body *strings.Reader

    // Note, has 3 conditions: text (json), form, or multipart file upload
    data := url.Values{}
    for key, value := range reqPayload {
        // value (type interface {}) as type string in argument to data.Set: need type assertion
        data.Set(key, fmt.Sprint(value))
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


func WriteOutputsDataToFile(testResult string, tcData testcase.TestCaseDataInfo, actualBody []byte) {
    var expOutputs []interface{}

    if testResult == "Success" {
        expOutputs = tcData.TestCase.Outputs()
        if len(expOutputs) > 0 {
            // get the actual value from actual body based on the fields in json outputs
            var keyStrList []string
            var valueStrList []string
            //
            for _, itemMap := range expOutputs {
                // item is {}
                for key, value := range itemMap.(map[string]interface{}) {
                    // for csv header
                    keyStrList = append(keyStrList, key)
                    //
                    if fmt.Sprint(value)[0:2] == "$." {
                        actualValue := GetActualValueBasedOnExpKeyAndActualBody(fmt.Sprint(value), actualBody)
                        if actualValue == nil {
                            valueStrList = append(valueStrList, "")
                        } else {
                            valueStrList = append(valueStrList, fmt.Sprint(actualValue))
                        }
                    } else {
                        valueStrList = append(valueStrList, fmt.Sprint(value))
                    }
                }
            }
            // get the full path of outputsfile
            jsonFileName := strings.TrimRight(filepath.Base(tcData.JsonFilePath), ".json")
            outputsFile := filepath.Join(filepath.Dir(tcData.JsonFilePath), jsonFileName + "_outputs.csv")
            // write csv header
            utils.GenerateFileBasedOnVarOverride(strings.Join(keyStrList, ",") + "\n", outputsFile)

            // write csv data
            utils.GenerateFileBasedOnVarAppend(strings.Join(valueStrList, ",") + "\n", outputsFile)
        }
    }
}

