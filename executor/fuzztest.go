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
    // "reflect"
    // "sync"
    "strings"
    "strconv"
    "go4api/types"
    "go4api/utils"
    "go4api/utils/mode"
    "path/filepath"
    // "go4api/logger"
    // "encoding/json"
)


func PrepFuzzTest(ch chan int, pStart_time time.Time, options map[string]string) [][]interface{} {
    fuzzFileList, _ := utils.WalkPath(options["testhome"] + "/testdata/", ".fuzz")
    fmt.Println("FuzzTest jsonFileList:", options["ifFuzzTestFirst"], fuzzFileList, "\n")

    // (1). generate the data tables based on the fuzz test, at least two dt files: positive and negative
    for _, fuzzFile := range fuzzFileList {
        fuzzData := GenerateFuzzData(fuzzFile)
        GenerateFuzzDataFiles(fuzzFile, fuzzData)
    }
    // (2). render the json using the fuzz dt(s)
    fuzzTcArray := GetFuzzTcArray(options)
    fmt.Println("fuzzTcArray:", fuzzTcArray, "\n")
    
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
func GenerateFuzzData(fuzzFile string) types.FuzzData {
    fuzzRowsByte := utils.GetContentFromFile(fuzzFile)

    fuzzRows := strings.Split(string(fuzzRowsByte), "\n")
    fmt.Println(fuzzRows) 

    var fuzzData types.FuzzData
    var validValueList []map[string][]string
    var invalidValueList []map[string][]string

    for i, fuzzLine := range fuzzRows {
        validValueMap := make(map[string][]string)
        invalidValueMap := make(map[string][]string)

        fmt.Println(fuzzLine)
        fieldName := "title" + strconv.Itoa(i) 
        fieldType := "RandStringRunes"
        fieldMin := 0
        fieldMax := 20

        // get the Boundary (valid, invalid), Equivalence, etc.
        var validLenList []int
        var invalidLenList []int

        validLenList = append(validLenList, fieldMin)
        validLenList = append(validLenList, fieldMin + 1)
        
        validLenList = append(validLenList, fieldMax)
        if fieldMax - 1 > fieldMin {
            validLenList = append(validLenList, fieldMax - 1)
        }

        if fieldMin - 1 > 0 {
           invalidLenList = append(invalidLenList, fieldMin - 1) 
        }
        invalidLenList = append(invalidLenList, fieldMax + 1) 

        //
        for _, validLen := range validLenList{
            validValue := mode.CallRands(fieldType, validLen)
            fmt.Println("validLen, validValue: ", validLen, validValue)

            validValueMap[fieldName] = append(validValueMap[fieldName], validValue)
        }

        for _, invalidLen := range invalidLenList{
            invalidValue := mode.CallRands(fieldType, invalidLen)
            fmt.Println("invalidLen, invalidValue: ", invalidLen, invalidValue)

            invalidValueMap[fieldName] = append(invalidValueMap[fieldName], invalidValue)
        } 

        validValueList = append(validValueList, validValueMap)
        invalidValueList = append(invalidValueList, validValueMap)

    }

    fuzzData = types.FuzzData {
        ValidData: invalidValueList,
        InvalidData: invalidValueList,
        ValidStatusCode: 200,
        InvalidStatusCode: 200,
    }

    fmt.Println("fuzzData: ", fuzzData)

    return fuzzData
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
    //
    combValid := mode.GetCombinationValid(fuzzData)
    //
    i := 1
    for subCombValid := range combValid {
        fmt.Println("subCombValid: ", subCombValid, len(subCombValid))
        combStr := strings.Join(subCombValid, ",")
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

    //     for i, value := range valueList {
    //         utils.GenerateFileBasedOnVarAppend("invalid" + strconv.Itoa(i + 1) + "," + value + "\n", outputsFile)
    // }
    // Get GetCombinationInvalid(fuzzData)
}


func GetFuzzTcArray(options map[string]string) [][]interface{} {
    var tcArray [][]interface{}

    jsonFileList, _ := utils.WalkPath(options["testhome"] + "/testdata/", ".json")
    // fmt.Println("jsonFileList:", jsonFileList, "\n")
    // to ge the json and related data file, then get tc from them
    for _, jsonFile := range jsonFileList {
        csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_fuzz_dt")
        // to get the json test data directly (if not template) based on template (if template)
        // tcInfos: [[casename, priority, parentTestCase, ], ...]
        var tcInfos [][]interface{}
        if len(csvFileList) > 0 {
            tcInfos = ConstructTcInfosBasedOnJsonTemplateAndDataTables(jsonFile, csvFileList)
        }
        // fmt.Println("tcInfos:", tcInfos, "\n")
        
        for _, tc := range tcInfos {
            tcArray = append(tcArray, tc)
        }
    }

    return tcArray
}





// func ConvertSliceInterfaceToSliceString(params []interface{}) []string {
//     var paramSlice []string
//     for _, param := range params {
//         paramSlice = append(paramSlice, (reflect.ValueOf(param)).String())
//         // fmt.Println("param: ", reflect.TypeOf(param), reflect.ValueOf(param))
//     }

//     return paramSlice
// }




