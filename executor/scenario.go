/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package executor

import (
    // "time"
    "fmt"
    "sync"
    "strings"
    "regexp"
    "encoding/json"
    
    "go4api/cmd"
    "go4api/lib/testcase"
    "go4api/lib/tree"
    "go4api/utils"
    "go4api/reports"
    "go4api/lib/csv"
)

func RunScenario (ch chan int, baseUrl string, resultsDir string, resultsLogFile string, 
        jsonFileList []string, tcArray []testcase.TestCaseDataInfo) tree.TcTreeStats {
    // --
    root, tcTree, tcTreeStats := InitRunScenario(tcArray)

    logFilePtr := reports.OpenExecutionResultsLogFile(resultsLogFile)
  
    miniLoop:
    for {
        resultsExeChan := make(chan testcase.TestCaseExecutionInfo, len(tcArray))
        var wg sync.WaitGroup
        //
        cReady := make(chan *tree.TcNode)
        go func(cReady chan *tree.TcNode) {
            defer close(cReady)
            tcTree.CollectNodeReadyByPriority(cReady, root, "1")
        }(cReady)

        ScheduleCases(cReady, &wg, resultsExeChan, baseUrl)
        //
        wg.Wait()

        close(resultsExeChan)

        for tcExecution := range resultsExeChan {
            tcTreeStats.DeductReadyCount("1")
            //
            BuildChilrenNodes(tcExecution, jsonFileList, root, tcTree)
            //
            c := make(chan *tree.TcNode)
            go func(c chan *tree.TcNode) {
                defer close(c)
                tcTree.SearchNode(c, root, tcExecution.TcName())
            }(c)
            //
            tcTree.RefreshNodeAndDirectChilrenTcResult(<-c, tcExecution.TestResult, tcExecution.StartTime, tcExecution.EndTime, 
                    tcExecution.HttpTestMessages, tcExecution.StartTimeUnixNano, tcExecution.EndTimeUnixNano)
            //
            tcReportResults := tcExecution.TcReportResults()
            repJson, _ := json.Marshal(tcReportResults)
            reports.WriteExecutionResults(string(repJson), logFilePtr)

            reports.ReportConsoleByTc(tcExecution)
        }
        tcTreeStats.CollectNodeStatusByPriority(root, "1")
        // if no more child cases can be added, then break
        if tcTreeStats.StatusCountByPriority["1"]["Ready"] == 0 {
            break miniLoop
        }
    }
    logFilePtr.Close()
    //
    RunConsoleOverallReport(tcArray, root, tcTreeStats)

    return tcTreeStats
}

func InitRunScenario (tcArray []testcase.TestCaseDataInfo) (*tree.TcNode, tree.TcTree, tree.TcTreeStats) {
    //
    tcTree := tree.CreateTcTree()
    root := tcTree.BuildTree(tcArray)

    // (3). then execute them, genrate the outputs if have
    prioritySet := []string{"1"}
    tcTreeStats := tree.CreateTcTreeStats(prioritySet)
    //
    tcTree.InitNodesRunResult(root, "Ready")
    tcTreeStats.CollectNodeStatusByPriority(root, "1")

    return root, tcTree, tcTreeStats
}

func BuildChilrenNodes (tcExecution testcase.TestCaseExecutionInfo, jsonFileList []string, root *tree.TcNode, tcTree tree.TcTree) {
    // render the child cases, using the previous outputs as the inputs
    // the case has inputs and its parent's runstatus == Success (i.e. not failed)
    if tcExecution.TestResult == "Success" {
        tcArrayT := ConstructChildTcInfosBasedOnParentTcName(jsonFileList, tcExecution.TcName(), "_outputs")
        for _, tcData := range tcArrayT {
            tcaseExecution := testcase.TestCaseExecutionInfo {
                TestCaseDataInfo: &tcData,
                TestResult: "",
                ActualStatusCode: 0,
                StartTime: "",
                EndTime: "",
                HttpTestMessages: []*testcase.TestMessage{},
                StartTimeUnixNano: 0,
                EndTimeUnixNano: 0,
                DurationUnixNano: 0,
            }

            ifAdded := tcTree.AddNode(root, tcaseExecution)
            if ifAdded && true {
                fmt.Println("-> Child added: ", tcData.TcName())
            } else {
                fmt.Println("-> Child not added: ", tcData.TcName())
            }
        }
    }
}


func GetJsonFiles () []string {
    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")

    return jsonFileList
}

func ConstructChildTcInfosBasedOnParentRoot (jsonFileList []string, parentTcName string, dataTableSuffix string) []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo
    var tcInfos []testcase.TestCaseDataInfo

    for _, jsonFile := range jsonFileList {
        tcNames := GetBasicParentTestCaseInfosPerFile(jsonFile, "root")
        if len(tcNames) > 0 {
            csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")

            if len(csvFileList) > 0 {
                tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList, parentTcName)
            } else {
                tcInfos = ConstructTcInfosBasedOnJson(jsonFile, parentTcName)
            }
            
            for _, tcData := range tcInfos {
                tcArray = append(tcArray, tcData)
            }
        }
    }

    return tcArray
}

func ConstructChildTcInfosBasedOnParentTcName (jsonFileList []string, parentTcName string, dataTableSuffix string) []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo
    var tcInfos []testcase.TestCaseDataInfo

    for _, jsonFile := range jsonFileList {
        parentTcNames := GetBasicParentTestCaseInfosPerFile(jsonFile, parentTcName)
        //
        if len(parentTcNames) > 0 {
            csvInputsFileList := GetBasicInputsFilesPerFile(jsonFile)
            // if has inputs -> if has *_dt -> use json
            if len(csvInputsFileList) > 0 && utils.CheckFilesExistence(csvInputsFileList) {
                tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvInputsFileList, parentTcName)
            } else if len(GetCsvDataFilesForJsonFile(jsonFile, "_dt")) > 0 {
                dtFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")
                tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, dtFileList, parentTcName)
            } else {
                tcInfos = ConstructTcInfosBasedOnJson(jsonFile, parentTcName)
            }

            for _, tcData := range tcInfos {
                tcArray = append(tcArray, tcData)
            }
        }
    }

    return tcArray
}

func GetBasicParentTestCaseInfosPerFile (filePath string, parentName string) []string {
    // as the raw Jsonfile itself is template, may not be valid json fomat, before rendered by data
    contentsBytes := utils.GetContentFromFile(filePath)
    contents := string(contentsBytes)
    // add some space a make regexp works well
    contents = strings.Replace(contents, `[`, "[ ", -1)
    contents = strings.Replace(contents, `{`, "{ ", -1)

    var tcNames []string
    var parentTcInfos []string

    reg := regexp.MustCompile(`[\p{L}\w\pP]+`)
    wordSlice := reg.FindAllString(contents, -1)

    parentPos := 0
    for i, value := range wordSlice {
        if strings.Contains(value, `"parentTestCase"`) {
            parentPos = i
            parentTcInfos = append(parentTcInfos, wordSlice[parentPos + 1])
        }
        if i == parentPos {
            continue
        }
    }
    //
    priorityPos := 0
    for i, value := range wordSlice {
        if strings.Contains(value, `"priority"`) {
            priorityPos = i
            tcNames = append(tcNames, wordSlice[priorityPos - 2])
        }
        if i == priorityPos {
            continue
        }
    }
    //
    for i, _ := range parentTcInfos {
        parentTcInfos[i] = strings.Replace(parentTcInfos[i], `"`, "", -1)
        parentTcInfos[i] = strings.Replace(parentTcInfos[i], `,`, "", -1)
        parentTcInfos[i] = strings.Replace(parentTcInfos[i], ` `, "", -1)
    }
    for i, _ := range tcNames {
        tcNames[i] = strings.Replace(tcNames[i], `"`, "", -1)
        tcNames[i] = strings.Replace(tcNames[i], `,`, "", -1)
        tcNames[i] = strings.Replace(tcNames[i], ` `, "", -1)
        tcNames[i] = strings.Replace(tcNames[i], `:`, "", -1)
    }

    var tcNamesMatchParent []string
    for i, _ := range tcNames {
        if parentTcInfos[i] == parentName {
            tcNamesMatchParent = append(tcNamesMatchParent, tcNames[i])
        }
    }
    // fmt.Println("parentTcInfos 2: ", parentTcInfos, tcNames, tcNamesMatchParent)
    return tcNamesMatchParent
}

func writeGcsvToCsv (gcsvPtr *gcsv.Gcsv, outFile string) {
    // header
    utils.GenerateCsvFileBasedOnVarOverride(gcsvPtr.Header, outFile)
    // data
    for i, _ := range gcsvPtr.DataRows {
        utils.GenerateCsvFileBasedOnVarAppend(gcsvPtr.DataRows[i], outFile)
    }
}

func writeMapToCsv (inputsMap map[string]interface{}, joinedFile string) {
    var keyStrList []string
    var valueStrList []string

    for key, value := range inputsMap {
        // for csv header
        keyStrList = append(keyStrList, key)
        // for cav data
        valueStrList = append(valueStrList, value.(string))
    }
    // write csv header
    utils.GenerateCsvFileBasedOnVarOverride(keyStrList, joinedFile)
    // write csv data
    utils.GenerateCsvFileBasedOnVarAppend(valueStrList, joinedFile)
}

