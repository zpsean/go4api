/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package testsuite

import (

)

type TestSuites []TestSuite

// test suite type,
type TestSuite map[string]*TestCaseBasics

type TestSuiteBasics struct {
    Priority string         `json:"priority"`   // is the highest of the testcases included (i.e. get p1 if has p1, p2, p9)
    TcRootPath string       `json:"tcRootPath"` // has highpriority than attribute TestCases
    TestCases []string      `json:"testCases"`
    Parameters map[string]interface{}   `json:"parameters"`
}

// for report format 
type TsReportResults struct { 
    TsName string 
    Priority string
    StartTime string
    EndTime string
    StartTimeUnixNano int64
    EndTimeUnixNano int64
    DurationUnixNano int64
    DurationUnixMillis int64
    TestResult string  // Success, Fail
}

type TsConsoleResults struct { 
    TsName string 
    Priority string
    StartTime string
    EndTime string
    StartTimeUnixNano int64
    EndTimeUnixNano int64
    DurationUnixNano int64
    DurationUnixMillis int64
    TestResult string  // Success, Fail
}

