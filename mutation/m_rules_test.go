/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package mutation

import (
    "fmt"
    "testing"
    "reflect"
    "encoding/json"
)

var (
    currArrayInt []interface{}
    currArrayString []interface{}
)


func init() {
    currArrayInt = []interface{}{1, 2, 3, 4}
    currArrayString = []interface{}{"a", "b", "c", "d"}
}

func Test_MArrayR1(t *testing.T) {
    mutatedValue := MArrayR1(currArrayInt, "", "")

    fmt.Println("------->", currArrayInt, mutatedValue)

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR1 failed")
    }
}

func Test_MArrayR2(t *testing.T) {
    mutatedValue := MArrayR2(currArrayInt, "", "")

    fmt.Println("------->", currArrayInt, mutatedValue)

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR2 failed")
    }
}

func Test_MArrayR3(t *testing.T) {
    mutatedValue := MArrayR3(currArrayInt, "", "")

    fmt.Println("------->", currArrayInt, mutatedValue)

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR3 failed")
    }
}

func Test_MArrayR4(t *testing.T) {
    mutatedValue := MArrayR4(currArrayInt, "", "")

    fmt.Println("------->", currArrayInt, mutatedValue)

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR4 failed")
    }
}

func Test_MArrayR5(t *testing.T) {
    mutatedValue := MArrayR5(currArrayInt, "", "")

    fmt.Println("------->", currArrayInt, mutatedValue)

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR5 failed")
    }
}

func Test_MArrayR5_2(t *testing.T) {
    mutatedValue := MArrayR5(currArrayString, "", "")

    fmt.Println("------->", currArrayString, mutatedValue)

    if reflect.DeepEqual(currArrayString, mutatedValue) != false {
        t.Fatalf("Error, currArrayInt failed")
    }
}

func Test_MArrayR6(t *testing.T) {
    mutatedValue := MArrayR6(currArrayInt, "", "")

    fmt.Println("------->", currArrayInt, mutatedValue)

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR6 failed")
    }
}

func Test_MArrayR6_2(t *testing.T) {
    mutatedValue := MArrayR6(currArrayString, "", "")

    fmt.Println("------->", currArrayString, mutatedValue)

    if reflect.DeepEqual(currArrayString, mutatedValue) != false {
        t.Fatalf("Error, MArrayR6 failed")
    }
}

func Test_MArrayR7(t *testing.T) {
    mutatedValue := MArrayR7(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR8(t *testing.T) {
    mutatedValue := MArrayR8(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR9(t *testing.T) {
    mutatedValue := MArrayR9(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR10(t *testing.T) {
    mutatedValue := MArrayR10(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR11(t *testing.T) {
    mutatedValue := MArrayR11(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR12(t *testing.T) {
    mutatedValue := MArrayR12(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR13(t *testing.T) {
    mutatedValue := MArrayR13(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR14(t *testing.T) {
    mutatedValue := MArrayR14(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR15(t *testing.T) {
    mutatedValue := MArrayR15(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR16(t *testing.T) {
    mutatedValue := MArrayR16(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR17(t *testing.T) {
    mutatedValue := MArrayR17(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR18(t *testing.T) {
    mutatedValue := MArrayR18(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}

func Test_MArrayR19(t *testing.T) {
    mutatedValue := MArrayR19(currArrayInt, "", "")

    r7Json, _ := json.Marshal(mutatedValue)

    fmt.Println("------->", currArrayInt, mutatedValue, string(r7Json))
    // [1 2 3 4] [1 <nil> 3 4] [1,null,3,4]

    if reflect.DeepEqual(currArrayInt, mutatedValue) != false {
        t.Fatalf("Error, MArrayR failed")
    }
}
