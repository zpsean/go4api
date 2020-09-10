/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2019.07
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package testcase

import (
    "fmt"
    // "time"
    "os"
    "path/filepath"
    "strings"
    "io/ioutil"
    "strconv"
    "encoding/json"

    "go4api/utils"
)


func InitFullTcSlice (filePathSlice []string) []*TestCaseDataInfo { 
    var fullTcSlice []*TestCaseDataInfo
    var jsonFileList []string

    // tend to support cmd.Opt.Testcase accepting comma delimited paths
    // path istself can be regular expression
    // for example: path1,path2,path3,path4*,...
    for i, _ := range filePathSlice {
        // to support pattern later
        // matches, _ := filepath.Glob(filePathSlice[i])

        jsonFileListTemp, _ := utils.WalkPath(filePathSlice[i], ".json")
        jsonFileList = append(jsonFileList, jsonFileListTemp[0:]...)
    }

    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_dt")
        //
        var tcInfos []*TestCaseDataInfo

        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosWithDt(jsonFile, csvFileList)
        } else {
            tcInfos = ConstructTcInfosWithoutDt(jsonFile)
        }

        for i, _ := range tcInfos {
            fullTcSlice = append(fullTcSlice, tcInfos[i])
        }
    }

    return fullTcSlice
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
        csvFileName := strings.TrimSuffix(filepath.Base(csvFile), ".csv")
        jsonFileName := strings.TrimSuffix(filepath.Base(jsonFile), ".json")
        // Note: the json file realted data table files is pattern: jsonFileName + "_dt[*]"
        if strings.HasPrefix(csvFileName, jsonFileName + suffix) {
            csvFileList = append(csvFileList, csvFile)
        }
    }

    return csvFileList
}

func ConstructTcInfosWithoutDt (jsonFile string) []*TestCaseDataInfo {
    var tcInfos []*TestCaseDataInfo
    var tcases TestCases

    jsonStr := utils.GetJsonFromFile(jsonFile)

    csvFile := ""
    csvRow := ""

    err := json.Unmarshal([]byte(jsonStr), &tcases)
    if err != nil {
        fmt.Println("!! Error, parse Json into cases failed: ", jsonFile, ": ", err)
        os.Exit(1)
    }
  
    for i, _ := range tcases {
        tcaseData := &TestCaseDataInfo {
            TestCase: &tcases[i],
            JsonFilePath: jsonFile,
            CsvFile: csvFile,
            CsvRow: csvRow,
        }
        tcInfos = append(tcInfos, tcaseData)
    }

    return tcInfos
}

// not using "text/template"
func ConstructTcInfosWithDt (jsonFile string, csvFileList []string) []*TestCaseDataInfo {
    var tcInfos []*TestCaseDataInfo

    for _, csvFile := range csvFileList {
        jsonStr := utils.GetJsonFromFile(jsonFile)

        csvRows := utils.GetCsvFromFile(csvFile)
        for i, csvRow := range csvRows {
            jsonR := jsonStr
            // csvRows[0] is the header row
            if i > 0 {
                for col, _ := range csvRow {
                    key := csvRows[0][col]
                    value := csvRows[i][col]

                    // Note: trial, introduce "Fn::ToRawJson" for key, value like:
                    // "jkey": {"Fn::ToRawJson": "${varableInCsv}"},
                    // it is mainly for json map {} or json array []
                    // !!! To be improved further using json key lookup method
                    jsonR = strings.Replace(jsonR, `{"Fn::ToRawJson": "${` + key + `}"}`, fmt.Sprint(value), -1)

                    jsonR = strings.Replace(jsonR, "${" + key + "}", fmt.Sprint(value), -1)  
                }
                
                var tcases TestCases
                err := json.Unmarshal([]byte(jsonR), &tcases)
                if err != nil {
                    fmt.Println("!! Error, parse Json into cases failed: ", jsonFile, ": ", csvFile, ": ", strconv.Itoa(i + 1), ": ", err)
                    os.Exit(1)
                }
    
                for tcI, _ := range tcases {
                    tcaseData := &TestCaseDataInfo {
                        TestCase: &tcases[tcI],
                        JsonFilePath: jsonFile,
                        CsvFile: csvFile,
                        CsvRow: strconv.Itoa(i + 1),
                    }
                    tcInfos = append(tcInfos, tcaseData)
                }
            }
        }
        if IfCaseNameDuplicated(tcInfos) {
            fmt.Println("!! Error, has duplicated case name, please fix before run")
            os.Exit(1)
        }
    }
    if IfCaseNameDuplicated(tcInfos) {
        fmt.Println("!! Error, has duplicated case name, please fix before run")
        os.Exit(1)
    }

    return tcInfos
}

//
func IfCaseNameDuplicated (tcInfos []*TestCaseDataInfo) bool {
    keys := make(map[string]bool)
    var caseNameSet []string

    for _, tcDataInfo := range tcInfos {
        entry := tcDataInfo.TcName()
        if _, ok := keys[entry]; !ok {
            keys[entry] = true
            caseNameSet = append(caseNameSet, entry)
        }
    }

    if len(tcInfos) > len(caseNameSet) {
        return true
    } else {
        return false
    }
}

//
func GetTcNameSet (tcArray []*TestCaseDataInfo) []string {
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

