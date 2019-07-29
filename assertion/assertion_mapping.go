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
    
)

type AssertionDetails struct {
    AssertionKey string
    AssertionFunc interface{}  // function
    ApplyTypes []string        // [string, number, bool]
}


var assertionMapping = make(map[string]*AssertionDetails)

func init() {
    assertionMapping["Equals"] = &AssertionDetails{"Equals", Equals, []string{"string", "int", "float64", "bool", "slice"}}
    assertionMapping["NotEquals"] = &AssertionDetails{"NotEquals", NotEquals, []string{"string", "int", "float64", "bool", "slice"}}

    assertionMapping["Contains"] = &AssertionDetails{"Contains", Contains, []string{"string"}}
    assertionMapping["StartsWith"] = &AssertionDetails{"StartsWith", StartsWith, []string{"string"}}
    assertionMapping["EndsWith"] = &AssertionDetails{"EndsWith", EndsWith, []string{"string"}}

    assertionMapping["Less"] = &AssertionDetails{"Less", Less, []string{"float64"}}
    assertionMapping["LessOrEquals"] = &AssertionDetails{"LessOrEquals", LessOrEquals, []string{"float64"}}
    assertionMapping["Greater"] = &AssertionDetails{"Greater", Greater, []string{"float64"}}
    assertionMapping["GreaterOrEquals"] = &AssertionDetails{"GreaterOrEquals", GreaterOrEquals, []string{"float64"}}

    assertionMapping["In"] = &AssertionDetails{"In", In, []string{"string", "int", "float64", "bool", "slice"}}
    assertionMapping["NotIn"] = &AssertionDetails{"NotIn", NotIn, []string{"string", "int", "float64", "bool", "slice"}}
    assertionMapping["Has"] = &AssertionDetails{"Has", Has, []string{"string", "int", "float64", "bool", "slice"}}
    assertionMapping["NotHas"] = &AssertionDetails{"NotHas", NotHas, []string{"string", "int", "float64", "bool", "slice"}}

    assertionMapping["Match"] = &AssertionDetails{"Match", Match, []string{"string"}}
}


