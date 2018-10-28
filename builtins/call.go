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
	// "math/rand"                                                                                                                                        
	// "time"
 	// "strings"
 	"reflect"
)

func CallBuiltinFunc (funcName string, funcParams ... interface{}) interface{} {
    f := reflect.ValueOf(funcName)

    in := make([]reflect.Value, len(funcParams))
    for k, param := range funcParams {
        in[k] = reflect.ValueOf(param)
    }

    result := f.Call(in)

    return result[0]
}
