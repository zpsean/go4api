/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gsession

import (
    // "os"
    // "fmt"
    // "strings"
    // "reflect"
    // "encoding/csv"
    // "encoding/json"
)

var Gsession map[string]map[string]interface{}

func init () {
    Gsession = make(map[string]map[string]interface{})
}


func LookupParentSession (parentTcName string) map[string]interface{} {
    var tcSession map[string]interface{}

    for k, v := range Gsession {
        if parentTcName == k {
            tcSession = v
        }
    } 

    return tcSession
}