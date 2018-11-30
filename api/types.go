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
    "reflect"

    "go4api/lib/testcase" 
    "go4api/utils" 

    gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"

)

type TcDataStore struct {
    TcData *testcase.TestCaseDataInfo

    TcLocalVariables map[string]interface{}

    HttpExpStatus map[string]interface{}
    HttpExpHeader map[string]interface{}
    HttpExpBody map[string]interface{}
    HttpActualStatusCode int
    HttpActualHeader map[string][]string
    HttpActualBody []byte

    CmdGroupLength int

    CmdSection string // setUp, tearDown
    CmdType string // sql, redis, etc.
    CmdExecStatus string
    CmdAffectedCount int
    CmdResults interface{}
}

func InitTcDataStore (tcData *testcase.TestCaseDataInfo) *TcDataStore {
    tcDataStore := &TcDataStore {
        tcData,

        map[string]interface{}{},

        map[string]interface{}{},
        map[string]interface{}{},
        map[string]interface{}{},
        -1,
        map[string][]string{},
        []byte{},

        0,

        "",
        "",
        "",
        -1,
        -1,
    }
    // aa, _ := json.Marshal(tcData)
    // fmt.Println(string(aa))
    
    return tcDataStore
}

func (tcDataStore *TcDataStore) RenderTcRequestVariables (path string) {
    var resTcData testcase.TestCaseDataInfo
    var resReq testcase.Request
    dataFeeder := tcDataStore.MergeTestData()

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()
    // fmt.Println("jsonStr 0: ", jsonStr)

    if strings.Contains(jsonStr, "${") {
        for key, value := range dataFeeder {
            // Note: the type of value may be: string, int, float64, etc. 
            // fmt.Sprint(value) can result in issues, need to fix
            // 
            var valueStr = ""

            if value != nil {
                switch reflect.TypeOf(value).Kind().String() {
                case "float64":
                    // fmt.Println("t type float64:", value)
                    valueStr = utils.FloatToString(value.(float64))
                default:
                    valueStr = fmt.Sprint(value)
                }
            }

            jsonStr = strings.Replace(jsonStr, "${" + key + "}", valueStr, -1)
        }
   
        json.Unmarshal([]byte(jsonStr), &resReq)
        tcDataJson, _  = sjson.Set(tcDataJson, path, resReq)

        json.Unmarshal([]byte(tcDataJson), &resTcData)
        tcDataStore.TcData = &resTcData
    }
} 

func (tcDataStore *TcDataStore) EvaluateTcRequestBuiltinFunctions (path string) {
    var resTcData testcase.TestCaseDataInfo
    var resReq testcase.Request
    var resReq2 testcase.Request

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()
    json.Unmarshal([]byte(jsonStr), &resReq)

    // Note: for function EvaluateBuiltinFunctions:
    // if the input is str, like "request":{"method":"POST","path":"... 
    // the returned str is: "{\"method\":\"POST\",\"path\":\"...
    // to be safe, using the underlying struct
    edReq := EvaluateBuiltinFunctions(resReq)
    switch edReq.(type) {
        case string:
            jsonStr = edReq.(string)

            json.Unmarshal([]byte(jsonStr), &resReq2)
            tcDataJson, _  = sjson.Set(tcDataJson, path, resReq2)
        default:
            tcDataJson, _  = sjson.Set(tcDataJson, path, resReq)
    }

    json.Unmarshal([]byte(tcDataJson), &resTcData)
    tcDataStore.TcData = &resTcData
}

func (tcDataStore *TcDataStore) RenderTcResponseVariables (path string) {
    var resTcData testcase.TestCaseDataInfo
    var resResp testcase.Response
    dataFeeder := tcDataStore.MergeTestData()

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()
    // fmt.Println("jsonStr 1: ", jsonStr)

    if strings.Contains(jsonStr, "${") {
        for key, value := range dataFeeder {
            jsonStr = strings.Replace(jsonStr, "${" + key + "}", fmt.Sprint(value), -1)
        }
   
        json.Unmarshal([]byte(jsonStr), &resResp)
        tcDataJson, _  = sjson.Set(tcDataJson, path, resResp)

        json.Unmarshal([]byte(tcDataJson), &resTcData)
        tcDataStore.TcData = &resTcData
    }
} 

func (tcDataStore *TcDataStore) EvaluateTcResponseBuiltinFunctions (path string) {
    var resTcData testcase.TestCaseDataInfo
    var resResp testcase.Response
    var resResp2 testcase.Response

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()
    json.Unmarshal([]byte(jsonStr), &resResp)

    edResp := EvaluateBuiltinFunctions(resResp)
    switch edResp.(type) {
        case string:
            jsonStr = edResp.(string)

            json.Unmarshal([]byte(jsonStr), &resResp2)
            tcDataJson, _  = sjson.Set(tcDataJson, path, resResp2)
        default:
            tcDataJson, _  = sjson.Set(tcDataJson, path, resResp)
    }
    
    json.Unmarshal([]byte(tcDataJson), &resTcData)
    tcDataStore.TcData = &resTcData
}


func (tcDataStore *TcDataStore) RenderTcVariables (path string) {
    var resTcData testcase.TestCaseDataInfo
    var res interface{}
    dataFeeder := tcDataStore.MergeTestData()

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()
    // fmt.Println("jsonStr 2: ", jsonStr)
 
    if strings.Contains(jsonStr, "${") {
        // Warning, this may have performance issues, need to improve, that is, get the Variables first, then replace
        for key, value := range dataFeeder {
            jsonStr = strings.Replace(jsonStr, "${" + key + "}", fmt.Sprint(value), -1)
        }
        // fmt.Println("jsonStr: ", jsonStr)
        json.Unmarshal([]byte(jsonStr), &res) // notice
        tcDataJson, _  = sjson.Set(tcDataJson, path, jsonStr)
        // fmt.Println("tcDataJson: ", tcDataJson)

        json.Unmarshal([]byte(tcDataJson), &resTcData)
        tcDataStore.TcData = &resTcData
    }
} 

func (tcDataStore *TcDataStore) EvaluateTcBuiltinFunctions (path string) {
    var resTcData testcase.TestCaseDataInfo
    var resMap interface{}

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)
    // fmt.Println(">>> tcDataJson: 0: ", tcDataJson)

    result := gjson.Get(tcDataJson, path)
    edResp := EvaluateBuiltinFunctions(result.Value())
    // fmt.Println(">>> edResp: 0: ", edResp)

    // to be noticed the special case: result.Value() is string, edResp is string
    if strings.Contains(result.String(), "Fn::") {
        switch edResp.(type) {
        case string:
            // fmt.Println(">>> ----------------->")
            jsonStr := edResp.(string)

            json.Unmarshal([]byte(jsonStr), &resMap)
            tcDataJson, _  = sjson.Set(tcDataJson, path, resMap)
        default:
            tcDataJson, _  = sjson.Set(tcDataJson, path, result.Value())
        }
    }
    
    // fmt.Println(">>> tcDataJson: 1: ", tcDataJson)
    json.Unmarshal([]byte(tcDataJson), &resTcData)
    tcDataStore.TcData = &resTcData
}


