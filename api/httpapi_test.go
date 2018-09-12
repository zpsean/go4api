/*
 * go4api - a api testing tool written in Go
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

