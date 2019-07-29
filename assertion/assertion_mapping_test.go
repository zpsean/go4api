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
    // "encoding/csv"
)

func Test_init(t *testing.T) {
    for _, value := range assertionMapping {
        fmt.Println("AssertionMapping: ", value)
    }
    
    if len(assertionMapping) != 14 {
        t.Fatalf("init failed")
    } else {
        t.Log("init test passed")
    }
}
