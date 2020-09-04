/*
 * go4api - an api testing tool written in Go
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

    // "go4api/cmd"
    "go4api/lib/testcase" 
)

func DispatchApi(wg *sync.WaitGroup, resultsExeChan chan testcase.TestCaseExecutionInfo, baseUrl string, tcData *testcase.TestCaseDataInfo) {
    // -----------
    defer wg.Done()

    tcDataStore := InitTcDataStore(tcData)

    tcSetUpResult, setUpTestMessages := tcDataStore.RunTcSetUp()
    //
    start_time := time.Now()
    start_str := start_time.Format("2006-01-02 15:04:05.999999999")
    
    var httpResult string
    var httpTestMessages []*testcase.TestMessage
    if IfValidHttp(tcData) == true {
        httpResult, httpTestMessages = tcDataStore.RunHttp(baseUrl)

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
        TestCaseDataInfo:  tcDataStore.TcData,
        SetUpResult:       tcSetUpResult,
        SetUpTestMessages: setUpTestMessages,
        HttpResult:        httpResult,
        ActualStatusCode:  tcDataStore.HttpActualStatusCode,
        StartTime:         start_str,
        EndTime:           end_str,
        HttpTestMessages:  httpTestMessages,
        StartTimeUnixNano: start_time.UnixNano(),
        EndTimeUnixNano:   end_time.UnixNano(),
        DurationUnixNano:  end_time.UnixNano() - start_time.UnixNano(),
        ActualBody:        tcDataStore.HttpActualBody,
        ActualHeader:      tcDataStore.HttpActualHeader,
        HttpUrl:           tcDataStore.HttpUrl,
        TearDownResult:    tcTearDownResult,
        TearDownTestMessages: tearDownTestMessages,
        TestResult:           testResult,
        LocalVariables:       tcDataStore.TcLocalVariables,
    }

    // (6). write the channel to executor for scheduler and log
    resultsExeChan <- tcExecution
}

func IfValidHttp (tcData *testcase.TestCaseDataInfo) bool {
    ifValidHttp := true

    if tcData.TestCase.Request() == nil {
        ifValidHttp = false
    }

    return ifValidHttp
}

// Note: for each SetUp, TesrDown, it may have more than one Command (including sql)
// for each Command, it may have more than one assertion
func (tcDataStore *TcDataStore) RunTcSetUp () (string, [][]*testcase.TestMessage) {
    var finalResults string
    var finalTestMessages = [][]*testcase.TestMessage{}

    cmdGroup := tcDataStore.TcData.TestCase.SetUp()

    if len(cmdGroup) > 0 {
        tcDataStore.CmdGroupLength = len(cmdGroup)
        tcDataStore.CmdSection = "setUp"
        finalResults, finalTestMessages = tcDataStore.CommandGroup(cmdGroup)
    } else {
        finalResults = "NoSetUp"
    }

    return finalResults, finalTestMessages
}

func (tcDataStore *TcDataStore) RunTcTearDown () (string, [][]*testcase.TestMessage) {
    var finalResults string
    var finalTestMessages = [][]*testcase.TestMessage{}

    cmdGroup := tcDataStore.TcData.TestCase.TearDown()

    if len(cmdGroup) > 0 {
        tcDataStore.CmdGroupLength = len(cmdGroup)
        tcDataStore.CmdSection = "tearDown"
        finalResults, finalTestMessages = tcDataStore.CommandGroup(cmdGroup)
    } else {
        finalResults = "NoTearDown"
    }

    return finalResults, finalTestMessages
}

