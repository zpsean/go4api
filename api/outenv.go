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
)


func (tcDataStore *TcDataStore) WriteOutEnvVariables (testResult string) {
    // ----
    var expEnvs map[string]interface{}

    tcData := tcDataStore.TcData

    actualStatusCode := tcDataStore.HttpActualStatusCode
    actualHeader := tcDataStore.HttpActualHeader
    actualBody := tcDataStore.HttpActualBody

    if testResult == "Success" {
        expEnvs = tcData.TestCase.OutEnvVariables()
  
        if len(expEnvs) > 0 {
            for k, v := range expEnvs {
                key := "go4_" + k
                value := GetResponseValue(v.(string), actualStatusCode, actualHeader, actualBody)

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





