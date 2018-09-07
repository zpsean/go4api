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
    "strings"
    "strconv"
    "encoding/json"
    "path/filepath"
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

// id1: Char(0, 20), Default("abc")
// id2: Char(0, 20)Numeric, Default("123")
// id3: Char(0, 20)Alpha, Default("abcd")
// id4: Char(0, 20)AlphaNumeric, Default("123abc")
// id5: Char(0, 20)Email
// id6: Char(0, 20)Time
// id7: Char(0, 20)IP, Default("127.0.0.1")
// id8: Int(0, 20), Default(123)
// id9: Int(0, 20)Time, Default(1533052800000)
// id10: Float(0, 20)Precision, Default(123.12)
// id11: Bool, Default(true)
// id12: Array([1, 2, 3, 4]), Default(1)

type FieldDefinitions []*FieldDefinition

type FieldDefinition struct {  
    FieldName string
    FieldType string
    FieldSubType string
    FieldMin int
    FieldMax int
    ArrayRaw []interface{}
    Default interface{}
}


func PrepFuzzTest(pStart_time time.Time) {
    fuzzFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".fuzz")
    // (1). generate the data tables based on the fuzz test, at least two dt files: positive and negative
    for _, fuzzFile := range fuzzFileList {
        fuzzData := GenerateFuzzData(fuzzFile)

        GenerateFuzzValidDataFiles(fuzzFile, fuzzData)
        // GenerateFuzzInvalidDataFiles(fuzzFile, fuzzData)
    }
    // (2). render the json using the fuzz dt(s)
    // fuzzTcArray := GetFuzzTcArray(options)
}

// to get the fuzz data table files with naming fuzzcase_fuzz_dt_valid.csv / fuzzcase_fuzz_dt_invalid.csv
func GenerateFuzzData(fuzzFile string) FuzzData {
    var fieldDefinitions FieldDefinitions
    defJson := utils.GetJsonFromFile(fuzzFile)
    json.Unmarshal([]byte(defJson), &fieldDefinitions)

    var fuzzData FuzzData
    var validValueList []map[string][]interface{}
    var invalidValueList []map[string][]interface{}

    for _, fieldDefinition := range fieldDefinitions {
        validValueMap := make(map[string][]interface{})
        invalidValueMap := make(map[string][]interface{})
        // call the rules to get values
        fuzzValidType := fieldDefinition.DetermineFuzzValidType()
        validValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules(fuzzValidType)
        // invalid 
        fuzzInvalidType := fieldDefinition.DetermineFuzzInvalidType()
        invalidValueMap[fieldDefinition.FieldName] = fieldDefinition.CallFuzzRules(fuzzInvalidType)
        // append to slice
        validValueList = append(validValueList, validValueMap)
        invalidValueList = append(invalidValueList, invalidValueMap)
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


func GenerateFuzzValidDataFiles(fuzzFile string, fuzzData FuzzData) {
    outputsFile := filepath.Join(filepath.Dir(fuzzFile), 
        strings.TrimRight(filepath.Base(fuzzFile), ".fuzz") + "_fuzz_dt_valid.csv")
    // write csv header, data
    validHeaderStr := ""
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
}

func GenerateFuzzInvalidDataFiles(fuzzFile string, fuzzData FuzzData) {
    outputsFile := filepath.Join(filepath.Dir(fuzzFile), 
        strings.TrimRight(filepath.Base(fuzzFile), ".fuzz") + "_fuzz_dt_invalid.csv")
    // write csv header
    invalidHeaderStr := ""
    invalidHeaderStr = invalidHeaderStr + "tcid"
    for _, invalidDataMap := range fuzzData.InvalidData {
        for key, _ := range invalidDataMap {
            invalidHeaderStr = invalidHeaderStr + "," + key
        }
    }
    utils.GenerateFileBasedOnVarOverride(invalidHeaderStr + "\n", outputsFile)

    combInvalid := GetCombinationInvalid(fuzzData)
    //
    i := 0
    tcid := ""
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



