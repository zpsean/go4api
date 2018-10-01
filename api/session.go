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
    // "os"
    // "fmt"
    // "strings"
    // "reflect"
    // "path/filepath"
    // "encoding/json"

    "go4api/lib/testcase"
    "go4api/lib/session"
    // "go4api/uti/ls"
)


func WriteSession (testResult string, tcData testcase.TestCaseDataInfo, actualBody []byte) {
    var tcSessionDef = make(map[string]interface{})

    if testResult == "Success" {
        // get its parent session
        parentTcSession := gsession.LookupParentSession(tcData.ParentTestCase())

        tcName := tcData.TcName()
        gsession.Gsession[tcName] = parentTcSession
        // get its session def
        tcSessionDef = tcData.TestCase.Session()

        if len(tcSessionDef) > 0 {
            for k, v := range tcSessionDef {
                gsession.Gsession[tcName][k] = v
            } 
        }
    } else {
        // fmt.Println("Warning: test execution failed, no OutEnvVariables set!")
    }
}


