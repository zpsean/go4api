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
    // "strings"

    "go4api/assertion"

    gjson "github.com/tidwall/gjson"
)

func GetSqlResponseValueff (searchPath string, count int, rows string) interface{} {
    // prefix = "$(sql)."
    var resValue interface{}

    prefix := "$(sql)."
    lenPrefix := len(prefix)

    if searchPath == "$(sql).Count" {
        return count
    }

    if len(searchPath) > lenPrefix && searchPath[0:lenPrefix] == prefix {
        value := gjson.Get(string(rows), searchPath[lenPrefix:])
        resValue = value.Value()
    } else {
        resValue = searchPath
    }
    
    return resValue
}

func SqLCompareCommoffn (assertionKey string, actualValue interface{}, expValue interface{}) {
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
    // neither nil
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

    fmt.Println(assertionResults, testRes)
}
