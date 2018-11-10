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

    "go4api/lib/testcase"
    "go4api/lib/tree"
    "go4api/reports"
)

type G4Store struct {
    OverallFail             int
    FullTcSlice             []*testcase.TestCaseDataInfo
    GlobalSetUpRunStore     *TcsRunStore
    NormalRunStore          *TcsRunStore
    GlobalTeardownRunStore  *TcsRunStore
}

type TcsRunStore struct {
    TcSlice     []*testcase.TestCaseDataInfo
    PrioritySet []string
    Root        *tree.TcNode
    TcTree      tree.TcTree
    TcTreeStats tree.TcTreeStats
    OverallFail int
}

func InitG4Store () *G4Store {
    fullTcSlice := InitFullTcSlice()

    globalSetUpTcSlice := InitGlobalSetUpTcSlice(fullTcSlice)
    globalSetUpRunStore := &TcsRunStore {
        TcSlice: globalSetUpTcSlice,
        PrioritySet: []string{},
        Root: &tree.TcNode{},
        TcTree: tree.TcTree{},
        TcTreeStats: tree.TcTreeStats{},
        OverallFail: 0,
    }

    normalTcSlice := InitNormalTcSlice(fullTcSlice)
    normalRunStore := &TcsRunStore {
        TcSlice: normalTcSlice,
        PrioritySet: []string{},
        Root: &tree.TcNode{},
        TcTree: tree.TcTree{},
        TcTreeStats: tree.TcTreeStats{},
        OverallFail: 0,
    }

    globalTeardownTcSlice := InitGlobalTeardownTcSlice(fullTcSlice)
    globalTeardownRunStore := &TcsRunStore {
        TcSlice: globalTeardownTcSlice,
        PrioritySet: []string{},
        Root: &tree.TcNode{},
        TcTree: tree.TcTree{},
        TcTreeStats: tree.TcTreeStats{},
        OverallFail: 0,
    }

    g4Store := &G4Store {
        OverallFail: 0,
        FullTcSlice: fullTcSlice,
        GlobalSetUpRunStore: globalSetUpRunStore,
        NormalRunStore: normalRunStore,
        GlobalTeardownRunStore: globalTeardownRunStore,
    }

    return g4Store
}

func (tcsRunStore *TcsRunStore) InitRun () { 
    tcTree := tree.CreateTcTree()
    root := tcTree.BuildTree(tcsRunStore.TcSlice)
    //
    prioritySet := GetPrioritySet(tcsRunStore.TcSlice)
    tcTreeStats := tree.CreateTcTreeStats(prioritySet)
    // 
    tcTree.InitNodesRunResult(root, "Ready")

    tcsRunStore.PrioritySet = prioritySet
    tcsRunStore.Root = root
    tcsRunStore.TcTree = tcTree
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
        reports.ReportConsoleByPriority(0, priority, tcsRunStore.TcTreeStats.StatusCountByPriority)

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

        tcsRunStore.TcTreeStats.CollectNodeStatusByPriority(tcsRunStore.Root, priority)

        for tcExecution := range resultsExeChan {
            //
            c := make(chan *tree.TcNode)
            go func(c chan *tree.TcNode) {
                defer close(c)
                tcsRunStore.TcTree.SearchNode(c, tcsRunStore.Root, tcExecution.TcName())
            }(c)
            //
            tcsRunStore.TcTree.RefreshNodeAndDirectChilrenTcResult(<-c, tcExecution.TestResult, tcExecution.StartTime, tcExecution.EndTime, 
                tcExecution.HttpTestMessages, tcExecution.StartTimeUnixNano, tcExecution.EndTimeUnixNano)

            tcsRunStore.TcTreeStats.DeductReadyCount(priority)
            // (3). <--> for log write to file
            tcReportResults := tcExecution.TcReportResults()

            repJson, _ := json.Marshal(tcReportResults)

            reports.WriteExecutionResults(string(repJson), logFilePtr)

            reports.ReportConsoleByTc(tcExecution)
        }
        // if tcTree has no node with "Ready" status, break the miniloop
        if tcsRunStore.TcTreeStats.StatusCountByPriority[priority]["Ready"] == 0 {
            break miniLoop
        }
    }
}

func (tcsRunStore *TcsRunStore) RunConsoleOverallReport () {
    tcsRunStore.TcTreeStats.CollectOverallNodeStatus(tcsRunStore.Root, "Overall")

    tcsRunStore.OverallFail = tcsRunStore.OverallFail + tcsRunStore.TcTreeStats.StatusCountByPriority["Overall"]["Fail"]

    reports.ReportConsoleOverall(len(tcsRunStore.TcSlice), "Overall", tcsRunStore.TcTreeStats.StatusCountByPriority)
}


func (g4Store *G4Store) RunFinalConsoleReport () {
    fmt.Println("")
    fmt.Println("Final Test Execution Statistics")

    totalTcCount := len(g4Store.FullTcSlice)

    reports.ReportConsoleOverall(totalTcCount, "Overall", 
        g4Store.GlobalSetUpRunStore.TcTreeStats.StatusCountByPriority, 
        g4Store.NormalRunStore.TcTreeStats.StatusCountByPriority, 
        g4Store.GlobalTeardownRunStore.TcTreeStats.StatusCountByPriority)
}

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


