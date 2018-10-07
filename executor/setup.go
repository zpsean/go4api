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
    "fmt"
    "time"
    // "os"
    // "sort"
    // "sync"
    // "encoding/json"

    // "go4api/utils"
    "go4api/lib/testcase"
    // "go4api/reports"
)

func GetSetupTcSlice (tcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var setUpTcSlice []testcase.TestCaseDataInfo
    for i, _ := range tcArray {
        if tcArray[i].TestCase.IfGlobalSetUpTestCase() == true {
            setUpTcSlice = append(setUpTcSlice, tcArray[i])
        }
    }

    return setUpTcSlice
}

func RunSetup(ch chan int, pStart_time time.Time, pStart string, baseUrl string, resultsDir string, tcArray []testcase.TestCaseDataInfo) { 
    prioritySet, root, tcTree, tcTreeStats := RunBefore(tcArray)

    fmt.Println("\n====> setup test cases execution starts!") 

    RunPriorities(ch, pStart, baseUrl, resultsDir, tcArray, prioritySet, root, tcTree, tcTreeStats)

    RunAfter(ch, pStart_time, pStart, resultsDir, tcArray, root, tcTree, tcTreeStats)
}

