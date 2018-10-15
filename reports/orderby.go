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

func (tcReportSlice TcReportSlice) SortByStartTime () TcReportSlice {
    var orderedTcReportSlice TcReportSlice
    From(tcReportSlice).
        OrderByT(
            func(item *testcase.TcReportResults) int64 { return item.StartTimeUnixNano },
        ).
        ToSlice(&orderedTcReportSlice)

    return orderedTcReportSlice
}

func (tcReportSlice TcReportSlice) SortByStartTimeDesc () TcReportSlice {
    var orderedTcReportSlice TcReportSlice
    From(tcReportSlice).
        OrderByDescendingT(
            func(item *testcase.TcReportResults) int64 { return item.StartTimeUnixNano },
        ).
        ToSlice(&orderedTcReportSlice)

    return orderedTcReportSlice
}


func (tcReportSlice TcReportSlice) SortByEndTime () TcReportSlice {
    var orderedTcReportSlice TcReportSlice
    From(tcReportSlice).
        OrderByT(
            func(item *testcase.TcReportResults) int64 { return item.EndTimeUnixNano },
        ).
        ToSlice(&orderedTcReportSlice)

    return orderedTcReportSlice
}

func (tcReportSlice TcReportSlice) SortByEndTimeDesc () TcReportSlice {
    var orderedTcReportSlice TcReportSlice
    From(tcReportSlice).
        OrderByDescendingT(
            func(item *testcase.TcReportResults) int64 { return item.EndTimeUnixNano },
        ).
        ToSlice(&orderedTcReportSlice)

    return orderedTcReportSlice
}


func (tcReportSlice TcReportSlice) SortByDuration () TcReportSlice {
    var orderedTcReportSlice TcReportSlice
    From(tcReportSlice).
        OrderByT(
            func(item *testcase.TcReportResults) int64 { return item.DurationUnixNano },
        ).
        ToSlice(&orderedTcReportSlice)

    return orderedTcReportSlice
}

func (tcReportSlice TcReportSlice) SortByDurationDesc () TcReportSlice {
    var orderedTcReportSlice TcReportSlice
    From(tcReportSlice).
        OrderByDescendingT(
            func(item *testcase.TcReportResults) int64 { return item.DurationUnixNano },
        ).
        ToSlice(&orderedTcReportSlice)

    return orderedTcReportSlice
}
