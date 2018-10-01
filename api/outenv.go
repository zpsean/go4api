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
    "os"
    // "fmt"
    // "syscall"
    // "strings"
    // "reflect"
    // "path/filepath"
    // "encoding/json"

    "go4api/lib/testcase"
    // "go4api/uti/ls"
)


func WriteOutEnvVariables (testResult string, tcData testcase.TestCaseDataInfo, actualBody []byte) {
    var expEnvs map[string]interface{}

    if testResult == "Success" {
        expEnvs = tcData.TestCase.OutEnvVariables()
        
        if len(expEnvs) > 0 {
            for k, v := range expEnvs {
                key := "go4_" + k
                value := GetActualValueByJsonPath(v.(string), actualBody)
                // fmt.Println("expEnvs: ", key, value.(string))

                err := os.Setenv(key, value.(string))
                // syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())
                if err != nil {
                    panic(err) 
                }

                // env := os.Getenv(key)
                // fmt.Println("----> env: ", env)
            } 
        }
    } else {
        // fmt.Println("Warning: test execution failed, no OutEnvVariables set!")
    }
}



