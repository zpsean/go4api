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
    "fmt"
    "strings"
    "encoding/json"

    gjson "github.com/tidwall/gjson"
)

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

func ValidateCallParams (name string, params []interface{}) bool {
    if len(params) != 2 {
        fmt.Println("!! Warning, the number of params is not adapted, false", len(params), params)
        return false
    }
    // (1). nil, if match
    if params[0] == nil || params[1] == nil {
        if params[0] != nil || params[1] != nil {
            // only one nil
            return false
        } else {
            // both nil
            return true
        }
    } 
    // (2). no nil, if two type match
    typeAct := GetType(params[0])
    typeExp := GetType(params[1])

    if typeAct != typeExp {
        return false
    }
    // (3). no nil, if type matches with the mapping
    ifMatch := false

    for _, value := range assertionMapping[name].ApplyTypes {
        if strings.ToLower(typeAct) == value {
            ifMatch = true
            break
        }
    }
 
    return ifMatch
}

func ifBothNil (params []interface{}) bool {
    // (1). if nil
    // Note: As get Go nil, for JSON null, need special care, two possibilities:
    // p1: expResult -> null, but can not find out actualValue, go set it to nil, i.e. null (assertion -> false)
    // p2: expResult -> null, actualValue can be founc, and its value --> null (assertion -> true)
    // but here can not distinguish them
    if params[0] == nil || params[1] == nil {
        if params[0] != nil || params[1] != nil {
            // only one nil
            return false
        } else {
            // both nil
            return true
        }
    } else {
        return false
    }
}

func GetValue (value interface{}) interface {} {
    val, _ := GetRawJsonResult(value)

    valResult := gjson.Parse(val)

    return valResult.Value()
}

func GetType (value interface{}) string {
    rawRes, _ := GetRawJsonResult(value)
    gjsonRes := gjson.Parse(rawRes)

    return fmt.Sprint(gjsonRes.Type)
}

func GetRawJsonResult (value interface{}) (string, error) {
    // to get the raw json string using json.Marshal
    byteValue, err := json.Marshal(value)
    if err != nil {
        return "", err
    }

    return string(byteValue), err
}
