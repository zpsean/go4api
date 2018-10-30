/*
 * go4api - a api testing tool written in Go
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

func BuiltinFunctionsMapping (key string) []interface{} {
    //
    FuncsMapping := map[string][]interface{} {
    	"NextInt": []interface{} {
            NextInt, 
            "",
        },
        "NextAlphaNumeric": []interface{} {
            NextAlphaNumeric, 
            "",
        },
        "CurrentTimeStampString": []interface{} {
            CurrentTimeStampString, 
            "ignoreParams",
        },
        "CurrentTimeStampMilliString": []interface{} {
            CurrentTimeStampMilliString, 
            "ignoreParams",
        },
        "CurrentTimeStampMicroString": []interface{} {
            CurrentTimeStampMicroString, 
            "ignoreParams",
        },
        "CurrentTimeStampNanoString": []interface{} {
            CurrentTimeStampNanoString, 
            "ignoreParams",
        },
        "CurrentTimeStampUnix": []interface{} {
            CurrentTimeStampUnix, 
            "ignoreParams",
        },
        "CurrentTimeStampUnixMilli": []interface{} {
            CurrentTimeStampUnixMilli, 
            "ignoreParams",
        },
        "CurrentTimeStampUnixMicro": []interface{} {
            CurrentTimeStampUnixMicro, 
            "ignoreParams",
        },
        "CurrentTimeStampUnixNano": []interface{} {
            CurrentTimeStampUnixNano, 
            "ignoreParams",
        },
    }

    return FuncsMapping[key]
}

func CallBuiltinFunc (funcName string, funcParams interface{}) interface{} {
    f := reflect.ValueOf(BuiltinFunctionsMapping(funcName)[0])

    var in []reflect.Value

    for _, param := range reflect.ValueOf(funcParams).Interface().([]interface{}) {
        in = append(in, reflect.ValueOf(param))
    }

    if BuiltinFunctionsMapping(funcName)[1] == "ignoreParams" {
    	ins := make([]reflect.Value, 0)
    	result := f.Call(ins)

    	return result[0].Interface()
    } else {
    	result := f.Call(in)

    	return result[0].Interface()
    }
}
