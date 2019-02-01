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


func (tcDataStore *TcDataStore) EvaluateBuiltinFunctions (value interface{}) interface{} {
    jsonBytes, _ := json.Marshal(value)
    jsonStr := string(jsonBytes)

    // check if has builtin function
    if !strings.Contains(jsonStr, "Fn::") {
        return value
    } else {
        // fmt.Println(">>...")
        builtinLeavesSlice := GetBuiltinLeavesSlice(value)

        maxLevel := g4json.GetJsonNodesLevel(builtinLeavesSlice)

        jsonStr = tcDataStore.IterateBuiltsins(jsonStr, builtinLeavesSlice, maxLevel)

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
// !! Warning: specail case, if the key is complex key, as contains ., \, ", etc., need specail handle

func (tcDataStore *TcDataStore) IterateBuiltsins (jsonStr string, builtinLeavesSlice []g4json.FieldDetails, maxLevel int) string {
    var evaluatedSlice []g4json.FieldDetails
    var evaluatedFuncPaths []string

    var replacerMap = make(map[string]string)

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

                        // if the key is complex key, as contains dot (.)
                        if strings.Contains(oKey, ".") {
                            rkey := "go4Api_efdvberipz_ReplacerKey_" + fmt.Sprint(i) + "_" + fmt.Sprint(j) + "_" + fmt.Sprint(ii)

                            // if the key is complex key, as contains \"
                            if strings.Contains(oKey, "\"") {
                                oKey = strings.Replace(oKey, "\"", "\\\"", -1)
                                replacerMap[rkey] = oKey
                            } else {
                                replacerMap[rkey] = oKey
                            }

                            tempSlice = append(tempSlice, rkey)
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

                    var funcParams_f interface{}
                    // Note: if funcParams is string, it has chance to be the json lookup path, like $(sql).xxx, $(body).xxx
                    switch funcParams.(type) {
                    case string:
                        funcParams_f = tcDataStore.GetResponseValue(funcParams.(string))
                    default:
                        funcParams_f = funcParams
                    }
                    
                    // call
                    resValue := builtins.CallBuiltinFunc(funcName, funcParams_f)

                    for key, _ := range replacerMap {
                        jsonStr = strings.Replace(jsonStr, replacerMap[key], key, -1)
                    }

                    jsonStr, _  = sjson.Set(jsonStr, nodePathStr, resValue)
      
                    evaluatedFuncPaths = append(evaluatedFuncPaths, funcParamsPath)
                }
            }
        }
    }

    for key, _ := range replacerMap {
        jsonStr = strings.Replace(jsonStr, key, replacerMap[key], -1)
    }

    return jsonStr
}

