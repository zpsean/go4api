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
 	"encoding/json"

	// "go4api/lib/testcase"
	"go4api/texttmpl"
    "go4api/ui/js" 

	. "github.com/ahmetb/go-linq"
)

type ReportsStats struct {
    ReportKey interface{}
    Count int
}

var (
	// tcStats TcStats
)

func GenerateStatsJs (tcReportSlice TcReportSlice, resultsDir string) {
    statsFile := resultsDir + "/js/stats.js"

    reJsons := tcReportSlice.GetStatsJson()
    texttmpl.GenerateStatsJs(js.Stats, statsFile, reJsons)
}

func (tcReportSlice TcReportSlice) GetStatsJson () []string {
    var reJsons []string

    // for index.html details stats
    statsGaugeJson := tcReportSlice.GetStatsGaugeJson()
    reJsons = append(reJsons, statsGaugeJson)

    reportsStatsTotalSlice, reJsonTotal := tcReportSlice.GetTotalStatsJson()
    reJsons = append(reJsons, reJsonTotal)

    reJsonSuccess := tcReportSlice.GetSuccessStatsJson(reportsStatsTotalSlice)
    reJsons = append(reJsons, reJsonSuccess)

    reJsonFail := tcReportSlice.GetFailStatsJson(reportsStatsTotalSlice)
    reJsons = append(reJsons, reJsonFail)

    // for pie chart
    reJsonOverallStatus := tcReportSlice.GetOverallStatusStatsJson()
    reJsons = append(reJsons, reJsonOverallStatus)

    return reJsons
}

// ---
func (tcReportSlice TcReportSlice) GetTotalStatsJson () ([]ReportsStats, string) {
    query := tcReportSlice.GroupByTotalStartTime()

    reportsStatsTotalSlice := PrintStatsGroup(query)
    reJson, _ := json.MarshalIndent(reportsStatsTotalSlice, "", "\t")

    return reportsStatsTotalSlice, string(reJson)
}


func (tcReportSlice TcReportSlice) GetSuccessStatsJson (reportsStatsTotalSlice []ReportsStats) string {
    query := tcReportSlice.GroupBySuccessStartTime()

    reportsStatsSuccessSlice := PrintStatsGroup(query)
    reportsStatsSuccessSliceOrdered := ToOrderStatsGroup(reportsStatsTotalSlice, reportsStatsSuccessSlice)

    reJson, _ := json.MarshalIndent(reportsStatsSuccessSliceOrdered, "", "\t")

    return string(reJson)
}


func (tcReportSlice TcReportSlice) GetFailStatsJson (reportsStatsTotalSlice []ReportsStats) string {
    query := tcReportSlice.GroupByFailStartTime()

    reportsStatsFailSlice := PrintStatsGroup(query)
    reportsStatsFailSliceOrdered := ToOrderStatsGroup(reportsStatsTotalSlice, reportsStatsFailSlice)

    reJson, _ := json.MarshalIndent(reportsStatsFailSliceOrdered, "", "\t")

    return string(reJson)
}

// --
func (tcReportSlice TcReportSlice) GetOverallStatusStatsJson () string {
    query := tcReportSlice.GroupByOverallStatus()

    reportsOverallStatusSlice := PrintStatsGroup(query)

    reJson, _ := json.MarshalIndent(reportsOverallStatusSlice, "", "\t")

    return string(reJson)
}

func PrintStatsGroup (query []Group) []ReportsStats {
    var reportsStatsSlice []ReportsStats

    for _, q := range query {
        reportsStats := ReportsStats {
            ReportKey: q.Key,
            Count: len(q.Group),
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




