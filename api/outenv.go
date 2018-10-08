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
    "fmt"
    // "syscall"
    "net/http"

    "go4api/lib/testcase"
)


func WriteOutEnvVariables (testResult string, tcData testcase.TestCaseDataInfo, actualStatusCode int, actualHeader http.Header, actualBody []byte) {
    // ----
    var expEnvs map[string]interface{}

    if testResult == "Success" {
        expEnvs = tcData.TestCase.OutEnvVariables()
  
        if len(expEnvs) > 0 {
            for k, v := range expEnvs {
                key := "go4_" + k
                value := GetResponseValue(v.(string), actualStatusCode, actualHeader, actualBody)

                fmt.Println("key, value: ", key, value)
                err := os.Setenv(key, fmt.Sprint(value))
                // syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())
                if err != nil {
                    panic(err) 
                }
            } 
        }
    } else {
        // fmt.Println("Warning: test execution failed, no OutEnvVariables set!")
    }
}





