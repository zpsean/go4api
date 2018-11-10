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
    "sync"
)

var Gsession sync.Map

func InitNewSession () sync.Map {
    var session = sync.Map{}

    return session
}

func LookupParentSession (parentTcName string) map[string]interface{} {
    tcSession := make(map[string]interface{})

    result, ok := Gsession.Load(parentTcName)
    if ok {
        tcSession = result.(map[string]interface{})
    }

    return tcSession
}

func WriteTcSession (tcName string, tcSession map[string]interface{}) {
    Gsession.Store(tcName, tcSession)
}
