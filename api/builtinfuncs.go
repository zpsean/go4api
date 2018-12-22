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

    "go4api/builtins"
    "go4api/lib/g4json"

    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

type BuiltinFieldDetails struct {
    FieldPath []string
    CurrValue interface{}
    FieldType string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
    EvaluatedValue interface{}
}


func EvaluateBuiltinFunctions (value interface{}) interface{} {
    jsonBytes, _ := json.Marshal(value)
    jsonStr := string(jsonBytes)

    // check if has builtin function
    if !strings.Contains(jsonStr, "Fn::") {
        return value
    } else {
        // fmt.Println(">>...")
        builtinLeavesSlice := GetBuiltinLeavesSlice(value)

        maxLevel := g4json.GetJsonNodesLevel(builtinLeavesSlice)

        jsonStr = IterateBuiltsins(jsonStr, builtinLeavesSlice, maxLevel)

        return jsonStr
    }
}

func GetBuiltinLeavesSlice (value interface{}) []g4json.FieldDetails {
    var builtinLeavesSlice []g4json.FieldDetails

    fieldDetailsSlice := g4json.GetFieldsDetails(value)
    leavesSlice := g4json.GetJsonLeaves(fieldDetailsSlice)

    for i, _ := range leavesSlice {
        nodePathStr := strings.Join(leavesSlice[i].FieldPath, ".")

        if strings.Contains(nodePathStr, "Fn::") {
            builtinLeavesSlice = append(builtinLeavesSlice, leavesSlice[i])
        }   
    }

    return builtinLeavesSlice
}

// need to consider the nested builtin functions, like:
// definition:  "field2": {"Fn::F2": [12, {"Fn::F3": ["aaa", "bbbb"]}]},
// leafpath:    "request.payload.text.field2.Fn::F2.0",
//              "request.payload.text.field2.Fn::F2.1.Fn::F3.1",
//              "request.payload.text.field2.Fn::F2.1.Fn::F3.0"
//
// !! Warning: specail case, if the key has dot itself, need specail handle:
// for example: 
// json origin:  {"body":{"$(body).msg_code":{"Equals":{"Fn::ToInt":"5109"}}}
// need mark the path as:  body.$(body)\.msg_code.Equals
// for convenience, here will use specila string (NdotReplacerN) to replace the dot: body.$(body)NdotReplacerNmsg_code.Equals
// otherwise, it will return unexpected results after sjson.Set()
func IterateBuiltsins (jsonStr string, builtinLeavesSlice []g4json.FieldDetails, maxLevel int) string {
    var evaluatedSlice []g4json.FieldDetails
    var evaluatedFuncPaths []string

    for i := maxLevel; i > 0; i-- {
        for j, _ := range builtinLeavesSlice {
            pathLength := len(builtinLeavesSlice[j].FieldPath)
            if pathLength >= i && i > 1 {
                // the last node (leaf), take its own CurrValue as the funcParams 
                if strings.Contains(builtinLeavesSlice[j].FieldPath[i - 1], "Fn::") {
                    var value interface{}
                    json.Unmarshal([]byte(jsonStr), &value)

                    evaluatedSlice = g4json.GetFieldsDetails(value)

                    var tempSlice []string
                    var nodePathStr string
                    for ii, _ := range builtinLeavesSlice[j].FieldPath[0:i - 1] {
                        oKey := builtinLeavesSlice[j].FieldPath[0:i - 1][ii]
                        
                        // reset the jsonStr with replacing NdotReplacerN
                        if strings.Contains(oKey, ".") {
                            cKey := strings.Replace(oKey, ".", "NdotReplacerN", -1)
                            jsonStr = strings.Replace(jsonStr, oKey, cKey, -1)

                            tempSlice = append(tempSlice, cKey)
                        } else {
                            tempSlice = append(tempSlice, oKey)
                        }
                    }
                    nodePathStr = strings.Join(tempSlice, ".")
                    
                    funcName := strings.TrimLeft(builtinLeavesSlice[j].FieldPath[i - 1], "Fn::")

                    var funcParams interface{}
                    ifExists := false
                    funcParamsPath := strings.Join(builtinLeavesSlice[j].FieldPath[0:i], ".")
                    for ind, _ := range evaluatedFuncPaths {
                        if funcParamsPath == evaluatedFuncPaths[ind] {
                            ifExists = true
                        }
                    }
                    if ifExists == true {
                        continue
                    }
                    for k, _ := range evaluatedSlice {
                        p := strings.Join(evaluatedSlice[k].FieldPath, ".")
                        if funcParamsPath == p {
                            funcParams = evaluatedSlice[k].CurrValue
                        }
                    }

                    resValue := builtins.CallBuiltinFunc(funcName, funcParams)
                    jsonStr, _  = sjson.Set(jsonStr, nodePathStr, resValue)
      
                    evaluatedFuncPaths = append(evaluatedFuncPaths, funcParamsPath)
                }
            }
        }
    }

    jsonStr = strings.Replace(jsonStr, "NdotReplacerN", ".", -1)

    return jsonStr
}

