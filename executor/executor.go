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
    "os"
    "sync"
    "encoding/json"

    "go4api/cmd"
    "go4api/api"
    "go4api/lib/testcase"
    "go4api/lib/tree"
    "go4api/reports"
)


func ScheduleCases (cReady chan *tree.TcNode, wg *sync.WaitGroup, resultsChan chan testcase.TestCaseExecutionInfo, baseUrl string) {
    //
    if cmd.Opt.IfConcurrency == true {
        tick := 0
        max := cmd.Opt.ConcurrencyLimit

        for tcNode := range cReady {
            wg.Add(1)
            // Note: to prevent reaching tcp connection limitation, here set a max, then sleep for a while
            if tick % max == 0 {
                time.Sleep(100 * time.Millisecond)
                go api.DispatchApi(wg, resultsChan, baseUrl, *(tcNode.TestCaseExecutionInfo.TestCaseDataInfo))
            } else {
                go api.DispatchApi(wg, resultsChan, baseUrl, *(tcNode.TestCaseExecutionInfo.TestCaseDataInfo))
            }

            tick = tick + 1
        }
    } else {
        for tcNode := range cReady {
            wg.Add(1)
            api.DispatchApi(wg, resultsChan, baseUrl, *(tcNode.TestCaseExecutionInfo.TestCaseDataInfo))
        }
    }   
}

func WriteNotNotExecutedToLog (priority string, logFilePtr *os.File, tcTreeStats tree.TcTreeStats) {
    notRunTime := time.Now()
    for i, _ := range tcTreeStats.TcNotExecutedByPriority[priority] {
        for _, tcExecution := range tcTreeStats.TcNotExecutedByPriority[priority][i] {
            // [casename, priority, parentTestCase, ...], tc, jsonFile, csvFile, row in csv
            if tcExecution.Priority() == priority {
                // set some dummy time for the tc not executed
                tcExecution.StartTimeUnixNano = notRunTime.UnixNano()
                tcExecution.EndTimeUnixNano = notRunTime.UnixNano()
                tcExecution.DurationUnixNano = notRunTime.UnixNano() - notRunTime.UnixNano()

                tcReportResults := tcExecution.TcReportResults()
                
                repJson, _ := json.Marshal(tcReportResults)
                //
                reports.WriteExecutionResults(string(repJson), logFilePtr)
            }
        }
    }
}



