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
	// "os"
    "strings"
    "path/filepath"
    "encoding/json"

    "go4api/utils"
    "go4api/lib/testcase"

)

// this function is called by cmd -report, to generate report from log file
func GenerateReportsFromLogFile(logResultsFile string) {
    // 1. retrieve the resultsDir form logResultsFile
    resultsDir := filepath.Dir(logResultsFile)

    // 2. get the content from log file, json lines, to ExecutionResultSlice
    jsonLinesBytes := utils.GetContentFromFile(logResultsFile)
    jsonLines := string(jsonLinesBytes)
    //
    jsonLines = strings.Replace(jsonLines, "\n", ",", strings.Count(jsonLines, "\n") - 1)
    tcReportStr := `[` + jsonLines + `]`

    json.Unmarshal([]byte(tcReportStr), &ExecutionResultSlice)

    // 2.1 sort the ExecutionResultSlice, by start time / end time
    orderedByStartTime := SortByStartTime()
    orderedByEndTime := SortByEndTime()
    orderedByDuration := SortByDuration()

    // 3. get the: pStart_time, pStart, pEnd_time, tcClassifedCountMap, totalTc, statusCountByPriority, 
    //      tcExecutedByPriority, tcNotExecutedByPriority
    totalTc := len(ExecutionResultSlice)
    pStart := orderedByStartTime[0].StartTime
    pEnd_time := orderedByEndTime[totalTc - 1].EndTime

    tcClassifedCountMap := map[string]int{}
    statusCountByPriority := map[string]map[string]int{} 
    tcExecutedByPriority := map[string]map[string][]*testcase.TestCaseExecutionInfo{}
    tcNotExecutedByPriority := map[string]map[string][]*testcase.TestCaseExecutionInfo{}

    fmt.Println("StartTime: ", totalTc, pStart)
    fmt.Println("EndTime: ", totalTc, pEnd_time)

    fmt.Println("Min Duration: ", orderedByDuration[0].DurationUnixNano)
    fmt.Println("Max Duration: ", orderedByDuration[totalTc - 1].DurationUnixNano)

    fmt.Println("tcClassifedCountMap: ", tcClassifedCountMap)
    fmt.Println("statusCountByPriority: ", statusCountByPriority)
    fmt.Println("tcExecutedByPriority: ", tcExecutedByPriority)
    fmt.Println("tcNotExecutedByPriority: ", tcNotExecutedByPriority)

    // 4. call the function GenerateTestReport()

    fmt.Println("Report Generated at: " + resultsDir + "/index.html")
    // fmt.Println("Execution Finished at: " + pEnd_time.String())
}



