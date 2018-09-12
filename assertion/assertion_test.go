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
    "testing"
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
    actualValue := true
    expValue := true
    res := Equals(actualValue, expValue)

    if res == false {
        t.Fatalf("Equals Bool failed")
    }  
}

func Test_Equals6(t *testing.T) {
    actualValue := true
    expValue := false
    res := Equals(actualValue, expValue)

    if res == true {
        t.Fatalf("Equals Bool failed")
    }
}

// func Test_Equals7(t *testing.T) {
//     aa := `
//         {
//             "a": 123,
//             "b": null
//         }`
//     var BB map[string]interface{}

//     json.Unmarshal([]byte(aa), &BB)

//     actualValue := nil // not work
//     expValue := BB.b

//     res := Equals(actualValue, expValue)

//     if res == true {
//         t.Fatalf("Equals null failed")
//     }
// }


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


