/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gsession

import (
    "sync"
)

var GlobalVariables sync.Map
var Session         sync.Map

func LookupTcSession (tcName string) map[string]interface{} {
    tcSession := make(map[string]interface{})

    result, ok := Session.Load(tcName)
    if ok {
        tcSession = result.(map[string]interface{})
    }

    return tcSession
}

func WriteTcSession (tcName string, tcSession map[string]interface{}) {
    Session.Store(tcName, tcSession)
}

// global
func LoopGlobalVariables () map[string]interface{} {
    var resMap = make(map[string]interface{})

    GlobalVariables.Range(func(key, value interface{}) bool {
        resMap[key.(string)] = value
        return true
    })

    return resMap
}

func LookupGlobalVariables (key string) interface{} {
    var value interface{}

    result, ok := GlobalVariables.Load(key)
    if ok {
        value = result
    }

    return value
}

func WriteGlobalVariables (key string, value interface{}) {
    GlobalVariables.Store(key, value)
}

