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
 	// "strconv"
 	"encoding/json"

	"go4api/lib/testcase"
	"go4api/texttmpl"

	. "github.com/ahmetb/go-linq"
)

type ReportsStats struct {
    ReportKey interface{}
    Count int
}

var (
	tcStats TcStats

    statusStats = map[string]map[string]int{}
    statusStatsPercentage = map[string]map[string]float64{}
)

func InitVariables(statusCountByPriority map[string]map[string]int) {
    for key, _ := range statusCountByPriority {
        statusStats[key] = map[string]int{}
        statusStatsPercentage[key] = map[string]float64{}
    }
}


func Get_Stats_1 (statusCountByPriority map[string]map[string]int) []byte {
	statsJsonBytes, _ := json.MarshalIndent(statusCountByPriority, "", "\t")

	return statsJsonBytes
}

func Get_Stats_2 () []Group {
    // type ReportsStuct struct {
    //     StartTimeUnixM int64
    // }
    var ExecutionStartSlice []int64
    for i, _ := range ExecutionResultSlice {
        ExecutionStartSlice = append(ExecutionStartSlice, (ExecutionResultSlice[i].StartTimeUnixNano / 1000 / 1000 / 1000) * 1000)
    }

    var query []Group

    From(ExecutionStartSlice).GroupByT(
        func(item int64) int64 { 
            return item
        },
        func(item int64) int { return 1 },
    ).OrderByT(
        func(g Group) int64 { return g.Key.(int64)},
    ).ToSlice(&query)

    return query
}


func Get_Stats_2_Success () []Group {
    var ExecutionStartSlice []int64
    for i, _ := range ExecutionResultSlice {
        if ExecutionResultSlice[i].TestResult == "Success" {
            ExecutionStartSlice = append(ExecutionStartSlice, (ExecutionResultSlice[i].StartTimeUnixNano / 1000 / 1000 / 1000) * 1000)
        }
    }

    var query []Group

    From(ExecutionStartSlice).GroupByT(
        func(item int64) int64 { 
            return item
        },
        func(item int64) int { return 1 },
    ).OrderByT(
        func(g Group) int64 { return g.Key.(int64)},
    ).ToSlice(&query)

    return query
}

func Get_Stats_2_Fail () []Group {
    var ExecutionStartSlice []int64
    for i, _ := range ExecutionResultSlice {
        if ExecutionResultSlice[i].TestResult == "Fail" {
            ExecutionStartSlice = append(ExecutionStartSlice, (ExecutionResultSlice[i].StartTimeUnixNano / 1000 / 1000 / 1000) * 1000)
        }
    }

    var query []Group

    From(ExecutionStartSlice).GroupByT(
        func(item int64) int64 { 
            return item
        },
        func(item int64) int { return 1 },
    ).OrderByT(
        func(g Group) int64 { return g.Key.(int64)},
    ).ToSlice(&query)

    return query
}

func Get_Stats_3 () {
    var orderedExecutionResultSlice []*testcase.TcReportResults
    From(ExecutionResultSlice).
        OrderByT(
            func(item *testcase.TcReportResults) int64 { return item.StartTimeUnixNano },
        ).
        ToSlice(&orderedExecutionResultSlice)

    for _, item := range orderedExecutionResultSlice {
        fmt.Println(item.StartTimeUnixNano, item.EndTimeUnixNano, item.DurationUnixNano)
    }
}


func PrintStatsGroup (query []Group) []ReportsStats {
    var reportsStatsSlice []ReportsStats

    for _, q := range query {
        ii := 0
        for range q.Group {
            ii += 1
        }

        reportsStats := ReportsStats {
            ReportKey: q.Key,
            Count: ii,
        }
        reportsStatsSlice = append(reportsStatsSlice, reportsStats)
    }
    return reportsStatsSlice
}

func ToOrderStatsGroup (reportsStatsTotalSlice []ReportsStats, reportsStatsSlice []ReportsStats) []ReportsStats {
    var reportsStatsOrdered []ReportsStats

    for i, _ := range reportsStatsTotalSlice {
        inx := -1
        for j, _ := range reportsStatsSlice {
            if reportsStatsTotalSlice[i].ReportKey == reportsStatsSlice[j].ReportKey {
                inx = j
                continue
            }
        }
        if inx != -1 {
            reportsStats := ReportsStats {
                ReportKey: reportsStatsTotalSlice[i].ReportKey,
                Count: reportsStatsSlice[inx].Count,
            }
            reportsStatsOrdered = append(reportsStatsOrdered, reportsStats)
        } else {
            reportsStats := ReportsStats {
                ReportKey: reportsStatsTotalSlice[i].ReportKey,
                Count: 0,
            }
            reportsStatsOrdered = append(reportsStatsOrdered, reportsStats)
        }
    }

    return reportsStatsOrdered
}

func GetStatsJson(statusCountByPriority map[string]map[string]int) []string {
    var reJsons []string

    reJson := Get_Stats_1(statusCountByPriority)
    reJsons = append(reJsons, string(reJson))

    query := Get_Stats_2()
    reportsStatsTotalSlice := PrintStatsGroup(query)
    reJson, _ = json.Marshal(reportsStatsTotalSlice)
    reJsons = append(reJsons, string(reJson))
    // fmt.Println("=====> reportsStatsSlice: ", string(reJson))

    query = Get_Stats_2_Success()
    reportsStatsSuccessSlice := PrintStatsGroup(query)

    reportsStatsSuccessSliceOrdered := ToOrderStatsGroup(reportsStatsTotalSlice, reportsStatsSuccessSlice)

    reJson, _ = json.Marshal(reportsStatsSuccessSliceOrdered)
    reJsons = append(reJsons, string(reJson))
    // fmt.Println("=====> reportsStatsSlice: ", string(reJson))

    query = Get_Stats_2_Fail()
    reportsStatsFailSlice := PrintStatsGroup(query)

    reportsStatsFailSliceOrdered := ToOrderStatsGroup(reportsStatsTotalSlice, reportsStatsFailSlice)

    reJson, _ = json.Marshal(reportsStatsFailSliceOrdered)
    reJsons = append(reJsons, string(reJson))

    return reJsons
}


func GetExecutedJson (tcExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) *texttmpl.DetailsJs {
	statsJsonBytes, _ := json.MarshalIndent(tcExecutedByPriority, "", "\t")

	tcStatsReport := texttmpl.DetailsJs {
		StatsStr: string(statsJsonBytes),
	}

	return &tcStatsReport
}

func GetNotExecutedJson (tcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) *texttmpl.DetailsJs {
	statsJsonBytes, _ := json.MarshalIndent(tcNotExecutedByPriority, "", "\t")

	tcStatsReport := texttmpl.DetailsJs {
		StatsStr: string(statsJsonBytes),
	}

	return &tcStatsReport
}

