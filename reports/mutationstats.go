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
	"go4api/texttmpl"

	"go4api/lib/testcase"
)


func GetMutationStats () []*MutationStats{
	var mutationStatsSlice []*MutationStats

	for i, _ := range ExecutionResultSlice {
	    tcReportRes := &MutationStats { 
	    	HttpUrl: ExecutionResultSlice[i].Path,
			HttpMethod: ExecutionResultSlice[i].Method,
			MutationPart: "Headers",
			MutationMethod: ExecutionResultSlice[i].MutationMethod,
			HttpStatus: ExecutionResultSlice[i].ActualStatusCode,
			TestStatus: ExecutionResultSlice[i].TestResult,
		}

		mutationStatsSlice = append(mutationStatsSlice, tcReportRes)
	}

	return mutationStatsSlice
}

// func GetMutationDetails () *MutationDetails {
// 	for i, _ := ExecutionResultSlice {
// 	    tcReportRes := &MutationDetails { 
// 	    	HttpUrl: ExecutionResultSlice[i].ReqPath(),
// 			HttpMethod: ExecutionResultSlice[i].ReqMethod(),
// 			MutationPart: "Headers",
// 			MutationMethod: ExecutionResultSlice[i].MutationMethod,
// 			HttpStatus: ExecutionResultSlice[i].ActualStatusCode,
// 			TestStatus: ExecutionResultSlice[i].TestResult,
// 			MutationMessage: ExecutionResultSlice[i].MutationInfo,
// 			TestMessages: ExecutionResultSlice[i].TestMessages,
// 		}

// 	return tcReportRes
// }


func GetMutationStatsJson (tcClassifedCountMap map[string]int, totalTc int, statusCountByPriority map[string]map[string]int, 
    tcExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo,
    tcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) *texttmpl.StatsJs {
	
	//-----
	mutationStatsSlice := GetMutationStats()

	statsJsonBytes, _ := json.MarshalIndent(mutationStatsSlice, "", "\t")

	tcStatsReport := texttmpl.StatsJs {
		StatsStr: string(statsJsonBytes),
	}

	return &tcStatsReport
}

