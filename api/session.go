/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (    
    // "fmt"
    // "reflect"

    gsession "go4api/lib/session"
)


func (tcDataStore *TcDataStore) WriteSession (expTcSession map[string]interface{}) {
    var tcSession = make(map[string]interface{})
 
    tcData := tcDataStore.TcData

    // if current tc has no out session yet, then init it from parent
    tcSessionTemp := gsession.LookupTcSession(tcData.TcName())

    if len(tcSessionTemp) == 0 {
        parentTcSession := gsession.LookupTcSession(tcData.ParentTestCase())

        // copy the parentTcSession to tcSession
        for k, v := range parentTcSession {
            tcSession[k] = v
        }
    } else {
        for k, v := range tcSessionTemp {
            tcSession[k] = v
        }
    }
    
    //
    if expTcSession != nil {
        for k, v := range expTcSession {
            var value interface{}

            switch v.(type) {
            case string:
                value = tcDataStore.GetResponseValue(v.(string))
            case int, int64, float64:
                value = v
            }

            tcSession[k] = value
        } 
    }
    tcName := tcData.TcName()
    gsession.WriteTcSession(tcName, tcSession)
}


