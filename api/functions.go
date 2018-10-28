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
    "encoding/json"

    "go4api/lib/testcase"
    "go4api/builtins"
    "go4api/lib/g4json"

    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

func EvaluateBuiltinFunctions (tcData testcase.TestCaseDataInfo) testcase.TestCaseDataInfo {
    tcJson, _ := json.Marshal(tcData)

    for key, value := range tcData.TestCase.ReqPayload() {
        if key == "text" {
            // to loop over the struct
            fieldDetailsSlice := g4json.GetFieldsDetails(value)

            for i, _ := range fieldDetailsSlice {
                plPath := key + "." + strings.Join(fieldDetailsSlice[i].FieldPath, ".")
                plFullPath := "TestCase." + tcData.TcName() + ".request.payload" + "." + plPath
                
                // check if field value has Fn::
                // e.g. "field": {"Fn::NextInt": ["min", "max"]}
                if strings.Contains(fieldDetailsSlice[i].CurrValue.(string), "Fn::") {
                    funcName := strings.TrimLeft(fieldDetailsSlice[i].CurrValue.(string), "Fn::")
                    // funcParams := gjson.Get(string(actualBody), key[lenPrefix:])

                    resValue := builtins.CallBuiltinFunc(funcName, "")

                    fmt.Println("resValue: ", resValue)

                    sjson.Set(string(tcJson), plFullPath, resValue)
                }
            }
        }
    }

    var fTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &fTcData)

    return fTcData
}




