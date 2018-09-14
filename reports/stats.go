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

)

// type ExecutionStatusCount map[string]*ByPriorityDetails

// type ByPriorityDetails struct {
// 	Priority string
// 	TcExecutedList []testcase.TestCaseExecutionInfo
// 	TcNotExecutedList []testcase.TestCaseExecutionInfo
// }


type ByPriorityReports struct {
	Priority string
	ByStatusCount map[string]int
	TcCount int
	TcExecutionDurationMax float64
	TcExecutionDurationMin float64
}

// each priority
// status count
