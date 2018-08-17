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
    // "go4api/types"
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
        GenerateFuzzDataFiles(fuzzFile)
        // fuzzDataFileList :=

    }

    // (2). render the json using the fuzz dt(s)
    fuzzTcArray := GetFuzzTcArray(options)
    fmt.Println("fuzzTcArray:", fuzzTcArray, "\n")
    // (3). execute the rendered json, 


    // JSON Schema defines the following basic types:
    // string
    // Numeric -> float64
    // boolean
    // null
    // object (raw)
    // array (raw)
    
    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    return fuzzTcArray
}

func GenerateFuzzDataFiles(fuzzFile string) {
    fuzzRowsByte := utils.GetContentFromFile(fuzzFile)

    fuzzRows := strings.Split(string(fuzzRowsByte), "\n")
    fmt.Println(fuzzRows)

    for _, fuzzLine := range fuzzRows {
        fmt.Println(fuzzLine)
        fieldName := "title"
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



        var validValueList []string
        var invalidValueList []string

        for _, validLen := range validLenList{
            validValue := mode.CallRands(fieldType, validLen)
            fmt.Println("validLen, validValue: ", validLen, validValue)
            validValueList = append(validValueList, validValue)
        }

        for _, invalidLen := range invalidLenList{
            invalidValue := mode.CallRands(fieldType, invalidLen)
            fmt.Println("invalidLen, invalidValue: ", invalidLen, invalidValue)
            invalidValueList = append(invalidValueList, invalidValue)  
        }

        // fmt.Println("validValueList: ", validValueList)
        // fmt.Println("invalidValueList: ", invalidValueList)


        outputsFile := filepath.Join(filepath.Dir(fuzzFile), 
            strings.TrimRight(filepath.Base(fuzzFile), ".fuzz") + "_fuzz_dt_valid.csv")

        utils.GenerateFileBasedOnVarOverride("tcid"+ "," + fieldName + "\n", outputsFile)

        // validValueStrList := ConvertSliceInterfaceToSliceString(validValueList)
        for i, value := range validValueList {
            utils.GenerateFileBasedOnVarAppend("valid" + strconv.Itoa(i) + "," + value + "\n", outputsFile)
        }
        

        // write csv data
        outputsFile = filepath.Join(filepath.Dir(fuzzFile), 
            strings.TrimRight(filepath.Base(fuzzFile), ".fuzz") + "_fuzz_dt_invalid.csv")

        utils.GenerateFileBasedOnVarOverride("tcid"+ "," + fieldName + "\n", outputsFile)

        // invalidValueStrList := ConvertSliceInterfaceToSliceString(invalidValueList)
        for i, value := range invalidValueList {
            utils.GenerateFileBasedOnVarAppend("invalid" + strconv.Itoa(i) + "," + value + "\n", outputsFile)
        }


    }
}



// func ConvertSliceInterfaceToSliceString(params []interface{}) []string {
//     var paramSlice []string
//     for _, param := range params {
//         paramSlice = append(paramSlice, (reflect.ValueOf(param)).String())
//         // fmt.Println("param: ", reflect.TypeOf(param), reflect.ValueOf(param))
//     }

//     return paramSlice
// }


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



