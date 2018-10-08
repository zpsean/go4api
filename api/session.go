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
    "net/http"

    "go4api/lib/testcase"
    "go4api/lib/session"
)


func WriteSession (testResult string, tcData testcase.TestCaseDataInfo, actualStatusCode int, actualHeader http.Header, actualBody []byte) {
    var tcSession = make(map[string]interface{})
    var tcSessionDef = make(map[string]interface{})

    if testResult == "Success" {
        // get its parent session
        parentTcSession := gsession.LookupParentSession(tcData.ParentTestCase())
        tcSession = parentTcSession

        // get its session def
        tcSessionDef = tcData.TestCase.Session()

        if len(tcSessionDef) > 0 {
            for k, v := range tcSessionDef {
                value := GetResponseValue(v.(string), actualStatusCode, actualHeader, actualBody)

                tcSession[k] = value
            } 
        }
        tcName := tcData.TcName()
        gsession.WriteTcSession(tcName, parentTcSession)

    } else {
        // fmt.Println("Warning: test execution failed, no OutEnvVariables set!")
    }
}


