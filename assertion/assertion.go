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
    // "fmt"
    // "os"
	"reflect"
    "strings"
    // "encoding/json"
    // "strconv"
    "regexp"
)

//
func CallAssertion (name string, params ... interface{}) bool {
    if !ValidateCallName(name) {
        return false
    }

    if ifBothNil(params) {
        return true
    } else if !ValidateCallParams(name, params) {
        var r bool

        n := strings.ToLower(name)
        if strings.Contains(n, "not") {
            r = true
        } else {
            r = false
        }

        return r
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
func Equals (actualValue interface{}, expValue interface{}) bool {
    switch reflect.TypeOf(actualValue).Kind().String() {
        case "int", "float64": {
            actualValueConverted, expValueConverted := convertIntToFloat64 (actualValue, expValue)

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

func convertIntToFloat64 (actualValue interface{}, expValue interface{}) (float64, float64) {
    var actualValueConverted float64
    var expValueConverted float64

    if reflect.TypeOf(actualValue).String() == "int" {
        actualValueConverted = float64(actualValue.(int))
    } else if reflect.TypeOf(actualValue).String() == "int64" {
        actualValueConverted = float64(actualValue.(int64))
    } else {
        actualValueConverted = actualValue.(float64)
    }

    if reflect.TypeOf(expValue).String() == "int" {
        expValueConverted = float64(expValue.(int))
    } else if reflect.TypeOf(expValue).String() == "int64" {
        expValueConverted = float64(expValue.(int64))
    } else {
        expValueConverted = expValue.(float64)
    }

    return actualValueConverted, expValueConverted
}

func NotEquals (actualValue interface{}, expValue interface{}) bool {
    return !reflect.DeepEqual(actualValue, expValue)
}

// string
func Contains (actualValue interface{}, expValue interface{}) bool {
    if strings.Contains(actualValue.(string), expValue.(string)) {
        return true
    } else {
        return false
    }
}

func StartsWith (actualValue interface{}, expValue interface{}) bool {
    if strings.HasPrefix(actualValue.(string), expValue.(string)) {
        return true
    } else {
        return false
    }
}

func EndsWith (actualValue interface{}, expValue interface{}) bool {
    if strings.HasSuffix(actualValue.(string), expValue.(string)) {
        return true
    } else {
        return false
    }
}

// float
func Less (actualValue interface{}, expValue interface{}) bool {
    actualValueConverted, expValueConverted := convertIntToFloat64 (actualValue, expValue)
    
    if actualValueConverted < expValueConverted {
        return true
    } else {
        return false
    }
}

func LessOrEquals (actualValue interface{}, expValue interface{}) bool {
    actualValueConverted, expValueConverted := convertIntToFloat64 (actualValue, expValue)
    
    if actualValueConverted <= expValueConverted {
        return true
    } else {
        return false
    }
}

func Greater (actualValue interface{}, expValue interface{}) bool {
    actualValueConverted, expValueConverted := convertIntToFloat64 (actualValue, expValue)
    
    if actualValueConverted > expValueConverted {
        return true
    } else {
        return false
    }
}

func GreaterOrEquals (actualValue interface{}, expValue interface{}) bool {
    actualValueConverted, expValueConverted := convertIntToFloat64 (actualValue, expValue)
    
    if actualValueConverted >= expValueConverted {
        return true
    } else {
        return false
    }
}

// In, NotIn, Has, NotHas, for item in item slice
func In (actualValue interface{}, expValue interface{}) bool {
    var ifIn bool

    for _, value := range reflect.ValueOf(expValue).Interface().([]interface{}) {
        if CallAssertion("Equals", actualValue, value) {
            ifIn = true
            break
        }
    }

    return ifIn
}

func NotIn (actualValue interface{}, expValue interface{}) bool {
    return !In(actualValue, expValue)
}


func Has (actualValue interface{}, expValue interface{}) bool {
    var ifHas bool

    for _, value := range reflect.ValueOf(actualValue).Interface().([]interface{}) {
        if CallAssertion("Equals", value, expValue) {
            ifHas = true
            break
        }
    }

    return ifHas
}

func NotHas (actualValue interface{}, expValue interface{}) bool {
    return !Has(actualValue, expValue)

}

// for regrex
func Match (actualValue interface{}, expPattern interface{}) bool {
    reg := regexp.MustCompile(expPattern.(string))
    resSlice := reg.FindAllString(actualValue.(string), -1)

    if resSlice != nil {
        return true
    } else {
        return false
    }
}

// special cases for map:
// map has key
func HasMapKey (actualValue interface{}, expValue interface{}) {
    
}

func NotHasMapKey (actualValue interface{}, expValue interface{}) {
    
}

// has key but value is null
func IsNull (actualValue interface{}, expValue interface{}) {
    
}

// has key but value is not null
func IsNotNull (actualValue interface{}, expValue interface{}) {
    
}


