/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
    "os"
    "fmt"
    "strings"
    "reflect"
    "path/filepath"
    "encoding/json"
    "net/http"

    "go4api/lib/testcase"
    "go4api/utils"
    
    gjson "github.com/tidwall/gjson"
)


func WriteOutputsDataToFile (testResult string, tcData testcase.TestCaseDataInfo, actualStatusCode int, actualHeader http.Header, actualBody []byte) {
    var expOutputs []*testcase.OutputsDetails

    if testResult == "Success" {
        expOutputs = tcData.TestCase.Outputs()
        if len(expOutputs) > 0 {
            // get the actual value from actual body based on the fields in json outputs
            for i, _ := range expOutputs {
                var keyStrList []string
                var valueStrList []string
                // (1). get the full path of outputsfile
                tempDir := utils.CreateTempDir(tcData.JsonFilePath)
                outputsFile := filepath.Join(tempDir, expOutputs[i].GetOutputsDetailsFileName())
                os.Remove(outputsFile)
                // (2). get the outputsfile format, may be csv, excel, etc.
                outputsFileFormat := expOutputs[i].GetOutputsDetailsFormat()
                // (3). get the outputsfile data
                outputsData := expOutputs[i].GetOutputsDetailsData()
                switch strings.ToLower(outputsFileFormat) {
                    case "csv":
                        keyStrList, valueStrList = GetOutputsCsvData(outputsData, actualStatusCode, actualHeader, actualBody)
                        // write csv header
                        utils.GenerateCsvFileBasedOnVarOverride(keyStrList, outputsFile)
                        // write csv data
                        utils.GenerateCsvFileBasedOnVarAppend(valueStrList, outputsFile)
                    case "xlsx":
                        SaveHttpRespFile(actualBody, outputsFile)
                }   
            } 
        }
    } else {
        // fmt.Println("Warning: test execution failed, no outputs file!")
    }
}

func GetOutputsCsvData (outputsData map[string][]interface{}, actualStatusCode int, actualHeader http.Header, actualBody []byte) ([]string, []string) {
    var keyStrList []string
    var valueStrList []string

    for key, valueSlice := range outputsData {
        // for csv header
        keyStrList = append(keyStrList, key)
        // for cav data
        if len(valueSlice) > 0 {
            // check if the valueSlice is [], or [[]], using the valueSlice[0]
            switch reflect.TypeOf(valueSlice[0]).Kind() {
                case reflect.Slice:
                    fieldStrList := GetOutputsDetailsDataForFieldSlice(valueSlice, actualStatusCode, actualHeader, actualBody)
                    fieldStr := convertSliceAsString(fieldStrList)
                    valueStrList = append(valueStrList, fieldStr)
                default: 
                    // Note, here may return array also
                    fieldStrList := GetOutputsDetailsDataForFieldString(valueSlice, actualStatusCode, actualHeader, actualBody)
                    valueStrList = append(valueStrList, strings.Join(fieldStrList, "")) 
            }
        }     
    }

    return keyStrList, valueStrList
}

func GetOutputsDetailsDataForFieldString (valueSlice []interface{}, actualStatusCode int, actualHeader http.Header, actualBody []byte) []string {
    var fieldStrList []string
    // check if the valueSlice is [], or [[]], using the valueSlice[0]
    for _, value := range valueSlice {
        // actualValue := GetActualValueByJsonPath(fmt.Sprint(value), actualBody)
        actualValue := GetResponseValue(fmt.Sprint(value), actualStatusCode, actualHeader, actualBody)
        
        if actualValue == nil {
            fieldStrList = append(fieldStrList, "")
        } else {
            switch reflect.TypeOf(actualValue).Kind() {
                case reflect.Slice:
                    str := convertSliceAsString(reflect.ValueOf(actualValue).Interface().([]interface{}))
                    fieldStrList = append(fieldStrList, fmt.Sprint(str))
                default:
                    fieldStrList = append(fieldStrList, fmt.Sprint(actualValue))
            }
        }
    }

    return fieldStrList
}


func GetOutputsDetailsDataForFieldSlice (valueSlice []interface{}, actualStatusCode int, actualHeader http.Header, actualBody []byte) []interface{} {
    var fieldStrList []interface{}
    // currently, suppose has only one sub slice
    firstSubSlice := valueSlice[0]

    for _, value := range reflect.ValueOf(firstSubSlice).Interface().([]interface{}) {
        // actualValue := GetActualValueByJsonPath(fmt.Sprint(value), actualBody)
        actualValue := GetResponseValue(fmt.Sprint(value), actualStatusCode, actualHeader, actualBody)

        if actualValue == nil {
            fieldStrList = append(fieldStrList, "")
        } else {
            fieldStrList = append(fieldStrList, actualValue)
        }
    }

    return fieldStrList
}

func convertSliceAsString (slice []interface{}) string {
    varStr := ""
    if len(slice) > 0 {
        varStrByte, _ := json.Marshal(slice)
        varStr = string(varStrByte)
    } else {
        varStr = "[]"
    }
    // fmt.Println("==>", slice, varStr)
    return varStr
}


// -- 
func GetResponseValue (searchPath string, actualStatusCode int, actualHeader http.Header, actualBody []byte) interface{} {
    // prefix = "$(status).", "$(headers).", "$(body)."
    var value interface{}
    if len(searchPath) > 1 {
        if strings.HasPrefix(searchPath, "$(status).") {
            value = GetStatusActualValue(searchPath, actualStatusCode)
        } else if strings.HasPrefix(searchPath, "$(headers).") {
            value = GetHeadersActualValue(searchPath, actualHeader)
        } else if strings.HasPrefix(searchPath, "$(body).") {
            value = GetActualValueByJsonPath(searchPath, actualBody)
        } else if strings.HasPrefix(searchPath, "$.") {
            value = GetActualValueByJsonPath(searchPath, actualBody)
        } else {
            value = searchPath
        }
    } else {
        value = searchPath
    }
    
    return value
}

func GetStatusActualValue (key string, actualStatusCode int) interface{} {
    var actualValue interface{}
    // leading "$(status)" is mandatory if want to retrive status
    if len(key) == 9 && key == "$(status)" {
        actualValue = actualStatusCode
    } else {
        actualValue = key
    }

    return actualValue
}

func GetHeadersActualValue (key string, actualHeader http.Header) interface{} { 
    var actualValue interface{}
    // leading "$(headers)" is mandatory if want to retrive headers value
    prefix := "$(headers)."
    lenPrefix := len(prefix)

    if len(key) > lenPrefix && key[0:lenPrefix] == prefix {
        actualValue = strings.Join(actualHeader[key[lenPrefix:]], ",")
    } else {
        actualValue = key
    }

    return actualValue
}

func GetActualValueByJsonPath (key string, actualBody []byte) interface{} {  
    var actualValue interface{}
    // leading "$." or "$(headers)." is mandatory if want to use path search
    prefix := "$(body)."
    lenPrefix := len(prefix)
    prefix2 := "$."
    lenPrefix2 := len(prefix2)

    if len(key) > lenPrefix && key[0:lenPrefix] == prefix {
        value := gjson.Get(string(actualBody), key[lenPrefix:])
        actualValue = value.Value()
    } else if len(key) > lenPrefix2 && key[0:lenPrefix2] == prefix2 {
        value := gjson.Get(string(actualBody), key[lenPrefix2:])
        actualValue = value.Value()
    } else {
        actualValue = key
    }

    return actualValue
}

