/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package fuzz

import (
    "fmt"
    "testing"
    // "encoding/csv"
)

func Test_getMaxLenVector(t *testing.T) {
    data := [][]interface{}{{1, 2, 3, 4}, {5, 6, 7}}

    res := getMaxLenVector(data)

    if res != 4 {
        t.Fatalf("getMaxLenVector test failed")
    } else {
        t.Log("getMaxLenVector test passed")
    }
}

func Test_GetCombinationInvalid(t *testing.T) {
    validVectors := [][]interface{}{{1, 2, 3, 4}, {5, 6, 7}}
    invalidVectors := [][]interface{}{{"a", "b", "c", "d"}, {"e", "f", "g"}}

    res := GetCombinationInvalid(validVectors, invalidVectors, 2)

    fmt.Println("validVectors: ", validVectors)
    fmt.Println("invalidVectors: ", invalidVectors)
    fmt.Println("res: ", res)

    if len(res) != 4 * 7 || len(res[0]) != 2 {
        t.Fatalf("getMaxLenVector test failed")
    } else {
        t.Log("getMaxLenVector test passed")
    }
}