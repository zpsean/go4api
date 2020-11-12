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
    "fmt"
    "strings" 
    "time"
    "strconv"
    "encoding/json"
    "reflect"

    "go4api/cmd"
    "go4api/lib/testcase"
    "go4api/utils"

    gjson "github.com/tidwall/gjson"
)

func (tcDataStore *TcDataStore) CommandGroup (cmdGroupOrigin []*testcase.CommandDetails) (string, [][]*testcase.TestMessage) {
    finalResults := "Success"
    var cmdsResults []bool
    var finalTestMessages [][]*testcase.TestMessage

    for i := 0; i < tcDataStore.CmdGroupLength; i ++ {
        var sResults []bool
        var sMessages [][]*testcase.TestMessage

        cmdType := cmdGroupOrigin[i].CmdType
        lc := strings.ToLower(cmdType)
        switch lc {
        case "sql", "mysql", "postgres", "postgresql":
            sResults, sMessages = tcDataStore.HandleSqlCmd(lc, i)

            cmdsResults = append(cmdsResults, sResults[0:]...)
            finalTestMessages = append(finalTestMessages, sMessages[0:]...)
        case "redis":
            sResults, sMessages = tcDataStore.HandleRedisCmd(i)

            cmdsResults = append(cmdsResults, sResults[0:]...)
            finalTestMessages = append(finalTestMessages, sMessages[0:]...)
        case "mongodb":
            sResults, sMessages = tcDataStore.HandleMongoDBCmd(i)

            cmdsResults = append(cmdsResults, sResults[0:]...)
            finalTestMessages = append(finalTestMessages, sMessages[0:]...)
        case "init":
            sResults, sMessages = tcDataStore.HandleInitCmd(i)

            cmdsResults = append(cmdsResults, sResults[0:]...)
            finalTestMessages = append(finalTestMessages, sMessages[0:]...)
        case "jsonfile":
            sResults, sMessages = tcDataStore.HandleJsonFile(i)

            cmdsResults = append(cmdsResults, sResults[0:]...)
            finalTestMessages = append(finalTestMessages, sMessages[0:]...)
        default:
            fmt.Println("!! warning, command ", cmdType, " can not be recognized.")
        }
    }

    for i, _ := range cmdsResults {
        if cmdsResults[i] == false {
            finalResults = "Fail"
            break
        }
    }

    return finalResults, finalTestMessages
}

// file
func (tcDataStore *TcDataStore) HandleJsonFile (i int) ([]bool, [][]*testcase.TestMessage) {
    var sResults []bool
    var sMessages [][]*testcase.TestMessage

    // cmd is always ""
    // cmdStr := ""

    tcDataJsonB, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonB)

    cmdTgtJsonFile := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmdSource"
    tgtJsonFile := gjson.Get(tcDataJson, cmdTgtJsonFile).String()

    tgtJsonFileFullPath := ""
    // check if tgtJsonFile is absolute path
    if len(tgtJsonFile) > 0 {
        if tgtJsonFile[0:1] == "/" {
            tgtJsonFileFullPath = tgtJsonFile
        } else {
            testResourcePath := cmd.Opt.Testresource

            if string(testResourcePath[len(testResourcePath) - 1]) != "/" {
                tgtJsonFileFullPath = testResourcePath + "/" + tgtJsonFile
            }
        }
    }

    jsonStr := utils.GetJsonFromFile(tgtJsonFileFullPath)

    // as no cmd is executed, the CmdExecStatus is always "cmdSuccess"
    tcDataStore.CmdType = "jsonFile"
    tcDataStore.CmdExecStatus = "cmdSuccess"
    tcDataStore.CmdAffectedCount = -1
    tcDataStore.CmdResults = jsonStr

    sResults, sMessages = tcDataStore.HandleSingleCmdResult(i)

    return sResults, sMessages
}


//mysql
func (tcDataStore *TcDataStore) HandleSqlCmd (lc string, i int) ([]bool, [][]*testcase.TestMessage) {
    var sResults []bool
    var sMessages [][]*testcase.TestMessage

    cmdStrPath := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmd"
    tcDataStore.PrepEmbeddedFunctions(cmdStrPath)

    tcDataJsonB, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonB)

    cmdStr := gjson.Get(tcDataJson, cmdStrPath).String()

    cmdTgtDb := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmdSource"
    tgtDb := gjson.Get(tcDataJson, cmdTgtDb).String()
    // init
    // tcDataStore.CmdType = "sql"
    tcDataStore.CmdExecStatus = ""
    tcDataStore.CmdAffectedCount = -1
    tcDataStore.CmdResults = -1

    // call sql
    if len(tgtDb) == 0 {
        fmt.Println("No target db provided, default to master")
        tgtDb = "master"
    }
    cmdAffectedCount := -1
    var cmdResults []map[string]interface{}
    cmdExecStatus    := ""
    //
    switch lc {
    case "sql", "mysql":
        tcDataStore.CmdType = "mysql"
        cmdAffectedCount, _, cmdResults, cmdExecStatus = RunSql(tgtDb, cmdStr)
    case "postgres", "postgresql":
        tcDataStore.CmdType = "postgresql"
        cmdAffectedCount, _, cmdResults, cmdExecStatus = RunPgSql(tgtDb, cmdStr)
    }
    tcDataStore.CmdExecStatus = cmdExecStatus
    tcDataStore.CmdAffectedCount = cmdAffectedCount
    tcDataStore.CmdResults = cmdResults

    sResults, sMessages = tcDataStore.HandleSingleCmdResult(i)

    return sResults, sMessages
}

// redis
func (tcDataStore *TcDataStore) HandleRedisCmd (i int) ([]bool, [][]*testcase.TestMessage) {
    var sResults []bool
    var sMessages [][]*testcase.TestMessage

    var cmdStr, cmdKey, cmdValue string

    cmdStrPath := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmd"
    tcDataStore.PrepEmbeddedFunctions(cmdStrPath)

    tcDataJsonB, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonB)

    // cmdTgtDb := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmdSource"
    // tgtDb := gjson.Get(tcDataJson, cmdTgtDb).String()

    cmdS := gjson.Get(tcDataJson, cmdStrPath).String()

    var mm []string
    m := strings.Split(cmdS, " ")
    for _, k := range m {
        if len(k) > 0 {
            mm = append(mm, k)
        }
    }

    switch len(mm) {
    case 1:
        cmdStr = mm[0]
    case 2:
        cmdStr = mm[0]
        cmdKey = mm[1]
    case 3:
        cmdStr = mm[0]
        cmdKey = mm[1]
        cmdValue = mm[2]
    }

    // init
    tcDataStore.CmdType = "redis"
    tcDataStore.CmdExecStatus = ""
    tcDataStore.CmdAffectedCount = -1
    tcDataStore.CmdResults = -1

    cmdAffectedCount, cmdResults, cmdExecStatus := RunRedis(cmdStr, cmdKey, cmdValue)
    
    tcDataStore.CmdExecStatus = cmdExecStatus
    tcDataStore.CmdAffectedCount = cmdAffectedCount
    tcDataStore.CmdResults = cmdResults

    sResults, sMessages = tcDataStore.HandleSingleCmdResult(i)

    return sResults, sMessages
}

// mongodb
func (tcDataStore *TcDataStore) HandleMongoDBCmd (i int) ([]bool, [][]*testcase.TestMessage) {
    var sResults []bool
    var sMessages [][]*testcase.TestMessage

    cmdStrPath := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmd"
    tcDataStore.PrepEmbeddedFunctions(cmdStrPath)

    tcDataJsonB, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonB)

    cmdStr := gjson.Get(tcDataJson, cmdStrPath).String()

    // init
    tcDataStore.CmdType = "mongodb"
    tcDataStore.CmdExecStatus = ""
    tcDataStore.CmdAffectedCount = -1
    tcDataStore.CmdResults = -1

    cmdAffectedCount, cmdResults, cmdExecStatus := RunMongoDB(cmdStr)
    
    tcDataStore.CmdExecStatus = cmdExecStatus
    tcDataStore.CmdAffectedCount = cmdAffectedCount
    tcDataStore.CmdResults = cmdResults

    sResults, sMessages = tcDataStore.HandleSingleCmdResult(i)

    return sResults, sMessages
}

// init
func (tcDataStore *TcDataStore) HandleInitCmd (i int) ([]bool, [][]*testcase.TestMessage) {
    var sResults []bool
    var sMessages [][]*testcase.TestMessage

    cmdStrPath := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmd"
    tcDataStore.PrepEmbeddedFunctions(cmdStrPath)

    tcDataJsonB, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonB)

    cmdStr := gjson.Get(tcDataJson, cmdStrPath).String()

    s := strings.ToLower(cmdStr)

    ss := strings.Fields(strings.TrimSpace(s))

    if len(ss) == 0 {
        // fmt.Println("No cmd is provided")
    } else {
        switch ss[0] {
        case "sleep":
            if len(ss) == 1 {
                fmt.Println("No sleep duration provided, not slept")
            } else {
                tm, err := strconv.Atoi(ss[1])
                if err != nil {
                    fmt.Println("Provided sleep duration is not number, please fix")
                } else {
                    time.Sleep(time.Duration(tm)*time.Second)
                }
            }
        }
    }
        
    // as maybe no cmd is executed, the CmdExecStatus is always "cmdSuccess"
    // init
    tcDataStore.CmdType = "init"
    tcDataStore.CmdExecStatus = "cmdSuccess"
    tcDataStore.CmdAffectedCount = -1
    tcDataStore.CmdResults = -1

    sResults, sMessages = tcDataStore.HandleSingleCmdResult(i)

    return sResults, sMessages
}


func (tcDataStore *TcDataStore) HandleSingleCmdResult (i int) ([]bool, [][]*testcase.TestMessage) {
    // --------
    var cmdsResults []bool
    var finalTestMessages = [][]*testcase.TestMessage{}
    var cmdGroup []*testcase.CommandDetails

    if tcDataStore.CmdExecStatus == "cmdSuccess" {
        // path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmdResponse"
        // tcDataStore.PrepEmbeddedFunctions(path)
        //
        switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
        }

        cmdExpResp := cmdGroup[i].CmdResponse

        singleCmdResults, testMessages := tcDataStore.CompareRespGroup(cmdExpResp)

        // HandleSingleCommandResults for out
        if singleCmdResults == true {
            tcDataStore.HandleCmdResultsForOut(i)
        }

        cmdsResults = append(cmdsResults, singleCmdResults)
        finalTestMessages = append(finalTestMessages, testMessages)
    } else {
        cmdsResults = append(cmdsResults, false)
    }

    return cmdsResults, finalTestMessages
}


// for trial
func (tcDataStore *TcDataStore) CompareRespGroup (cmdExpResp []map[string]interface{}) (bool, []*testcase.TestMessage){
    //-----------
    singleCmdResults := true
    var testResults []bool
    var testMessages []*testcase.TestMessage
    for _, v := range cmdExpResp {
        // path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".cmdResponse"
        // tcDataStore.PrepEmbeddedFunctions(path)

        testRes, msg := tcDataStore.CompareRespGroupSingleAssertion(v)

        testMessages = append(testMessages, msg)
        testResults = append(testResults, testRes)
    }

    for key := range testResults {
        if testResults[key] == false {
            singleCmdResults = false
            break
        }
    }

    return singleCmdResults, testMessages
}

// for trial
func (tcDataStore *TcDataStore) CompareRespGroupSingleAssertion (v map[string]interface{}) (bool, *testcase.TestMessage){
    //-----------
    var testResult bool
    var testMessage *testcase.TestMessage

    for actualOrigin, value := range v {
        cmdExpResp_sub := value.(map[string]interface{})
        for assertionKey, expValueOrigin := range cmdExpResp_sub {
            switch assertionKey {
            case "HasMapKey", "NotHasMapKey":
                
            case "IsNull", "IsNotNull":

            default:
                var actualValue interface{}
                if strings.Contains(actualOrigin, "$(") {
                    l := tcDataStore.RenderExpresionA(actualOrigin)

                    actualValue = tcDataStore.GetResponseValue(l)
                } else {
                    actualValue = tcDataStore.RenderExpresionB(actualOrigin)
                }

                var expValue interface{}
                // expValueOrigin may have "$(...)" or "Fn::"
                // !!! currently, suppose the "$(...)" and "Fn::" do not coexist
                t := reflect.TypeOf(expValueOrigin).Kind().String()
                switch t {
                case "float64":
                    expValue = expValueOrigin
                case "string":
                    e := expValueOrigin.(string)
                    if strings.Contains(e, "$(") {
                        l := tcDataStore.RenderExpresionA(e)

                        expValue = tcDataStore.GetResponseValue(l)
                    } else {
                        expValue = tcDataStore.RenderExpresionB(e)
                    }
                case "map":
                    b, _ := json.Marshal(expValueOrigin)
                    s := tcDataStore.GetRenderTcVariables(string(b))

                    expValue = s

                    if strings.Contains(s, "Fn::") {
                        f := tcDataStore.EvaluateEmbeddedFunctions(cmdExpResp_sub)

                        var vv interface{}
                        json.Unmarshal([]byte(f.(string)), &vv)

                        for _, v1 := range vv.(map[string]interface{}) {
                            expValue = v1
                        }
                    }
                case "slice":
                    expValue = expValueOrigin
                default:
                    expValue = expValueOrigin
                }
                
                // fmt.Println("CompareRespGroupSingleAssertion: ", expValueOrigin, t, assertionKey, actualValue, expValue)
                testResult, testMessage = compareCommon(tcDataStore.CmdType, actualOrigin, assertionKey, actualValue, expValue)
            }  
        }
    }

    return testResult, testMessage
}

//
func (tcDataStore *TcDataStore) HandleCmdResultsForOut (i int) {
    var cmdGroup []*testcase.CommandDetails
    
    // write out session if has
    path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".session"
    tcDataStore.PrepEmbeddedFunctions(path)

    switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
    }

    expTcSession := cmdGroup[i].Session
    tcDataStore.WriteSession(expTcSession)

    // write out global variables if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".outGlobalVariables"
    tcDataStore.PrepEmbeddedFunctions(path)

    switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
    }
    
    expOutGlobalVariables := cmdGroup[i].OutGlobalVariables
    tcDataStore.WriteOutGlobalVariables(expOutGlobalVariables)

    // write out tc local variables if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".outLocalVariables"
    tcDataStore.PrepEmbeddedFunctions(path)

    switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
    }

    expOutLocalVariables := cmdGroup[i].OutLocalVariables
    tcDataStore.WriteOutTcLocalVariables(expOutLocalVariables)

    // write out files if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + tcDataStore.CmdSection + "." + fmt.Sprint(i) + ".outFiles"
    tcDataStore.PrepEmbeddedFunctions(path)

    switch tcDataStore.CmdSection {
        case "setUp":
            cmdGroup = tcDataStore.TcData.TestCase.SetUp()
        case "tearDown":
            cmdGroup = tcDataStore.TcData.TestCase.TearDown()
    }

    expOutFiles := cmdGroup[i].OutFiles
    tcDataStore.HandleOutFiles(expOutFiles)
}

