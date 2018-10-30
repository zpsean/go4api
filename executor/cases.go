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
    // "os"
    "sort"
    // "sync"
    "path/filepath"
    "strings"
    "io/ioutil"
    "strconv"
    "encoding/json"

    "go4api/cmd"
    "go4api/utils"
    "go4api/lib/testcase"
    // "go4api/lib/tree"
    "go4api/texttmpl"
    // "go4api/reports"
)

func GetTcArray () []testcase.TestCaseDataInfo { 
    var tcArray []testcase.TestCaseDataInfo

    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")
        //
        var tcInfos []testcase.TestCaseDataInfo

        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList, "")
        } else {
            tcInfos = ConstructTcInfosBasedOnJson(jsonFile, "")
        }

        for _, tcData := range tcInfos {
            // fmt.Println("\n tcData:", tcData.TcName(), tcData.TestCase.IfGlobalSetUpTestCase())
            tcArray = append(tcArray, tcData)
        }
    }

    return tcArray
}

func GetNormalTcSlice (tcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var tcSlice []testcase.TestCaseDataInfo
    for i, _ := range tcArray {
        if tcArray[i].TestCase.IfGlobalSetUpTestCase() != true && tcArray[i].TestCase.IfGlobalTearDownTestCase() != true {
            tcSlice = append(tcSlice, tcArray[i])
        }
    }
    
    return tcSlice
}


func GetCsvDataFilesForJsonFile (jsonFile string, suffix string) []string {
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


func ConstructTcInfosBasedOnJsonTemplateAndDataTables (jsonFile string, csvFileList []string, parentTcName string) []testcase.TestCaseDataInfo {
    var tcInfos []testcase.TestCaseDataInfo

    for _, csvFile := range csvFileList {
        // to check the csv file's existence
        csvRows := utils.GetCsvFromFile(csvFile)
        for i, csvRow := range csvRows {
            // starting with data row
            if i > 0 {
                // note: here pass the csvRows[0], csvRow, but they can be replaced by map[string]interface{} for later enhancement
                var cvsRowInterface []interface{}
                for rowI, _ := range csvRow {
                    cvsRowInterface = append(cvsRowInterface, csvRow[rowI])
                }
                mergedTestData := MergeTestData(csvRows[0], cvsRowInterface, parentTcName)

                outTempJson := texttmpl.GenerateJsonBasedOnTemplateAndCsv(jsonFile, mergedTestData)

                var tcases testcase.TestCases
                resJson, _ := ioutil.ReadAll(outTempJson)
                err := json.Unmarshal([]byte(resJson), &tcases)
                if err != nil {
                    fmt.Println("!! error, Json format validation : ", jsonFile, ": ", csvFile, ": ", strconv.Itoa(i + 1), ": ", err)
                }
                // as the json is generated based on templated dynamically, so that, to cache all the resulted json in array
                for tcI, _ := range tcases {
                    // populate the testcase.TestCaseDataInfo
                    tcaseData := testcase.TestCaseDataInfo {
                        TestCase: &tcases[tcI],
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

func ConstructTcInfosBasedOnJson (jsonFile string, parentTcName string) []testcase.TestCaseDataInfo {
    var tcInfos []testcase.TestCaseDataInfo

    csvFile := ""
    csvRow := ""
    // mergedTestData := map[string]interface{}{}
    mergedTestData := MergeTestData([]string{""}, []interface{}{}, parentTcName)

    outTempJson := texttmpl.GenerateJsonBasedOnTemplateAndCsv(jsonFile, mergedTestData)
    
    var tcases testcase.TestCases
    resJson, _ := ioutil.ReadAll(outTempJson)
    err := json.Unmarshal([]byte(resJson), &tcases)
    if err != nil {
        fmt.Println("!! error, Json format validation : ", jsonFile, ": ", err)
    }
    // fmt.Println("resJson: ", string(resJson), tcases)
    // tJson, _ := json.Marshal(tcases)
    // fmt.Println("tJson: ", string(tJson))
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


func GetTcNameSet (tcArray []testcase.TestCaseDataInfo) []string {
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


func GetPrioritySet (tcArray []testcase.TestCaseDataInfo) []string {
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

    prioritySet = SortPrioritySet(prioritySet)

    return prioritySet
}

func SortPrioritySet (prioritySet []string) []string {
    // Note: here is a bug, if sort as string, as the sort results is 1, 10, 11, 2, 3, etc. => fixed
    prioritySet_Int := utils.ConvertStringArrayToIntArray(prioritySet)
    sort.Ints(prioritySet_Int)
    prioritySet = utils.ConvertIntArrayToStringArray(prioritySet_Int)

    return prioritySet
}



func GetTestCasesByPriority (prioritySet []string, tcArray []testcase.TestCaseDataInfo) (map[string][]testcase.TestCaseDataInfo, map[string]int) {
    // build the map
    tcClassifedMap := make(map[string][]testcase.TestCaseDataInfo)
    tcClassifedCountMap := make(map[string]int)

    for _, entry := range prioritySet {
        for _, tcaseData := range tcArray {
            if entry == tcaseData.Priority() {
                tcClassifedMap[entry] = append(tcClassifedMap[entry], tcaseData)
                tcClassifedCountMap[entry] += 1
            }
        }
    }

    return tcClassifedMap, tcClassifedCountMap
}


func GetFuzzTcArray () []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo

    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_fuzz_dt")
        // to get the json test data directly (if not template) based on template (if template)
        // tcInfos: [[casename, priority, parentTestCase, ], ...]
        var tcInfos []testcase.TestCaseDataInfo
        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList, "")
        }

        for _, tcData := range tcInfos {
            tcArray = append(tcArray, tcData)
        }
    }

    return tcArray
}


func GetOriginMutationTcArray () []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo
    jsonFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".json")

    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")

        var tcInfos []testcase.TestCaseDataInfo
        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList, "")
        } else {
            tcInfos = ConstructTcInfosBasedOnJson(jsonFile, "")
        }

        for _, tcData := range tcInfos {
            tcArray = append(tcArray, tcData)
        }
    }

    // tJson, _ := json.MarshalIndent(tcArray, "", "\t")
    // fmt.Println("tcArray: ", string(tJson))

    return tcArray
}


