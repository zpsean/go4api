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

func Test_VerifyTypes(t *testing.T) {
    actualValue := 1234.4
    expValue := 1234
    res := VerifyTypes(actualValue, expValue)

    if res == "false" {
        t.Fatalf("Test_VerifyTypes failed")
    } else {
        t.Log("Test_VerifyTypes passed")
    }
}

func Test_VerifyTypes2(t *testing.T) {
    actualValue := 1234
    expValue := "1234"
    res := VerifyTypes(actualValue, expValue)

    if res == "true" {
        t.Fatalf("Test_VerifyTypes failed")
    } else {
        t.Log("Test_VerifyTypes passed")
    }
}

func Test_VerifyTypes3(t *testing.T) {
    actualValue := true
    expValue := ""
    res := VerifyTypes(actualValue, expValue)

    if res == "true" {
        t.Fatalf("Test_VerifyTypes failed")
    } else {
        t.Log("Test_VerifyTypes passed")
    }
}

