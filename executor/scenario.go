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
    "os"
    "time"
    "fmt"
    "sync"
    "strings"
    "encoding/json"
    "path/filepath"
    "go4api/cmd"
    "go4api/testcase"
    "go4api/utils"
    "go4api/logger"
)


func RunScenario(ch chan int, pStart_time time.Time, pStart string, baseUrl string, resultsDir string) {
    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")
    fmt.Println("Scenario jsonFileList:", cmd.Opt.IfScenario, jsonFileList, "")

    var tcArray []testcase.TestCaseDataInfo

    // (1). get the root cases in json (but maybe the json has notation, not valid json)
    // => the json has parentTestCase = root, or the the data table has parentTestCase = root
    tcArray = ConstructChildTcInfosBasedOnParentRoot(jsonFileList, "root" , "_dt") 

    // (2). render them, get the rendered cases
    root, _ := BuildTree(tcArray)

    // (3). then execute them, genrate the outputs if have
    InitNodesRunResult(root, "Ready")
  
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
                        TestMessages: "",
                        StartTimeUnixNano: 0,
                        EndTimeUnixNano: 0,
                        DurationUnixNano: 0,
                    }

                    ifAdded := AddNode(tcaseExecution)
                    if ifAdded && true {
                        fmt.Println("----- Child added: ", tcData.TcName())
                    } else {
                        fmt.Println("----- Child not added: ", tcData.TcName())
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
            // set Priority to 1 for all
            repJson, _ := json.Marshal(tcExecution)
            // (4). put the execution log into resultstestResult
            logger.WriteExecutionResults(string(repJson), pStart, resultsDir)
        }

        // (5). execute the chilren, and so on
        statusReadyCount = 0
        CollectNodeReadyStatus(root, "1")

        // no more child cases can be added, then break
        if statusReadyCount == 0 {
            break miniLoop
        }
    }

    // ShowNodes(root)

    // CollectOverallNodeStatus(root, len(prioritySet))

    // generate the html report based on template, and results data
    // time.Sleep(1 * time.Second)
    pEnd_time := time.Now()
    //
    GenerateTestReport(resultsDir, pStart_time, pStart, pEnd_time)
    //
    fmt.Println("---------------------------------------------------------------------------")
    fmt.Println("Report Generated at: " + resultsDir + "index.html")
    fmt.Println("Execution Finished at: " + pEnd_time.String())
    
    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    ch <- 1
}


func ConstructChildTcInfosBasedOnParentRoot(jsonFileList []string, parentTcName string, dataTableSuffix string) []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo
    var tcInfos []testcase.TestCaseDataInfo

    for _, jsonFile := range jsonFileList {
        tcNames := GetTestCaseBasicBasedOnParentFromJsonFile(jsonFile, "root")
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
        tcNames := GetTestCaseBasicBasedOnParentFromJsonFile(jsonFile, parentTcName)
        //
        if len(tcNames) > 0 {
            csvFileList := GetTestCaseBasicInputsFileNameFromJsonFile(jsonFile)
            // if has inputs -> if has *_dt -> use json
            if len(csvFileList) > 0 && CheckFilesExistence(csvFileList) {
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

func CheckFilesExistence(csvFileList []string) bool {
    ifExist := true

    for _, csvFile := range csvFileList {
        _, err := os.Stat(csvFile)
        if err != nil {
            fmt.Println("!!Error: " + csvFile + " does not exist.\n")
            ifExist = false
            break
        }
    }

    return ifExist
}


func GetTestCaseBasicBasedOnParentFromJsonFile(filePath string, parentName string) []string {
    // as the raw Jsonfile itself is template, may not be valid json fomat, before rendered by data
    contents := utils.GetContentFromFile(filePath)

    var tcNames []string
    // Note: as we can not ensure if the field inputs and its value will on the same line, so use : as delimiter
    strList := strings.Split(string(contents), ":")
    for ii, value := range strList {
        if strings.Contains(value, `"parentTestCase"`) {
            parentStr := strings.Split(strList[ii + 1], ",")[0]
            parentStr = strings.TrimSpace(strings.Replace(parentStr, `"`, "", -1))
            // here do not use Contains, use equal exactly
            if strings.EqualFold(parentStr, parentName) {
                parentStr := strings.Split(strList[ii - 2], ",")[0]
                str := strings.Split(parentStr, `"`)[len(strings.Split(parentStr, `"`)) - 2]
                tcNames = append(tcNames, str)
            } 
        }
    }
    return tcNames
}


func GetTestCaseBasicInputsFileNameFromJsonFile(filePath string) []string {
    // as the raw Jsonfile itself is template, may not be valid json fomat, before rendered by data
    contents := utils.GetContentFromFile(filePath)

    var inputsFiles []string
    // Note: as we can not ensure if the field inputs and its value will on the same line, so use : as delimiter
    strList := strings.Split(string(contents), ":")
    for ii, value := range strList {
        if strings.Contains(value, `"inputs"`) {
            fileStr := strings.Split(strList[ii + 1], ",")[0]
            inputsFileBaseName := strings.TrimSpace(strings.Replace(fileStr, `"`, "", -1))
            if inputsFileBaseName != "" {
                inputsFiles = append(inputsFiles, filepath.Join(filepath.Dir(filePath), inputsFileBaseName))
            }
        }
    }
    return inputsFiles
}
