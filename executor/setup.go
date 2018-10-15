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

    "go4api/lib/testcase"
    "go4api/lib/tree"
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


func RunGlobalSetup (ch chan int, baseUrl string, resultsDir string, resultsLogFile string, tcArray []testcase.TestCaseDataInfo) tree.TcTreeStats { 
    //-----
    prioritySet, root, tcTree, tcTreeStats := RunInit(tcArray)

    fmt.Println("\n====> Global SetUp test cases execution starts!") 

    RunPriorities(baseUrl, resultsDir, resultsLogFile, tcArray, prioritySet, root, tcTree, tcTreeStats)
    RunConsoleOverallReport(tcArray, root, tcTreeStats)

    return tcTreeStats
}
