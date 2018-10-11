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
    "os"
    // "sort"
    "sync"
    // "path/filepath"
    // "strings"
    // "io/ioutil"
    // "strconv"
    "encoding/json"

    "go4api/cmd"
    // "go4api/utils"
    "go4api/api"
    "go4api/lib/testcase"
    "go4api/lib/tree"
    // "go4api/texttmpl"
    "go4api/reports"
)


func Run (ch chan int, pStart_time time.Time, pStart string, baseUrl string, resultsDir string, tcArray []testcase.TestCaseDataInfo) { 
    //-----
    prioritySet, root, tcTree, tcTreeStats := RunBefore(tcArray)

    fmt.Println("\n====> test cases execution starts!") 

    RunPriorities(ch, pStart, baseUrl, resultsDir, tcArray, prioritySet, root, tcTree, tcTreeStats)
    RunConsoleOverallReport(ch, pStart_time, pStart, resultsDir, tcArray, root, tcTree, tcTreeStats)
    RunFinalReport(ch, pStart_time, pStart, resultsDir, tcArray, root, tcTree, tcTreeStats)
}

func RunBefore (tcArray []testcase.TestCaseDataInfo) ([]string, *tree.TcNode, tree.TcTree, tree.TcTreeStats) { 
    // check the tcArray, if the case not distinct, report it to fix
    if len(tcArray) != len(GetTcNameSet(tcArray)) {
        fmt.Println("\n!! There are duplicated test case names, please make them distinct")
        os.Exit(1)
    }
    //
    tcTree := tree.CreateTcTree()
    root := tcTree.BuildTree(tcArray)
    //
    prioritySet := GetPrioritySet(tcArray)
    tcTreeStats := tree.CreateTcTreeStats(prioritySet)
    // Init
    tcTree.InitNodesRunResult(root, "Ready")

    return prioritySet, root, tcTree, tcTreeStats
}

func RunPriorities (ch chan int, pStart string, baseUrl string, resultsDir string, tcArray []testcase.TestCaseDataInfo, prioritySet []string, 
        root *tree.TcNode, tcTree tree.TcTree, tcTreeStats tree.TcTreeStats) {
    // -------
    logFilePtr := reports.OpenExecutionResultsLogFile(resultsDir + pStart + ".log")

    for _, priority := range prioritySet {
        fmt.Println("====> Priority " + priority + " starts!")
        //
        RunEachPriority(ch, pStart, baseUrl, resultsDir, tcArray, priority, root, tcTree, logFilePtr, tcTreeStats)

        // Put out the cases which has not been executed (i.e. not Success or Fail)
        WriteNotNotExecutedToLog(priority, logFilePtr, tcTreeStats)

        // report to console
        reports.ReportConsoleByPriority(0, priority, tcTreeStats.StatusCountByPriority)

        fmt.Println("====> Priority " + priority + " ended!")
        fmt.Println("")
        // sleep for debug
        // time.Sleep(500 * time.Millisecond)
    }

    logFilePtr.Close()
}


func RunEachPriority (ch chan int, pStart string, baseUrl string, resultsDir string, tcArray []testcase.TestCaseDataInfo, 
        priority string, root *tree.TcNode, tcTree tree.TcTree, logFilePtr *os.File, tcTreeStats tree.TcTreeStats) {
    // ----------
    miniLoop:
    for {
        //
        resultsExeChan := make(chan testcase.TestCaseExecutionInfo, len(tcArray))
        var wg sync.WaitGroup
        //
        cReady := make(chan *tree.TcNode)
        go func(cReady chan *tree.TcNode) {
            defer close(cReady)
            tcTree.CollectNodeReadyByPriority(cReady, root, priority)
        }(cReady)

        ScheduleCases(cReady, &wg, resultsExeChan, pStart, baseUrl, resultsDir)
        //
        wg.Wait()

        close(resultsExeChan)

        tcTreeStats.CollectNodeStatusByPriority(root, priority)

        for tcExecution := range resultsExeChan {
            // (1). tcName, testResult, the search result is saved to *findNode
            c := make(chan *tree.TcNode)
            go func(c chan *tree.TcNode) {
                defer close(c)
                tcTree.SearchNode(c, root, tcExecution.TcName())
            }(c)
            // (2). 
            tcTree.RefreshNodeAndDirectChilrenTcResult(<-c, tcExecution.TestResult, tcExecution.StartTime, tcExecution.EndTime, 
                tcExecution.TestMessages, tcExecution.StartTimeUnixNano, tcExecution.EndTimeUnixNano)

            tcTreeStats.DeductReadyCount(priority)
            // (3). <--> for log write to file
            tcReportResults := tcExecution.TcReportResults()
            reports.ExecutionResultSlice = append(reports.ExecutionResultSlice, tcReportResults)

            repJson, _ := json.Marshal(tcReportResults)
            reports.WriteExecutionResults(string(repJson), logFilePtr)

            reports.ReportConsoleByTc(tcExecution)
        }
        // if tcTree has no node with "Ready" status, break the miniloop
        if tcTreeStats.StatusCountByPriority[priority]["Ready"] == 0 {
            break miniLoop
        }
    }
}

func RunConsoleOverallReport (ch chan int, pStart_time time.Time, pStart string, resultsDir string, tcArray []testcase.TestCaseDataInfo, 
        root *tree.TcNode, tcTree tree.TcTree, tcTreeStats tree.TcTreeStats) {
    // -------
    tcTreeStats.CollectOverallNodeStatus(root, "Overall")
    reports.ReportConsoleOverall(len(tcArray), "Overall", tcTreeStats.StatusCountByPriority)
}

func RunFinalReport (ch chan int, pStart_time time.Time, pStart string, resultsDir string, tcArray []testcase.TestCaseDataInfo, 
        root *tree.TcNode, tcTree tree.TcTree, tcTreeStats tree.TcTreeStats) {
    // -------
    // generate the html report based on template, and results data
    // time.Sleep(1 * time.Second)
    pEnd_time := time.Now()
    //
    reports.GenerateTestReport(resultsDir, pStart_time, pStart, pEnd_time, 
        "", len(tcArray), tcTreeStats.StatusCountByPriority)
    //
    fmt.Println("Report Generated at: " + resultsDir + "index.html")
    fmt.Println("Execution Finished at: " + pEnd_time.String())

    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    // ch <- tcTreeStats.StatusCountByPriority["Overall"]["Fail"]

    // repJson, _ := json.Marshal(tcTree)
    // fmt.Println(string(repJson))
}


func ScheduleCases (cReady chan *tree.TcNode, wg *sync.WaitGroup, resultsChan chan testcase.TestCaseExecutionInfo, 
        pStart string, baseUrl string, resultsDir string) {
    //
    tick := 0
    max := cmd.Opt.ConcurrencyLimit

    for tcNode := range cReady {
        wg.Add(1)
        // Note: to prevent reaching tcp connection limitation, here set a max, then sleep for a while
        if tick % max == 0 {
            time.Sleep(100 * time.Millisecond)
            go api.HttpApi(wg, resultsChan, pStart, baseUrl, *(tcNode.TestCaseExecutionInfo.TestCaseDataInfo), resultsDir)
        } else {
            go api.HttpApi(wg, resultsChan, pStart, baseUrl, *(tcNode.TestCaseExecutionInfo.TestCaseDataInfo), resultsDir)
        }

        tick = tick + 1
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
                reports.ExecutionResultSlice = append(reports.ExecutionResultSlice, tcReportResults)
                
                repJson, _ := json.Marshal(tcReportResults)
                //
                reports.WriteExecutionResults(string(repJson), logFilePtr)
            }
        }
    }
}



