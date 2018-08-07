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
    // "time"
    "os"
    "sort"
    "sync"
    // "go4api/api"  
    "go4api/ui"     
    "go4api/ui/js"  
    "go4api/ui/style"                                                                                                                                
    "go4api/utils"
    "go4api/utils/texttmpl"
    "path/filepath"
    "strings"
    "io/ioutil"
    "strconv"
    // simplejson "github.com/bitly/go-simplejson"
)


func Run(ch chan int, pStart string, options map[string]string) { //client
    //
    testenv := options["testEnv"]
    baseUrl := ""
    if options["baseUrl"] != "" {
        baseUrl = options["baseUrl"]
    } else {
        _, err := os.Stat(options["testhome"] + "/testconfig/testconfig.json")
        // fmt.Println("err: ", err)
        if err == nil {
            baseUrl = utils.GetBaseUrlFromConfig(options["testhome"] + "/testconfig/testconfig.json", testenv) 
        }
    }
    if baseUrl == "" {
        fmt.Println("Warning: baseUrl is not set")
    } else {
        fmt.Println("baseUrl set to: " + baseUrl)
    }
    // get results dir
    resultsDir := GetResultsDir(pStart, options)
    //
    // (1), get the text path, default is ../data/*, then search all the sub-folder to get the test scripts
    //
    tcArray := GetTcArray(options)
    // fmt.Println("tcArray:", tcArray, "\n")
    // myabe there needs a scheduler, for priority 1 (w or w/o dependency) -> priority 2 (w or w/o dependency), ...
    // --
    // How to impliment the case Dependency???
    // Two big categories: 
    // (1) case has No parent Dependency or successor Dependency, which can be scheduled concurrently
    // (2) case has parent Dependency or successor Dependency, which has rules to be scheduled concurrently
    // 
    // need a tree to track and schedule the run dynamiclly, but need a dummy root test case
 
    // dummy root tc => {"root", "0", "0", rooTC, "", "", ""}
    root, _ := BuildTree(tcArray)
    fmt.Println("------------------")
    //
    prioritySet := GetPrioritySet(tcArray)
    // classifications := GetTestCasesByPriority(prioritySet, tcArray)
    // Note, before starting execution, needs to sort the priorities_set first by priority
    // Note: here is a bug, as the sort results is 1, 10, 11, 2, 3, etc. => fixed
    prioritySet_Int := utils.ConvertStringArrayToIntArray(prioritySet)
    sort.Ints(prioritySet_Int)
    prioritySet = utils.ConvertIntArrayToStringArray(prioritySet_Int)

    // If need to set the Concurrency MAX?
    // fmt.Println("------------------", root, &root)
    // ShowNodes(root)
    // fmt.Println("------------------", root, &root)
    InitNodesRunResult(root, "Ready")
    // fmt.Println("------------------", root, &root)
    // ShowNodes(root)
    // fmt.Println("------------------", root, &root)

    //
    fmt.Println("\n====> test cases execution starts!\n")
    tcTotalCount = 0
    statusReadyCount = 0
    statusSuccessCount = 0
    statusFailCount = 0
    statusOtherCount = 0
    //
    for _, priority := range prioritySet {
        // tcArrayPriority := classifications[priority]
        fmt.Println("====> Priority " + priority + " starts!")
        
        miniLoop:
        for {
            //
            resultsChan := make(chan []interface{}, len(tcArray))
            var wg sync.WaitGroup
            //
            ScheduleNodes(root, &wg, options, priority, resultsChan, pStart, baseUrl, resultsDir)
            //
            wg.Wait()

            close(resultsChan)

            for tcRunResults := range resultsChan {
                // tcName, testResult
                SearchNode(&root, tcRunResults[0].(string))
                // fmt.Println("---- the found node: ", tcRunResults[0].(string), *findNode)
                RefreshNodeAndDirectChilrenTcResult(*findNode, tcRunResults[1].(string))
                // fmt.Println("------------------")
                // ShowNodes(root)
                }
            // if tcTree has no node with "Ready" status, break the miniloop
            statusReadyCount = 0
            //
            CollectNodeStatus(root, priority)
            // fmt.Println("------------------ statusReadyCount: ", statusReadyCount)
            if statusReadyCount == 0 {
                break miniLoop
            }
        }

        fmt.Println("====> Priority " + priority + " ended! \n")
        // sleep for debug
        // time.Sleep(1 * time.Second)
    }

    //
    fmt.Println("---------------------------------------------------------------------------")
    fmt.Println("----- Total " + strconv.Itoa(len(tcArray)) + " Cases in template -----")
    fmt.Println("----- Total " + strconv.Itoa(tcTotalCount) + " Cases put onto tcTree -----")
    fmt.Println("----- Total " + strconv.Itoa(statusOtherCount) + " Cases Skipped (Not Executed, due to Parent Failed) -----")
    fmt.Println("----- Total " + strconv.Itoa(tcTotalCount - statusOtherCount) + " Cases Executed -----")
    fmt.Println("----- Total " + strconv.Itoa(statusSuccessCount) + " Cases Success -----")
    fmt.Println("----- Total " + strconv.Itoa(statusFailCount) + " Cases Fail -----")
    fmt.Println("---------------------------------------------------------------------------\n\n")


    // generate the html report based on template, and results data
    // time.Sleep(1 * time.Second)
    GenerateTestReport(resultsDir, pStart)
    // fmt.Println("====> Report Generated!\n")

    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    ch <- 1

}



func GetTcArray(options map[string]string) [][]interface{} {
    var tcArray [][]interface{}
    jsonFileList, _ := utils.WalkPath(options["testhome"] + "/testdata/", ".json")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        // here search out the csv files under the same dir, not to use utils.WalkPath as it is recursively
        var csvFileListTemp []string
        infos, err := ioutil.ReadDir(filepath.Dir(jsonFile))
        if err != nil {
          panic(err)
        }
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
            if strings.Contains(csvFileName, jsonFileName + "_dt") {
                csvFileList = append(csvFileList, csvFile)
            }
        }
        // to get the json test data directly (if not template) based on template (if template)
        // tcInfos: [[casename, priority, parentTestCase, ], ...]
        var tcInfos [][]interface{}
        if len(csvFileList) > 0 {
            for _, csvFile := range csvFileList {
                csvRows := utils.GetCsvFromFile(csvFile)
                for i, csvRow := range csvRows {
                    // starting with data row
                    if i > 0 {
                        // outTempFile := texttmpl.GenerateJsonFileBasedOnTemplateAndCsv(jsonFile, csvRows[0], csvRow, tmpJsonDir)
                        // tcJsonsTemp := utils.GetTestCaseJsonFromTestDataFile(outTempFile)
                        outTempJson := texttmpl.GenerateJsonBasedOnTemplateAndCsv(jsonFile, csvRows[0], csvRow)
                        tcJsonsTemp := utils.GetTestCaseJsonFromTestData(outTempJson)
                        // as the json is generated based on templated dynamically, so that, to cache all the resulted json in array
                        var tcInfo []interface{}
                        for _, tc := range tcJsonsTemp {
                            // to get the case info like [casename, priority, parentTestCase, ...], tc, jsonFile, csvFile, row in csv
                            // Note: row in csv = i + 1 (i.e. plus csv header line)
                            tcInfo = utils.GetTestCaseBasicInfoFromTestData(tc)
                            // append last field: tcRunResult, it is tc[7]
                            tcInfo = append(tcInfo, tc, jsonFile, csvFile, strconv.Itoa(i + 1), "")
                            tcInfos = append(tcInfos, tcInfo)
                        }
                    }
                }
            }
        } else {
            csvFile := ""
            csvRow := ""
            // outTempFile := texttmpl.GenerateJsonFileBasedOnTemplateAndCsv(jsonFile, []string{""}, []string{""}, tmpJsonDir)
            // tcJsonsTemp := utils.GetTestCaseJsonFromTestDataFile(outTempFile)
            outTempJson := texttmpl.GenerateJsonBasedOnTemplateAndCsv(jsonFile, []string{""}, []string{""})
            tcJsonsTemp := utils.GetTestCaseJsonFromTestData(outTempJson)
            // as the json is generated based on templated dynamically, so that, to cache all the resulted json in array
            var tcInfo []interface{}
            for _, tc := range tcJsonsTemp {
                // to get the case info like [casename, priority, parentTestCase, ...]
                tcInfo = utils.GetTestCaseBasicInfoFromTestData(tc)
                // append last field: tcRunResult, it is tc[7]
                tcInfo = append(tcInfo, tc, jsonFile, csvFile, csvRow, "")
                tcInfos = append(tcInfos, tcInfo)
            }
        }

        // fmt.Println("tcInfos:", tcInfos, "\n")
        
        for _, tc := range tcInfos {
            tcArray = append(tcArray, tc)
        }
    }

    return tcArray
}


func GetPrioritySet(tcArray [][]interface{}) []string {
    // get the priorities
    var priorities []interface{}
    for _, tc := range tcArray {
        priorities = append(priorities, tc[1])
    }
    // go get the distinct key in priorities
    keys := make(map[string]bool)
    prioritySet := []string{}
    for _, entry := range priorities {
        // uses 'value, ok := map[key]' to determine if map's key exists, if ok, then true
        if _, value := keys[entry.(string)]; !value {
            keys[entry.(string)] = true
            prioritySet = append(prioritySet, entry.(string))
        }
    }

    return prioritySet
}

func GetTestCasesByPriority(prioritySet []string, tcArray [][]interface{}) map[string][][]interface{} {
    // build the map
    classifications := make(map[string][][]interface{})
    for _, entry := range prioritySet {
        for _, tc := range tcArray {
            // tc[1] represents the priority
            if entry == tc[1] {
                classifications[entry] = append(classifications[entry], tc)
            }
        }
    }
    // fmt.Println("classifications: ", classifications)
    return classifications
}


func GenerateTestReport(resultsDir string, pStart string) {
    // read the resource under /ui/*
    // fmt.Println("ui: ", ui.Index_template)

    texttmpl.GenerateHtmlReportFromTemplateAndVar(ui.Index_template, resultsDir, resultsDir + pStart + ".log")
    //
    err := os.MkdirAll(resultsDir + "js", 0777)
    if err != nil {
      panic(err) 
    }
    utils.GenerateFileBasedOnVar(js.Js, resultsDir + "js/go4api.js")
    //
    err = os.MkdirAll(resultsDir + "style", 0777)
    if err != nil {
      panic(err) 
    }
    utils.GenerateFileBasedOnVar(style.Style, resultsDir + "style/go4api.css")
}

func GetResultsDir(pStart string, options map[string]string) string {
    var resultsDir string
    err := os.MkdirAll(options["testresults"] + "/" + pStart + "/", 0777)
    if err != nil {
      panic(err) 
    } else {
        resultsDir = options["testresults"] + "/" + pStart + "/"
    }

    return resultsDir
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

