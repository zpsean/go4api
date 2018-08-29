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
    "path/filepath"
    "strings"
    "strconv"
    "go4api/utils"  
)

// valid, invalid data may have more than one field, but the map itself can not ensure the key sequence
// so that, here use slice
type FuzzData struct {  
    ValidData []map[string][]interface{}
    InvalidData []map[string][]interface{}
    ValidStatusCode int
    InvalidStatusCode int
}

type FieldDefinition struct {  
    FieldName string
    FieldType string
    FieldSubType string
    FieldMin int
    FieldMax int
}


func PrepFuzzTest(pStart_time time.Time, options map[string]string) {
    fuzzFileList, _ := utils.WalkPath(options["testhome"] + "/testdata/", ".fuzz")
    fmt.Println("FuzzTest jsonFileList:", options["ifFuzzTestFirst"], fuzzFileList, "")

    // (1). generate the data tables based on the fuzz test, at least two dt files: positive and negative
    for _, fuzzFile := range fuzzFileList {
        fuzzData := GenerateFuzzData(fuzzFile)
        GenerateFuzzDataFiles(fuzzFile, fuzzData)
    }
    // (2). render the json using the fuzz dt(s)
    // fuzzTcArray := GetFuzzTcArray(options)
    // fmt.Println("fuzzTcArray:", fuzzTcArray, "\n")
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

    var fuzzData FuzzData
    var validValueList []map[string][]interface{}
    var invalidValueList []map[string][]interface{}

    for _, fuzzLine := range fuzzRows {
        if len(strings.TrimSpace(fuzzLine)) > 0 {
            validValueMap := make(map[string][]interface{})
            invalidValueMap := make(map[string][]interface{})

            // fmt.Println("\nfuzzLine: ", fuzzLine)
            fieldDefinition := parseLine(fuzzLine)

            switch strings.ToLower(fieldDefinition.FieldType) {
                case "char", "varchar", "string": {
                    fmt.Println("\n------ char -")
                    validValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules("FuzzCharValid")
                    invalidValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules("FuzzCharInvalid")
                }
                case "int", "int64": {
                    fmt.Println("\n------ int -")
                    validValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules("FuzzCharValid")
                    invalidValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules("FuzzCharInvalid")
                }
                default: {
                    fmt.Println("\n------ default -")
                    validValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules("FuzzCharValid")
                    invalidValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules("FuzzCharInvalid")
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

    fuzzData = FuzzData {
        ValidData: validValueList,
        InvalidData: invalidValueList,
        ValidStatusCode: 200,
        InvalidStatusCode: 200,
    }

    fmt.Println("fuzzData: ", fuzzData)

    return fuzzData
}


func GenerateFuzzDataFiles(fuzzFile string, fuzzData FuzzData) {
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
    for _, subCombValid := range combValid {
        fmt.Println("subCombValid -- : ", subCombValid, len(subCombValid))
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
    utils.GenerateFileBasedOnVarOverride(invalidHeaderStr + "\n", outputsFile)

    combInvalid := GetCombinationInvalid(fuzzData)
    //
    i = 1
    for _, subCombInvalid := range combInvalid {
        fmt.Println("subCombInvalid: ", subCombInvalid, len(subCombInvalid))
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




func parseLine(fuzzLine string) FieldDefinition {
    var fieldDefinition FieldDefinition

    line := strings.Split(fuzzLine, ":")

    fieldDefinition.FieldName = strings.TrimSpace(line[0])

    if strings.Index(line[1], "(") > 0 {
        fieldDefinition.FieldType = strings.TrimSpace(line[1][0:strings.Index(line[1], "(")])
    } else {
        fieldDefinition.FieldType = "float64"
    }

    fieldDefinition.FieldSubType = ""
    fieldDefinition.FieldMin = 0
    fieldDefinition.FieldMax = 20

    return fieldDefinition
}



