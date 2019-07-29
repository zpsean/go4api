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
    // "fmt"
	// "math/rand"                                                                                                                                        
	// "time"
 	// "strings"
 	"reflect"
)

func BuiltinFunctionsMapping (key string) interface{} {
    //
    FuncsMapping := map[string]interface{} {
    	"NextInt": NextInt,
        "NextAlphaNumeric": NextAlphaNumeric,
        "NextStringNumeric": NextStringNumeric,
        "Substitute": Substitute,
        "Select": Select,
        "Join": Join,
        "Split": Split,
        "Replace": Replace,
        "SubString": SubString,
        "ToString": ToString,
        "Length": Length,
        "ToInt": ToInt,
        "ToBool": ToBool,
        "CurrentTimeStampString": CurrentTimeStampString,
        "CurrentTimeStampUnix": CurrentTimeStampUnix,
        "DayStart": DayStart,
        "DayEnd": DayEnd,
        "ConvertTimeToUnix": ConvertTimeToUnix,
        "ConvertTimeToStr": ConvertTimeToStr,
        "TimeStampUnixOffset": TimeStampUnixOffset,
    }

    return FuncsMapping[key]
}

func CallBuiltinFunc (funcName string, funcParams interface{}) interface{} {
    f := reflect.ValueOf(BuiltinFunctionsMapping(funcName))

    var in []reflect.Value
    in = append(in, reflect.ValueOf(funcParams))

    result := f.Call(in)

    return result[0].Interface()
}
