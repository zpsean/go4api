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
    "fmt"
    "strings"
    "reflect"
    "encoding/json"
    "strconv"
    
    gjson "github.com/tidwall/gjson"
)

func (tcDataStore *TcDataStore) GetResponseValue (searchPath string) interface{} {
    // searchPath has prefix = "$(status).", "$(headers).", "$(body).", "$(xx).", etc.
    var value interface{}

    z := strings.SplitN(searchPath, ".", 2)

    if len(z) == 2 {
        if len(z[0]) > 0 && len(z[1]) > 0 {
            switch z[0] {
            case "$(status)":
                value = tcDataStore.GetStatusActualValue()
            case "$(headers)":
                value = tcDataStore.GetHeadersActualValue(z[1])
            case "$(body)":
                s := string(tcDataStore.HttpActualBody)

                value = tcDataStore.GetContentByPath(s, z[1])
            case "$(sql)", "$(mysql)", "$(postgresql)", "$(mongodb)":
                s := tcDataStore.CmdResults

                value = tcDataStore.GetContentByPath(s, z[1])
            case "$(file)":
                s := tcDataStore.CmdResults

                value = tcDataStore.GetContentByPath(s, z[1])
            case "$(redis)":
                s := tcDataStore.CmdResults
                
                value = tcDataStore.GetContentByPath(s, z[1])
                // value = tcDataStore.GetRedisActualValueByPath(searchPath)
            default:
                // If the value from from declared variable: $(${variable}).xx.yy
                if strings.Contains(z[0], "$(") {
                    s := tcDataStore.GetVariableContent(z[0])
                    value = tcDataStore.GetContentByPath(s, z[1])
                } else {
                    value = searchPath
                }
            } 
        } else {
            value = searchPath
        }
    } else {
        value = searchPath
    }
 
    return value
}

// http response status code
func (tcDataStore *TcDataStore) GetStatusActualValue () interface{} {
    actualStatusCode := tcDataStore.HttpActualStatusCode

    return actualStatusCode
}

// http response headers
func (tcDataStore *TcDataStore) GetHeadersActualValue (key string) interface{} { 
    var actualValue interface{}
    actualHeader := tcDataStore.HttpActualHeader
 
    actualValue = strings.Join(actualHeader[key], ",")

    return actualValue
}

// ------------------------
func (tcDataStore *TcDataStore) GetContentByPath (res interface{}, jsonPath string) interface{} {  
    var r interface{}

    if res == nil {
        return "_null_key_"
    }

    if jsonPath == "*" {
        return res
    }

    t := reflect.TypeOf(res).Kind().String()
    switch t {
    case "string":
        if strings.HasSuffix(jsonPath, "_keys_count_") {
            r = tcDataStore.GetKeysCount(res.(string), jsonPath)
        } else {
            r = tcDataStore.GetRes(res.(string), jsonPath)
        }
    case "float64":
        r = res
    case "bool":
        r = res
    // case "map":
    //     fmt.Println("type is: ", t)
    case "slice", "map":
        s, err := json.Marshal(res)
        if err != nil {
            panic(err)
        }

        if strings.HasSuffix(jsonPath, "_keys_count_") {
            r = tcDataStore.GetKeysCount(string(s), jsonPath)
        } else {
            r = tcDataStore.GetRes(string(s), jsonPath)
        }
    default:
        r = res
    }

    return r
}

func (tcDataStore *TcDataStore) GetRes (s string, jsonPath string) interface{} {
    var resValue interface{}
    
    value := gjson.Get(s, jsonPath)

    // if the path (key) exists
    if !value.Exists() {
        // the key does not exist, set the actualValue = _null_key_
        resValue = "_null_key_"
    } else {
        vv := value.Value()

        switch vv.(type) {
        case nil:
            resValue = "_null_value_"
        case string:
            resValue = value.String()
        case bool:
            resValue = value.Bool()
        case float64:
            // for big int, to avoid the automatical convertion to 
            // scientific notation convertion for float64
            s := value.String()
            intI, err := strconv.Atoi(s)
            if err == nil {
                resValue = intI
            } else {
                resValue = value.Float()
            }
        case map[string]interface{}, []interface{}:
            r := value.Value()

            rs, err := json.Marshal(r)
            if err != nil {
                panic(err)
            }
            resValue = string(rs)
        default:
            resValue = value.Value()
        }
    }

    return resValue
}


func (tcDataStore *TcDataStore) GetKeysCount (s string, jsonPath string) interface{} {
    var resValue int
    var result gjson.Result

    if jsonPath == "_keys_count_" {
        result = gjson.Result {
            Type:  gjson.JSON,
            Raw:   s,
        }
    } else {
        subPath := jsonPath[0 : len(jsonPath) - len("_keys_count_") - 1]

        result = gjson.Get(s, subPath)
    }
    //
    i := 0
    result.ForEach(func(key, value gjson.Result) bool {
        i = i + 1
        return true // keep iterating
    })

    resValue = i

    return resValue
}

// -----
// Variable's Content, supports json path expression
func (tcDataStore *TcDataStore) GetVariableContent (s string) interface{} {  
    // var resValue interface{}

    s = strings.Replace(s, "$(", "", -1)
    s = strings.Replace(s, ")", "", -1)

    f := tcDataStore.RenderExpresionB(s)

    return f
}

// ------------------------

// redis
func (tcDataStore *TcDataStore) GetRedisActualValueByPath (searchPath string) interface{} {
    var resValue interface{}
 
    prefix := "$(redis)."
    lenPrefix := len(prefix)

    fmt.Println("==> tcDataStore.CmdResults: ", tcDataStore.CmdResults)

    cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    cmdResultsJson := string(cmdResultsB)

    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        if searchPath == "$(redis).#" {
            resValue = tcDataStore.CmdAffectedCount
        } else if searchPath == "$(redis).*" {
            resValue = tcDataStore.CmdResults
        } else if tcDataStore.IfCmdResultsPrimitive() {
            resValue = tcDataStore.CmdResults
        } else {
            value := gjson.Get(string(cmdResultsJson), searchPath[lenPrefix:])
            resValue = value.Value()
        }
    } else {
        resValue = searchPath
    }

    return resValue
}

//
func (tcDataStore *TcDataStore) IfCmdResultsPrimitive () bool {
    var cmdResultsJson string

    switch tcDataStore.CmdResults.(type) {
    case string:
        cmdResultsJson = tcDataStore.CmdResults.(string)
    default:
        cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
        cmdResultsJson = string(cmdResultsB)
    }

    // remove left \n
    ss := strings.TrimLeft(cmdResultsJson, "\n")
    // remove space
    ss = strings.TrimSpace(ss)

    if len(ss) == 0 {
        return true
    } else {
        // if read content from file, there may have " at left-most, so use ss[1:2]
        // to be fixed later !!! 
        if ss[0:1] == "[" || ss[0:1] == "{" || ss[1:2] == "[" || ss[1:2] == "{" {
            return false
        } else {
            return true
        }
    }
}

