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
    "testing"
    // "encoding/csv"
)

// <------ GetType() ------>
func Test_GetType(t *testing.T) {
    value := "abcde"
    res := GetType(value)

    fmt.Println("res: ", res, fmt.Sprint(res))
    if res != "String" {
        t.Fatalf("GetType failed")
    }
}

func Test_GetType2(t *testing.T) {
    value := 1234
    res := GetType(value)

    fmt.Println("res: ", res, fmt.Sprint(res))
    if res != "Number" {
        t.Fatalf("GetType failed")
    }
}

func Test_GetType3(t *testing.T) {
    value := 1234.12
    res := GetType(value)

    fmt.Println("res: ", res, fmt.Sprint(res))
    if res != "Number" {
        t.Fatalf("GetType failed")
    }
}

func Test_GetType4(t *testing.T) {
    value := true
    res := GetType(value)

    fmt.Println("res: ", res, fmt.Sprint(res))
    if res != "True" {
        t.Fatalf("GetType failed")
    }
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


