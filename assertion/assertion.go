package assertion

import (
    // "fmt"
	"reflect"
    "strings"
    "encoding/json"
    // simplejson "github.com/bitly/go-simplejson"
)

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