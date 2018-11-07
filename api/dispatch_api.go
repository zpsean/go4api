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
    // "fmt"
    "time"
    "sync"
    // "encoding/json"

    "go4api/cmd"
    "go4api/lib/testcase" 
)

func DispatchApi(wg *sync.WaitGroup, resultsExeChan chan testcase.TestCaseExecutionInfo, baseUrl string, oTcData testcase.TestCaseDataInfo) {
    //
    defer wg.Done()
    //--- TBD: here to identify and call the builtin functions in Body, then modify the tcData
    tcData := oTcData
    if !cmd.Opt.IfMutation {
        tcData = EvaluateBuiltinFunctions(oTcData)
    }
    //
    var actualStatusCode int
    var actualHeader = make(map[string][]string)
    var actualBody []byte
    // setUp
    tcSetUpResult, setUpTestMessages := RunTcSetUp(tcData, actualStatusCode, actualHeader, actualBody)
    //
    var httpResult string
    var httpTestMessages []*testcase.TestMessage

    start_time := time.Now()
    start_str := start_time.Format("2006-01-02 15:04:05.999999999")

    if IfValidHttp(tcData) == true {
        expStatus := tcData.TestCase.RespStatus()
        expHeader := tcData.TestCase.RespHeaders()
        expBody := tcData.TestCase.RespBody()
        //
        actualStatusCode, actualHeader, actualBody = CallHttp(baseUrl, tcData)
        // (3). compare
        tcName := tcData.TcName()
        httpResult, httpTestMessages = Compare(tcName, actualStatusCode, actualHeader, actualBody, expStatus, expHeader, expBody)

        // (4). here to generate the outputs file if the Json has "outputs" field
        WriteOutputsDataToFile(httpResult, tcData, actualStatusCode, actualHeader, actualBody)
        WriteOutEnvVariables(httpResult, tcData, actualStatusCode, actualHeader, actualBody)
        WriteSession(httpResult, tcData, actualStatusCode, actualHeader, actualBody)
    } else {
        httpResult = "NoHttp"
        actualStatusCode = 999
    }
    end_time := time.Now()
    end_str := end_time.Format("2006-01-02 15:04:05.999999999")

    // tearDown
    tcTearDownResult, tearDownTestMessages := RunTcTearDown(tcData, actualStatusCode, actualHeader, actualBody)

    testResult := "Success"
    if tcSetUpResult == "Fail" || httpResult == "Fail" || tcTearDownResult == "Fail" {
        testResult = "Fail"
    }

    // get the TestCaseExecutionInfo
    tcExecution := testcase.TestCaseExecutionInfo {
        TestCaseDataInfo: &tcData,
        SetUpResult: tcSetUpResult,
        SetUpTestMessages: setUpTestMessages,
        HttpResult: httpResult,
        ActualStatusCode: actualStatusCode,
        StartTime: start_str,
        EndTime: end_str,
        HttpTestMessages: httpTestMessages,
        StartTimeUnixNano: start_time.UnixNano(),
        EndTimeUnixNano: end_time.UnixNano(),
        DurationUnixNano: end_time.UnixNano() - start_time.UnixNano(),
        ActualBody: actualBody,
        TearDownResult: tcTearDownResult,
        TearDownTestMessages: tearDownTestMessages,
        TestResult: testResult,
    }

    // (6). write the channel to executor for scheduler and log
    resultsExeChan <- tcExecution
}

func IfValidHttp (tcData testcase.TestCaseDataInfo) bool {
    ifValidHttp := true

    if tcData.TestCase.Request() == nil {
        ifValidHttp = false
    }

    return ifValidHttp
}




