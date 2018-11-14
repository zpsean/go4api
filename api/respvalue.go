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

    "go4api/assertion" 
    "go4api/lib/testcase" 
    
    gjson "github.com/tidwall/gjson"
)

func (tcDataStore *TcDataStore) GetResponseValue (searchPath string) interface{} {
    // prefix = "$(status).", "$(headers).", "$(body)."
    var value interface{}

    switch {
        case strings.HasPrefix(searchPath, "$(status)."):
            value = tcDataStore.GetStatusActualValue(searchPath)
        case strings.HasPrefix(searchPath, "$(headers)."):
            value = tcDataStore.GetHeadersActualValue(searchPath)
        case strings.HasPrefix(searchPath, "$(body)."):
            value = tcDataStore.GetBodyActualValueByPath(searchPath)
        case strings.HasPrefix(searchPath, "$(sql)."):
            value = tcDataStore.GetSqlActualValueByPath(searchPath)
        case strings.HasPrefix(searchPath, "$(redis)."):
            value = tcDataStore.GetRedisActualValueByPath(searchPath)
        case strings.HasPrefix(searchPath, "$."):
            value = tcDataStore.GetBodyActualValueByPath(searchPath)
        default:
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

func (tcDataStore *TcDataStore) GetBodyActualValueByPath (key string) interface{} {  
    var actualValue interface{}
    actualBody := tcDataStore.HttpActualBody

    prefix := "$(body)."
    lenPrefix := len(prefix)
    prefix2 := "$."
    lenPrefix2 := len(prefix2)

    if len(key) > lenPrefix && key[0:lenPrefix] == prefix {
        value := gjson.Get(string(actualBody), key[lenPrefix:])
        actualValue = value.Value()
    } else if len(key) > lenPrefix2 && key[0:lenPrefix2] == prefix2 {
        value := gjson.Get(string(actualBody), key[lenPrefix2:])
        actualValue = value.Value()
    } else {
        actualValue = key
    }

    return actualValue
}


func (tcDataStore *TcDataStore) GetSqlActualValueByPath (searchPath string) interface{} {
    var resValue interface{}
 
    prefix := "$(sql)."
    lenPrefix := len(prefix)

    cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    cmdResultsJson := string(cmdResultsB)

    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        if searchPath == "$(sql).*" {
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

func (tcDataStore *TcDataStore) GetRedisActualValueByPath (searchPath string) interface{} {
    var resValue interface{}
 
    prefix := "$(redis)."
    lenPrefix := len(prefix)

    cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    cmdResultsJson := string(cmdResultsB)

    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        if searchPath == "$(redis).*" {
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

func (tcDataStore *TcDataStore) IfCmdResultsPrimitive () bool {
    cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    cmdResultsJson := string(cmdResultsB)

    ss := strings.TrimSpace(cmdResultsJson)
    if len(ss) == 0 {
        return true
    } else {
        if ss[0:1] == "[" || ss[0:1] == "{" {
            return false
        } else {
            return true
        }
    }
}

func compareCommon (reponsePart string, key string, assertionKey string, actualValue interface{}, expValue interface{}) (bool, *testcase.TestMessage) {
    // Note: As get Go nil, for JSON null, need special care, two possibilities:
    // p1: expResult -> null, but can not find out actualValue, go set it to nil, i.e. null (assertion -> false)
    // p2: expResult -> null, actualValue can be founc, and its value --> null (assertion -> true)
    // but here can not distinguish them
    assertionResults := ""
    var testRes bool

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
        ReponsePart: reponsePart,
        FieldName: key,
        AssertionKey:  assertionKey,
        ActualValue: actualValue,
        ExpValue: expValue,   
    }

    return testRes, &msg
}
