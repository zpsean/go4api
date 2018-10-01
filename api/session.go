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
    // "strings"
    // "reflect"
    // "path/filepath"
    // "encoding/json"

    "go4api/lib/testcase"
    // "go4api/uti/ls"
)


func WriteSession (testResult string, tcData testcase.TestCaseDataInfo, actualBody []byte) {
    var expEnvs []map[string]interface{}

    if testResult == "Success" {
        expEnvs = tcData.TestCase.OutEnvVariables()
        if len(expEnvs) > 0 {
            // 
            for i, _ := range expEnvs {
                for k, v := range expEnvs[i] {
                    key := "go4_" + k
                    value := GetActualValueByJsonPath(v.(string), actualBody)

                    os.Setenv(key, value.(string))
                } 
            }
        }
    } else {
        // fmt.Println("Warning: test execution failed, no OutEnvVariables set!")
    }
}

func RetriveSession (testResult string, tcData testcase.TestCaseDataInfo, actualBody []byte) {
    var expEnvs []map[string]interface{}

    if testResult == "Success" {
        expEnvs = tcData.TestCase.OutEnvVariables()
        if len(expEnvs) > 0 {
            // 
            for i, _ := range expEnvs {
                for k, v := range expEnvs[i] {
                    key := "go4_" + k
                    value := GetActualValueByJsonPath(v.(string), actualBody)

                    os.Setenv(key, value.(string))
                } 
            }
        }
    } else {
        // fmt.Println("Warning: test execution failed, no OutEnvVariables set!")
    }
}

func MergeSessionWithTcJson (testResult string, tcData testcase.TestCaseDataInfo, actualBody []byte) {
    var expEnvs []map[string]interface{}

    if testResult == "Success" {
        expEnvs = tcData.TestCase.OutEnvVariables()
        if len(expEnvs) > 0 {
            // 
            for i, _ := range expEnvs {
                for k, v := range expEnvs[i] {
                    key := "go4_" + k
                    value := GetActualValueByJsonPath(v.(string), actualBody)

                    os.Setenv(key, value.(string))
                } 
            }
        }
    } else {
        // fmt.Println("Warning: test execution failed, no OutEnvVariables set!")
    }
}

