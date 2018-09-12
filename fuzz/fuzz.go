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
    "time"
    "fmt"
    "strings"
    "strconv"
    "encoding/json"
    "path/filepath"
    
    "go4api/cmd"
    "go4api/utils"
    "go4api/lib/pairwise"
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
        GenerateFuzzInvalidDataFiles(fuzzFile, fuzzData)
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
        fmt.Println("fuzzInvalidType: ", fuzzInvalidType)
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

    // this is to get the combinations, set use pairwise length = 2
    tcDataSlice := GetValidTcData(fuzzData, 2)
    //
    i := 0
    tcid := ""
    for _, tcData := range tcDataSlice {
        i = i + 1
        tcid = "valid" + strconv.Itoa(i)

        combStr := ""
        for ii, item := range tcData {
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

    tcDataSlice := GetInvalidTcData(fuzzData, 2)
    //
    i := 0
    tcid := ""
    for _, tcData := range tcDataSlice {
        i = i + 1
        tcid = "invalid" + strconv.Itoa(i)

        combStr := ""
        for ii, item := range tcData {
            if ii == 0 {
                combStr = combStr + fmt.Sprint(item)
            } else{
                combStr = combStr + "," + fmt.Sprint(item)
            }
        }
        utils.GenerateFileBasedOnVarAppend(tcid + "," + combStr + "\n", outputsFile)  
    }
}

func GetValidVectors(fuzzData FuzzData) [][]interface{} {
    var validVectors [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            validVectors = append(validVectors, validList)
        }
    }

    return validVectors
}

func GetInvalidVectors(fuzzData FuzzData) [][]interface{} {
    var invalidVectors [][]interface{}
    for _, invalidDataMap := range fuzzData.InvalidData {
        for _, invalidList := range invalidDataMap {
            invalidVectors = append(invalidVectors, invalidList)
        }
    }
    return invalidVectors
}


func GetValidTcData(fuzzData FuzzData, pwLength int) [][]interface{} {
    validVectors := GetValidVectors(fuzzData)

    validTcData := GetPairWiseValid(validVectors, pwLength)

    return validTcData
}


func GetInvalidTcData(fuzzData FuzzData, pwLength int) [][]interface{} {
    validVectors := GetValidVectors(fuzzData)
    invalidVectors := GetInvalidVectors(fuzzData)

    fmt.Println("--> validVectors: ", validVectors)
    fmt.Println("--> invalidVectors: ", invalidVectors)
    invalidTcData := GetCombinationInvalid(validVectors, invalidVectors, pwLength)

    return invalidTcData
}


//
func GetPairWiseValid(validVectors [][]interface{}, pwLength int) [][]interface{} {
    var validTcData [][]interface{}

    // need to consiber the len(combins) = 1 / = 2 / > 2
    if len(validVectors) >= pwLength {
        c := make(chan []interface{})

        go func(c chan []interface{}) {
            defer close(c)
            pairwise.GetPairWise(c, validVectors, 2)
        }(c)

        for tcData := range c {
            validTcData = append(validTcData, tcData)
        }
    } else if len(validVectors) == 1{
        for _, item := range validVectors[0] {
            var itemSlice []interface{}
            itemSlice = append(itemSlice, item)
            validTcData = append(validTcData, itemSlice)
        }
    }

    return validTcData
}

// -- for the fuzz data
func GetCombinationInvalid(validVectors [][]interface{}, invalidVectors [][]interface{}, pwLength int) [][]interface{} {
    // to ensure each negative value will be combined with each positive value(s)
    var invalidTcData [][]interface{}

    max := getMaxLenVector(validVectors)

    for i, _ := range invalidVectors {
        for j, _ := range invalidVectors[i] { 
            // loop the validVectors
            for jj := 0; jj < max; jj++ {
                tcData := make([]interface{}, len(validVectors))
                for k := 0; k < len(validVectors); k++ {
                    if i != k {
                        if jj <= len(validVectors[k]) - 1 {
                            tcData[k] = validVectors[k][jj]
                        } else {
                            // using the first one for valid vector
                            tcData[k] = validVectors[k][0]
                        }
                    }
                }
                tcData[i] = invalidVectors[i][j]  
                invalidTcData = append(invalidTcData, tcData)
            }
        }
    }
    
    return invalidTcData
}

func getMaxLenVector (vectors [][]interface{}) int {
    max := 0
    for i, _ := range vectors {
        if len(vectors[i]) > max {
            max = len(vectors[i])
        }
    }
    return max
}




