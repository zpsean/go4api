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
    "encoding/json"
    
    gjson "github.com/tidwall/gjson"
)

// (tcDataStore *TcDataStore)
func GetResponseValue (searchPath string, actualStatusCode int, actualHeader map[string][]string, actualBody []byte) interface{} {
    // prefix = "$(status).", "$(headers).", "$(body)."
    var value interface{}
    if len(searchPath) > 1 {
        if strings.HasPrefix(searchPath, "$(status).") {
            value = GetStatusActualValue(searchPath, actualStatusCode)
        } else if strings.HasPrefix(searchPath, "$(headers).") {
            value = GetHeadersActualValue(searchPath, actualHeader)
        } else if strings.HasPrefix(searchPath, "$(body).") {
            value = GetActualValueByJsonPath(searchPath, actualBody)
        } else if strings.HasPrefix(searchPath, "$.") {
            value = GetActualValueByJsonPath(searchPath, actualBody)
        } else {
            value = searchPath
        }
    } else {
        value = searchPath
    }
    
    return value
}

func GetStatusActualValue (key string, actualStatusCode int) interface{} {
    var actualValue interface{}
    // leading "$(status)" is mandatory if want to retrive status
    if len(key) == 9 && key == "$(status)" {
        actualValue = actualStatusCode
    } else {
        actualValue = key
    }

    return actualValue
}

func GetHeadersActualValue (key string, actualHeader map[string][]string) interface{} { 
    var actualValue interface{}
    // leading "$(headers)" is mandatory if want to retrive headers value
    prefix := "$(headers)."
    lenPrefix := len(prefix)

    if len(key) > lenPrefix && key[0:lenPrefix] == prefix {
        actualValue = strings.Join(actualHeader[key[lenPrefix:]], ",")
    } else {
        actualValue = key
    }

    return actualValue
}

func GetActualValueByJsonPath (key string, actualBody []byte) interface{} {  
    var actualValue interface{}
    // leading "$." or "$(headers)." is mandatory if want to use path search
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


func GetSqlActualRespValue (searchPath string, rowsCount int, rowsData []map[string]interface{}) interface{} {
    // prefix = "$(sql)."
    var resValue interface{}
 
    prefix := "$(sql)."
    lenPrefix := len(prefix)

    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        rowsDataB, _ := json.Marshal(rowsData)
        rowsDataJson := string(rowsDataB)

        value := gjson.Get(string(rowsDataJson), searchPath[lenPrefix:])
        resValue = value.Value()
    } else {
        resValue = searchPath
    }

    return resValue
}

