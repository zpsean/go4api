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
    "fmt"
    "strings"
    "reflect"
    "path/filepath"
    "encoding/json"

    "go4api/testcase"
    "go4api/utils"
    
    gjson "github.com/tidwall/gjson"
)


func WriteOutputsDataToFile(testResult string, tcData testcase.TestCaseDataInfo, actualBody []byte) {
    var expOutputs []*testcase.OutputsDetails

    if testResult == "Success" {
        expOutputs = tcData.TestCase.Outputs()
        if len(expOutputs) > 0 {
            // get the actual value from actual body based on the fields in json outputs
            for i, _ := range expOutputs {
                var keyStrList []string
                var valueStrList []string
                // item is {}
                // (1). get the full path of outputsfile
                outputsFile := filepath.Join(filepath.Dir(tcData.JsonFilePath), expOutputs[i].GetOutputsDetailsFileName())
                // (2). get the outputsfile format
                // outputsFileFormat := expOutputs[i].GetOutputsDetailsFormat()
                // (3). get the outputsfile data
                for key, valueSlice := range expOutputs[i].GetOutputsDetailsData() {
                    // for csv header
                    keyStrList = append(keyStrList, key)
                    // for cav data
                    if len(valueSlice) > 0 {
                        // check if the valueSlice is [], or [[]], using the valueSlice[0]
                        switch reflect.TypeOf(valueSlice[0]).Kind() {
                            case reflect.String, reflect.Float64: 
                                fieldStrList := GetOutputsDetailsDataForFieldString(valueSlice, actualBody)
                                valueStrList = append(valueStrList, strings.Join(fieldStrList, "")) 
                            case reflect.Slice:
                                fieldStrList := GetOutputsDetailsDataForFieldSlice(valueSlice, actualBody)
                                fieldStr := convertSliceAsString(fieldStrList)
                                valueStrList = append(valueStrList, fieldStr)
                        }
                    }     
                }
                // write csv header
                // utils.GenerateFileBasedOnVarOverride(strings.Join(keyStrList, ",") + "\n", outputsFile)
                utils.GenerateCsvFileBasedOnVarOverride(keyStrList, outputsFile)

                // write csv data
                // utils.GenerateFileBasedOnVarAppend(strings.Join(valueStrList, ",") + "\n", outputsFile)
                utils.GenerateCsvFileBasedOnVarAppend(valueStrList, outputsFile)
            } 
        }
    }
}

func GetOutputsDetailsDataForFieldString (valueSlice []interface{}, actualBody []byte) []string {
    var fieldStrList []string
    // check if the valueSlice is [], or [[]], using the valueSlice[0]
    for _, value := range valueSlice {
        if fmt.Sprint(value)[0:2] == "$." {
            actualValue := GetActualValueBasedOnExpKeyAndActualBody(fmt.Sprint(value), actualBody)
            if actualValue == nil {
                fieldStrList = append(fieldStrList, "")
                // valueStrList = append(valueStrList, "")
            } else {
                fieldStrList = append(fieldStrList, fmt.Sprint(actualValue))
                // valueStrList = append(valueStrList, fmt.Sprint(actualValue))
            }
        } else {
            fieldStrList = append(fieldStrList, fmt.Sprint(value))
            // valueStrList = append(valueStrList, fmt.Sprint(value))
        }
    }

    return fieldStrList
}


func GetOutputsDetailsDataForFieldSlice (valueSlice []interface{}, actualBody []byte) []interface{} {
    var fieldStrList []interface{}
    // currently, deal with the first sub slice only
    firstSubSlice := valueSlice[0]
    for _, value := range reflect.ValueOf(firstSubSlice).Interface().([]interface{}) {
        if fmt.Sprint(value)[0:2] == "$." {
            actualValue := GetActualValueBasedOnExpKeyAndActualBody(fmt.Sprint(value), actualBody)
            if actualValue == nil {
                fieldStrList = append(fieldStrList, "")
                // valueStrList = append(valueStrList, "")
            } else {
                fieldStrList = append(fieldStrList, actualValue)
                // valueStrList = append(valueStrList, fmt.Sprint(actualValue))
            }
        } else {
            fieldStrList = append(fieldStrList, value)
            // valueStrList = append(valueStrList, fmt.Sprint(value))
        }
    }

    return fieldStrList
}

func convertSliceAsString (slice []interface{}) string {
    varStr := ""
    if len(slice) > 0 {
        switch reflect.TypeOf(slice[0]).Kind() {
            case reflect.String, reflect.Float64:
                varStrByte, _ := json.Marshal(slice)
                fmt.Println(string(varStrByte))
                varStr = string(varStrByte)
            case reflect.Float32:
                varFloat := 0.0
                for _, v := range slice {
                    varFloat = varFloat + v.(float64)
                }
                varStr = fmt.Sprint(varFloat)
        }
    } else {
        varStr = ""
    }

    return varStr
}


func GetActualValueBasedOnExpKeyAndActualBody(key string, actualBody []byte) interface{} {
    var actualValue interface{}
    // if key starts with "$.", it represents the path, for xml, json
    // if key == "text", it is plain text, represents its valu is the whole returned body
    //
    // parse it based on the json by default, need add logic for xml, and other format
    if key[0:2] == "$." {
        value := gjson.Get(string(actualBody), key[2:])
        actualValue = value.Value()
    } else {
        value := gjson.Get(string(actualBody), key)
        actualValue = value.Value()
    }

    // fmt.Println("actualValue: ", actualValue)
    return actualValue
}
