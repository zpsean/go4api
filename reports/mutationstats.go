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
 	"encoding/json"
	"go4api/texttmpl"

	"go4api/lib/testcase"

    . "github.com/ahmetb/go-linq"
)


func GetMutationStats () []*MutationStats{
	var mutationStatsSlice []*MutationStats

	for i, _ := range ExecutionResultSlice {
	    tcReportRes := &MutationStats { 
	    	HttpUrl: ExecutionResultSlice[i].Path,
			HttpMethod: ExecutionResultSlice[i].Method,
			MutationPart: "Headers",
			MutationRule: ExecutionResultSlice[i].MutationRule,
			HttpStatus: ExecutionResultSlice[i].ActualStatusCode,
			TestStatus: ExecutionResultSlice[i].TestResult,
		}

		mutationStatsSlice = append(mutationStatsSlice, tcReportRes)
	}

	return mutationStatsSlice
}


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


func GetMutationDetailsJson (tcClassifedCountMap map[string]int, totalTc int, statusCountByPriority map[string]map[string]int, 
    tcExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo,
    tcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) {
	
	//-----
	//Exclude duplicates.
    var noduplicates []*testcase.TcReportResults

    type AA struct {
    	HttpUrl string
    	HttpMethod string
    	MutationPart string
    	MutationCategory string
    	// MutationRule string
    	HttpStatus int
    	TestStatus string
    }


    From(ExecutionResultSlice).
        DistinctByT(
            func(item *testcase.TcReportResults) AA { return AA{item.Path, item.Method, item.MutationArea, item.MutationCategory, item.ActualStatusCode, item.TestResult } },
        ).
        ToSlice(&noduplicates)

    for _, item := range noduplicates {
        // fmt.Println(product.Path, product.Method, product.ActualStatusCode, product.TestResult)
        fmt.Println(item.Path, item.Method, item.MutationArea, item.MutationCategory, item.ActualStatusCode, item.TestResult )
    }
}


func GetMutationStatsJson2 (tcClassifedCountMap map[string]int, totalTc int, statusCountByPriority map[string]map[string]int, 
    tcExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo,
    tcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo) {
	
	//-----
	//Exclude duplicates.
    // var noduplicates []*testcase.TcReportResults
    var query []Group

// HttpUrl	HttpMethod	MutationPart	MutationCategory	MutationRule	HttpStatus	TestStatus	Count	MutationMessage

    type AA struct {
    	HttpUrl string
    	HttpMethod string
    	MutationPart string
    	MutationCategory string
    	// MutationRule string
    	HttpStatus int
    	TestStatus string
    }


    // From(ExecutionResultSlice).
    //     DistinctByT(
    //         func(item *testcase.TcReportResults) int { return item.ActualStatusCode },
    //     ).
    //     ToSlice(&noduplicates)

    From(ExecutionResultSlice).GroupByT(
    	func(item *testcase.TcReportResults) AA { return AA{item.Path, item.Method, item.MutationArea, item.MutationCategory, item.ActualStatusCode, item.TestResult }},
	    func(item *testcase.TcReportResults) int { return 1 },
	    
	    
	).ToSlice(&query)

    // for _, product := range noduplicates {
    //     fmt.Println(product.Path, product.Method, product.ActualStatusCode, product.TestResult)
    // }
    for _, petGroup := range query {
	    // fmt.Printf("%d\n", petGroup.Key)
	    fmt.Println(petGroup.Key)
	    ii := 0
	    for range petGroup.Group {
	    	ii += 1
	        // fmt.Printf("  %s\n", petName)

	    }
	    fmt.Println("ii: ", ii)
	}
}



