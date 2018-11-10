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
    TcData testcase.TestCaseDataInfo

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

func InitTcDataStore (tcData testcase.TestCaseDataInfo) *TcDataStore {
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


func (tcDataStore *TcDataStore) RenderTcVariables (path string) {
    var resTcData testcase.TestCaseDataInfo
    dataFeeder := tcDataStore.MergeTestData()

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()

    if strings.Contains(jsonStr, "${") {
        // Warning, this may have performance issues, need to improve, that is, get the Variables first, then replace
        for key, value := range dataFeeder {
            jsonStr = strings.Replace(jsonStr, "${" + key + "}", fmt.Sprint(value), -1)
        }
        tcDataJson, _  = sjson.Set(tcDataJson, path, jsonStr)

        json.Unmarshal([]byte(tcDataJson), &resTcData)
        tcDataStore.TcData = resTcData
    }
} 


func (tcDataStore *TcDataStore) EvaluateTcBuiltinFunctions (path string) {
    var resTcData testcase.TestCaseDataInfo

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()
    jsonStr = EvaluateBuiltinFunctions(jsonStr)
    // path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + ".setUp"

    tcDataJson, _  = sjson.Set(tcDataJson, path, jsonStr)

    json.Unmarshal([]byte(tcDataJson), &resTcData)
    tcDataStore.TcData = resTcData
}

func (tcDataStore *TcDataStore) RenderHttpSectionVariables () {

} 

func (tcDataStore *TcDataStore) EvaluateHttpSectionBuiltinFunctions () {

} 
