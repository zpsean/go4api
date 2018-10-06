/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package executor

import (
    // "fmt"
    "time"
    // "os"
    // "sort"
    // "sync"
    // "encoding/json"

    // "go4api/utils"
    "go4api/lib/testcase"
    // "go4api/reports"
)

func GetTeardownTcSlice (tcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var teardownTcSlice []testcase.TestCaseDataInfo
    for i, _ := range tcArray {
        if tcArray[i].TestCase.IfGlobalTearDownTestCase() == true {
            teardownTcSlice = append(teardownTcSlice, tcArray[i])
        }
    }
    
    return teardownTcSlice
}


func RunTeardown(ch chan int, pStart_time time.Time, pStart string, baseUrl string, resultsDir string, tcArray []testcase.TestCaseDataInfo) { 
    prioritySet, root, tcTree := RunBefore(tcArray)

    RunPriorities(ch, pStart, baseUrl, resultsDir, tcArray, prioritySet, root, tcTree)

    RunAfter(ch, pStart_time, pStart, resultsDir, tcArray, root, tcTree)
}

