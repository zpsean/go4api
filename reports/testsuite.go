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
    "encoding/json"

    "go4api/lib/testcase"

    . "github.com/ahmetb/go-linq"
)



func (tcReportSlice TcReportSlice) GroupByTestSuite () []Group {
    type ReportsStuct struct {
        TestSuite string
    }

    var query []Group

    From(tcReportSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.TestSuite}
        },
        func(item *testcase.TcReportResults) int64 { return 1 },
    ).ToSlice(&query)

    return query
}

func (tcReportSlice TcReportSlice) GetOverallTestSuiteStatusStatsJson () string {
    query := tcReportSlice.GroupByTestSuite()

    reportsOverallStatusSlice := PrintStatsGroup(query)

    reJson, _ := json.MarshalIndent(reportsOverallStatusSlice, "", "\t")

    return string(reJson)
}