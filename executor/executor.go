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
    "go4api/testcase"
    "go4api/ui"     
    "go4api/ui/js"  
    "go4api/ui/style"                                                                                                                                
    "go4api/utils"
    "go4api/utils/texttmpl"
    "path/filepath"
    "strings"
    "io/ioutil"
    "strconv"
    "go4api/logger"
    "encoding/json"
)


func Run(ch chan int, pStart_time time.Time, options map[string]string, pStart string, baseUrl string, resultsDir string, tcArray []testcase.TestCaseDataInfo) { //client
    // (1), get the text path, default is ../data/*, then search all the sub-folder to get the test scripts
    // to check the tcArray, if the case not distinct, report it to fix
    if len(tcArray) != len(GetTcNameSet(tcArray)) {
        fmt.Println("\n!! There are duplicated test case names, please make them distinct\n")
        os.Exit(1)
    }
    //
    // fmt.Println("tcArray:", tcArray, "\n")
    // myabe there needs a scheduler, for priority 1 (w or w/o dependency) -> priority 2 (w or w/o dependency), ...
    // --
    // How to impliment the case Dependency???
    // Two big categories: 
    // (1) case has No parent Dependency or successor Dependency, which can be scheduled concurrently
    // (2) case has parent Dependency or successor Dependency, which has rules to be scheduled concurrently
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

    // If need to set the Concurrency MAX?
    InitNodesRunResult(root, "Ready")
    // fmt.Println("------------------", root, &root)

    //
    fmt.Println("\n====> test cases execution starts!\n")
    statusReadyCount = 0
    // init the status count list
    statusCountList = make([][]int, len(prioritySet) + 1)
    for i := range statusCountList {
        statusCountList[i] = make([]int, 5)
    }
    //
    for p_index, priority := range prioritySet {
        tcArrayPriority := classifications[priority]
        fmt.Println("====> Priority " + priority + " starts!")
        
        miniLoop:
        for {
            //
            resultsExeChan := make(chan testcase.TestCaseExecutionInfo, len(tcArray))
            var wg sync.WaitGroup
            //
            ScheduleNodes(root, &wg, options, priority, resultsExeChan, pStart, baseUrl, resultsDir)
            //
            wg.Wait()

            close(resultsExeChan)

            for tcExecution := range resultsExeChan {
                // (1). tcName, testResult, the search result is saved to *findNode
                SearchNode(&root, tcExecution.TcName())
                // (2). 
                RefreshNodeAndDirectChilrenTcResult(*findNode, tcExecution.TestResult, tcExecution.StartTime, tcExecution.EndTime, 
                    tcExecution.TestMessages, tcExecution.StartTimeUnixNano, tcExecution.EndTimeUnixNano)
                // fmt.Println("------------------")
                // (3). <--> for log write to file
                tcReportResults := tcExecution.TcReportResults()
                repJson, _ := json.Marshal(tcReportResults)
                // fmt.Println(string(repJson))
                // (4). put the execution log into results
                logger.WriteExecutionResults(string(repJson), pStart, resultsDir)
                // fmt.Println("------!!!------")
            }
            // if tcTree has no node with "Ready" status, break the miniloop
            statusReadyCount = 0
            CollectNodeReadyStatus(root, priority)
            // fmt.Println("------------------ statusReadyCount: ", statusReadyCount)
            if statusReadyCount == 0 {
                break miniLoop
            }
            // ShowNodes(root)
        }
        //
        CollectNodeStatusByPriority(root, p_index, priority)

        // (5). also need to put out the cases which has not been executed (i.e. not Success, Fail)
        notRunTime := time.Now()
        for _, tcExecution := range tcNotExecutedList {
            // [casename, priority, parentTestCase, ...], tc, jsonFile, csvFile, row in csv
            if tcExecution.Priority() == priority {
                // set some dummy time for the tc not executed
                tcExecution.StartTimeUnixNano =notRunTime.UnixNano()
                tcExecution.EndTimeUnixNano = notRunTime.UnixNano()
                tcExecution.DurationUnixNano = notRunTime.UnixNano() - notRunTime.UnixNano()

                tcReportResults := tcExecution.TcReportResults()
                repJson, _ := json.Marshal(tcReportResults)
                //
                logger.WriteExecutionResults(string(repJson), pStart, resultsDir)
                // to console
            }
        }
        
        //
        var successCount = statusCountList[p_index][2]
        var failCount = statusCountList[p_index][3]
        //
        fmt.Println("---------------------------------------------------------------------------")
        fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(len(tcArrayPriority)) + " Cases in template -----")
        fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(statusCountList[p_index][0]) + " Cases put onto tcTree -----")
        fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(statusCountList[p_index][0] - successCount - failCount) + " Cases Skipped (Not Executed, due to Parent Failed) -----")
        fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(successCount + failCount) + " Cases Executed -----")
        fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(successCount) + " Cases Success -----")
        fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(failCount) + " Cases Fail -----")
        fmt.Println("---------------------------------------------------------------------------")

        fmt.Println("====> Priority " + priority + " ended! \n")
        // sleep for debug
        time.Sleep(500 * time.Millisecond)
    }
    // ShowNodes(root)
    CollectOverallNodeStatus(root, len(prioritySet))
    // fmt.Println("====> statusCountList final: ", statusCountList)
    //
    var successCount = statusCountList[len(prioritySet)][2]
    var failCount = statusCountList[len(prioritySet)][3]
    //
    fmt.Println("---------------------------------------------------------------------------")
    fmt.Println("----- Total " + strconv.Itoa(len(tcArray)) + " Cases in template -----")
    fmt.Println("----- Total " + strconv.Itoa(statusCountList[len(prioritySet)][0]) + " Cases put onto tcTree -----")
    fmt.Println("----- Total " + strconv.Itoa(statusCountList[len(prioritySet)][0] - successCount - failCount) + " Cases Skipped (Not Executed, due to Parent Failed) -----")
    fmt.Println("----- Total " + strconv.Itoa(successCount + failCount) + " Cases Executed -----")
    fmt.Println("----- Total " + strconv.Itoa(successCount) + " Cases Success -----")
    fmt.Println("----- Total " + strconv.Itoa(failCount) + " Cases Fail -----")
    fmt.Println("---------------------------------------------------------------------------\n\n")


    // generate the html report based on template, and results data
    // time.Sleep(1 * time.Second)
    pEnd_time := time.Now()
    //
    GenerateTestReport(resultsDir, pStart_time, pStart, pEnd_time)
    //
    fmt.Println("Report Generated at: " + resultsDir + "index.html")
    fmt.Println("Execution Finished at: " + pEnd_time.String())

    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    ch <- 1

}



func GetTcArray(options map[string]string) []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo

    jsonFileList, _ := utils.WalkPath(options["testhome"] + "/testdata/", ".json")
    // fmt.Println("jsonFileList:", jsonFileList, "\n")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")
        // fmt.Println("csvFileList:", csvFileList, "\n")
        //
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
        
        // if             
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
                for _, tcase := range tcases {
                    // populate the testcase.TestCaseDataInfo
                    tcaseData := testcase.TestCaseDataInfo {
                        TestCase: tcase,
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
    // as the json is generated based on templated dynamically, so that, to cache all the resulted json in array
     for _, tcase := range tcases {
        // populate the testcase.TestCaseDataInfo
        tcaseData := testcase.TestCaseDataInfo {
            TestCase: tcase,
            JsonFilePath: jsonFile,
            CsvFile: csvFile,
            CsvRow: csvRow,
        }
        tcInfos = append(tcInfos, tcaseData)
    }

    return tcInfos
}




func GetTcNameSet(tcArray []testcase.TestCaseDataInfo) []string {
    // get the tcNames
    var tcNames []string
    for _, tcaseInfo := range tcArray {
        tcNames = append(tcNames, tcaseInfo.TcName())
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
        if _, value := keys[entry]; !value {
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
    // fmt.Println("classifications: ", classifications)
    return classifications
}



func GenerateTestReport(resultsDir string, pStart_time time.Time, pStart string, pEnd_time time.Time) {
    // read the resource under /ui/*
    // fmt.Println("ui: ", ui.Index_template)

    // copy the value of var Index to file
    utils.GenerateFileBasedOnVarOverride(ui.Index, resultsDir + "index.html")

    //
    err := os.MkdirAll(resultsDir + "js", 0777)
    if err != nil {
      panic(err) 
    }
    // copy the value of var js.Js to file
    texttmpl.GenerateHtmlJsCSSFromTemplateAndVar(js.Results, pStart_time, pEnd_time, resultsDir, resultsDir + pStart + ".log")
    //
    utils.GenerateFileBasedOnVarOverride(js.Js, resultsDir + "js/go4api.js")
    //
    err = os.MkdirAll(resultsDir + "style", 0777)
    if err != nil {
      panic(err) 
    }
    // copy the value of var style.Style to file
    utils.GenerateFileBasedOnVarOverride(style.Style, resultsDir + "style/go4api.css")
}



func GetTmpJsonDir(path string) string {
    // check if the /tmp/go4api_wfasf exists, if exists, then rm first
    os.RemoveAll("/tmp/" + path)
    //
    var resultsDir string
    err := os.Mkdir("/tmp/" + path + "/", 0777)
    if err != nil {
      panic(err) 
    } else {
        resultsDir = "/tmp/" + path + "/"
    }

    return resultsDir
}



func GetFuzzTcArray(options map[string]string) []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo

    jsonFileList, _ := utils.WalkPath(options["testhome"] + "/testdata/", ".json")
    // fmt.Println("jsonFileList:", jsonFileList, "\n")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_fuzz_dt")
        // to get the json test data directly (if not template) based on template (if template)
        // tcInfos: [[casename, priority, parentTestCase, ], ...]
        var tcInfos []testcase.TestCaseDataInfo
        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
        }
        // fmt.Println("tcInfos:", tcInfos, "\n")
        
        for _, tcData := range tcInfos {
            tcArray = append(tcArray, tcData)
        }
    }

    return tcArray
}
