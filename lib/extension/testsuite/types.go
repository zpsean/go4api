/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package testsuite

import (
    "go4api/lib/testcase"
)

type TestSuites []*TestSuite

// test suite type,
type TestSuite map[string]*TestSuiteBasics

type TestSuiteBasics struct {
    Priority string            `json:"priority"`      // is the highest of the testcases included (i.e. get p1 if has p1, p2, p9)
    Description string         `json:"description"`   
    TestCasePaths []string     `json:"testCasePaths"` // has highp riority than attribute TestCases
    OriginalTestCases []string      `json:"originalTestCases"`
    AnalyzedTestCases []*testcase.TestCaseDataInfo      `json:"analyzedTestCases"`
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

