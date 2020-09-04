/*
 * go4api - an api testing tool written in Go
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
    "os"
    "sync"
    "encoding/json"

    "go4api/cmd"
    "go4api/api"
    "go4api/lib/testcase"
    "go4api/lib/tree"
    "go4api/reports"
)

func (tcsRunStore *TcsRunStore) InitRun () { 
    tcTree := tree.CreateTcTree()
    root   := tcTree.BuildTree(tcsRunStore.TcSlice)
    //
    prioritySet := GetPrioritySet(tcsRunStore.TcSlice)
    tcTreeStats := tree.CreateTcTreeStats(prioritySet)
    // 
    tcTree.InitNodesRunResult(root, "Ready")

    // aa, _ := json.Marshal(tcTree)
    // fmt.Println(string(aa))

    tcsRunStore.PrioritySet = prioritySet
    tcsRunStore.Root        = root
    tcsRunStore.TcTree      = tcTree
    tcsRunStore.TcTreeStats = tcTreeStats
}

func (tcsRunStore *TcsRunStore) RunPriorities (baseUrl string, resultsLogFile string) {
    logFilePtr := reports.OpenExecutionResultsLogFile(resultsLogFile)

    for _, priority := range tcsRunStore.PrioritySet {
        fmt.Println("====> Priority " + priority + " starts!")
        //
        tcsRunStore.RunEachPriority(baseUrl, priority, logFilePtr)

        // Put out the cases which has not been executed (i.e. not Success or Fail)
        WriteNotNotExecutedToLog(priority, logFilePtr, tcsRunStore.TcTreeStats)

        // report to console
        reports.ReportConsoleByPriority(tcsRunStore.TcTreeStats.StatusCountByPriority[priority]["Total"], priority, tcsRunStore.TcTreeStats.StatusCountByPriority)

        fmt.Println("====> Priority " + priority + " ended!")
        fmt.Println("")
        // sleep for debug
        // time.Sleep(500 * time.Millisecond)
    }

    logFilePtr.Close()
}

func (tcsRunStore *TcsRunStore) RunEachPriority (baseUrl string, priority string, logFilePtr *os.File) {
    miniLoop:
    for {
        resultsExeChan := make(chan testcase.TestCaseExecutionInfo, len(tcsRunStore.TcSlice))
        var wg sync.WaitGroup
        //
        cReady := make(chan *tree.TcNode)
        go func(cReady chan *tree.TcNode) {
            defer close(cReady)
            tcsRunStore.TcTree.CollectNodeReadyByPriority(cReady, tcsRunStore.Root, priority)
        }(cReady)

        ScheduleCases(cReady, &wg, resultsExeChan, baseUrl)
        //
        wg.Wait()

        close(resultsExeChan)
        
        for tcExecution := range resultsExeChan {
            //
            c := make(chan *tree.TcNode)
            go func(c chan *tree.TcNode) {
                defer close(c)
                tcsRunStore.TcTree.SearchNode(c, tcsRunStore.Root, tcExecution.TcName())
            }(c)
            //
            // tcsRunStore.TcTreeStats.ResetTcTreeStats(priority)
            // tcsRunStore.TcTreeStats.CollectNodeStatusByPriority(tcsRunStore.Root, priority)

            // tcsRunStore.TcTree.RefreshNodeAndDirectChilrenTcResult(<-c, tcExecution.TestResult, tcExecution.StartTime, tcExecution.EndTime, 
            //     tcExecution.HttpTestMessages, tcExecution.StartTimeUnixNano, tcExecution.EndTimeUnixNano)

            tcsRunStore.TcTree.RefreshNodeAndChilrenTcResult(<-c, tcExecution.TestResult, tcExecution.StartTime, tcExecution.EndTime, 
                tcExecution.HttpTestMessages, tcExecution.StartTimeUnixNano, tcExecution.EndTimeUnixNano)

            tcsRunStore.TcTreeStats.ResetTcTreeStats(priority)
            tcsRunStore.TcTreeStats.CollectNodeStatusByPriority(tcsRunStore.Root, priority)

            // console 1. to print status to console (i.e. executed cases: "Success", "Fail")
            tcReportResults := tcExecution.TcReportResults()

            repJson, _ := json.Marshal(tcReportResults)

            tcsRunStore.TcDs = append(tcsRunStore.TcDs, tcExecution.TestCaseDataInfo)

            reports.WriteExecutionResults(string(repJson), logFilePtr)
            reports.ReportConsoleByTc(tcExecution)
        }
        // if tcTree has no node with "Ready" status, break the miniloop
        if tcsRunStore.TcTreeStats.StatusCountByPriority[priority]["Ready"] == 0 {
            break miniLoop
        }
    }

    // fixed report issue: to display report for cases if have parent-child relationship, but not same priority
    tcsRunStore.TcTreeStats.ResetTcTreeStats(priority)
    tcsRunStore.TcTreeStats.CollectNodeStatusByPriority(tcsRunStore.Root, priority)
}

// for each global setup, normal, global teardown
func (tcsRunStore *TcsRunStore) RunConsoleOverallReport () {
    tcsRunStore.TcTreeStats.CollectOverallNodeStatus(tcsRunStore.Root, "Overall")

    tcsRunStore.OverallFail = tcsRunStore.OverallFail + tcsRunStore.TcTreeStats.StatusCountByPriority["Overall"]["Fail"]

    reports.ReportConsoleOverall(len(tcsRunStore.TcSlice), "Overall", tcsRunStore.TcTreeStats.StatusCountByPriority)
}

// for all (global setup, normal, global teardown)
func (g4Store *G4Store) RunFinalConsoleReport () {
    fmt.Println("")
    fmt.Println("---------------------------------------------------------------------------------")
    fmt.Println("Final Test Case Execution Statistics - Overall")

    totalTcCount := len(g4Store.FullTcSlice)

    reports.ReportConsoleOverall(totalTcCount, "Overall", 
        g4Store.GlobalSetUpRunStore.TcTreeStats.StatusCountByPriority, 
        g4Store.NormalRunStore.TcTreeStats.StatusCountByPriority, 
        g4Store.MutationRunStore.TcTreeStats.StatusCountByPriority, 
        g4Store.GlobalTeardownRunStore.TcTreeStats.StatusCountByPriority)
}

// for file report
func (g4Store *G4Store) RunFinalReport (ch chan int, gStart_str string, resultsDir string, resultsLogFile string) {
    gEnd_time := time.Now()
    gEnd_str := gEnd_time.Format("2006-01-02 15:04:05.000000000 +0800 CST")

    g4Store.OverallFail = g4Store.GlobalSetUpRunStore.OverallFail + g4Store.NormalRunStore.OverallFail + g4Store.GlobalTeardownRunStore.OverallFail

    reports.GenerateTestReport(gStart_str, gEnd_str, resultsDir, resultsLogFile)
    //
    fmt.Println("")
    fmt.Println("Report Generated at: " + resultsDir + "index.html")
    fmt.Println("Execution Finished at: " + gEnd_str)

    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    ch <- g4Store.OverallFail
}


func ScheduleCases (cReady chan *tree.TcNode, wg *sync.WaitGroup, resultsChan chan testcase.TestCaseExecutionInfo, baseUrl string) {
    // ------
    if cmd.Opt.IfConcurrency == true {
        tick := 0
        max := cmd.Opt.ConcurrencyLimit

        for tcNode := range cReady {
            wg.Add(1)
            // Note: to prevent reaching tcp connection limitation, here set a max, then sleep for a while
            if tick % max == 0 {
                time.Sleep(100 * time.Millisecond)
                go api.DispatchApi(wg, resultsChan, baseUrl, tcNode.TestCaseExecutionInfo.TestCaseDataInfo)
            } else {
                go api.DispatchApi(wg, resultsChan, baseUrl, tcNode.TestCaseExecutionInfo.TestCaseDataInfo)
            }

            tick = tick + 1
        }
    } else {
        for tcNode := range cReady {
            wg.Add(1)
            api.DispatchApi(wg, resultsChan, baseUrl, tcNode.TestCaseExecutionInfo.TestCaseDataInfo)
        }
    }   
}

func WriteNotNotExecutedToLog (priority string, logFilePtr *os.File, tcTreeStats tree.TcTreeStats) {
    notRunTime := time.Now()
    for i, _ := range tcTreeStats.TcNotExecutedByPriority[priority] {
        for _, tcExecution := range tcTreeStats.TcNotExecutedByPriority[priority][i] {
            //
            if tcExecution.Priority() == priority {
                // set some dummy time for the tc not executed
                tcExecution.StartTimeUnixNano = notRunTime.UnixNano()
                tcExecution.EndTimeUnixNano = notRunTime.UnixNano()
                tcExecution.DurationUnixNano = notRunTime.UnixNano() - notRunTime.UnixNano()

                tcReportResults := tcExecution.TcReportResults()
                
                repJson, _ := json.Marshal(tcReportResults)
                //
                reports.WriteExecutionResults(string(repJson), logFilePtr)

                // console 2. to print status to console (i.e. not executed cases: "ParentFailed", "ParentSkipped")
                reports.ReportConsoleByTc(*tcExecution)
            }
        }
    }
}



