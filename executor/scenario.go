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
    // "os"
    "time"
    "fmt"
    "sync"
    "strings"
    "go4api/types"
    "go4api/utils"
    "path/filepath"
    // "go4api/utils/texttmpl"
    // simplejson "github.com/bitly/go-simplejson"
    "strconv"
    "go4api/logger"
)


func RunScenario(ch chan int, pStart_time time.Time, options map[string]string) {
    jsonFileList, _ := utils.WalkPath(options["testhome"] + "/Scenarios/", ".json")
    fmt.Println("Scenario jsonFileList:", options["ifScenario"], jsonFileList, "\n")

    var tcArray [][]interface{}
    // var tcNames []string

    baseUrl := GetBaseUrl(options)
    pStart := pStart_time.String()
    resultsDir := GetResultsDir(pStart, options)

    // (1). get the root cases in json (but maybe the json has notation, not valid json)
    // => the json has parentTestCase = root, or the the data table has parentTestCase = root
    tcArray = ConstructChildTcInfosBasedOnParentRoot(jsonFileList, "root" , "_dt") 

    fmt.Println("tcArray:", tcArray, "\n")

    
    // (2). render them, get the rendered cases
    // => need to build a tree???
    root, _ := BuildTree(tcArray)
    fmt.Println("------------------")
    fmt.Println("------------------", root, &root)
    ShowNodes(root)

    // (3). then execute them, genrate the outputs if have
    InitNodesRunResult(root, "Ready")

    

    miniLoop:
    for {
        resultsChan := make(chan types.TcRunResults, len(tcArray))
        var wg sync.WaitGroup
        //
        ScheduleNodes(root, &wg, options, "1", resultsChan, pStart, baseUrl, resultsDir)
        //
        wg.Wait()

        close(resultsChan)

        for tcRunResults := range resultsChan {
            fmt.Println("tcRunResults: ", tcRunResults)
            // here can refactor to struct
            tcName := tcRunResults.TcName
            parentTestCase := tcRunResults.ParentTestCase
            testResult := tcRunResults.TestResult
            actualStatusCode := tcRunResults.ActualStatusCode
            jsonFile_Base := tcRunResults.JsonFile_Base
            csvFileBase := tcRunResults.CsvFileBase
            rowCsv := tcRunResults.RowCsv
            start := tcRunResults.Start
            end := tcRunResults.End
            testMessages := tcRunResults.TestMessages
            start_time_UnixNano := tcRunResults.Start_time_UnixNano
            end_time_UnixNano := tcRunResults.End_time_UnixNano
            duration_UnixNano := tcRunResults.Duration_UnixNano
            //
            //(4). render the child cases, using the previous outputs as the inputs
            fmt.Println("----- testResult: ", testResult)
            if testResult == "Success" {
                tcArrayT := ConstructChildTcInfosBasedOnParentTcName(jsonFileList, tcName, "_outputs")
                fmt.Println("----- tcArrayT: ", tcArrayT)
                for _, tc := range tcArrayT {
                    ifAdded := AddNode(tc[0].(string), tc[2].(string), "", tc, "", "", "")
                    if ifAdded && true {
                        fmt.Println("----- child added")
                    } else {
                        fmt.Println("----- child not added")
                    }
                }
            }
                

            // (1). tcName, testResult, the search result is saved to *findNode
            SearchNode(&root, tcName)
            // (2). 
            RefreshNodeAndDirectChilrenTcResult(*findNode, testResult, start, end, 
                testMessages, start_time_UnixNano, end_time_UnixNano)
            // fmt.Println("------------------")
            // (3). <--> for log write to file
            resultReportString1 := "1" + "," + tcName + "," + parentTestCase + "," + testResult + "," + actualStatusCode + "," + jsonFile_Base + "," + csvFileBase
            resultReportString2 := "," + rowCsv + "," + start + "," + end + "," + "`" + "d" + "`" + "," + strconv.FormatInt(start_time_UnixNano, 10)
            resultReportString3 := "," + strconv.FormatInt(end_time_UnixNano, 10) + "," +  strconv.FormatInt(duration_UnixNano, 10)
            resultReportString :=  resultReportString1 + resultReportString2 + resultReportString3
            // (4). put the execution log into results
            logger.WriteExecutionResults(resultReportString, pStart, resultsDir)
            // fmt.Println("------!!!------")
        }

        // (4). render the child cases, using the previous outputs as the inputs
        // tcInputsFiles := utils.GetTestCaseBasicInputsFileNameFromJsonFile(jsonFile)
        // fmt.Println("tcInputsFiles: ", jsonFile, tcInputsFiles)
        
        // the case has inputs and its parent's runstatus == Success (i.e. not failed)


        // (5). execute the chilren, and so on
        statusReadyCount = 0
        CollectNodeReadyStatus(root, "1")

        // no more child cases can be added, then break
        if statusReadyCount == 0 {
            break miniLoop
        }
    }

    ShowNodes(root)

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



func ConstructChildTcInfosBasedOnParentRoot(jsonFileList []string, parentTcName string, dataTableSuffix string) [][]interface{} {
    var tcArray [][]interface{}
    var tcInfos [][]interface{}

    for _, jsonFile := range jsonFileList {
        tcNames := GetTestCaseBasicBasedOnParentFromJsonFile(jsonFile, "root")
        if len(tcNames) > 0 {
            csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")
            fmt.Println("tcNames: ", jsonFile, tcNames)

            if len(csvFileList) > 0 {
                tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
            } else {
                tcInfos = ConstructTcInfosBasedOnJson(jsonFile)
            }

            // fmt.Println("tcInfos:", tcInfos, "\n")
            
            for _, tc := range tcInfos {
                tcArray = append(tcArray, tc)
            }
        }
    }

    return tcArray
}


func ConstructChildTcInfosBasedOnParentTcName(jsonFileList []string, parentTcName string, dataTableSuffix string) [][]interface{} {
    var tcArray [][]interface{}
    var tcInfos [][]interface{}

    for _, jsonFile := range jsonFileList {
        tcNames := GetTestCaseBasicBasedOnParentFromJsonFile(jsonFile, parentTcName)
        if len(tcNames) > 0 {
            csvFileList := GetTestCaseBasicInputsFileNameFromJsonFile(jsonFile)
            fmt.Println("tcInputsFiles: ", jsonFile, csvFileList)
            // csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_outputs")
            fmt.Println("tcNames: ", jsonFile, tcNames)

            if len(csvFileList) > 0 {
                tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
            } else {
                tcInfos = ConstructTcInfosBasedOnJson(jsonFile)
            }

            // fmt.Println("tcInfos:", tcInfos, "\n")
            
            for _, tc := range tcInfos {
                tcArray = append(tcArray, tc)
            }
        }
    }

    return tcArray
}



func GetTestCaseBasicBasedOnParentFromJsonFile(filePath string, parentName string) []string {
    // as the raw Jsonfile itself is template, may not be valid json fomat, before rendered by data
    contents := utils.GetContentFromFile(filePath)

    var tcNames []string
    // Note: as we can not ensure if the field inputs and its value will on the same line, so use : as delimiter
    strList := strings.Split(string(contents), ":")
    // fmt.Println("strList - root: ", strList[0], strList)
    for ii, value := range strList {
        if strings.Contains(value, `"parentTestCase"`) {
            // "ParentTestCase-001": {
            //     "priority": "10",
            //     "parentTestCase": "root",
            parentStr := strings.Split(strList[ii + 1], ",")[0]
            if strings.Contains(parentStr, parentName) {
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
    // fmt.Println("strList - inputs: ", strList)
    for ii, value := range strList {
        if strings.Contains(value, `"inputs"`) {
            fileStr := strings.Split(strList[ii + 1], ",")[0]
            inputsFileBaseName := strings.TrimSpace(strings.Replace(fileStr, `"`, "", -1))
            inputsFiles = append(inputsFiles, filepath.Join(filepath.Dir(filePath), inputsFileBaseName))
        }
    }
    return inputsFiles
}
