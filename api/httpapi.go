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
    "io"
    "net/http"     
    "net/url"  
    "strings"
    "encoding/json"

    "go4api/cmd"
    "go4api/lib/testcase"                                                                                                                               
    "go4api/assertion"
    "go4api/protocal/http"
)

type TestMessage struct {  
    FieldName interface{}
    AssertionKey  interface{}
    ExpValue interface{}
    ActualValue  interface{}
}

const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_N = "\x1b[0m"

func HttpApi(wg *sync.WaitGroup, resultsExeChan chan testcase.TestCaseExecutionInfo, pStart string, baseUrl string, 
        tcData testcase.TestCaseDataInfo, resultsDir string) {
    //
    defer wg.Done()
    //
    start_time := time.Now()
    start := start_time.String()
    //
    actualStatusCode, actualHeader, actualBody := CallHttp(baseUrl, tcData)
    //
    // (2). Expected response
    expStatus := tcData.TestCase.RespStatus()
    expHeader := tcData.TestCase.RespHeaders()
    expBody := tcData.TestCase.RespBody()

    // (3). compare
    tcName := tcData.TcName()
    testResult, TestMessages := Compare(tcName, actualStatusCode, actualHeader, actualBody, expStatus, expHeader, expBody)
    //
    end_time := time.Now()
    end := end_time.String()
    // fmt.Println(tcName + " end: ", end)

    // (4). here to generate the outputs file if the Json has "outputs" field
    WriteOutputsDataToFile(testResult, tcData, actualBody)

    // get the TestCaseExecutionInfo
    tcExecution := testcase.TestCaseExecutionInfo {
        TestCaseDataInfo: &tcData,
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
    ReportConsole(tcExecution, actualBody)

    // (6). write the channel to executor for scheduler and log
    resultsExeChan <- tcExecution
}


func CallHttp(baseUrl string, tcData testcase.TestCaseDataInfo) (int, http.Header, []byte) {
    urlStr := tcData.TestCase.UrlEncode(baseUrl)
    //
    apiMethodSelector, apiMethod, bodyText, bodyMultipart, boundary := GetPayloadInfo(tcData)
    //
    reqHeaders := make(map[string]interface{})
    reqHeaders = tcData.TestCase.ReqHeaders()
    // set the boundary to headers, if multipart
    if boundary != "" {
        reqHeaders["Content-Type"] = boundary
    }

    // < !! ----------- !! >
    // (1). Actual response
    var actualStatusCode int
    var actualHeader http.Header
    var actualBody []byte
    // 
    httpRequest := protocal.HttpRestful{}
    if apiMethodSelector == "POSTMultipart" {
        actualStatusCode, actualHeader, actualBody = httpRequest.Request(urlStr, apiMethod, reqHeaders, bodyMultipart)    
    } else {
        actualStatusCode, actualHeader, actualBody = httpRequest.Request(urlStr, apiMethod, reqHeaders, bodyText)
        }

    return actualStatusCode, actualHeader, actualBody
}


func GetPayloadInfo (tcData testcase.TestCaseDataInfo) (string, string, *strings.Reader, *bytes.Buffer, string) {
    apiMethod := tcData.TestCase.ReqMethod()
    // request payload(body)
    reqPayload := tcData.TestCase.ReqPayload()
    //
    var bodyText *strings.Reader // init body
    bodyMultipart := &bytes.Buffer{}
    boundary := ""
    //
    apiMethodSelector := apiMethod
    // Note, has 3 conditions: text (json), form, or multipart file upload
    for key, value := range reqPayload {
        // case 1: multipart upload
        if key == "filename" {
            if string(cmd.Opt.Testresource[len(cmd.Opt.Testresource) - 1]) == "/" {
                bodyMultipart, boundary, _ = PrepMultipart(cmd.Opt.Testresource + value.(string), "excel")
            } else {
                bodyMultipart, boundary, _ = PrepMultipart(cmd.Opt.Testresource + "/" + value.(string), "excel")
            }
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

    return apiMethodSelector, apiMethod, bodyText, bodyMultipart, boundary
}


func Compare(tcName string, actualStatusCode int, actualHeader http.Header, actualBody []byte, 
        expStatus map[string]interface{}, expHeader map[string]interface{}, expBody map[string]interface{}) (string, string) {
    //
    var testResults []bool
    var testMessages []TestMessage
    // status
    testResultsS, testMessagesS := CompareStatus(actualStatusCode, expStatus)
    testResults = append(testResults, testResultsS[0:]...)
    testMessages = append(testMessages, testMessagesS[0:]...)
    // headers
    testResultsH, testMessagesH := CompareHeaders(actualHeader, expHeader)
    testResults = append(testResults, testResultsH[0:]...)
    testMessages = append(testMessages, testMessagesH[0:]...)
    // body
    testResultsB, testMessagesB := CompareBody(actualBody, expBody)
    testResults = append(testResults, testResultsB[0:]...)
    testMessages = append(testMessages, testMessagesB[0:]...)

    // default finalResults
    finalResults := "Success"

    for key := range testResults {
        if testResults[key] == false {
            finalResults = "Fail"
            break
        }
    }

    testMessagesJson, _ := json.Marshal(testMessages)
    testMessagesJsonStr := string(testMessagesJson)
    
    return finalResults, testMessagesJsonStr
} 


func CompareStatus(actualStatusCode int, expStatus map[string]interface{}) ([]bool, []TestMessage) {
    var testResults []bool
    var testMessages []TestMessage
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
            testMessages = append(testMessages, msg)
        }
        testResults = append(testResults, testResult)
    }

    return testResults, testMessages
} 

func CompareHeaders(actualHeader http.Header, expHeader map[string]interface{}) ([]bool, []TestMessage) {
    var testResults []bool
    var testMessages []TestMessage
    // headers
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
                testMessages = append(testMessages, msg)
            }
            testResults = append(testResults, testResult)
        } 
    }

    return testResults, testMessages
} 

func CompareBody(actualBody []byte, expBody map[string]interface{}) ([]bool, []TestMessage) {
    var testResults []bool
    var testMessages []TestMessage
    // body
    for key, value := range expBody {
        // Note, the below statement does not work, if the key starts with $, such as $.#, maybe bug for gjson???
        expBody_sub := value.(map[string]interface{})
        for assertionKey, expValue := range expBody_sub {
            // if path, then value - value, otherwise, key - value
            actualValue := GetActualValueByJsonPath(key, actualBody)
            // check the value gotten
            if actualValue == nil {
                msg := TestMessage {
                    FieldName: key,
                    AssertionKey:  assertionKey,
                    ExpValue: expValue,
                    ActualValue: actualValue,
                }
                testMessages = append(testMessages, msg)

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
                    testMessages = append(testMessages, msg)
                }
                testResults = append(testResults, testResult)
            } 
        }
    }

    return testResults, testMessages
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


func ReportConsole (tcExecution testcase.TestCaseExecutionInfo, actualBody []byte) {
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

        fmt.Printf("\n%s%-40s%-3s%-30s%-10s%-30s%-30s%-4s%d%s\n", CLR_R, tcReportResults.TcName , tcReportResults.Priority, tcReportResults.ParentTestCase, 
            tcReportResults.TestResult, tcReportResults.JsonFilePath, tcReportResults.CsvFile, tcReportResults.CsvRow,
            tcReportResults.ActualStatusCode, CLR_N)

        if tcReportResults.MutationInfo != nil {
            fmt.Println(tcReportResults.MutationInfo)
        }
        
        // fmt.Println(tcReportResults.MutationInfo)
        fmt.Println(tcReportResults.TestMessages)
        fmt.Println(string(actualBody)[0:out_len], "...")
    } else {
        fmt.Printf("\n%s%-40s%-3s%-30s%-10s%-30s%-30s%-4s%d%s\n", CLR_G, tcReportResults.TcName, tcReportResults.Priority, tcReportResults.ParentTestCase, 
            tcReportResults.TestResult, tcReportResults.JsonFilePath, tcReportResults.CsvFile, tcReportResults.CsvRow,
            tcReportResults.ActualStatusCode, CLR_N)

        if tcReportResults.MutationInfo != nil {
            fmt.Println(tcReportResults.MutationInfo)
        }
    }
}
