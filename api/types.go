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

    gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"

)

type TcDataStore struct {
    TcData *testcase.TestCaseDataInfo

    TcLocalVariables map[string]interface{}
    SetUpStore []map[string]interface{}

    HttpExpStatus map[string]interface{}
    HttpExpHeader map[string]interface{}
    HttpExpBody map[string]interface{}
    HttpActualStatusCode int
    HttpActualHeader map[string][]string
    HttpActualBody []byte

    HttpStore map[string]interface{}
    TearDownStore []map[string]interface{}
}

func InitTcDataStore (tcData *testcase.TestCaseDataInfo) *TcDataStore {
    tcDataStore := &TcDataStore {
        tcData,

        map[string]interface{}{},
        []map[string]interface{}{},

        map[string]interface{}{},
        map[string]interface{}{},
        map[string]interface{}{},
        -1,
        map[string][]string{},
        []byte{},

        map[string]interface{}{},
        []map[string]interface{}{},
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

    if strings.Contains(jsonStr, "${") {
        for key, value := range dataFeeder {
            jsonStr = strings.Replace(jsonStr, "${" + key + "}", fmt.Sprint(value), -1)
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
    jsonStr = EvaluateBuiltinFunctions(resReq)
    json.Unmarshal([]byte(jsonStr), &resReq2)
 
    tcDataJson, _  = sjson.Set(tcDataJson, path, resReq2)

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

    jsonStr = EvaluateBuiltinFunctions(resResp)
    json.Unmarshal([]byte(jsonStr), &resResp2)
 
    tcDataJson, _  = sjson.Set(tcDataJson, path, resResp2)

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
 
    if strings.Contains(jsonStr, "${") {
        // Warning, this may have performance issues, need to improve, that is, get the Variables first, then replace
        for key, value := range dataFeeder {
            jsonStr = strings.Replace(jsonStr, "${" + key + "}", fmt.Sprint(value), -1)
        }
        
        json.Unmarshal([]byte(jsonStr), &res)
        tcDataJson, _  = sjson.Set(tcDataJson, path, res)

        json.Unmarshal([]byte(tcDataJson), &resTcData)
        tcDataStore.TcData = &resTcData
    }
} 

func (tcDataStore *TcDataStore) EvaluateTcBuiltinFunctions (path string) {
    var resTcData testcase.TestCaseDataInfo
    var res interface{}

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()
    jsonStr = EvaluateBuiltinFunctions(jsonStr)

    json.Unmarshal([]byte(jsonStr), &res)
    tcDataJson, _  = sjson.Set(tcDataJson, path, res)

    json.Unmarshal([]byte(tcDataJson), &resTcData)
    tcDataStore.TcData = &resTcData
}


