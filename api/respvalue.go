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
    // "reflect"
    "encoding/json"

    "go4api/assertion" 
    "go4api/lib/testcase" 
    
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
                value = tcDataStore.GetStatusActualValue(searchPath)
            case "$(headers)":
                value = tcDataStore.GetHeadersActualValue(searchPath)
            case "$(body)":
                value = tcDataStore.GetBodyActualValueByPath(z[1])
            case "$(sql)", "$(mysql)", "$(postgresql)", "$(mongodb)":
                value = tcDataStore.GetSqlActualValueByPath(z[1])
            case "$(file)":
                value = tcDataStore.GetFileActualValueByPath(z[1])
            case "$(redis)":
                value = tcDataStore.GetRedisActualValueByPath(searchPath)
            default:
                value = searchPath
            } 
        } else {
            value = searchPath
        }
    } else {
        value = searchPath
    }
    
    return value
}

func (tcDataStore *TcDataStore) GetStatusActualValue (key string) interface{} {
    var actualValue interface{}
    actualStatusCode := tcDataStore.HttpActualStatusCode
    
    prefix := "$(status)."
    lenPrefix := len(prefix)

    if len(key) == lenPrefix && key[0:lenPrefix] == prefix {
        actualValue = actualStatusCode
    } else {
        actualValue = key
    }

    return actualValue
}

func (tcDataStore *TcDataStore) GetHeadersActualValue (key string) interface{} { 
    var actualValue interface{}
    actualHeader := tcDataStore.HttpActualHeader
 
    prefix := "$(headers)."
    lenPrefix := len(prefix)

    if len(key) > lenPrefix && key[0:lenPrefix] == prefix {
        actualValue = strings.Join(actualHeader[key[lenPrefix:]], ",")
    } else {
        actualValue = key
    }

    return actualValue
}


// http response body
func (tcDataStore *TcDataStore) GetBodyActualValueByPath (jsonPath string) interface{} {  
    var resValue interface{}

    actualBodyString := string(tcDataStore.HttpActualBody)

    switch {
    case jsonPath == "*":
        resValue = actualBodyString
    default:
        b := strings.HasSuffix(jsonPath, "__keys_count_")

        if b == true {
            var result gjson.Result

            if jsonPath == "__keys_count_" {
                result = gjson.Result {
                    Type:  gjson.JSON,
                    Raw:   actualBodyString,
                }
            } else {
                subPath := jsonPath[0 : len(jsonPath) - len("__keys_count_") - 1]

                result = gjson.Get(actualBodyString, subPath)
            }
            //
            i := 0
            result.ForEach(func(key, value gjson.Result) bool {
                i = i + 1
                return true // keep iterating
            })

            resValue = i

        } else {
            value := gjson.Get(actualBodyString, jsonPath)
            // if the path (key) exists
            if !value.Exists() {
                // the key does not exist, set the actualValue = _null_key_
                resValue = "_null_key_"
            } else {
                resValue = value.Value()
            }
        }   
    }

    if resValue == nil {
        resValue = "_null_value_"
    }

    return resValue
}


// for RDBMS sql
func (tcDataStore *TcDataStore) GetSqlActualValueByPath (jsonPath string) interface{} {
    var resValue interface{}

    if len(tcDataStore.CmdResults.([]map[string]interface {})) == 0 {
        resValue = "_null_key_"

        return resValue
    }

    resValue = tcDataStore.CommonGetActualValueByPath(jsonPath)

    return resValue
}


// redis
func (tcDataStore *TcDataStore) GetRedisActualValueByPath (searchPath string) interface{} {
    var resValue interface{}
 
    prefix := "$(redis)."
    lenPrefix := len(prefix)

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


// file
func (tcDataStore *TcDataStore) GetFileActualValueByPath (jsonPath string) interface{} {
    var resValue interface{}

    resValue = tcDataStore.CommonGetActualValueByPath(jsonPath)

    return resValue
}


// common json compatible value search
func (tcDataStore *TcDataStore) CommonGetActualValueByPath (jsonPath string) interface{} {
    var resValue interface{}

    cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    cmdResultsJson := string(cmdResultsB)

    switch {
    // case jsonPath == "affectedCount":
    //     resValue = tcDataStore.CmdAffectedCount
    case jsonPath == "*":
        resValue = tcDataStore.CmdResults
    case tcDataStore.IfCmdResultsPrimitive():
        resValue = tcDataStore.CmdResults
    default:
        value := gjson.Get(string(cmdResultsJson), jsonPath)
        // if the path (key) exists
        if !value.Exists() {
            // the key does not exist, set the actualValue = _null_key_
            resValue = "_null_key_"
        } else {
            resValue = value.Value()
        }
    }

    if resValue == nil {
        resValue = "_null_value_"
    }

    return resValue
}


//
func (tcDataStore *TcDataStore) IfCmdResultsPrimitive () bool {
    cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    cmdResultsJson := string(cmdResultsB)

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

//
func compareCommon (reponsePart string, key string, assertionKey string, actualValue interface{}, expValue interface{}) (bool, *testcase.TestMessage) {
    // Note: As get Go nil, for JSON null, need special care, two possibilities:
    // p1: expResult -> null, but can not find out actualValue, go set it to nil, i.e. null (assertion -> false)
    // p2: expResult -> null, actualValue can be founc, and its value --> null (assertion -> true)
    // but here can not distinguish them
    assertionResults := ""
    var testRes bool

    // reserved word: _ignore_assertion_
    if fmt.Sprint(expValue) == "_ignore_assertion_" {
        msg := testcase.TestMessage {
            AssertionResults: "Success",
            ReponsePart:      reponsePart,
            FieldName:        key,
            AssertionKey:     assertionKey,
            ActualValue:      actualValue,
            ExpValue:         expValue,   
        }
        
        testRes = true

        return testRes, &msg
    } 

    if actualValue == nil || expValue == nil {
        // if only one nil
        if actualValue != nil || expValue != nil {
            assertionResults = "Failed"
            testRes = false
        // both nil
        } else {
            assertionResults = "Success"
            testRes = true
        }
    // no nil
    } else {
        // call the assertion function
        testResult := assertion.CallAssertion(assertionKey, actualValue, expValue)
        
        if testResult == false {
            assertionResults = "Failed"
            testRes = false
        } else {
            assertionResults = "Success"
            testRes = true
        }
    }
    //
    msg := testcase.TestMessage {
        AssertionResults: assertionResults,
        ReponsePart:      reponsePart,
        FieldName:        key,
        AssertionKey:     assertionKey,
        ActualValue:      actualValue,
        ExpValue:         expValue,   
    }

    return testRes, &msg
}
