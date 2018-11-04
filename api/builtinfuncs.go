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

    "go4api/lib/testcase"
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

func EvaluateBuiltinFunctions (tcData testcase.TestCaseDataInfo) testcase.TestCaseDataInfo {
    tcJsonBytes, _ := json.Marshal(tcData)
    tcJson := string(tcJsonBytes)

    // fmt.Println(tcJson)

    // check if has builtin function
    if !strings.Contains(tcJson, "Fn::") {
        return tcData
    } else {
        builtinLeavesSlice := GetBuiltinLeavesSlice(tcData.TestCase.TestCaseBasics())
        maxLevel := g4json.GetJsonNodesLevel(builtinLeavesSlice)
        fmt.Println("maxLevel: ", maxLevel)

        tcJson = IterateBuiltsins(tcJson, builtinLeavesSlice, maxLevel)

        var fTcData testcase.TestCaseDataInfo
        json.Unmarshal([]byte(tcJson), &fTcData)

        return fTcData
    }
}

func GetBuiltinLeavesSlice (tcBasics *testcase.TestCaseBasics) []g4json.FieldDetails {
    var value interface{}
    var builtinLeavesSlice []g4json.FieldDetails

    tcBasicsJsonBytes, _ := json.Marshal(tcBasics)
    tcBasicsJson := string(tcBasicsJsonBytes)
    json.Unmarshal([]byte(tcBasicsJson), &value)

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

// scan from end to front
// create new shadow slice to record Fn value
// merge the values
// then sjson to update the field
func IterateBuiltsins (tcJson string, builtinLeavesSlice []g4json.FieldDetails, maxLevel int) string {
    // ---
    var evaluatedSlice []g4json.FieldDetails
    var evaluatedFuncPaths []string

    for i := maxLevel; i > 0; i-- {
        for j, _ := range builtinLeavesSlice {
            pathLength := len(builtinLeavesSlice[j].FieldPath)
            if pathLength >= i && i > 1 {
                // the last node (leaf), take its own CurrValue as the funcParams 
                if strings.Contains(builtinLeavesSlice[j].FieldPath[i - 1], "Fn::") {
                    //
                    var evalData testcase.TestCaseDataInfo
                    var value interface{}
                    json.Unmarshal([]byte(tcJson), &evalData)

                    tcBasicsJsonBytes, _ := json.Marshal(evalData.TestCase.TestCaseBasics())
                    tcBasicsJson := string(tcBasicsJsonBytes)
                    json.Unmarshal([]byte(tcBasicsJson), &value)

                    evaluatedSlice = g4json.GetFieldsDetails(value)
                    // zz, _ := json.MarshalIndent(evaluatedSlice, "", "\t")
                    // fmt.Println("evaluatedSlice: ", string(zz))

                    
                    // evaluatedSlice = GetBuiltinLeavesSlice(evalData.TestCase.TestCaseBasics())
                    // zz, _ := json.MarshalIndent(evaluatedSlice, "", "\t")
                    // fmt.Println("evaluatedSlice: ", string(zz))

                    //
                    nodePathStr := strings.Join(builtinLeavesSlice[j].FieldPath[0:i - 1], ".")
                    plFullPath := "TestCase." + evalData.TestCase.TcName() + "." + nodePathStr

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
                    
                    fmt.Println("---> vv i, j: ", fmt.Sprint(i), fmt.Sprint(j), ": ", funcName, funcParams)
                    fmt.Println("---> vv plFullPath: ", plFullPath)

                    resValue := builtins.CallBuiltinFunc(funcName, funcParams)
                    
                    fmt.Println("---> vv i, j: ", fmt.Sprint(i), fmt.Sprint(j), ": ", resValue)
                    

                    tcJson, _  = sjson.Set(tcJson, plFullPath, resValue)

                    evaluatedFuncPaths = append(evaluatedFuncPaths, funcParamsPath)
                }
            }
        }
    }

    return tcJson
}

