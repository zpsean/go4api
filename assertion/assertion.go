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
    // "os"
	"reflect"
    "strings"
    // "encoding/json"
    // "strconv"
    "regexp"
)

//
func CallAssertion(name string, params ... interface{}) bool {
    if !ValidateCallName(name) {
        return false
    }

    if ifBothNil(params) {
        return true
    } else if !ValidateCallParams(name, params) {
        return false
    } else {
        f := reflect.ValueOf(assertionMapping[name].AssertionFunc)

        in := make([]reflect.Value, len(params))
        for k, param := range params {
            in[k] = reflect.ValueOf(param)
        }
        
        result := f.Call(in)

        return result[0].Interface().(bool)
    }
}


// for both String and Numeric
func Equals(actualValue interface{}, expValue interface{}) bool {
    switch reflect.TypeOf(actualValue).Kind().String() {
        case "int", "float64": {
            var actualValueConverted float64
            var expValueConverted float64

            if reflect.TypeOf(actualValue).String() == "int" {
                actualValueConverted = float64(actualValue.(int))
            } else {
                actualValueConverted = actualValue.(float64)
            }

            if reflect.TypeOf(expValue).String() == "int" {
                expValueConverted = float64(expValue.(int))
            } else {
                expValueConverted = expValue.(float64)
            }
            //
            if actualValueConverted == expValueConverted {
                return true
            } else {
                return false
            }
        }
        case "string": {
            if actualValue.(string) == expValue.(string) {
                return true
            } else {
                return false
            }
        }
        case "bool": {
            if actualValue.(bool) == expValue.(bool) {
                return true
            } else {
                return false
            }
        }
        default:
            if reflect.DeepEqual(actualValue, expValue) {
                return true
            } else {
                return false
            }
    } 
}   
    

func Contains(actualValue interface{}, expValue interface{}) bool {
    if strings.Contains(actualValue.(string), expValue.(string)) {
        return true
    } else {
        return false
    }
}

func StartsWith(actualValue interface{}, expValue interface{}) bool {
    if strings.HasPrefix(actualValue.(string), expValue.(string)) {
        return true
    } else {
        return false
    }
}

func EndsWith(actualValue interface{}, expValue interface{}) bool {
    if strings.HasSuffix(actualValue.(string), expValue.(string)) {
        return true
    } else {
        return false
    }
}

func NotEquals(actualValue interface{}, expValue interface{}) bool {
    return !reflect.DeepEqual(actualValue, expValue)
}

func Less(actualValue interface{}, expValue interface{}) bool {
    if actualValue.(float64) < expValue.(float64) {
        return true
    } else {
        return false
    }
}

func LessOrEquals(actualValue interface{}, expValue interface{}) bool {
    if actualValue.(float64) <= expValue.(float64) {
        return true
    } else {
        return false
    }
}

func Greater(actualValue interface{}, expValue interface{}) bool {
    if actualValue.(float64) > expValue.(float64) {
        return true
    } else {
        return false
    }
}

func GreaterOrEquals(actualValue interface{}, expValue interface{}) bool {
    if actualValue.(float64) >= expValue.(float64) {
        return true
    } else {
        return false
    }
}

// for regrex
func Match(actualValue interface{}, expPattern interface{}) bool {
    reg := regexp.MustCompile(expPattern.(string))
    resSlice := reg.FindAllString(actualValue.(string), -1)

    if resSlice != nil {
        return true
    } else {
        return false
    }
}



