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
    "strconv"
    "encoding/json"

    "go4api/cmd"
    "go4api/lib/testcase"
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

func ReportConsoleByTc (tcExecution testcase.TestCaseExecutionInfo) {
    tcReportResults := tcExecution.TcConsoleResults()
    // repJson, _ := json.Marshal(tcReportResults)

    if tcReportResults.TestResult == "Fail" {
        length := len(string(tcExecution.ActualBody))
        out_len := 0
        if length > 300 {
            out_len = 300
        } else {
            out_len = length
        }

        fmt.Printf("\n%s%-40s%-3s%-30s%-10s%-30s%-30s%-4s%d%s\n", CLR_R, tcReportResults.TcName , tcReportResults.Priority, tcReportResults.ParentTestCase, 
            tcReportResults.TestResult, tcReportResults.JsonFilePath, tcReportResults.CsvFile, tcReportResults.CsvRow,
            tcReportResults.ActualStatusCode, CLR_N)

        if cmd.Opt.IfMutation {
            fmt.Println(tcReportResults.MutationInfoStr)
        }
        
        // fmt.Println(tcReportResults.MutationInfo)

        // by default, print failed field in testMessages
        failedTM := filterTestMessages(tcReportResults.TestMessages)
        failedTMBytes, _ := json.Marshal(failedTM)
        fmt.Println(string(failedTMBytes))

        fmt.Println(string(tcExecution.ActualBody)[0:out_len], "...")
    } else {
        fmt.Printf("\n%s%-40s%-3s%-30s%-10s%-30s%-30s%-4s%d%s\n", CLR_G, tcReportResults.TcName, tcReportResults.Priority, tcReportResults.ParentTestCase, 
            tcReportResults.TestResult, tcReportResults.JsonFilePath, tcReportResults.CsvFile, tcReportResults.CsvRow,
            tcReportResults.ActualStatusCode, CLR_N)

        if cmd.Opt.IfMutation {
            fmt.Println(tcReportResults.MutationInfoStr)
        }
    }
}


func ReportConsoleByPriority (totalTc int, priority string, statusCountByPriority map[string]map[string]int) {
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

func ReportConsoleOverall (totalTc int, key string, params ... map[string]map[string]int) {
    // ---
    var totalCount = 0
    var successCount = 0
    var failCount = 0
    var skipCount = 0

    for _, param := range params {
        totalCount = totalCount + param[key]["Total"]
        successCount = successCount + param[key]["Success"]
        failCount = failCount + param[key]["Fail"]
        skipCount = skipCount + param[key]["ParentFailed"]
    }
    
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
