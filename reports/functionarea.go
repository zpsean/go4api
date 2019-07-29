/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package reports

import (
 	// "fmt"

	"go4api/lib/testcase"

	. "github.com/ahmetb/go-linq"
)



func (tcReportSlice TcReportSlice) GroupByFunctionArea () []Group {
    type ReportsStuct struct {
        FunctionArea string
    }

    var query []Group

    From(tcReportSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.FunctionAreas[0]}
        },
        func(item *testcase.TcReportResults) int64 { return 1 },
    ).ToSlice(&query)

    return query
}
