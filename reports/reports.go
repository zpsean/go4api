/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package reports

import (
    "fmt"
	"os"
    "time"
    "strconv"
    "strings"
    "encoding/json"

    "go4api/lib/testcase"
    "go4api/ui"     
    "go4api/ui/js"  
    "go4api/ui/style"                                                                                                                                
    "go4api/utils"
    "go4api/texttmpl"
)

const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_N = "\x1b[0m"


func GenerateTestReport(resultsDir string, pStart_time time.Time, pStart string, pEnd_time time.Time,
        tcClassifedCountMap map[string]int, totalTc int, statusCountByPriority map[string]map[string]int, 
        tcExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo,
        tcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) {
    // --------

    utils.GenerateFileBasedOnVarOverride(ui.Index, resultsDir + "index.html")
    // (0)
    err := os.MkdirAll(resultsDir + "js", 0777)
    if err != nil {
      panic(err) 
    }
    logResultsFile := resultsDir + pStart + ".log"

    statsFile := resultsDir + "/js/stats.js"
    tcStats := GetStats(tcClassifedCountMap, totalTc, statusCountByPriority, tcExecutedByPriority, tcNotExecutedByPriority)
    statsJs := tcStats.PrepStatsReport()

    texttmpl.GenerateStatsJs(js.Stats, statsFile, statsJs, logResultsFile)

    // (1). get js/reslts.js
    resultsFile := resultsDir + "/js/reslts.js"
    resultsJs := GetResultsJs(pStart_time, pEnd_time, logResultsFile)

    texttmpl.GenerateResultsJs(js.Results, resultsFile, resultsJs, logResultsFile)
    
    // (2). get js/go4api.js
    utils.GenerateFileBasedOnVarOverride(js.Js, resultsDir + "js/go4api.js")
    
    // (3). get style/go4api.css
    err = os.MkdirAll(resultsDir + "style", 0777)
    if err != nil {
      panic(err) 
    }
    utils.GenerateFileBasedOnVarOverride(style.Style, resultsDir + "style/go4api.css")
}

func GetResultsJs (pStart_time time.Time, pEnd_time time.Time, logResultsFile string) *texttmpl.ResultsJs {
    // get the data from the log results file, used for ui
    var tcReportStr string

    jsonLinesBytes := utils.GetContentFromFile(logResultsFile)
    jsonLines := string(jsonLinesBytes)
    //
    jsonLines = strings.Replace(jsonLines, "\n", ",", strings.Count(jsonLines, "\n") - 1)
    tcReportStr = `[` + jsonLines + `]`        
    //
    resultsJs := texttmpl.ResultsJs {
        PStart_time: pStart_time.UnixNano(), 
        PStart: `"` + pStart_time.String() + `"`, 
        PEnd_time: pEnd_time.UnixNano(), 
        PEnd: `"` + pEnd_time.String() + `"`, 
        TcReportStr: tcReportStr,
    }
     
    return &resultsJs
}

func ReportConsoleByTc (tcExecution testcase.TestCaseExecutionInfo, actualBody []byte) {
    tcReportResults := tcExecution.TcConsoleResults()
    // repJson, _ := json.Marshal(tcReportResults)

    if tcReportResults.TestResult == "Fail" {
        length := len(string(actualBody))
        out_len := 0
        if length > 300 {
            out_len = 300
        } else {
            out_len = length
        }

        fmt.Printf("\n%s%-40s%-3s%-30s%-10s%-30s%-30s%-4s%d%s\n", CLR_R, tcReportResults.TcName , tcReportResults.Priority, tcReportResults.ParentTestCase, 
            tcReportResults.TestResult, tcReportResults.JsonFilePath, tcReportResults.CsvFile, tcReportResults.CsvRow,
            tcReportResults.ActualStatusCode, CLR_N)

        if tcReportResults.MutationInfo != nil {
            fmt.Println(tcReportResults.MutationInfo)
        }
        
        // fmt.Println(tcReportResults.MutationInfo)

        // by default, print failed field in testMessages
        failedTM := filterTestMessages(tcReportResults.TestMessages)
        failedTMBytes, _ := json.Marshal(failedTM)
        fmt.Println(string(failedTMBytes))

        fmt.Println(string(actualBody)[0:out_len], "...")
    } else {
        fmt.Printf("\n%s%-40s%-3s%-30s%-10s%-30s%-30s%-4s%d%s\n", CLR_G, tcReportResults.TcName, tcReportResults.Priority, tcReportResults.ParentTestCase, 
            tcReportResults.TestResult, tcReportResults.JsonFilePath, tcReportResults.CsvFile, tcReportResults.CsvRow,
            tcReportResults.ActualStatusCode, CLR_N)

        if tcReportResults.MutationInfo != nil {
            fmt.Println(tcReportResults.MutationInfo)
        }
    }
}


func ReportConsoleByPriority (totalTc int, priority string, statusCountByPriority map[string]map[string]int, 
        tcExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo,
        tcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) {
    // ---
    var totalCount = statusCountByPriority[priority]["Total"]
    var successCount = statusCountByPriority[priority]["Success"]
    var failCount = statusCountByPriority[priority]["Fail"]
    var skipCount = statusCountByPriority[priority]["ParentFailed"]
    //
    fmt.Println("---------------------------------------------------------------------------------")
    fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(totalTc) + " Cases in Source")
    fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(totalCount) + " Cases recognized from template")
    fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(successCount + failCount) + " Cases Executed")
    fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(successCount) + " Cases Success")
    fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(failCount) + " Cases Fail")
    fmt.Println("----- Priority " + priority + ": " + strconv.Itoa(skipCount) + " Cases Skipped (Not Executed, due to Parent Failed)")
    fmt.Println("---------------------------------------------------------------------------------")
}

func ReportConsoleOverall (totalTc int, key string, statusCountByPriority map[string]map[string]int, 
        tcExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo,
        tcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) {
    // ---
    var totalCount = statusCountByPriority[key]["Total"]
    var successCount = statusCountByPriority[key]["Success"]
    var failCount = statusCountByPriority[key]["Fail"]
    var skipCount = statusCountByPriority[key]["ParentFailed"]
    //
    fmt.Println("---------------------------------------------------------------------------------")
    fmt.Println("----- " + key + ": " + strconv.Itoa(totalTc) + " Cases in Source")
    fmt.Println("----- " + key + ": " + strconv.Itoa(totalCount) + " Cases recognized from template")
    fmt.Println("----- " + key + ": " + strconv.Itoa(successCount + failCount) + " Cases Executed")
    fmt.Println("----- " + key + ": " + strconv.Itoa(successCount) + " Cases Success")
    fmt.Println("----- " + key + ": " + strconv.Itoa(failCount) + " Cases Fail")
    fmt.Println("----- " + key + ": " + strconv.Itoa(skipCount) + " Cases Skipped (Not Executed, due to Parent Failed)")
    fmt.Println("---------------------------------------------------------------------------------")
}

func filterTestMessages (testMessages []*testcase.TestMessage) []*testcase.TestMessage {
    var failedTM []*testcase.TestMessage
    for i, _ := range testMessages {
        if testMessages[i].AssertionResults == "Failed" {
            failedTM = append(failedTM, testMessages[i])
        }
    }

    return failedTM
}


