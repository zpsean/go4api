/*
 * go4api - an api testing tool written in Go
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
    "reflect"
    "encoding/json"
)

// use the pakcage reflect to get the type, and determine if they are comparable, then compare
// ----
// JSON Schema defines the following basic types:
// string
// Numeric -> float64
// boolean
// null
// object (raw)
// array (raw)

func ValidateCallName (name string) bool {
    if _, ok := assertionMapping[name]; ok {
        return true
    } else {
        return false
    }
}

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
    // (2). no nil, if two type assertable
    typeAct := reflect.TypeOf(params[0]).Kind().String()
    typeExp := reflect.TypeOf(params[1]).Kind().String()

    // fmt.Println("typeAct, typeExp: ", typeAct, typeExp)

    // consider the type int, float64, they are comparable
    if typeAct == "int" && typeExp == "float64" {
        return true
    } else if typeAct == "int64" && typeExp == "float64" {
        return true
    } else if typeAct == "float64" && typeExp == "int" {
        return true
    } else if typeAct == "float64" && typeExp == "int64" {
        return true
    }

    // consider slice
    switch typeAct {
        case "slice":
            switch typeExp {
                case "slice": 
                    return true
                case "string", "int", "int64", "float64", "bool":
                    lowerName := strings.ToLower(name)
                    if lowerName == "has" || lowerName == "nothas" {
                        return true
                    } else {
                        return false
                    }
                default:
                    return false
            }
        case "string", "int", "int64", "float64", "bool":
            switch typeExp {
                case "slice": 
                    lowerName := strings.ToLower(name)
                    if lowerName == "in" || lowerName == "notin" {
                        return true
                    } else {
                        return false
                    }
                case "string", "int", "int64", "float64", "bool":
                    if typeAct != typeExp {
                        return false
                    }
                default:
                    return false
            }
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

func GetRawJsonResult (value interface{}) string {
    // to get the raw json string using json.Marshal
    byteValue, err := json.Marshal(value)
    if err != nil {
        return ""
    }

    return string(byteValue)
}
