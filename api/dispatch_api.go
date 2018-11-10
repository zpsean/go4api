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
    "strings"
    "encoding/json"

    "go4api/cmd"
    "go4api/lib/testcase" 

    sjson "github.com/tidwall/sjson"
)

type TcDataStore struct {
    TcData testcase.TestCaseDataInfo

    TcLocalStore map[string]interface{}
    SetUpStore []map[string]interface{}

    HttpExpStatus map[string]interface{}
    HttpExpHeader map[string]interface{}
    HttpExpBody map[string]interface{}
    HttpActualStatusCode int
    HttpActualHeader map[string][]string
    HttpActualBody []byte

    HttpStore map[string]interface{}
    TearDownStore []map[string]interface{}
}

func InitTcDataStore (tcData testcase.TestCaseDataInfo) *TcDataStore {
    tcDataStore := &TcDataStore {
        tcData,
        map[string]interface{}{},
        []map[string]interface{}{},

        map[string]interface{}{},
        map[string]interface{}{},
        map[string]interface{}{},
        -1,
        map[string][]string{},
        []byte{},

        map[string]interface{}{},
        []map[string]interface{}{},
    }
    // aa, _ := json.Marshal(tcData)
    // fmt.Println(string(aa))
    
    return tcDataStore
}

func DispatchApi(wg *sync.WaitGroup, resultsExeChan chan testcase.TestCaseExecutionInfo, baseUrl string, tcData testcase.TestCaseDataInfo) {
    //
    defer wg.Done()

    tcDataStore := InitTcDataStore(tcData)
    // setUp
    if !cmd.Opt.IfMutation {
        // Stage 1: deal with Setup().Cmd before run
        // tcDataStore.RenderSetUpCmdVariables()
        // tcDataStore.EvaluateSetUpCmdBuiltinFunctions()
        // bb, _ := json.Marshal(tcDataStore.TcData)
        // fmt.Println("bb>>>>>>>>: ", string(bb))
    }
    // Stage 2: deal with the Setup().Assertion
    tcSetUpResult, setUpTestMessages := tcDataStore.RunTcSetUp()
    if !cmd.Opt.IfMutation {
        // Stage 3: deal with the Setup().Out* after run

        // tcDataStore.RenderSetUpResultsVariables()
        // tcDataStore.EvaluateSetUpResultsBuiltinFunctions()

        // tcDataStore.WriteOutEnvVariables(tcSetUpResult)
        // tcDataStore.WriteSession(tcSetUpResult)
    }
    //
    var httpResult string
    var httpTestMessages []*testcase.TestMessage

    start_time := time.Now()
    start_str := start_time.Format("2006-01-02 15:04:05.999999999")

    if IfValidHttp(tcData) == true {
        tcDataStore.CallHttp(baseUrl)

        httpResult, httpTestMessages = tcDataStore.Compare()

        // (4). here to generate the outputs file if the Json has "outputs" field
        // tcDataStore.WriteOutputsDataToFile(httpResult)
        // tcDataStore.WriteOutEnvVariables(httpResult)
        // tcDataStore.WriteSession(httpResult)
    } else {
        httpResult = "NoHttp"
        tcDataStore.HttpActualStatusCode = 999
    }
    end_time := time.Now()
    end_str := end_time.Format("2006-01-02 15:04:05.999999999")

    // tearDown
    if !cmd.Opt.IfMutation {
        tcDataStore.EvaluateTearDownSectionBuiltinFunctions()
        tcData = tcDataStore.TcData
        // fmt.Println(tcData)
    }
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


func (tcDataStore *TcDataStore) RenderSetUpCmdVariables () {
    dataFeeder := make(map[string]interface{})

    tcSetup := tcDataStore.TcData.TestCase.SetUp()
    for i, _ := range tcSetup {
        cmdStr := tcSetup[i].Cmd

        for key, value := range dataFeeder{
            cmdStr = strings.Replace(cmdStr, "${" + key + "}", fmt.Sprint(value), -1)
        }
        tcSetup[i].Cmd = cmdStr
    }
} 

func (tcDataStore *TcDataStore) RenderSetUpResultsVariables () {
    dataFeeder := make(map[string]interface{})

    tcSetup := tcDataStore.TcData.TestCase.SetUp()
    for i, _ := range tcSetup {
        cmdStr := tcSetup[i].Cmd

        for key, value := range dataFeeder {
            cmdStr = strings.Replace(cmdStr, "${" + key + "}", fmt.Sprint(value), -1)
        }
        tcSetup[i].Cmd = cmdStr
    }
} 

func (tcDataStore *TcDataStore) EvaluateSetUpCmdBuiltinFunctions () {
    var resTcData testcase.TestCaseDataInfo
    var tcTempSetup []*testcase.CommandDetails

    tcSetup := tcDataStore.TcData.TestCase.SetUp()
    if tcSetup != nil {
        if len(tcSetup) > 0 {
            jsonStr := EvaluateBuiltinFunctions(tcSetup)
            json.Unmarshal([]byte(jsonStr), &tcTempSetup)
   
            tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
            tcDataJson := string(tcDataJsonBytes)
            path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + ".setUp"

            tcDataJson, _  = sjson.Set(tcDataJson, path, tcTempSetup)

            json.Unmarshal([]byte(tcDataJson), &resTcData)

            tcDataStore.TcData = resTcData
        }
    }
}

func (tcDataStore *TcDataStore) RenderHttpSectionVariables () {

} 

func (tcDataStore *TcDataStore) EvaluateHttpSectionBuiltinFunctions () {

} 


func (tcDataStore *TcDataStore) RenderTearDownSectionVariables () {

} 

func (tcDataStore *TcDataStore) EvaluateTearDownSectionBuiltinFunctions () {
    var resTcData testcase.TestCaseDataInfo
    var tcTempTearDown []*testcase.CommandDetails

    tcTearDown := tcDataStore.TcData.TestCase.SetUp()
    if tcTearDown != nil {
        if len(tcTearDown) > 0 {
            jsonStr := EvaluateBuiltinFunctions(tcTearDown)
            json.Unmarshal([]byte(jsonStr), &tcTempTearDown)
   
            tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
            tcDataJson := string(tcDataJsonBytes)
            path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + ".tearDown"

            tcDataJson, _  = sjson.Set(tcDataJson, path, tcTempTearDown)

            json.Unmarshal([]byte(tcDataJson), &resTcData)

            tcDataStore.TcData = resTcData
        }
    }
} 
