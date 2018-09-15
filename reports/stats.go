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
 	"encoding/json"

	"go4api/lib/testcase"
	"go4api/texttmpl"
)


type TcStats struct {
	StatusStats map[string]map[string]int
	StatusStatsPercentage map[string]map[string]float64
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


func (tcStats *TcStats) PrepStatsReport () *texttmpl.StatsJs {
	statsJsonBytes, _ := json.MarshalIndent(tcStats, "", "\t")

	tcStatsReport := texttmpl.StatsJs {
		StatsStr: string(statsJsonBytes),
	}

	return &tcStatsReport
}


func GetStats (tcClassifedCountMap map[string]int, totalTc int, statusCountByPriority map[string]map[string]int, 
    tcExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo,
    tcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) TcStats {

	InitVariables(statusCountByPriority)

	for key, _ := range statusCountByPriority {
		for k, v := range statusCountByPriority[key] {
		    statusStats[key][k] = v

		    if len(tcClassifedCountMap) == 0 {
		    	statusStats[key]["TotalInSource"] = 0
		    } else {
		    	statusStats[key]["TotalInSource"] = statusCountByPriority[key]["Total"]
		    }
		    
		    if statusCountByPriority[key]["Total"] > 0 {
		    	statusStatsPercentage[key][k] = float64(v) / float64(statusCountByPriority[key]["Total"])
		    } else {
		    	statusStatsPercentage[key][k] = 0
			}
		}
    }

    tcStats = TcStats {
    	StatusStats: statusStats,
    	StatusStatsPercentage: statusStatsPercentage,
    }

    return tcStats
}

