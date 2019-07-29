/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
    // "fmt"
    "testing"
    // "reflect"
    // "encoding/csv"
)

func Test_compareCommon(t *testing.T) {
    actualValue := "123"
    expValue := "123"
    res, _ := compareCommon("", "", "Equals", actualValue, expValue)

    if res != true {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}

func Test_compareCommon2(t *testing.T) {
    actualValue := "123"
    expValue := 123
    res, _ := compareCommon("", "", "Equals", actualValue, expValue)

    if res != false {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}


func Test_compareCommon3(t *testing.T) {
    actualValue := []interface{}{1,2,3}
    expValue := []interface{}{"1","2","3"}

    res, _ := compareCommon("", "", "Equals", actualValue, expValue)

    if res != false {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}

func Test_compareCommon4(t *testing.T) {
    actualValue := []interface{}{"1","2","3"}
    expValue := []interface{}{"1","2","3"}

    res, _ := compareCommon("", "", "Equals", actualValue, expValue)

    if res != true {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}

func Test_compareCommon5(t *testing.T) {
    actualValue := []interface{}{"1","2","3"}
    expValue := []interface{}{"1","2","3"}

    res, _ := compareCommon("", "", "NotEquals", actualValue, expValue)

    if res != false {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}


func Test_compareCommon6(t *testing.T) {
    actualValue := []interface{}{"1","2","3"}
    expValue := []interface{}{"1","2","3"}

    res, _ := compareCommon("", "", "NNNotEquals", actualValue, expValue)

    if res != false {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}

// in, notin, has, nothas
func Test_compareCommon7(t *testing.T) {
    var actualValue int
    actualValue = 1

    expValue := []interface{}{1, 2, 3}

    res, _ := compareCommon("", "", "In", actualValue, expValue)

    if res != true {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}

func Test_compareCommon8(t *testing.T) {
    var actualValue string
    actualValue = "1"

    expValue := []interface{}{1, 2, 3}

    res, _ := compareCommon("", "", "In", actualValue, expValue)

    if res != false {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}

func Test_compareCommon9(t *testing.T) {
    actualValue := []interface{}{1, 2, 3}

    var expValue int
    expValue = 1

    res, _ := compareCommon("", "", "Has", actualValue, expValue)

    if res != true {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}


func Test_compareCommon10(t *testing.T) {
    actualValue := []interface{}{1, 2, 3}

    var expValue string
    expValue = "1"

    res, _ := compareCommon("", "", "Has", actualValue, expValue)

    if res != false {
        t.Fatalf("compareCommon test failed")
    } else {
        t.Log("compareCommon test passed")
    }
}
