/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package keyword

import (
    "go4api/lib/testcase"
)


type GKeyWords []*GKeyWord

// currently, supports TestCases, Settings, Keywords, Variables
type GKeyWord struct {
    Settings  *Settings
    TestCases *TestCases
    // Keywords  Keywords
    Variables *Variables
}

//
type Settings struct {
    StartLine          int
    EndLine            int
    OriginalContent    []string
    ID                 string
    TestSuitePaths     []string  // paths
    BasicTestCasePaths []string  // paths
    JsFuncPaths        []string  // paths
}

type TestCases struct {
    StartLine        int
    EndLine          int
    OriginalContent  []string
    ParsedTestCases  []*KWTestCase
}

type KWTestCase struct {
    OriginalLine            int
    OriginalTestCase        string
    KWTestCaseName          string   
    ParsedTestCase          []string   // format: ["tsName / tcNmae", "arg1 / v", "arg2 / v", ...]
    MappingToTestSuiteId    string     // 
    MappingToTestSuiteFile  string     // 
    MappingToBasicTestCase  *testcase.TestCaseDataInfo     // 
}

type Variables struct {
    StartLine        int
    EndLine          int
    OriginalContent  []string
    ParsedContent    interface{}
}

// for report format 
type KWTcReportResults struct { 
    KWName             string 
    StartTime          string
    EndTime            string
    StartTimeUnixNano  int64
    EndTimeUnixNano    int64
    DurationUnixNano   int64
    DurationUnixMillis int64
    TestResult         string  // Success, Fail
}

type KWTcConsoleResults struct { 
    KWName             string 
    StartTime          string
    EndTime            string
    StartTimeUnixNano  int64
    EndTimeUnixNano    int64
    DurationUnixNano   int64
    DurationUnixMillis int64
    TestResult         string  // Success, Fail
}

