/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package assertion

import (
    // "fmt"
    "encoding/json"
    gjson "github.com/tidwall/gjson"
)

// To get the type of concrete type, the underlying type for *simplejson.Json, json.Number, ets.
// so that, he possible combinations for actualtype and exptype are:
// ==> status
// int json.Number
// ==> header
// string string
// ==> body
// *simplejson.Json json.Number
// *simplejson.Json string
// *simplejson.Json bool
// int json.Number
// ...

// ==> there are two options to deal with the types and values:
// Option 1: use the pakcage reflect to get the type, and determine if they are comparable, then compare
// Option 2: to get the raw data first, then determine if they are: string, number, bool, null and Raw (json), then compare

// after trying the Option 1, now prefer to use Option 2

// JSON Schema defines the following basic types:
// string
// Numeric -> float64
// boolean
// null
// object (raw)
// array (raw)


func GetRawJsonResult(value interface{}) (string, error) {
    // to get the raw json string using json.Marshal
    byteValue, err := json.Marshal(value)
    if err != nil {
        return "", err
    }

    return string(byteValue), err
}

func VerifyTypes(actualValue interface{}, expValue interface{}) string {
    act, _ := GetRawJsonResult(actualValue)
    exp, _ := GetRawJsonResult(expValue)

    actResult := gjson.Parse(act)
    expResult := gjson.Parse(exp)

    // as get Go nil, for JSON null, need special care, two possibilities:
    // p1: expResult -> null, but can not find out actualValue, go set it to nil, i.e. null (assertion -> false)
    // p2: expResult -> null, actualValue can be founc, and its value --> null (assertion -> true)
    // but here can not distinguish them

    if actualValue == nil || expValue == nil {
        if actualValue != nil || expValue != nil {
            return "false"
        } else {
            return "true"
        }
    } else if actResult.Type == expResult.Type {
        return "true"
    } else {
        return "false"
    }
}

func GetValue(value interface{}) interface {} {
    val, _ := GetRawJsonResult(value)

    valResult := gjson.Parse(val)

    return valResult.Value()
}

