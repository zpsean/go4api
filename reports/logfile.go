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
    // "fmt"
	// "os"
    // "strings"
    "path/filepath"
    "encoding/json"

    "go4api/utils"
    "go4api/ui/js"
    "go4api/lib/testcase"
    "go4api/texttmpl"

)

// this function is called by cmd -report, to generate report from log file
func GenerateReportsFromLogFile(resultsLogFile string) {
    var executionResultSlice []*testcase.TcReportResults
    //
    resultsDir := filepath.Dir(resultsLogFile)
    //
    jsonLinesBytes := utils.GetContentFromFile(resultsLogFile)
    json.Unmarshal(jsonLinesBytes, &executionResultSlice)

    setUpResultSlice, normalResultSlice, tearDownResultSlice := ClassifyResults(ExecutionResultSlice)
    // 
    setUpResultsJs := GetResultsJs(setUpResultSlice)
    normalResultsJs := GetResultsJs(normalResultSlice)
    tearDownResultsJs := GetResultsJs(tearDownResultSlice)
    //
    resultsFile := resultsDir + "/js/setUpResults.js"
    texttmpl.GenerateResultsJs(js.Results, resultsFile, setUpResultsJs, resultsLogFile)

    resultsFile = resultsDir + "/js/Results.js"
    texttmpl.GenerateResultsJs(js.Results, resultsFile, normalResultsJs, resultsLogFile)

    resultsFile = resultsDir + "/js/tearDownResults.js"
    texttmpl.GenerateResultsJs(js.Results, resultsFile, tearDownResultsJs, resultsLogFile)
    //
    // tcClassifedCountMap := map[string]int{}
    // statusCountByPriority := map[string]map[string]int{} 
    // tcExecutedByPriority := map[string]map[string][]*testcase.TestCaseExecutionInfo{}
    // tcNotExecutedByPriority := map[string]map[string][]*testcase.TestCaseExecutionInfo{}

    // fmt.Println("StartTime: ", totalTc, gStart)
    // fmt.Println("EndTime: ", totalTc, pEnd_time)

    // fmt.Println("Min Duration: ", orderedByDuration[0].DurationUnixNano)
    // fmt.Println("Max Duration: ", orderedByDuration[totalTc - 1].DurationUnixNano)

    // fmt.Println("tcClassifedCountMap: ", tcClassifedCountMap)
    // fmt.Println("statusCountByPriority: ", statusCountByPriority)
    // fmt.Println("tcExecutedByPriority: ", tcExecutedByPriority)
    // fmt.Println("tcNotExecutedByPriority: ", tcNotExecutedByPriority)

    // // 4. call the function GenerateTestReport()

    // fmt.Println("Report Generated at: " + resultsDir + "/index.html")
    // // fmt.Println("Execution Finished at: " + pEnd_time.String())
}

func ClassifyResults (ExecutionResultSlice []*testcase.TcReportResults) ([]*testcase.TcReportResults, []*testcase.TcReportResults, []*testcase.TcReportResults) {
    var setUpResultSlice []*testcase.TcReportResults
    var normalResultSlice []*testcase.TcReportResults
    var tearDownResultSlice []*testcase.TcReportResults

    for i, _ := range ExecutionResultSlice {
        switch ExecutionResultSlice[i].IfGlobalSetUpTearDown {
            case "SetUp":
                setUpResultSlice = append(setUpResultSlice, ExecutionResultSlice[i])
            case "TearDown":
                tearDownResultSlice = append(tearDownResultSlice, ExecutionResultSlice[i])
            default:
                normalResultSlice = append(normalResultSlice, ExecutionResultSlice[i])
        }
    }

    return setUpResultSlice, normalResultSlice, tearDownResultSlice
}

func GetResultsJs (resultSlice []*testcase.TcReportResults) *texttmpl.ResultsJs {
    var resultsJs texttmpl.ResultsJs
    //
    if len(resultSlice) > 0 {
        tcReportBytes, _ := json.Marshal(resultSlice)
        tcReportStr := string(tcReportBytes)
        //
        orderedByStartTime := SortByStartTime(resultSlice)
        orderedByEndTime := SortByEndTime(resultSlice)
        // orderedByDuration := SortByDuration(resultSlice)
        //
        totalTc := len(resultSlice)
        pStart_time := orderedByStartTime[0].StartTime
        pEnd_time := orderedByEndTime[totalTc - 1].EndTime
        //
        resultsJs = texttmpl.ResultsJs {
            // GStart_time: pStart_time.UnixNano(), 
            GStart: `"` + pStart_time + `"`, 
            // PEnd_time: pEnd_time.UnixNano(), 
            PEnd: `"` + pEnd_time + `"`, 
            TcReportStr: tcReportStr,
        }
    }
    
     
    return &resultsJs
}


