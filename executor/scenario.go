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
    "time"
    "fmt"
    "sync"
    "strings"
    "regexp"
    "encoding/json"
    
    "go4api/cmd"
    "go4api/lib/testcase"
    "go4api/utils"
    "go4api/reports"
    "go4api/lib/csv"
)

func RunScenario(ch chan int, pStart_time time.Time, pStart string, baseUrl string, resultsDir string) {
    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")
    // fmt.Println("Scenario jsonFileList:", cmd.Opt.IfScenario, jsonFileList, "")

    var tcArray []testcase.TestCaseDataInfo

    // (1). get the root cases in json (but maybe the json has notation, not valid json)
    // => the json has parentTestCase = root, or the the data table has parentTestCase = root
    tcArray = ConstructChildTcInfosBasedOnParentRoot(jsonFileList, "root" , "_dt") 

    // (2). render them, get the rendered cases
    root, _ := BuildTree(tcArray)

    // (3). then execute them, genrate the outputs if have
    prioritySet := []string{"1"}
    InitVariables(prioritySet)
    InitNodesRunResult(root, "Ready")
    logFilePtr := reports.OpenExecutionResultsLogFile(resultsDir + pStart + ".log")
  
    miniLoop:
    for {
        resultsExeChan := make(chan testcase.TestCaseExecutionInfo, len(tcArray))
        var wg sync.WaitGroup
        //
        ScheduleNodes(root, &wg, "1", resultsExeChan, pStart, baseUrl, resultsDir)
        //
        wg.Wait()

        close(resultsExeChan)

        for tcExecution := range resultsExeChan {
            //(). render the child cases, using the previous outputs as the inputs
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
                        TestMessages: []*testcase.TestMessage{},
                        StartTimeUnixNano: 0,
                        EndTimeUnixNano: 0,
                        DurationUnixNano: 0,
                    }
                    ifAdded := AddNode(tcaseExecution)
                    if ifAdded && true {
                        fmt.Println("-> Child added: ", tcData.TcName())
                    } else {
                        fmt.Println("-> Child not added: ", tcData.TcName())
                    }
                }
            }
            // (1). tcName, testResult, the search result is saved to *findNode
            c := make(chan *tcNode)
            go func(c chan *tcNode) {
                defer close(c)
                SearchNode(c, root, tcExecution.TcName())
            }(c)
            // (2). 
            RefreshNodeAndDirectChilrenTcResult(<-c, tcExecution.TestResult, tcExecution.StartTime, tcExecution.EndTime, 
                    tcExecution.TestMessages, tcExecution.StartTimeUnixNano, tcExecution.EndTimeUnixNano)
            // (3). <--> for log write to file
            tcReportResults := tcExecution.TcReportResults()
            repJson, _ := json.Marshal(tcReportResults)
            reports.WriteExecutionResults(string(repJson), logFilePtr)

            reports.ReportConsoleByTc(tcExecution)
        }
        // (4). execute the chilren, and so on
        statusReadyCount = 0
        CollectNodeReadyStatusByPriority(root, "1")

        // no more child cases can be added, then break
        if statusReadyCount == 0 {
            break miniLoop
        }
    }
    logFilePtr.Close()
    
    CollectOverallNodeStatus(root, "Overall")
    reports.ReportConsoleOverall(statusCountByPriority["Overall"]["Total"], "Overall", statusCountByPriority, tcExecutedByPriority, tcNotExecutedByPriority)

    // generate the html report based on template, and results data
    // time.Sleep(1 * time.Second)
    pEnd_time := time.Now()
    //
    reports.GenerateTestReport(resultsDir, pStart_time, pStart, pEnd_time, 
        map[string]int{}, statusCountByPriority["Overall"]["Total"], statusCountByPriority, tcExecutedByPriority, tcNotExecutedByPriority)
    //
    fmt.Println("---------------------------------------------------------------------------")
    fmt.Println("Report Generated at: " + resultsDir + "index.html")
    fmt.Println("Execution Finished at: " + pEnd_time.String())
    
    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    ch <- statusCountByPriority["Overall"]["Fail"]
}


func ConstructChildTcInfosBasedOnParentRoot(jsonFileList []string, parentTcName string, dataTableSuffix string) []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo
    var tcInfos []testcase.TestCaseDataInfo

    for _, jsonFile := range jsonFileList {
        tcNames := GetBasicParentTestCaseInfosPerFile(jsonFile, "root")
        if len(tcNames) > 0 {
            csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")

            if len(csvFileList) > 0 {
                tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
            } else {
                tcInfos = ConstructTcInfosBasedOnJson(jsonFile)
            }
            
            for _, tcData := range tcInfos {
                tcArray = append(tcArray, tcData)
            }
        }
    }

    return tcArray
}

func ConstructChildTcInfosBasedOnParentTcName(jsonFileList []string, parentTcName string, dataTableSuffix string) []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo
    var tcInfos []testcase.TestCaseDataInfo

    for _, jsonFile := range jsonFileList {
        parentTcNames := GetBasicParentTestCaseInfosPerFile(jsonFile, parentTcName)
        //
        if len(parentTcNames) > 0 {
            csvFileList := GetBasicInputsFilesPerFile(jsonFile)
            // if has inputs -> if has *_dt -> use json
            if len(csvFileList) > 0 && utils.CheckFilesExistence(csvFileList) {
                tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
            } else if len(GetCsvDataFilesForJsonFile(jsonFile, "_dt")) > 0 {
                tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, GetCsvDataFilesForJsonFile(jsonFile, "_dt"))
            } else {
                tcInfos = ConstructTcInfosBasedOnJson(jsonFile)
            }

            for _, tcData := range tcInfos {
                tcArray = append(tcArray, tcData)
            }
        }
    }

    return tcArray
}

func GetBasicParentTestCaseInfosPerFile(filePath string, parentName string) []string {
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





