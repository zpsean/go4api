/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package builtins

import (
    "fmt"
    "testing"
    // "encoding/csv"
)

func Test_CurrentTimeStampString (t *testing.T) {
    res := CurrentTimeStampString("milli")

    fmt.Println("res: ", res)
    // if res != 4 {
    //     t.Fatalf("getMaxLenVector test failed")
    // } else {
    //     t.Log("getMaxLenVector test passed")
    // }
}

func Test_CurrentTimeStampUnix (t *testing.T) {
    res := CurrentTimeStampUnix("milli")

    fmt.Println("res: ", res)
    // if res != 4 {
    //     t.Fatalf("getMaxLenVector test failed")
    // } else {
    //     t.Log("getMaxLenVector test passed")
    // }
}

func Test_DayStart (t *testing.T) {
    res := DayStart(CurrentTimeStampUnix("milli"))

    fmt.Println("res: ", res)
}

func Test_DayEnd (t *testing.T) {
    res := DayEnd(CurrentTimeStampUnix("milli"))

    fmt.Println("res: ", res)
}

func Test_TimeStampUnixOffset (t *testing.T) {
    var oset = []interface{}{1543593599000 , "-1", "day"}

    res := TimeStampUnixOffset(oset)

    fmt.Println("res: ", res)
}

