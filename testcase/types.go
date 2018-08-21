/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package testcase

import (
)

// test case data type, includes testcase
type TestCaseDataInfo struct {
    TestCase TestCase
    JsonFilePath string
    CsvFile string
    CsvRow string
}

// test case execution type, includes testdata
type TestCaseExecutionInfo struct {
    TestCaseDataInfo TestCaseDataInfo
    TestResult string
    ActualStatusCode int
    StartTime string
    EndTime string
    TestMessages string
    StartTimeUnixNano int64
    EndTimeUnixNano int64
    DurationUnixNano int64
}
type TestCases struct {
    TestCases []TestCase
}

// test case type,
type TestCase map[string]TestCaseBasics

type TestCaseBasics struct {
    Priority string
    ParentTestCase string
    Request Request
    Response Response
}

type Request struct {  
    Method string
    Path string
    Inputs string
    Headers map[string]interface{}
    QueryString map[string]interface{}
    Payload map[string]interface{}
}


type Response struct {  
    Status map[string]interface{}
    Headers map[string]interface{}
    Body map[string]interface{}
    Outputs []interface{}
}



