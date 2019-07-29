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
    "testing"
    "reflect"
)

// <--------- GetRawJsonResult() ------------>
func Test_GetRawJsonResult(t *testing.T) {
    value := "abcd"
    res := GetRawJsonResult(value)

    fmt.Println("res: ", res, reflect.TypeOf(res), fmt.Sprint(res), value, reflect.TypeOf(value), fmt.Sprint(value))
    if res != `"abcd"` {
        t.Fatalf("GetRawJsonResult failed")
    }
}

func Test_GetRawJsonResult2(t *testing.T) {
    value := 1234
    res := GetRawJsonResult(value)

    fmt.Println("res: ", res, reflect.TypeOf(res), fmt.Sprint(res), value, reflect.TypeOf(value), fmt.Sprint(value))
    if res != "1234" {
        t.Fatalf("GetRawJsonResult failed")
    }
}

func Test_GetRawJsonResult3(t *testing.T) {
    aa := make([]interface{}, 1)
    aa[0] = 2

    a := aa[0]

    fmt.Println(a)
    fmt.Println(float64(a.(int)))
}

// <--------- ValidateCallParams() ------------>
func Test_ValidateCallParams(t *testing.T) {
    name := "Contains"
    values := []interface{}{"1234", 12}
    res := ValidateCallParams(name, values)

    fmt.Println("res: ", res, fmt.Sprint(res))
    if res != false {
        t.Fatalf("ValidateCallParams failed")
    }
}

func Test_ValidateCallParams2(t *testing.T) {
    name := "Contains"
    values := []interface{}{"1234", "12"}
    res := ValidateCallParams(name, values)

    fmt.Println("res: ", res, fmt.Sprint(res))
    if res != true {
        t.Fatalf("ValidateCallParams failed")
    }
}


func Test_ValidateCallParams3(t *testing.T) {
    name := "Contains"
    values := []interface{}{1234.12, 12}
    res := ValidateCallParams(name, values)

    fmt.Println("res: ", res, fmt.Sprint(res))
    if res != true {
        t.Fatalf("ValidateCallParams failed")
    }
}

