/*
 * go4api - a api testing tool written in Go
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

    gsession "go4api/lib/session"
)


func (tcDataStore *TcDataStore) WriteSession (expTcSession map[string]interface{}) {
    var tcSession = make(map[string]interface{})
 
    tcData := tcDataStore.TcData

    // get its parent session
    parentTcSession := gsession.LookupTcSession(tcData.ParentTestCase())
    tcSession = parentTcSession

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


