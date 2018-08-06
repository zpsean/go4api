package assertion

import (
    // "fmt"
	"reflect"
    "strings"
    // simplejson "github.com/bitly/go-simplejson"
)

func Equals(a interface{}, b interface{}) bool {
    if a == b {
        return true
    } else {
        return false
    }
}

func Contains(a interface{}, b interface{}) bool {
    if strings.Contains(b.(string), a.(string)) {
        return true
    } else {
        return false
    }
}

// func LargerThan(a interface{}, b interface{}) bool {
//     if a > b {
//         return true
//     } else {
//         return false
//     }
// }

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