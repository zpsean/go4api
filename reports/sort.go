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
 	// "strconv"
 	// "encoding/json"

	"go4api/lib/testcase"

	. "github.com/ahmetb/go-linq"
)

func SortByStartTime (executionResultSlice []*testcase.TcReportResults) []*testcase.TcReportResults {
    var orderedExecutionResultSlice []*testcase.TcReportResults
    From(executionResultSlice).
        OrderByT(
            func(item *testcase.TcReportResults) int64 { return item.StartTimeUnixNano },
        ).
        ToSlice(&orderedExecutionResultSlice)

    // for _, item := range orderedExecutionResultSlice {
    //     fmt.Println(item.StartTimeUnixNano, item.EndTimeUnixNano, item.DurationUnixNano)
    // }

    return orderedExecutionResultSlice
}

func SortByEndTime (executionResultSlice []*testcase.TcReportResults) []*testcase.TcReportResults {
    var orderedExecutionResultSlice []*testcase.TcReportResults
    From(executionResultSlice).
        OrderByT(
            func(item *testcase.TcReportResults) int64 { return item.EndTimeUnixNano },
        ).
        ToSlice(&orderedExecutionResultSlice)

    // for _, item := range orderedExecutionResultSlice {
    //     fmt.Println(item.StartTimeUnixNano, item.EndTimeUnixNano, item.DurationUnixNano)
    // }

    return orderedExecutionResultSlice
}

func SortByDuration (executionResultSlice []*testcase.TcReportResults) []*testcase.TcReportResults {
    var orderedExecutionResultSlice []*testcase.TcReportResults
    From(executionResultSlice).
        OrderByT(
            func(item *testcase.TcReportResults) int64 { return item.DurationUnixNano },
        ).
        ToSlice(&orderedExecutionResultSlice)

    // for _, item := range orderedExecutionResultSlice {
    //     fmt.Println(item.StartTimeUnixNano, item.EndTimeUnixNano, item.DurationUnixNano)
    // }

    return orderedExecutionResultSlice
}

