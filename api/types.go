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
    "encoding/json"

    "go4api/utils" 
    "go4api/lib/testcase" 

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

    HttpUrl string

    CmdSection string // setUp, tearDown
    CmdGroupLength int
    
    CmdType string // sql, redis, init, etc.
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

        "",

        "",
        0,

        "",
        "",
        -1,
        -1,
    }
    // aa, _ := json.Marshal(tcData)
    // fmt.Println(string(aa))
    
    return tcDataStore
}


// for http: .request, .response, .session, .outGlobalVariables, .outLocalVariables, .outFiles
// for cmd (setUp, tearDown): .cmd, .cmdResponse, .session, .outGlobalVariables, .outLocalVariables, .outFiles
func (tcDataStore *TcDataStore) PrepEmbeddedFunctions (path string) {
    pathSlice := strings.Split(path, ".")
    pathLength := len(pathSlice)
    pathType := pathSlice[pathLength - 1]

    switch pathType {
    case "request":
        var res testcase.Request

        tcDataStore.RenderTcVariables(path, res)
        tcDataStore.EvaluateTcEmbeddedFunctions(path, res)
    case "response":
        var res testcase.Response

        tcDataStore.RenderTcVariables(path, res)
        tcDataStore.EvaluateTcEmbeddedFunctions(path, res)

    case "session", "outGlobalVariables", "outLocalVariables":
        var res map[string]interface{}

        tcDataStore.RenderTcVariables(path, res)
        tcDataStore.EvaluateTcEmbeddedFunctions(path, res)
    case "outFiles":
        var res []*testcase.OutFilesDetails

        tcDataStore.RenderTcVariables(path, res)
        tcDataStore.EvaluateTcEmbeddedFunctions(path, res)
    case "cmd":
        var res string

        tcDataStore.RenderTcVariables(path, res)
        tcDataStore.EvaluateTcEmbeddedFunctions(path, res)
    case "cmdResponse":
        var res map[string]interface{}

        tcDataStore.RenderTcVariables(path, res)
        tcDataStore.EvaluateTcEmbeddedFunctions(path, res)
    }
}


func (tcDataStore *TcDataStore) RenderTcVariables (path string, res interface{}) {
    var resTcData testcase.TestCaseDataInfo

    dataFeeder := tcDataStore.MergeTestData()

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    jsonStr := gjson.Get(tcDataJson, path).String()
  
    if strings.Contains(jsonStr, "${") {
        // Warning, there may have performance issues
        for key, value := range dataFeeder {
            var valueStr = ""
            if value != nil {
                switch reflect.TypeOf(value).Kind().String() {
                case "float64":
                    valueStr = utils.FloatToString(value.(float64))
                case "string":
                    valueStr = value.(string)
                case "slice":
                    // for slice, []string or []float64, may have better solution later
                    // for example:
                    // valueB, _ := json.Marshal(value)
                    // valueStr = "`" + string(valueB) + "`"

                    valueStr = fmt.Sprint(value)
                default:
                    valueStr = fmt.Sprint(value)
                }
            }

            jsonStr = strings.Replace(jsonStr, "${" + key + "}", valueStr, -1)
        }
        // Note: if the jsonStr is string, like "request":{"method":"POST","path":"... 
        // the returned string tcDataJson is: "{\"method\":\"POST\",\"path\":\"...
        // for this issue, be kind to use the right struct but not string
        switch res.(type) {
        case string:
            tcDataJson, _  = sjson.Set(tcDataJson, path, jsonStr)
        default:
            json.Unmarshal([]byte(jsonStr), &res) 

            tcDataJson, _  = sjson.Set(tcDataJson, path, res)
        }
        
        json.Unmarshal([]byte(tcDataJson), &resTcData)
        tcDataStore.TcData = &resTcData
    }
} 

func (tcDataStore *TcDataStore) EvaluateTcEmbeddedFunctions (path string, res interface{}) {
    var resTcData testcase.TestCaseDataInfo

    tcDataJsonBytes, _ := json.Marshal(tcDataStore.TcData)
    tcDataJson := string(tcDataJsonBytes)

    result := gjson.Get(tcDataJson, path)
    edResp := tcDataStore.EvaluateEmbeddedFunctions(result.Value())

    // to be noticed the special case: result.Value() is string, edResp is string
    if strings.Contains(result.String(), "Fn::") {
        switch edResp.(type) {
        case string:
            jsonStr := edResp.(string)

            json.Unmarshal([]byte(jsonStr), &res)
            tcDataJson, _  = sjson.Set(tcDataJson, path, res)
        default:
            tcDataJson, _  = sjson.Set(tcDataJson, path, result.Value())
        }
    }

    json.Unmarshal([]byte(tcDataJson), &resTcData)
    tcDataStore.TcData = &resTcData
}


