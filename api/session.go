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
    gsession "go4api/lib/session"
)


func (tcDataStore *TcDataStore) WriteSession (expTcSession map[string]interface{}, rowsCount int, rowsData interface{}) {
    var tcSession = make(map[string]interface{})
 
    tcData := tcDataStore.TcData

    // get its parent session
    parentTcSession := gsession.LookupTcSession(tcData.ParentTestCase())
    tcSession = parentTcSession

    if expTcSession != nil {
        for k, v := range expTcSession {
            value := tcDataStore.GetResponseValue(v.(string), rowsCount, rowsData)

            tcSession[k] = value
        } 
    }
    tcName := tcData.TcName()
    gsession.WriteTcSession(tcName, tcSession)

}


