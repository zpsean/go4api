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
    "go4api/cmd"
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


func PrepFuzzTest(pStart_time time.Time) {
    fuzzFileList, _ := utils.WalkPath(cmd.Opt.Testcase + "/testdata/", ".fuzz")
    // fmt.Println("FuzzTest jsonFileList:", options["ifFuzzTestFirst"], fuzzFileList, "")

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

            fType := fieldDefinition.DetermineFuzzType()
            // call the rules to get values
            validValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules(fType)
            invalidValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules("FuzzCharInvalid")
            // append to slice
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
    // (1). for valid data ------------------------
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

    // this is to get the combinations, maybe use pairwise
    combValid := GetCombinationValid(fuzzData)
    //
    i := 0
    tcid := ""
    for _, subCombValid := range combValid {
        i = i + 1
        tcid = "valid" + strconv.Itoa(i)

        fmt.Println("subCombValid -- : ", subCombValid, len(subCombValid))
        combStr := ""
        for ii, item := range subCombValid {
            if ii == 0 {
                combStr = combStr + fmt.Sprint(item)
            } else{
                combStr = combStr + "," + fmt.Sprint(item)
            }
        }
        utils.GenerateFileBasedOnVarAppend(tcid + "," + combStr + "\n", outputsFile)
    }
    

    // (2). for invalid data ------------------------
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
    i = 0
    tcid = ""
    for _, subCombInvalid := range combInvalid {
        i = i + 1
        tcid = "invalid" + strconv.Itoa(i)

        fmt.Println("subCombInvalid: ", subCombInvalid, len(subCombInvalid))
        combStr := ""
        for ii, item := range subCombInvalid {
            if ii == 0 {
                combStr = combStr + fmt.Sprint(item)
            } else{
                combStr = combStr + "," + fmt.Sprint(item)
            }
        }
        utils.GenerateFileBasedOnVarAppend(tcid + "," + combStr + "\n", outputsFile)  
    }
}



func parseLine(fuzzLine string) FieldDefinition {
    var fieldDefinition FieldDefinition

    // id: char(10, 10)
    lineNS := strings.Replace(fuzzLine, " ", "", -1)
    line := strings.Split(lineNS, ":")
    // [id, char(10, 10)]

    fieldDefinition.FieldName = line[0]

    if strings.Index(line[1], "(") > 0 {
        fieldDefinition.FieldType = line[1][0:strings.Index(line[1], "(")]
    } else {
        fieldDefinition.FieldType = "float64"
    }

    fieldDefinition.FieldSubType = strings.Split(lineNS, ")")[1]

    line2 := strings.Split(line[1], ",")
    // [char(10, 10)]

    fieldDefinition.FieldMin, _ = strconv.Atoi(line2[0][strings.Index(line[1], "("):])
    fieldDefinition.FieldMax, _ = strconv.Atoi(line2[1][0:strings.Index(line2[1], ")")])

    return fieldDefinition
}



