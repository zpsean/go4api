/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package combins

import (
    "fmt"
    "testing"
    // "encoding/csv"
)

func Test_CombinationsInt(t *testing.T) {
    var res [][]int

    list := []int{1, 2, 3, 4}
    length := 2

    c := CombinationsInt(list, length)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(list, res)
    if len(res) != 6 || len(res[0]) != length {
        t.Fatalf("combinationsInt test failed")
    } else {
        t.Log("combinationsInt test passed")
    }
}


func Test_CombinationsInt2(t *testing.T) {
    var res [][]int

    list := []int{1, 2, 3, 4}
    length := 4

    c := CombinationsInt(list, length)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(list, res)
    if len(res) != 1 || len(res[0]) != length {
        t.Fatalf("combinationsInt test failed")
    } else {
        t.Log("combinationsInt test passed")
    }
}

func Test_CombinationsInt3(t *testing.T) {
    var res [][]int

    list := []int{1, 2, 3, 4}
    length := 1

    c := CombinationsInt(list, length)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(list, res)
    if len(res) != 4 || len(res[0]) != length {
        t.Fatalf("combinationsInt test failed")
    } else {
        t.Log("combinationsInt test passed")
    }
}

//
func Test_CombinationsInterface(t *testing.T) {
    var res [][]interface{}

    list := []interface{}{1, 2, 3, 4}
    length := 3

    c := CombinationsInterface(list, length)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(list, res)
    if len(res) != 4 || len(res[0]) != length {
        t.Fatalf("combinationsInterface test failed")
    } else {
        t.Log("combinationsInterface test passed")
    }
}


func Test_GenerateProductString(t *testing.T) {
    var res [][]string

    list := []string{"a", "b", "c", "d"}
    length := 2

    c := GenerateProductString(list, length)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(list, res)
    if len(res) != 16 || len(res[0]) != length {
        t.Fatalf("GenerateCombinationsString test failed")
    } else {
        t.Log("GenerateCombinationsString test passed")
    }
}

func Test_GenerateProductInt(t *testing.T) {
    var res [][]int

    list := []int{1, 2, 3, 4}
    length := 2

    c := GenerateProductInt(list, length)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(list, res)
    if len(res) != 16 || len(res[0]) != length {
        t.Fatalf("GenerateCombinationsInt test failed")
    } else {
        t.Log("GenerateCombinationsInt test passed")
    }
}

// [[1 2 3 4] [5 6 7 8]] ==> [[1 5] [1 6] [1 7] [1 8] [2 5] [2 6] [2 7] [2 8] [3 5] [3 6] [3 7] [3 8] [4 5] [4 6] [4 7] [4 8]]
func Test_combinsSliceInterface(t *testing.T) {
    // var combin []interface{}
    // var data [][]interface{}
    var res [][]interface{}

    data := [][]interface{}{{1, 2, 3, 4}, {5, 6, 7, 8}}

    c := make(chan []interface{})
    go func(c chan []interface{}) {
        defer close(c)
        CombinsSliceInterface(c, []interface{}{}, data)
    }(c)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(data, res)
    if len(res) != 4 * 4 || len(res[0]) != len(data) {
        t.Fatalf("combinsSliceInterface test failed")
    } else {
        t.Log("combinsSliceInterface test passed")
    }
}


// [[1 2 3 4] [5 6 7]] ==> [[1 5] [1 6] [1 7] [2 5] [2 6] [2 7] [3 5] [3 6] [3 7] [4 5] [4 6] [4 7]]
func Test_combinsSliceInterface2(t *testing.T) {
    // var combin []interface{}
    // var data [][]interface{}
    var res [][]interface{}

    data := [][]interface{}{{1, 2, 3, 4}, {5, 6, 7}}

    c := make(chan []interface{})
    go func(c chan []interface{}) {
        defer close(c)
        CombinsSliceInterface(c, []interface{}{}, data)
    }(c)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(data, res)
    if len(res) != 4 * 3 || len(res[0]) != len(data) {
        t.Fatalf("combinsSliceInterface test failed")
    } else {
        t.Log("combinsSliceInterface test passed")
    }
}


func Test_combinsSliceInterface3(t *testing.T) {
    // var combin []interface{}
    // var data [][]interface{}
    var res [][]interface{}

    data := [][]interface{}{{1, 2, 3, 4}, {5, 6, 7}, {8, 9, 10}}

    c := make(chan []interface{})
    go func(c chan []interface{}) {
        defer close(c)
        CombinsSliceInterface(c, []interface{}{}, data)
    }(c)

    for value := range c {
        res = append(res, value)
    }

    fmt.Println(data, res)
    if len(res) != 4 * 3 * 3 || len(res[0]) != len(data) {
        t.Fatalf("combinsSliceInterface test failed")
    } else {
        t.Log("combinsSliceInterface test passed")
    }
}


