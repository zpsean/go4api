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
    "encoding/json"
    // "encoding/xml"
    // "encoding/json"
    "strconv"
    "regexp"
    simplejson "github.com/bitly/go-simplejson"
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
    // fmt.Println("Contains", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if actualValue == expValue {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("Contains", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if act == exp {
            return true
        } else {
            return false
        }
    }
}   
    

func Contains(actualValue interface{}, expValue interface{}) bool {
    // fmt.Println("Contains", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if strings.Contains(actualValue.(string), expValue.(string)) {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("Contains", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if strings.Contains(act.(string), exp.(string)) {
            return true
        } else {
            return false
        }
    }
}

func StartsWith(actualValue interface{}, expValue interface{}) bool {
    // fmt.Println("Contains", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if strings.HasPrefix(actualValue.(string), expValue.(string)) {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("Contains", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if strings.HasPrefix(act.(string), exp.(string)) {
            return true
        } else {
            return false
        }
    }
}

func EndsWith(actualValue interface{}, expValue interface{}) bool {
    // fmt.Println("Contains", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if strings.HasSuffix(actualValue.(string), expValue.(string)) {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("Contains", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if strings.HasSuffix(act.(string), exp.(string)) {
            return true
        } else {
            return false
        }
    }
}

func NotEquals(actualValue interface{}, expValue interface{}) bool {
    // fmt.Println("GreaterOrEquals", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if actualValue.(float64) != expValue.(float64) {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("GreaterOrEquals", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if act.(float64) != exp.(float64) {
            return true
        } else {
            return false
        }
    }
}

func Less(actualValue interface{}, expValue interface{}) bool {
    // fmt.Println("GreaterOrEquals", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if actualValue.(float64) < expValue.(float64) {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("GreaterOrEquals", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if act.(float64) < exp.(float64) {
            return true
        } else {
            return false
        }
    }
}

func LessOrEquals(actualValue interface{}, expValue interface{}) bool {
    // fmt.Println("GreaterOrEquals", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if actualValue.(float64) <= expValue.(float64) {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("GreaterOrEquals", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if act.(float64) <= exp.(float64) {
            return true
        } else {
            return false
        }
    }
}

func Greater(actualValue interface{}, expValue interface{}) bool {
    // fmt.Println("GreaterOrEquals", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if actualValue.(float64) > expValue.(float64) {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("GreaterOrEquals", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if act.(float64) > exp.(float64) {
            return true
        } else {
            return false
        }
    }
}

func GreaterOrEquals(actualValue interface{}, expValue interface{}) bool {
    // fmt.Println("GreaterOrEquals", actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue)) 
    if reflect.TypeOf(actualValue) == reflect.TypeOf(expValue) {
        if actualValue.(float64) >= expValue.(float64) {
            return true
        } else {
            return false
        }
    } else {
        act, exp := CovertValuesBasedOnTypes(actualValue, expValue)
        // fmt.Println("GreaterOrEquals", act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 
        if act.(float64) >= exp.(float64) {
            return true
        } else {
            return false
        }
    }
}

func CovertValuesBasedOnTypes(actualValue interface{}, expValue interface{}) (interface{}, interface{}) {
    // Note: the possible combinations for actualtype and exptype are:
    // ==> status
    // string json.Number
    // ==> header
    // string string
    // ==> body
    // *simplejson.Json json.Number
    // *simplejson.Json string
    // int json.Number


    typeActualValue := reflect.TypeOf(actualValue)
    typeExpValue := reflect.TypeOf(expValue)

    // valueActualValue := reflect.ValueOf(actualValue)
    // valueExpValue := reflect.ValueOf(expValue)

    // fmt.Println("Convert types: ", typeActualValue, valueActualValue, typeExpValue, valueExpValue)

    // to check the valueExpValue first, it may be string, number, boolean, null, array, json, etc.
    var act, exp interface{}

    switch typeExpValue.String() {
        case "json.Number": {
            switch typeActualValue.String(){
                case "*simplejson.Json": {
                    act, _ = actualValue.(*simplejson.Json).Float64()
                    exp, _ = expValue.(json.Number).Float64()
                }
                case "string": {
                    act = actualValue

                    expF, _ := expValue.(json.Number).Float64()
                    exp = strconv.FormatFloat(expF, 'f', -1, 64)
                }
                case "int": {
                    act = float64(actualValue.(int))
                    exp, _ = expValue.(json.Number).Float64()
                }
            }
        }
        case "string": {
            act, _ = actualValue.(*simplejson.Json).String()
            exp = expValue
        }
        case "bool": {
            act, _ = actualValue.(*simplejson.Json).Bool()
            exp = expValue
        }
    }

    // fmt.Println("Convert types: ", actualValue, expValue, typeActualValue, typeExpValue, act, exp, reflect.TypeOf(act), reflect.TypeOf(exp)) 

    return act, exp
}


// For regrex, Match function, for value - value match 
// a is the key, wold be path, like: $.headers.Content-Type, $.body.resource[0], $.body.resource.count, etc. 
// a may be a simple concrete type liek string, number, boolean, null, etc. or other complex type like array, json, etc.
// b is the value, wold be regrex expression, like: application\\/json, ^\\d{4}-\\d{2}-\\d{2}$, etc.
// b may be a simple concrete type liek string, number, boolean, null, etc. or other complex type like array, json, etc. 
func Match(actualValue interface{}, expPattern interface{}) bool {
    act, expP := CovertValuesBasedOnTypes(actualValue, expPattern)   
    // fmt.Println("GreaterOrEquals", act, exp)

    ind, _ := regexp.MatchString(expP.(string), act.(string))

    if ind {
        return true
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


