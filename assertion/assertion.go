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
    // "regexp"
    // simplejson "github.com/bitly/go-simplejson"
)

// What to support assertion here:
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

// for general regrex
// Match


// for both String and Numeric
func Equals(a interface{}, b interface{}) bool {
    // fmt.Println("Equals", a, b, reflect.TypeOf(a), reflect.TypeOf(b))

    fb := b.(json.Number).String()

    if a == fb {
        return true
    } else {
        return false
    }
}

func Contains(a interface{}, b interface{}) bool {
    // fmt.Println("Contains", a, b, reflect.TypeOf(a), reflect.TypeOf(b))
    if strings.Contains(a.(string), b.(string)) {
        return true
    } else {
        return false
    }
}

func LargerThan(a interface{}, b interface{}) bool {
    // fmt.Println("LargerThan", a, b, reflect.TypeOf(a), reflect.TypeOf(b))

    fa := int64(a.(int))
    fb, _ := b.(json.Number).Int64()
    
    // fmt.Println("LargerThan", fa, fb)
    if fa > fb {
        return true
    } else {
        return false
    }
}


// For regrex, Match function, 
// a is the key, wold be path, like: $.headers.Content-Type, $.body.resource[0], $.body.resource.count, etc. 
// a may be a simple concrete type liek string, Int, etc. or other complex type like struct, map
// b is the value, wold be regrex expression, like: application\\/json, ^\\d{4}-\\d{2}-\\d{2}$, etc.
// b may be a simple concrete type liek string, Int, etc. or other complex type like struct, map
func Match(a interface{}, b interface{}) bool {
    // fmt.Println("LargerThan", a, b, reflect.TypeOf(a), reflect.TypeOf(b))

    fa := int64(a.(int))
    fb, _ := b.(json.Number).Int64()
    
    // fmt.Println("LargerThan", fa, fb)
    if fa > fb {
        return true
    } else {
        return false
    }
}

func CallAssertion(m map[string]interface{}, name string, params ... interface{}) (result []reflect.Value, err error) {
    f := reflect.ValueOf(m[name])
    // if len(params) != f.Type().NumIn() {
    //     err = errors.New("The number of params is not adapted.")
    //     return
    // }
    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }
    result = f.Call(in)

    return
}


