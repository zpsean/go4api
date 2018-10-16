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
	// "go4api/texttmpl"

	. "github.com/ahmetb/go-linq"
)

// Stamp      = "Jan _2 15:04:05"
// StampMilli = "Jan _2 15:04:05.000"
// StampMicro = "Jan _2 15:04:05.000000"
// StampNano  = "Jan _2 15:04:05.000000000"

func (tcReportSlice TcReportSlice) GroupByTotalStartTime () []Group {
    var tcReportSubSlice []int64
    for i, _ := range tcReportSlice {
        tcReportSubSlice = append(tcReportSubSlice, (tcReportSlice[i].StartTimeUnixNano / 1000 / 1000 / 1000) * 1000)
    }

    query := GroupByStartTime(tcReportSubSlice)

    return query
}


func (tcReportSlice TcReportSlice) GroupBySuccessStartTime () []Group {
    var tcReportSubSlice []int64
    for i, _ := range tcReportSlice {
        if tcReportSlice[i].TestResult == "Success" {
            tcReportSubSlice = append(tcReportSubSlice, (tcReportSlice[i].StartTimeUnixNano / 1000 / 1000 / 1000) * 1000)
        }
    }

    query := GroupByStartTime(tcReportSubSlice)

    return query
}

func (tcReportSlice TcReportSlice) GroupByFailStartTime () []Group {
    var tcReportSubSlice []int64
    for i, _ := range tcReportSlice {
        if tcReportSlice[i].TestResult == "Fail" {
            tcReportSubSlice = append(tcReportSubSlice, (tcReportSlice[i].StartTimeUnixNano / 1000 / 1000 / 1000) * 1000)
        }
    }

    query := GroupByStartTime(tcReportSubSlice)

    return query
}

func GroupByStartTime (execStartSlice []int64) []Group {
    var query []Group

    From(execStartSlice).GroupByT(
        func(item int64) int64 { 
            return item
        },
        func(item int64) int { return 1 },
    ).OrderByT(
        func(g Group) int64 { return g.Key.(int64)},
    ).ToSlice(&query)

    return query
}


func (tcReportSlice TcReportSlice) GroupByOverallStatus () []Group {
    type ReportsStuct struct {
        TestResult string
    }

    var query []Group

    From(tcReportSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.TestResult}
        },
        func(item *testcase.TcReportResults) int64 { return 1 },
    ).ToSlice(&query)

    return query
}
