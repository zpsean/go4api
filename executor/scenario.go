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
    "regexp"
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
        tcNames := GetBasicParentTestCaseInfosPerFile(jsonFile, parentTcName)
        //
        if len(tcNames) > 0 {
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


func GetBasicInputsFilesPerFile(filePath string) []string {
    fileInputsInfos := GetBasicInputsInfos(filePath)
    inputsFiles := GenerateInputsFiles(filePath, fileInputsInfos[0])

    GenerateInputsFileWithConsolidatedData(filePath, fileInputsInfos[0])

    fmt.Println("-- :", filePath, inputsFiles)
    return inputsFiles
}

func GetBasicInputsInfos(filePath string) [][]string {
    // as the raw Jsonfile itself is template, may not be valid json fomat, before rendered by data
    contentsBytes := utils.GetContentFromFile(filePath)
    contents := string(contentsBytes)
    // add some space a make regexp works well
    contents = strings.Replace(contents, `[`, "[ ", -1)
    contents = strings.Replace(contents, `{`, "{ ", -1)

    var inputsInfos []string
    var fileInputsInfos [][]string
    var inputsPosSlice []int

    reg := regexp.MustCompile(`[\p{L}\w\pP]+`)
    wordSlice := reg.FindAllString(contents, -1)

    for i, value := range wordSlice {
        if strings.Contains(value, `"inputs"`) {
            inputsPosSlice = append(inputsPosSlice, i)
        }
    }
    for _, inputsPos := range inputsPosSlice {
        for ii := inputsPos + 1; ; ii ++ {
            v := strings.Replace(wordSlice[ii], `"`, "", -1)
            v = strings.Replace(v, `[`, "", -1)
            v = strings.Replace(v, `]`, "", -1)
            v = strings.Replace(v, `,`, "", -1)
            v = strings.Replace(v, ` `, "", -1)
            if len(v) > 0 {
                inputsInfos = append(inputsInfos, v)
                // fmt.Println(v)
            }
            if strings.Contains(wordSlice[ii], `]`) {
                break
            }
        }
        fileInputsInfos = append(fileInputsInfos, inputsInfos)
    }

    return fileInputsInfos
}


func GenerateInputsFiles (filePath string, inputsInfos []string) []string {
    var inputsFiles []string
    
    for _, value := range inputsInfos {  
        tempDir := filepath.Dir(filePath) + "/temp"

        inputsFiles = append(inputsFiles, filepath.Join(tempDir, value))
        break
    }
    fmt.Println("inputsInfos: ", inputsInfos, inputsFiles)
    return inputsFiles
}


// to implements the operator: union, join, append for inputs files    
func GenerateInputsFileWithConsolidatedData (filePath string, inputsInfos []string) {
    // inputsInfos => ["s1ParentTestCase_out.csv", "join", "s1ParentTestCase_out.csv"]
    // var inputsFiles []string
    // 1. len(inputsInfos)
    if len(inputsInfos) > 0 {
        if len(inputsInfos) % 2 != 1 {
            fmt.Println("!! Error, inputs contents error, please check")
            os.Exit(1)
        }
        for i := 1; i <= len(inputsInfos) / 2; i ++ {
            operator := strings.ToLower(inputsInfos[2 * (i - 1) + 1])
            if operator != "union" && operator != "join" && operator != "append" {
                fmt.Println("!! Error, inputs operator error, please check")
                os.Exit(1)
            }
        }
        // loop the inputsInfos and apply operator
        inputsMap := make(map[string]interface{})
        // init inputsMap
        tempDir := filepath.Dir(filePath) + "/temp"
        csvRows := utils.GetCsvFromFile(filepath.Join(tempDir, inputsInfos[0]))
        // suppose we have only one row here
        for _, csvRow := range csvRows {
            for i, _ := range csvRow {
                inputsMap[csvRows[0][i]] = csvRows[1][i]
            }
        }
        //
        for i := 1; i <= len(inputsInfos) / 2; i ++ {
            operator := strings.ToLower(inputsInfos[2 * (i - 1) + 1])
            tempDir := filepath.Dir(filePath) + "/temp"

            switch operator {
                // case "union": 
                case "join": 
                    inputsMapLatter := make(map[string]interface{})
                    csvRowsLatter := utils.GetCsvFromFile(filepath.Join(tempDir, inputsInfos[i + 1]))
                    // suppose we have only one row here
                    for _, csvRow := range csvRowsLatter {
                        for i, _ := range csvRow {
                            inputsMapLatter[csvRowsLatter[0][i]] = csvRowsLatter[1][i]
                        }
                    }
                    inputsMap = mergeMaps(inputsMap, inputsMapLatter)
                    // ==> 
                    joinedFile := filepath.Join(tempDir, "_join.csv")
                    writeMapToCsv(inputsMap, joinedFile)

                // case "append":
            }
        }
    }
    // return inputsFiles
}


func mergeMaps (inputsMap map[string]interface{}, inputsMapLatter map[string]interface{}) map[string]interface{} {
    for k, v := range inputsMapLatter {
        inputsMap[k] = v
    }

    return inputsMap
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





