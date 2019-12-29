/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gsql

import (
    "strings"
)


//
func renderValue (jsonStr string, feeder map[string]string) string {
    s := jsonStr
    
    for key, value := range feeder {
        k := "${" + key + "}"
        if k == s {
            s = strings.Replace(s, k, value, -1)
        }
    }

    return s
}

