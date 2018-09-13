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


func GenerateTestReport(resultsDir string, pStart_time time.Time, pStart string, pEnd_time time.Time) {
    // read the resource under /ui/*
    // copy the value of var Index to file
    utils.GenerateFileBasedOnVarOverride(ui.Index, resultsDir + "index.html")
    // (0)
    err := os.MkdirAll(resultsDir + "js", 0777)
    if err != nil {
      panic(err) 
    }
    // (1). get js/reslts.js
    logResultsFile := resultsDir + pStart + ".log"
    jsResults := resultsDir + "/js/reslts.js"
    //
    texttmpl.GenerateHtmlJsCSSFromTemplateAndVar(js.Results, pStart_time, pEnd_time, jsResults, logResultsFile)
    // (2). get js/go4api.js
    utils.GenerateFileBasedOnVarOverride(js.Js, resultsDir + "js/go4api.js")
    //
    err = os.MkdirAll(resultsDir + "style", 0777)
    if err != nil {
      panic(err) 
    }
    // (3). get style/go4api.css
    utils.GenerateFileBasedOnVarOverride(style.Style, resultsDir + "style/go4api.css")
}


func ReportConsole (tcExecution testcase.TestCaseExecutionInfo, actualBody []byte) {
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

func filterTestMessages (testMessages []*testcase.TestMessage) []*testcase.TestMessage {
    var failedTM []*testcase.TestMessage
    for i, _ := range testMessages {
        if testMessages[i].AssertionResults == "Failed" {
            failedTM = append(failedTM, testMessages[i])
        }
    }

    return failedTM
}


