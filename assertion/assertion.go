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
	"reflect"
    "strings"
    // "encoding/json"
    // "strconv"
    "regexp"
)


// To support assertion here:
// if response body is xml: [key, using xpath] [operator, like Equals, ...] [value, can use regrex]
// if response body is html: [key, using xpath, css] [operator, like Equals, ...] [value, can use regrex]
// if response body is json: [key] [operator, like Equals, ...] [value, can use regrex]

// for String:
// Equals
// Contains
// StartsWith
// EndsWith

// for Numeric:
// Equals
// NotEquals
// Less
// LessOrEquals
// Greater
// GreaterOrEquals

// for Bool (true, false):
// Equals
// NotEquals

// for general regrex
// Match


// for both String and Numeric
func Equals(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        switch reflect.TypeOf(GetValue(actualValue)).String() {
            case "float64": {
                if GetValue(actualValue).(float64) == GetValue(expValue).(float64) {
                    return true
                } else {
                    return false
                }
            }
            case "string": {
                if GetValue(actualValue).(string) == GetValue(expValue).(string) {
                    return true
                } else {
                    return false
                }
            }
            case "bool": {
                if GetValue(actualValue).(bool) == GetValue(expValue).(bool) {
                    return true
                } else {
                    return false
                }
            }
            default:
                return false
            }
    } else {
        return false
    }
}   
    

func Contains(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        if strings.Contains(GetValue(actualValue).(string), GetValue(expValue).(string)) {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}

func StartsWith(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        if strings.HasPrefix(GetValue(actualValue).(string), GetValue(expValue).(string)) {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}

func EndsWith(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        if strings.HasSuffix(GetValue(actualValue).(string), GetValue(expValue).(string)) {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}

func NotEquals(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        switch reflect.TypeOf(GetValue(actualValue)).String() {
            case "float64": {
                if GetValue(actualValue).(float64) != GetValue(expValue).(float64) {
                    return true
                } else {
                    return false
                }
            }
            case "string": {
                if GetValue(actualValue).(string) != GetValue(expValue).(string) {
                    return true
                } else {
                    return false
                }
            }
            case "bool": {
                if GetValue(actualValue).(bool) != GetValue(expValue).(bool) {
                    return true
                } else {
                    return false
                }
            }
            default:
                return false
        }
    } else {
        return false
    }
}

func Less(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        if GetValue(actualValue).(float64) < GetValue(expValue).(float64) {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}

func LessOrEquals(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        if GetValue(actualValue).(float64) <= GetValue(expValue).(float64) {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}

func Greater(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        if GetValue(actualValue).(float64) > GetValue(expValue).(float64) {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}

func GreaterOrEquals(actualValue interface{}, expValue interface{}) bool {
    if VerifyTypes(actualValue, expValue) == "true" {
        if GetValue(actualValue).(float64) >= GetValue(expValue).(float64) {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}


// For regrex, Match function, for value - value match 
// a is the key, wold be path, like: $.headers.Content-Type, $.body.resource[0], $.body.resource.count, etc. 
// a may be a simple concrete type liek string, number, boolean, null, etc. or other complex type like array, json, etc.
// b is the value, wold be regrex expression, like: application\\/json, ^\\d{4}-\\d{2}-\\d{2}$, etc.
// b may be a simple concrete type liek string, number, boolean, null, etc. or other complex type like array, json, etc. 
func Match(actualValue interface{}, expPattern interface{}) bool {
    if VerifyTypes(actualValue, expPattern) == "true" {
        ind, _ := regexp.MatchString(GetValue(actualValue).(string), GetValue(expPattern).(string))

        if ind {
            return true
        } else {
            return false
        }
    } else {
        return false
    }
}


//
func CallAssertion(name string, params ... interface{}) bool {
    funcs := map[string]interface{} {
        "Equals": Equals,
        "Contains": Contains,
        "StartsWith": StartsWith,
        "EndsWith": EndsWith,
        "NotEquals": NotEquals,
        "Less": Less,
        "LessOrEquals": LessOrEquals,
        "Greater": Greater,
        "GreaterOrEquals": GreaterOrEquals,
        "Match": Match,
    }

    f := reflect.ValueOf(funcs[name])
    // if len(params) != f.Type().NumIn() {
    //     err: = errors.New("The number of params is not adapted.")
    //     return
    // }
    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }
    result := f.Call(in)

    return result[0].Interface().(bool)
}


