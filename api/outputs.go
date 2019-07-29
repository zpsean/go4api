/*
 * go4api - an api testing tool written in Go
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

    "go4api/lib/testcase"
    "go4api/utils"
)


func (tcDataStore *TcDataStore) WriteOutputsDataToFile () {
    var expOutputs []*testcase.OutputsDetails

    tcData := tcDataStore.TcData

    // actualStatusCode := tcDataStore.HttpActualStatusCode
    // actualHeader := tcDataStore.HttpActualHeader
    actualBody := tcDataStore.HttpActualBody

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
                    keyStrList, valueStrList = tcDataStore.GetOutputsCsvData(outputsData)
                    // write csv header
                    utils.GenerateCsvFileBasedOnVarOverride(keyStrList, outputsFile)
                    // write csv data
                    utils.GenerateCsvFileBasedOnVarAppend(valueStrList, outputsFile)
                case "xlsx":
                    SaveHttpRespFile(actualBody, outputsFile)
            }   
        } 
    }

}

func (tcDataStore *TcDataStore) GetOutputsCsvData (outputsData map[string][]interface{}) ([]string, []string) {
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
                    fieldStrList := tcDataStore.GetOutputsDetailsDataForFieldSlice(valueSlice)
                    fieldStr := convertSliceAsString(fieldStrList)
                    valueStrList = append(valueStrList, fieldStr)
                default: 
                    // Note, here may return array also
                    fieldStrList := tcDataStore.GetOutputsDetailsDataForFieldString(valueSlice)
                    valueStrList = append(valueStrList, strings.Join(fieldStrList, "")) 
            }
        }     
    }

    return keyStrList, valueStrList
}

func (tcDataStore *TcDataStore) GetOutputsDetailsDataForFieldString (valueSlice []interface{}) []string {
    var fieldStrList []string
    // check if the valueSlice is [], or [[]], using the valueSlice[0]
    for _, value := range valueSlice {
        // actualValue := GetActualValueByJsonPath(fmt.Sprint(value), actualBody)
        actualValue := tcDataStore.GetResponseValue(fmt.Sprint(value))
        
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


func (tcDataStore *TcDataStore) GetOutputsDetailsDataForFieldSlice (valueSlice []interface{}) []interface{} {
    var fieldStrList []interface{}
    // currently, suppose has only one sub slice
    firstSubSlice := valueSlice[0]

    for _, value := range reflect.ValueOf(firstSubSlice).Interface().([]interface{}) {
        // actualValue := GetActualValueByJsonPath(fmt.Sprint(value), actualBody)
        actualValue := tcDataStore.GetResponseValue(fmt.Sprint(value))

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

