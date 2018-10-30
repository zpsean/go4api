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
    // "fmt"
    "strings"
    // "reflect"
    "encoding/json"

    "go4api/lib/testcase"
    "go4api/builtins"
    "go4api/lib/g4json"

    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

func EvaluateBuiltinFunctions (tcData testcase.TestCaseDataInfo) testcase.TestCaseDataInfo {
    tcJsonBytes, _ := json.Marshal(tcData)
    tcJson := string(tcJsonBytes)

    var value interface{}

    tcBasicsJsonBytes, _ := json.Marshal(tcData.TestCase.TestCaseBasics())
    tcBasicsJson := string(tcBasicsJsonBytes)
    json.Unmarshal([]byte(tcBasicsJson), &value)

    fieldDetailsSlice := g4json.GetFieldsDetails(value)
    // tJson, _ := json.MarshalIndent(fieldDetailsSlice, "", "\t")
    // fmt.Println("=======>11: ", string(tJson))

    for i, _ := range fieldDetailsSlice {
        // e.g. "field": {"Fn::NextInt": ["min", "max"]}
        for j, _ := range fieldDetailsSlice[i].FieldPath {
            if strings.Contains(fieldDetailsSlice[i].FieldPath[j], "Fn::") {
                plPath := strings.Join(fieldDetailsSlice[i].FieldPath[0:j], ".")
                plFullPath := "TestCase." + tcData.TcName() + "." + plPath

                funcName := strings.TrimLeft(fieldDetailsSlice[i].FieldPath[j], "Fn::")
                funcParams := fieldDetailsSlice[i].CurrValue

                resValue := builtins.CallBuiltinFunc(funcName, funcParams)

                // switch reflect.TypeOf(resValue).Kind() {
                //     case reflect.Int, reflect.Int64:
                //         resValue = resValue
                //     case reflect.String:
                //         resValue = fmt.Sprint(resValue)
                // }

                tcJson, _  = sjson.Set(tcJson, plFullPath, resValue)
            }
        }    
    }

    var fTcData testcase.TestCaseDataInfo
    json.Unmarshal([]byte(tcJson), &fTcData)

    return fTcData
}




