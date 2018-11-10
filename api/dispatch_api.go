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
    "time"
    "sync"

    // "go4api/cmd"
    "go4api/lib/testcase" 
)

func DispatchApi(wg *sync.WaitGroup, resultsExeChan chan testcase.TestCaseExecutionInfo, baseUrl string, tcData testcase.TestCaseDataInfo) {
    // -----------
    defer wg.Done()

    tcDataStore := InitTcDataStore(tcData)
    // setUp
    // if !cmd.Opt.IfMutation {
    // }
    tcSetUpResult, setUpTestMessages := tcDataStore.RunTcSetUp()
    //
    var httpResult string
    var httpTestMessages []*testcase.TestMessage

    start_time := time.Now()
    start_str := start_time.Format("2006-01-02 15:04:05.999999999")

    if IfValidHttp(tcData) == true {
        tcDataStore.CallHttp(baseUrl)
        httpResult, httpTestMessages = tcDataStore.Compare()

        if httpResult == "Success" {
            tcDataStore.HandleHttpResultsForOut()
        }
    } else {
        httpResult = "NoHttp"
        tcDataStore.HttpActualStatusCode = 999
    }
    end_time := time.Now()
    end_str := end_time.Format("2006-01-02 15:04:05.999999999")

    // tearDown
    tcTearDownResult, tearDownTestMessages := tcDataStore.RunTcTearDown()

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
        ActualStatusCode: tcDataStore.HttpActualStatusCode,
        StartTime: start_str,
        EndTime: end_str,
        HttpTestMessages: httpTestMessages,
        StartTimeUnixNano: start_time.UnixNano(),
        EndTimeUnixNano: end_time.UnixNano(),
        DurationUnixNano: end_time.UnixNano() - start_time.UnixNano(),
        ActualBody: tcDataStore.HttpActualBody,
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


