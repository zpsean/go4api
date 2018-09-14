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
    "fmt"
    "time"
    "os"
    "sort"
    "sync"
    "path/filepath"
    "strings"
    "io/ioutil"
    "strconv"
    "encoding/json"

    "go4api/cmd"
    "go4api/utils"
    "go4api/lib/testcase"
    "go4api/texttmpl"
    "go4api/reports"
)


func Run(ch chan int, pStart_time time.Time, pStart string, baseUrl string, resultsDir string, tcArray []testcase.TestCaseDataInfo) { 
    // check the tcArray, if the case not distinct, report it to fix
    if len(tcArray) != len(GetTcNameSet(tcArray)) {
        fmt.Println("\n!! There are duplicated test case names, please make them distinct")
        os.Exit(1)
    }
    //
    root, _ := BuildTree(tcArray)
    fmt.Println("------------------")
    //
    prioritySet := GetPrioritySet(tcArray)
    classifications := GetTestCasesByPriority(prioritySet, tcArray)
    // Note, before starting execution, needs to sort the priorities_set first by priority
    // Note: here is a bug, as the sort results is 1, 10, 11, 2, 3, etc. => fixed
    prioritySet_Int := utils.ConvertStringArrayToIntArray(prioritySet)
    sort.Ints(prioritySet_Int)
    prioritySet = utils.ConvertIntArrayToStringArray(prioritySet_Int)
    // Init
    InitVariables(prioritySet)
    InitNodesRunResult(root, "Ready")
    //
    fmt.Println("\n====> test cases execution starts!")

    logFilePtr := reports.OpenExecutionResultsLogFile(resultsDir + pStart + ".log")
    
    for _, priority := range prioritySet {
        tcArrayPriority := classifications[priority]
        fmt.Println("====> Priority " + priority + " starts!")
        
        miniLoop:
        for {
            //
            resultsExeChan := make(chan testcase.TestCaseExecutionInfo, len(tcArray))
            var wg sync.WaitGroup
            //
            ScheduleNodes(root, &wg, priority, resultsExeChan, pStart, baseUrl, resultsDir)
            //
            wg.Wait()

            close(resultsExeChan)

            for tcExecution := range resultsExeChan {
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
                // (4). put the execution log into results
                reports.WriteExecutionResults(string(repJson), logFilePtr)
            }
            // if tcTree has no node with "Ready" status, break the miniloop
            statusReadyCount = 0
            CollectNodeReadyStatusByPriority(root, priority)
            //
            if statusReadyCount == 0 {
                break miniLoop
            }
        }
        CollectNodeStatusByPriority(root, priority)

        // (5). also need to put out the cases which has not been executed (i.e. not Success, Fail)
        notRunTime := time.Now()
        for i, _ := range tcNotExecutedByPriority[priority] {
            for _, tcExecution := range tcNotExecutedByPriority[priority][i] {
                // [casename, priority, parentTestCase, ...], tc, jsonFile, csvFile, row in csv
                if tcExecution.Priority() == priority {
                    // set some dummy time for the tc not executed
                    tcExecution.StartTimeUnixNano =notRunTime.UnixNano()
                    tcExecution.EndTimeUnixNano = notRunTime.UnixNano()
                    tcExecution.DurationUnixNano = notRunTime.UnixNano() - notRunTime.UnixNano()

                    tcReportResults := tcExecution.TcReportResults()
                    repJson, _ := json.Marshal(tcReportResults)
                    //
                    reports.WriteExecutionResults(string(repJson), logFilePtr)
                    // to console
                }
            }
        }
        // report to console
        reports.ReportConsoleByPriority(len(tcArrayPriority), priority, statusCountByPriority, tcExecutedByPriority, tcNotExecutedByPriority)

        fmt.Println("====> Priority " + priority + " ended!")
        fmt.Println("")
        // sleep for debug
        // time.Sleep(500 * time.Millisecond)
    }
    logFilePtr.Close()

    CollectOverallNodeStatus(root, "Overall")
    reports.ReportConsoleByPriority(len(tcArray), "Overall", statusCountByPriority, tcExecutedByPriority, tcNotExecutedByPriority)
    
    // generate the html report based on template, and results data
    // time.Sleep(1 * time.Second)
    pEnd_time := time.Now()
    //
    reports.GenerateTestReport(resultsDir, pStart_time, pStart, pEnd_time)
    //
    fmt.Println("Report Generated at: " + resultsDir + "index.html")
    fmt.Println("Execution Finished at: " + pEnd_time.String())

    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    ch <- statusCountByPriority["Overall"]["Fail"]
}


func GetTcArray() []testcase.TestCaseDataInfo { 
    var tcArray []testcase.TestCaseDataInfo

    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")
        //
        var tcInfos []testcase.TestCaseDataInfo

        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
        } else {
            tcInfos = ConstructTcInfosBasedOnJson(jsonFile)
        }

        for _, tcData := range tcInfos {
            // fmt.Println("\n tcData:", tcData.TcName())
            tcArray = append(tcArray, tcData)
        }
    }

    return tcArray
}


func GetCsvDataFilesForJsonFile(jsonFile string, suffix string) []string {
    // here search out the csv files under the same dir, not to use utils.WalkPath as it is recursively
    var csvFileListTemp []string
    infos, err := ioutil.ReadDir(filepath.Dir(jsonFile))
    if err != nil {
      panic(err)
    }

    // get the csv file, ignore the fields "inputs", "outputs"
    for _, info := range infos {
      if filepath.Ext(info.Name()) == ".csv" {
        csvFileListTemp = append(csvFileListTemp, filepath.Join(filepath.Dir(jsonFile), info.Name()))
      }
    }
    // 
    var csvFileList []string
    for _, csvFile := range csvFileListTemp {
        csvFileName := strings.TrimRight(filepath.Base(csvFile), ".csv")
        jsonFileName := strings.TrimRight(filepath.Base(jsonFile), ".json")
        // Note: the json file realted data table files is pattern: jsonFileName + "_dt[*]"
        if strings.Contains(csvFileName, jsonFileName + suffix) {
            csvFileList = append(csvFileList, csvFile)
        }
    }

    return csvFileList
}


func ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile string, csvFileList []string) []testcase.TestCaseDataInfo {
    var tcInfos []testcase.TestCaseDataInfo

    for _, csvFile := range csvFileList {
        // to check the csv file's existence
        csvRows := utils.GetCsvFromFile(csvFile)
        for i, csvRow := range csvRows {
            // starting with data row
            if i > 0 {
                outTempJson := texttmpl.GenerateJsonBasedOnTemplateAndCsv(jsonFile, csvRows[0], csvRow)

                var tcases testcase.TestCases
                resJson, _ := ioutil.ReadAll(outTempJson)
                json.Unmarshal([]byte(resJson), &tcases)
                // as the json is generated based on templated dynamically, so that, to cache all the resulted json in array
                for i, _ := range tcases {
                    // populate the testcase.TestCaseDataInfo
                    tcaseData := testcase.TestCaseDataInfo {
                        TestCase: &tcases[i],
                        JsonFilePath: jsonFile,
                        CsvFile: csvFile,
                        CsvRow: strconv.Itoa(i + 1),
                    }
                    tcInfos = append(tcInfos, tcaseData)
                }
            }
        }
    }
    return tcInfos
}

func ConstructTcInfosBasedOnJson(jsonFile string) []testcase.TestCaseDataInfo {
    var tcInfos []testcase.TestCaseDataInfo

    csvFile := ""
    csvRow := ""
    outTempJson := texttmpl.GenerateJsonBasedOnTemplateAndCsv(jsonFile, []string{""}, []string{""})
    
    var tcases testcase.TestCases
    resJson, _ := ioutil.ReadAll(outTempJson)
    json.Unmarshal([]byte(resJson), &tcases)
    // fmt.Println("resJson: ", string(resJson), tcases)
    // as the json is generated based on templated dynamically, so that, to cache all the resulted json in array
     for i, _ := range tcases {
        // populate the testcase.TestCaseDataInfo
        tcaseData := testcase.TestCaseDataInfo {
            TestCase: &tcases[i],
            JsonFilePath: jsonFile,
            CsvFile: csvFile,
            CsvRow: csvRow,
        }
        tcInfos = append(tcInfos, tcaseData)
    }

    return tcInfos
}


func GetTcNameSet(tcArray []testcase.TestCaseDataInfo) []string {
    var tcNames []string

    for _, tcaseInfo := range tcArray {
        var ifExists bool
        ifExists = false
        for _, tcN := range tcNames {
            if tcaseInfo.TcName() == tcN {
                ifExists = true
                break
            }
        } 
        if ifExists == false {
            tcNames = append(tcNames, tcaseInfo.TcName())
        }   
    }
    return tcNames
}


func GetPrioritySet(tcArray []testcase.TestCaseDataInfo) []string {
    // get the priorities
    var priorities []string
    for _, tcaseInfo := range tcArray {
        priorities = append(priorities, tcaseInfo.Priority())
    }
    // go get the distinct key in priorities
    keys := make(map[string]bool)
    var prioritySet []string
    for _, entry := range priorities {
        // uses 'value, ok := map[key]' to determine if map's key exists, if ok, then true
        if _, ok := keys[entry]; !ok {
            keys[entry] = true
            prioritySet = append(prioritySet, entry)
        }
    }

    return prioritySet
}



func GetTestCasesByPriority(prioritySet []string, tcArray []testcase.TestCaseDataInfo) map[string][]testcase.TestCaseDataInfo {
    // build the map
    classifications := make(map[string][]testcase.TestCaseDataInfo)
    for _, entry := range prioritySet {
        for _, tcaseData := range tcArray {
            if entry == tcaseData.Priority() {
                classifications[entry] = append(classifications[entry], tcaseData)
            }
        }
    }

    return classifications
}


func GetFuzzTcArray() []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo

    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_fuzz_dt")
        // to get the json test data directly (if not template) based on template (if template)
        // tcInfos: [[casename, priority, parentTestCase, ], ...]
        var tcInfos []testcase.TestCaseDataInfo
        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
        }

        for _, tcData := range tcInfos {
            tcArray = append(tcArray, tcData)
        }
    }

    return tcArray
}


func GetOriginMutationTcArray() []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo

    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_mutation_dt")
        // to get the json test data directly (if not template) based on template (if template)
        // tcInfos: [[casename, priority, parentTestCase, ], ...]
        var tcInfos []testcase.TestCaseDataInfo
        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
        } else {
            tcInfos = ConstructTcInfosBasedOnJson(jsonFile)
        }

        for _, tcData := range tcInfos {
            tcArray = append(tcArray, tcData)
        }
    }

    return tcArray
}


