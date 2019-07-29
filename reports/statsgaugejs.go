/*
 * go4api - an api testing tool written in Go
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
    "math"
 	// "strconv"
 	"encoding/json"

	"go4api/lib/testcase"
	// "go4api/texttmpl"
    // "go4api/ui/js" 

	. "github.com/ahmetb/go-linq"
)

type StatsGauge struct {
    ReportKey interface{}
    Count int
    PerformanceGauge *PerformanceGauge
}

type PerformanceGauge struct {
    Min int64
    P50 int64
    P75 int64
    P95 int64
    P99 int64
    Max int64
    Mean int64
    StdDev int64
}

// func GenerateStatsGaugeJs (tcReportSlice TcReportSlice, resultsDir string) {
//     statsFile := resultsDir + "/js/statsgauge.js"

//     reJsons := tcReportSlice.GeStatsGaugeJson()
//     texttmpl.GenerateStatsGaugeJs(js.StatsGauge, statsFile, reJsons)
// }

func (tcReportSlice TcReportSlice) GetStatsGaugeJson () string {
    var finalStatsGaugeSlice []StatsGauge

    reportsStatsGaugeSliceL1 := tcReportSlice.GetStatsGaugeJsonL1()
    finalStatsGaugeSlice = append(finalStatsGaugeSlice, reportsStatsGaugeSliceL1...)

    reportsStatsGaugeSliceL2 := tcReportSlice.GetStatsGaugeJsonL2()
    finalStatsGaugeSlice = append(finalStatsGaugeSlice, reportsStatsGaugeSliceL2...)

    reportsStatsGaugeSliceL3 := tcReportSlice.GetStatsGaugeJsonL3()
    finalStatsGaugeSlice = append(finalStatsGaugeSlice, reportsStatsGaugeSliceL3...)

    reportsStatsGaugeSliceL4 := tcReportSlice.GetStatsGaugeJsonL4()
    finalStatsGaugeSlice = append(finalStatsGaugeSlice, reportsStatsGaugeSliceL4...)

    reJson, _ := json.MarshalIndent(finalStatsGaugeSlice, "", "\t")
    
    return string(reJson)
}


func (tcReportSlice TcReportSlice) GetStatsGaugeJsonL1 () []StatsGauge {
    query := tcReportSlice.GroupByStatsGaugeDetailsL1()
    reportsStatsGaugeSliceL1 := PrintGroupStatsGauge(query)
    
    return reportsStatsGaugeSliceL1
}

func (tcReportSlice TcReportSlice) GetStatsGaugeJsonL2 () []StatsGauge {
    query := tcReportSlice.GroupByStatsGaugeDetailsL2()
    reportsStatsGaugeSliceL2 := PrintGroupStatsGauge(query)

    return reportsStatsGaugeSliceL2
}

func (tcReportSlice TcReportSlice) GetStatsGaugeJsonL3 () []StatsGauge {
    query := tcReportSlice.GroupByStatsGaugeDetailsL3()
    reportsStatsGaugeSliceL3 := PrintGroupStatsGauge(query)
 
    return reportsStatsGaugeSliceL3
}

func (tcReportSlice TcReportSlice) GetStatsGaugeJsonL4 () []StatsGauge {
    var tcReportGaugeSlice []StatsGauge
    var durationSlice []int64

    performanceGauge := PerformanceGauge {
        Min: 0,
        P50: 0,
        P75: 0,
        P95: 0,
        P99: 0,
        Max: 0,
        Mean: 0,
        StdDev: 0,
    }
    totalTc := len(tcReportSlice)
    if len(tcReportSlice) > 0 {
        orderedByDuration := tcReportSlice.SortByDuration()
        totalTcF := float64(totalTc)

        for i, _ := range orderedByDuration {
            durationSlice = append(durationSlice, orderedByDuration[i].DurationUnixNano / 1000000)
        }
        mean := From(durationSlice).SumInts() / int64(totalTc)

        var v int64
        v = 0
        for i, _ := range durationSlice {
            v = v + (durationSlice[i] - mean) * (durationSlice[i] - mean)
        }
        // try later to use v / (int64(totalTc) - 1)
        variance := v / int64(totalTc)
        stddev := int64(math.Sqrt(float64(variance)))

        performanceGauge = PerformanceGauge {
            Min: orderedByDuration[0].DurationUnixNano / 1000000,
            P50: orderedByDuration[int(math.Floor(totalTcF * 0.5))].DurationUnixNano / 1000000,
            P75: orderedByDuration[int(math.Floor(totalTcF * 0.75))].DurationUnixNano / 1000000,
            P95: orderedByDuration[int(math.Floor(totalTcF * 0.95))].DurationUnixNano / 1000000,
            P99: orderedByDuration[int(math.Floor(totalTcF * 0.99))].DurationUnixNano / 1000000,
            Max: orderedByDuration[totalTc - 1].DurationUnixNano / 1000000,
            Mean: mean,
            StdDev: stddev,
        }
    }

    statsGauge := StatsGauge {
        ReportKey: map[string]string{
            "IfGlobalSetUpTearDown": "ALL",
            "Priority": "ALL",
            "TestResult": "ALL",
        },
        Count: totalTc,
        PerformanceGauge: &performanceGauge,
    }

    tcReportGaugeSlice = append(tcReportGaugeSlice, statsGauge)

    return tcReportGaugeSlice
}

// ----
func (tcReportSlice TcReportSlice) GroupByStatsGaugeDetailsL1 () []Group {
    type ReportsStuct struct {
        IfGlobalSetUpTearDown string
        Priority string
        TestResult string
    }

    var query []Group

    From(tcReportSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.IfGlobalSetUpTearDown, item.Priority, item.TestResult}
        },
        func(item *testcase.TcReportResults) int64 { return item.DurationUnixNano / 1000000 },
    ).ToSlice(&query)

    return query
}

func (tcReportSlice TcReportSlice) GroupByStatsGaugeDetailsL2 () []Group {
    type ReportsStuct struct {
        IfGlobalSetUpTearDown string
        Priority string
        TestResult string
    }

    var query []Group

    From(tcReportSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.IfGlobalSetUpTearDown, item.Priority, "ALL"}
        },
        func(item *testcase.TcReportResults) int64 { return item.DurationUnixNano / 1000000 },
    ).ToSlice(&query)

    return query
}

func (tcReportSlice TcReportSlice) GroupByStatsGaugeDetailsL3 () []Group {
    type ReportsStuct struct {
        IfGlobalSetUpTearDown string
        Priority string
        TestResult string
    }

    var query []Group

    From(tcReportSlice).GroupByT(
        func(item *testcase.TcReportResults) ReportsStuct { 
            return ReportsStuct{item.IfGlobalSetUpTearDown, "ALL", "ALL"}
        },
        func(item *testcase.TcReportResults) int64 { return item.DurationUnixNano / 1000000 },
    ).ToSlice(&query)

    return query
}

// ---
func PrintGroupStatsGauge (query []Group) []StatsGauge {
    // []Group = [ {"Key": {,,}, "Group": [,,]}, ]
    var tcReportGaugeSlice []StatsGauge

    for _, q := range query {
        performanceGauge := GetPerformanceGauge(q.Group)

        statsGauge := StatsGauge {
            ReportKey: q.Key,
            Count: len(q.Group),
            PerformanceGauge: performanceGauge,
        }
        tcReportGaugeSlice = append(tcReportGaugeSlice, statsGauge)
    }
    return tcReportGaugeSlice
}

func GetPerformanceGauge (group []interface{}) *PerformanceGauge {
    performanceGauge := PerformanceGauge {
        Min: 0,
        P50: 0,
        P75: 0,
        P95: 0,
        P99: 0,
        Max: 0,
        Mean: 0,
        StdDev: 0,
    }
    if len(group) > 0 {
        totalTc := len(group)
        totalTcF := float64(totalTc)

        mean := From(group).SumInts() / int64(totalTc)

        var v int64
        v = 0
        for i, _ := range group {
            v = v + (group[i].(int64) - mean) * (group[i].(int64) - mean)
        }
        variance := v / int64(totalTc)
        stddev := int64(math.Sqrt(float64(variance)))

        performanceGauge = PerformanceGauge {
            Min: group[0].(int64),
            P50: group[int(math.Floor(totalTcF * 0.5))].(int64),
            P75: group[int(math.Floor(totalTcF * 0.75))].(int64),
            P95: group[int(math.Floor(totalTcF * 0.95))].(int64),
            P99: group[int(math.Floor(totalTcF * 0.99))].(int64),
            Max: group[totalTc - 1].(int64),
            Mean: mean,
            StdDev: stddev,
        }
    }
    
    return &performanceGauge
}



