/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package statechart

import (
    "fmt"
    "os"
    // "strings"
    // "bufio"
    // "io"
    // "path/filepath"
    "encoding/json"

    "go4api/utils"
    "go4api/lib/testcase"
)

func InitFullScTcSlice (scfilePathSlice []string) []*testcase.TestCaseDataInfo {
    var fullScTcSlice []*testcase.TestCaseDataInfo
    // var fullKwJsPathSlice []string

    fmt.Println(scfilePathSlice)

    for i, _ := range scfilePathSlice {
        // scFileListTemp, _ := utils.WalkPath(scfilePathSlice[i], ".scxml")
        scFileListTemp, _ := utils.WalkPath(scfilePathSlice[i], ".xstate")

        for _, path := range scFileListTemp {
            // content := utils.GetContentFromFile(path)
            // XmlDecode(content) 

            ConstructXstate(path)
        }
    }

    return fullScTcSlice
}    

func ConstructXstate (xstateFile string) {
    var xstate State

    jsonStr := utils.GetJsonFromFile(xstateFile)

    err := json.Unmarshal([]byte(jsonStr), &xstate)
    if err != nil {
        fmt.Println("!! Error, parse xstate into xstate failed: ", xstateFile, ". Cause: ", err)
        os.Exit(1)
    }
        
    b, _ := json.Marshal(xstate)
    fmt.Println(string(b))

    xstate.SetStateIds()
    xstate.GetStateIds()

    b, _ = json.Marshal(xstate)
    fmt.Println(string(b))

    xstate.GetStateTransitions()

}





