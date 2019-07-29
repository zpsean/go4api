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
    // "encoding/csv"
)

func Test_Equals(t *testing.T) {
    actualValue := "abc"
    expValue := "abc"
    res := Equals(actualValue, expValue)

    if res == false {
        t.Fatalf("Equals String failed")
    } else {
        t.Log("Equals String test passed")
    }
}

func Test_Equals2(t *testing.T) {
    actualValue := "abc"
    expValue := "abcd"
    res := Equals(actualValue, expValue)

    if res == true {
        t.Fatalf("Equals String failed")
    } else {
        t.Log("Equals String test passed")
    }
}

func Test_Equals3(t *testing.T) {
    actualValue := 123
    expValue := 123
    res := Equals(actualValue, expValue)

    if res == false {
        t.Fatalf("Equals float64 failed")
    } else {
        t.Log("Equals float64 test passed")
    } 
}

func Test_Equals4(t *testing.T) {
    actualValue := 123
    expValue := 1234
    res := Equals(actualValue, expValue)

    if res == true {
        t.Fatalf("Equals float64 failed")
    } else {
        t.Log("Equals float64 test passed")
    } 
}

func Test_Equals5(t *testing.T) {
    var actualValue int
    actualValue = 123
    var expValue float64
    expValue = 123
    res := Equals(actualValue, expValue)

    if res != true {
        t.Fatalf("Equals float64 failed")
    } else {
        t.Log("Equals float64 test passed")
    } 
}

func Test_Equals6(t *testing.T) {
    actualValue := true
    expValue := true
    res := Equals(actualValue, expValue)

    if res == false {
        t.Fatalf("Equals Bool failed")
    }  
}

func Test_Equals7(t *testing.T) {
    actualValue := true
    expValue := false
    res := Equals(actualValue, expValue)

    if res == true {
        t.Fatalf("Equals Bool failed")
    }
}

func Test_Equals8(t *testing.T) {
    actualValue := []int{1,2,3}
    expValue := []int{1,2,3}

    res := Equals(actualValue, expValue)

    if res != true {
        t.Fatalf("Equals Slice failed")
    }
}

func Test_Equals9(t *testing.T) {
    actualValue := []int{1,2,33}
    expValue := []int{1,2,3}

    res := Equals(actualValue, expValue)

    if res != false {
        t.Fatalf("Equals Slice failed")
    }
}

func Test_Equals10(t *testing.T) {
    actualValue := []interface{}{1,2,3}
    expValue := []interface{}{"1","2","3"}

    // res := Equals(actualValue, expValue)
    res := reflect.DeepEqual(actualValue, expValue)
    fmt.Println("actualValue, expValue: ", res, reflect.TypeOf(res), actualValue, expValue)

    if res != false {
        t.Fatalf("Equals Slice failed")
    }
}

func Test_Equals11(t *testing.T) {
    actualValue := []interface{}{"1","2","3"}
    expValue := []interface{}{"1","2","3"}

    // res := Equals(actualValue, expValue)
    res := reflect.DeepEqual(actualValue, expValue)
    fmt.Println("actualValue, expValue: ", res, reflect.TypeOf(res), actualValue, expValue)

    if res != true {
        t.Fatalf("Equals Slice failed")
    }
}


// <------->
func Test_Contains(t *testing.T) {
    actualValue := "abcde"
    expValue := "abc"
    res := Contains(actualValue, expValue)

    if res == false {
        t.Fatalf("Contains failed")
    }  
}

func Test_Contains2(t *testing.T) {
    actualValue := "abcde"
    expValue := "abcf"
    res := Contains(actualValue, expValue)

    if res == true {
        t.Fatalf("Contains failed")
    }
}

func Test_StartsWith(t *testing.T) {
    actualValue := "abcde"
    expValue := "abc"
    res := StartsWith(actualValue, expValue)

    if res == false {
        t.Fatalf("StartsWith failed")
    }  
}

func Test_StartsWith2(t *testing.T) {
    actualValue := "abcde"
    expValue := "bcd"
    res := StartsWith(actualValue, expValue)

    if res == true {
        t.Fatalf("StartsWith failed")
    }
}

func Test_EndsWith(t *testing.T) {
    actualValue := "abcde"
    expValue := "de"
    res := EndsWith(actualValue, expValue)

    if res == false {
        t.Fatalf("EndsWith failed")
    }  
}

func Test_EndsWith2(t *testing.T) {
    actualValue := "abcde"
    expValue := "cd"
    res := EndsWith(actualValue, expValue)

    if res == true {
        t.Fatalf("EndsWith failed")
    }
}



func Test_GreaterOrEquals(t *testing.T) {
    actualValue := 1234.4
    expValue := 1234.0
    res := GreaterOrEquals(actualValue, expValue)

    if res == false {
        t.Fatalf("GreaterOrEquals failed")
    }
}

func Test_GreaterOrEquals2(t *testing.T) {
    actualValue := 1234.0
    expValue := 1234.4
    res := GreaterOrEquals(actualValue, expValue)

    if res == true {
        t.Fatalf("GreaterOrEquals failed")
    }
}


func Test_GreaterOrEquals3(t *testing.T) {
     var actualValue int
    actualValue = 1234.0
    var expValue float64
    expValue = 1234.4

    res := GreaterOrEquals(actualValue, expValue)

    if res == true {
        t.Fatalf("GreaterOrEquals failed")
    }
}

func Test_Greater(t *testing.T) {
     var actualValue int
    actualValue = 1234.0
    var expValue float64
    expValue = 1234.4

    res := Greater(actualValue, expValue)

    if res == true {
        t.Fatalf("Greater failed")
    }
}

// item in slice
func Test_In(t *testing.T) {
    var actualValue int
    actualValue = 1

    expValue := []interface{}{1, 2, 3}

    res := In(actualValue, expValue)

    if res != true {
        t.Fatalf("In failed")
    }
}

func Test_In2(t *testing.T) {
    var actualValue int
    actualValue = 11

    expValue := []interface{}{1, 2, 3}

    res := In(actualValue, expValue)

    if res != false {
        t.Fatalf("In failed")
    }
}

func Test_NotIn(t *testing.T) {
    var actualValue int
    actualValue = 1

    expValue := []interface{}{1, 2, 3}

    res := NotIn(actualValue, expValue)

    if res != false {
        t.Fatalf("NotIn failed")
    }
}

func Test_NotIn2(t *testing.T) {
    var actualValue int
    actualValue = 11

    expValue := []interface{}{1, 2, 3}

    res := NotIn(actualValue, expValue)

    if res != true {
        t.Fatalf("NotIn failed")
    }
}


func Test_Has(t *testing.T) {
    actualValue := []interface{}{1, 2, 3}

    var expValue int
    expValue = 1

    res := Has(actualValue, expValue)

    if res != true {
        t.Fatalf("In failed")
    }
}

func Test_Has2(t *testing.T) {
    actualValue := []interface{}{1, 2, 3}

    var expValue int
    expValue = 11

    res := Has(actualValue, expValue)

    if res != false {
        t.Fatalf("In failed")
    }
}

func Test_NotHas(t *testing.T) {
    actualValue := []interface{}{1, 2, 3}

    var expValue int
    expValue = 1

    res := NotHas(actualValue, expValue)

    if res != false {
        t.Fatalf("In failed")
    }
}

func Test_NotHas2(t *testing.T) {
    actualValue := []interface{}{1, 2, 3}

    var expValue int
    expValue = 11

    res := NotHas(actualValue, expValue)

    if res != true {
        t.Fatalf("In failed")
    }
}

// match
func Test_Match(t *testing.T) {
    actualValue := "abcde"
    expPattern := `[a-z]+`
    res := Match(actualValue, expPattern)

    if res == false {
        t.Fatalf("Match failed")
    }  
}

func Test_Match2(t *testing.T) {
    actualValue := "abcde"
    expPattern := `[0-9]+`
    res := Match(actualValue, expPattern)

    if res == true {
        t.Fatalf("Match failed")
    }
}


