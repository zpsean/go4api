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
        case strings.HasPrefix(searchPath, "$(postgresql)."):
            value = tcDataStore.GetPgSqlActualValueByPath(searchPath)
        case strings.HasPrefix(searchPath, "$(redis)."):
            value = tcDataStore.GetRedisActualValueByPath(searchPath)
        case strings.HasPrefix(searchPath, "$(mongodb)."):
            value = tcDataStore.GetMongoDBActualValueByPath(searchPath)
        case strings.HasPrefix(searchPath, "$(file)."):
            value = tcDataStore.GetFileActualValueByPath(searchPath)
        // case strings.HasPrefix(searchPath, "$."):
        //     value = tcDataStore.GetBodyActualValueByPath(searchPath)
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

// http response body
func (tcDataStore *TcDataStore) GetBodyActualValueByPath (key string) interface{} {  
    var actualValue interface{}
    actualBody := tcDataStore.HttpActualBody

    prefix := "$(body)."
    lenPrefix := len(prefix)
    prefix2 := "$."
    lenPrefix2 := len(prefix2)

    if key == prefix + "*" || key == prefix2 + "*" {
        actualValue = string(actualBody)
    } else if len(key) > lenPrefix && key[0:lenPrefix] == prefix {
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

// mysql
func (tcDataStore *TcDataStore) GetSqlActualValueByPath (searchPath string) interface{} {
    var resValue interface{}
 
    prefix := "$(sql)."
    lenPrefix := len(prefix)

    cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    cmdResultsJson := string(cmdResultsB)

    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        switch {
        case searchPath == "$(sql).affectedCount":
            resValue = tcDataStore.CmdAffectedCount
        case searchPath == "$(sql).*":
            resValue = tcDataStore.CmdResults
        case tcDataStore.IfCmdResultsPrimitive():
            resValue = tcDataStore.CmdResults
        default:
            value := gjson.Get(string(cmdResultsJson), searchPath[lenPrefix:])
            resValue = value.Value()
        }

        // if searchPath == "$(sql).#" {
        //     resValue = tcDataStore.CmdAffectedCount
        // } else if searchPath == "$(sql).*" {
        //     resValue = tcDataStore.CmdResults
        // } else if tcDataStore.IfCmdResultsPrimitive() {
        //     resValue = tcDataStore.CmdResults
        // } else {
        //     value := gjson.Get(string(cmdResultsJson), searchPath[lenPrefix:])
        //     resValue = value.Value()
        // }
    } else {
        resValue = searchPath
    }

    return resValue
}

// postgresql
func (tcDataStore *TcDataStore) GetPgSqlActualValueByPath (searchPath string) interface{} {
    var resValue interface{}
 
    prefix := "$(postgresql)."
    lenPrefix := len(prefix)

    cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    cmdResultsJson := string(cmdResultsB)

    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        switch {
        case searchPath == "$(postgresql).affectedCount":
            resValue = tcDataStore.CmdAffectedCount
        case searchPath == "$(postgresql).*":
            resValue = tcDataStore.CmdResults
        case tcDataStore.IfCmdResultsPrimitive():
            resValue = tcDataStore.CmdResults
        default:
            value := gjson.Get(string(cmdResultsJson), searchPath[lenPrefix:])
            resValue = value.Value()
        }
    } else {
        resValue = searchPath
    }

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

// MongoDB
func (tcDataStore *TcDataStore) GetMongoDBActualValueByPath (searchPath string) interface{} {
    var resValue interface{}
 
    prefix := "$(mongodb)."
    lenPrefix := len(prefix)

    // cmdResultsB, _ := json.Marshal(tcDataStore.CmdResults)
    // cmdResultsJson := string(cmdResultsB)

    cmdResultsJson := tcDataStore.CmdResults.(string)

    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        switch {
        case searchPath == "$(mongodb).affectedCount":
            resValue = tcDataStore.CmdAffectedCount
        case searchPath == "$(mongodb).*":
            resValue = tcDataStore.CmdResults
        case tcDataStore.IfCmdResultsPrimitive():
            resValue = tcDataStore.CmdResults
        default:
            value := gjson.Get(string(cmdResultsJson), searchPath[lenPrefix:])
            resValue = value.Value()
        }

        // if searchPath == "$(sql).#" {
        //     resValue = tcDataStore.CmdAffectedCount
        // } else if searchPath == "$(sql).*" {
        //     resValue = tcDataStore.CmdResults
        // } else if tcDataStore.IfCmdResultsPrimitive() {
        //     resValue = tcDataStore.CmdResults
        // } else {
        //     value := gjson.Get(string(cmdResultsJson), searchPath[lenPrefix:])
        //     resValue = value.Value()
        // }
    } else {
        resValue = searchPath
    }

    return resValue
}

// file
func (tcDataStore *TcDataStore) GetFileActualValueByPath (searchPath string) interface{} {
    var resValue interface{}
 
    prefix := "$(file)."
    lenPrefix := len(prefix)

    cmdResultsJson := tcDataStore.CmdResults.(string)
 
    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        switch {
        case tcDataStore.IfCmdResultsPrimitive():
            resValue = tcDataStore.CmdResults
        default:
            value := gjson.Get(cmdResultsJson, searchPath[lenPrefix:])
            resValue = value.Value()
        }
    } else {
        resValue = searchPath
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
