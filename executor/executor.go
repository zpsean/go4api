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
    "sync"
    "encoding/json"

    "go4api/cmd"
    "go4api/api"
    "go4api/lib/testcase"
    "go4api/lib/tree"
    "go4api/reports"
)

var overallFail = 0

func Run (ch chan int, baseUrl string, resultsDir string, resultsLogFile string, tcArray []testcase.TestCaseDataInfo) tree.TcTreeStats { 
    //-----
    prioritySet, root, tcTree, tcTreeStats := RunInit(tcArray)

    fmt.Println("\n====> test cases execution starts!") 

    RunPriorities(baseUrl, resultsDir, resultsLogFile, tcArray, prioritySet, root, tcTree, tcTreeStats)
    RunConsoleOverallReport(tcArray, root, tcTreeStats)

    return tcTreeStats
}

func RunInit (tcArray []testcase.TestCaseDataInfo) ([]string, *tree.TcNode, tree.TcTree, tree.TcTreeStats) { 
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

func RunPriorities (baseUrl string, resultsDir string, resultsLogFile string, tcArray []testcase.TestCaseDataInfo, prioritySet []string, 
        root *tree.TcNode, tcTree tree.TcTree, tcTreeStats tree.TcTreeStats) {
    // -------
    logFilePtr := reports.OpenExecutionResultsLogFile(resultsLogFile)

    for _, priority := range prioritySet {
        fmt.Println("====> Priority " + priority + " starts!")
        //
        RunEachPriority(baseUrl, tcArray, priority, root, tcTree, logFilePtr, tcTreeStats)

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

func RunEachPriority (baseUrl string, tcArray []testcase.TestCaseDataInfo, 
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

        ScheduleCases(cReady, &wg, resultsExeChan, baseUrl)
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

func RunConsoleOverallReport (tcArray []testcase.TestCaseDataInfo, root *tree.TcNode, tcTreeStats tree.TcTreeStats) {
    tcTreeStats.CollectOverallNodeStatus(root, "Overall")

    overallFail = overallFail + tcTreeStats.StatusCountByPriority["Overall"]["Fail"]

    reports.ReportConsoleOverall(len(tcArray), "Overall", tcTreeStats.StatusCountByPriority)
}

func RunFinalConsoleReport (totalTcCount int, setUpTcTreeStats tree.TcTreeStats, normalTcTreeStats tree.TcTreeStats, teardownTcTreeStats tree.TcTreeStats) {
    fmt.Println("")
    fmt.Println("Final Test Execution Statistics")

    reports.ReportConsoleOverall(totalTcCount, "Overall", setUpTcTreeStats.StatusCountByPriority, 
        normalTcTreeStats.StatusCountByPriority, teardownTcTreeStats.StatusCountByPriority)
}

func RunFinalReport (ch chan int, gStart_str string, resultsDir string, resultsLogFile string) {
    gEnd_time := time.Now()
    gEnd_str := gEnd_time.Format("2006-01-02 15:04:05.000000000 +0800 CST")

    reports.GenerateTestReport(gStart_str, gEnd_str, resultsDir, resultsLogFile)
    //
    fmt.Println("")
    fmt.Println("Report Generated at: " + resultsDir + "index.html")
    fmt.Println("Execution Finished at: " + gEnd_str)

    // channel code, can be used for the overall success or fail indicator, especially for CI/CD
    ch <- overallFail
}


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
                go api.HttpApi(wg, resultsChan, baseUrl, *(tcNode.TestCaseExecutionInfo.TestCaseDataInfo))
            } else {
                go api.HttpApi(wg, resultsChan, baseUrl, *(tcNode.TestCaseExecutionInfo.TestCaseDataInfo))
            }

            tick = tick + 1
        }
    } else {
        for tcNode := range cReady {
            wg.Add(1)
            api.HttpApi(wg, resultsChan, baseUrl, *(tcNode.TestCaseExecutionInfo.TestCaseDataInfo))
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



