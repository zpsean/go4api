/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package fuzz

import (                                                                                                                                             
    // "os"
    "time"
    "fmt"
    // "reflect"
    // "sync"
    "path/filepath"
    "strings"
    "strconv"
    "go4api/testcase"
    "go4api/utils"  
    // "go4api/logger"
    // "encoding/json"
)

// valid, invalid data may have more than one field, but the map itself can not ensure the key sequence
// so that, here use slice
type FuzzData struct {  
    ValidData []map[string][]interface{}
    InvalidData []map[string][]interface{}
    ValidStatusCode int
    InvalidStatusCode int
}


func PrepFuzzTest(ch chan int, pStart_time time.Time, options map[string]string) []testcase.TestCaseDataInfo {
    fuzzFileList, _ := utils.WalkPath(options["testhome"] + "/testdata/", ".fuzz")
    fmt.Println("FuzzTest jsonFileList:", options["ifFuzzTestFirst"], fuzzFileList, "\n")

    // (1). generate the data tables based on the fuzz test, at least two dt files: positive and negative
    for _, fuzzFile := range fuzzFileList {
        fuzzData := GenerateFuzzData(fuzzFile)
        GenerateFuzzDataFiles(fuzzFile, fuzzData)
    }
    // (2). render the json using the fuzz dt(s)
    fuzzTcArray := GetFuzzTcArray(options)
    // fmt.Println("fuzzTcArray:", fuzzTcArray, "\n")
    
    return fuzzTcArray
}

// JSON Schema defines the following basic types:
// string
// Numeric -> float64
// boolean
// null
// object (raw)
// array (raw)

// to get the fuzz data table files with naming fuzzcase_fuzz_dt_valid.csv / fuzzcase_fuzz_dt_invalid.csv
func GenerateFuzzData(fuzzFile string) FuzzData {
    fuzzRowsByte := utils.GetContentFromFile(fuzzFile)

    fuzzRows := strings.Split(string(fuzzRowsByte), "\n")

    var fuzzData types.FuzzData
    var validValueList []map[string][]interface{}
    var invalidValueList []map[string][]interface{}

    for _, fuzzLine := range fuzzRows {
        if len(strings.TrimSpace(fuzzLine)) > 0 {
            validValueMap := make(map[string][]interface{})
            invalidValueMap := make(map[string][]interface{})

            // fmt.Println("\nfuzzLine: ", fuzzLine)

            fieldName, fieldType, fieldMin, fieldMax := parseLine(fuzzLine)

            switch strings.ToLower(fieldType) {
                case "char", "varchar", "string": {
                    fmt.Println("\n------ char -")
                    validValueMap, invalidValueMap = getChar(fieldName, fieldType, fieldMin, fieldMax)
                }
                case "int", "int64": {
                    fmt.Println("\n------ int -")
                    validValueMap, invalidValueMap = getInt(fieldName, fieldType, fieldMin, fieldMax)
                }
                default: {
                    fmt.Println("\n------ default -")
                    validValueMap, invalidValueMap = getChar(fieldName, fieldType, fieldMin, fieldMax)
                }
                // case numeric
                // case email
                // case float
                // case list
            }

            validValueList = append(validValueList, validValueMap)
            invalidValueList = append(invalidValueList, invalidValueMap)
        }
    }

    fuzzData = types.FuzzData {
        ValidData: validValueList,
        InvalidData: invalidValueList,
        ValidStatusCode: 200,
        InvalidStatusCode: 200,
    }

    fmt.Println("fuzzData: ", fuzzData)

    return fuzzData
}


func getChar(fieldName string, fieldType string, fieldMin int, fieldMax int) (map[string][]interface{}, map[string][]interface{}) {
    validValueMap := make(map[string][]interface{})
    invalidValueMap := make(map[string][]interface{})
    // get the Boundary (valid, invalid), Equivalence, etc.
    var validLenList []int
    var invalidLenList []int
    //
    validLenList = append(validLenList, fieldMin)
    validLenList = append(validLenList, fieldMin + 1)

    validLenList = append(validLenList, fieldMax)
    if fieldMax - 1 > fieldMin {
        validLenList = append(validLenList, fieldMax - 1)
    }
    //
    if fieldMin - 1 > 0 {
       invalidLenList = append(invalidLenList, fieldMin - 1) 
    }
    invalidLenList = append(invalidLenList, fieldMax + 1) 
    //

    fieldRands := []string{"RandStringRunes", "RandStringCNRunes"}
    //
    for _, validLen := range validLenList{
        for _, randType := range fieldRands {
            validValue := CallRands(randType, validLen)
            // fmt.Println("validLen, validValue: ", validLen, validValue)

            validValueMap[fieldName] = append(validValueMap[fieldName], validValue)
        }        
    }
    //
    for _, invalidLen := range invalidLenList{
        for _, randType := range fieldRands {
            invalidValue := CallRands(randType, invalidLen)
            // fmt.Println("invalidLen, invalidValue: ", invalidLen, invalidValue)

            invalidValueMap[fieldName] = append(invalidValueMap[fieldName], invalidValue)
        }
    }

    return validValueMap, invalidValueMap
}


func getInt(fieldName string, fieldType string, fieldMin int, fieldMax int) (map[string][]interface{}, map[string][]interface{}) {
    validValueMap := make(map[string][]interface{})
    invalidValueMap := make(map[string][]interface{})
    // get the Boundary (valid, invalid), Equivalence, etc.
    validValueMap[fieldName] = append(validValueMap[fieldName], fieldMin)
    validValueMap[fieldName] = append(validValueMap[fieldName], fieldMin + 1)
    validValueMap[fieldName] = append(validValueMap[fieldName], fieldMax)
    validValueMap[fieldName] = append(validValueMap[fieldName], fieldMax - 1)
    //
    invalidValueMap[fieldName] = append(invalidValueMap[fieldName], fieldMin - 1)
    invalidValueMap[fieldName] = append(invalidValueMap[fieldName], fieldMax + 1)
    //

    return validValueMap, invalidValueMap
}


func GenerateFuzzDataFiles(fuzzFile string, fuzzData types.FuzzData) {
    // fmt.Println("validValueList: ", validValueList)
    // fmt.Println("invalidValueList: ", invalidValueList)

    // for valid data
    outputsFile := filepath.Join(filepath.Dir(fuzzFile), 
        strings.TrimRight(filepath.Base(fuzzFile), ".fuzz") + "_fuzz_dt_valid.csv")
    // write csv header, data
    var validHeaderStr string
    validHeaderStr = validHeaderStr + "tcid"
    for _, validDataMap := range fuzzData.ValidData {
        for key, _ := range validDataMap {
            validHeaderStr = validHeaderStr + "," + key
        }
    }
    utils.GenerateFileBasedOnVarOverride(validHeaderStr + "\n", outputsFile)
    
    combValid := GetCombinationValid(fuzzData)
    //
    i := 1
    for subCombValid := range combValid {
        // fmt.Println("subCombValid -- : ", subCombValid, len(subCombValid))
        combStr := ""
        for ii, item := range subCombValid {
            if ii == 0 {
                combStr = combStr + fmt.Sprint(item)
            } else{
                combStr = combStr + "," + fmt.Sprint(item)
            }
        }
        utils.GenerateFileBasedOnVarAppend("valid" + strconv.Itoa(i) + "," + combStr + "\n", outputsFile)  
        i = i + 1  
    }
    

    // for invalid data
    outputsFile = filepath.Join(filepath.Dir(fuzzFile), 
        strings.TrimRight(filepath.Base(fuzzFile), ".fuzz") + "_fuzz_dt_invalid.csv")
    var invalidHeaderStr string
    invalidHeaderStr = invalidHeaderStr + "tcid"
    // write csv header, data
    for _, invalidDataMap := range fuzzData.InvalidData {
        for key, _ := range invalidDataMap {
            invalidHeaderStr = invalidHeaderStr + "," + key
        }
    }
    utils.GenerateFileBasedOnVarOverride(invalidHeaderStr, outputsFile)

    combInvalid := GetCombinationInvalid(fuzzData)
    //
    i = 1
    for subCombInvalid := range combInvalid {
        // fmt.Println("subCombInvalid: ", subCombInvalid, len(subCombInvalid))
        combStr := ""
        for ii, item := range subCombInvalid {
            if ii == 0 {
                combStr = combStr + fmt.Sprint(item)
            } else{
                combStr = combStr + "," + fmt.Sprint(item)
            }
        }
        utils.GenerateFileBasedOnVarAppend("invalid" + strconv.Itoa(i) + "," + combStr + "\n", outputsFile)  
        i = i + 1  
    }
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


func parseLine(fuzzLine string) (string, string, int, int) {
    var fieldName, fieldType string


    line := strings.Split(fuzzLine, ":")

    fieldName = strings.TrimSpace(line[0])

    if strings.Index(line[1], "(") > 0 {
        fieldType = strings.TrimSpace(line[1][0:strings.Index(line[1], "(")])
    }

    fmt.Print("fieldName, fieldType: ", fieldName, fieldType)

    return fieldName, fieldType, 0, 20
}


// func ConvertSliceInterfaceToSliceString(params []interface{}) []string {
//     var paramSlice []string
//     for _, param := range params {
//         paramSlice = append(paramSlice, (reflect.ValueOf(param)).String())
//         // fmt.Println("param: ", reflect.TypeOf(param), reflect.ValueOf(param))
//     }

//     return paramSlice
// }




