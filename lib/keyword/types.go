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
    // "go4api/lib/testcase"
)


type KWBlocks []*KWBlock

type KWBlock struct {
    StartLine        int
    EndLine          int
    OriginalContent  []string
    ParsedContent    interface{}
    BlockType        string  // currently, supports TestCases, Settings, Keywords, Variables
}

type KWTestCase struct {
    OriginalLine     int
    OriginalTestCase string
    ParsedTestCase   []string   // format: ["tsName / tcNmae", "arg1 / v", "arg2 / v", ...]
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

