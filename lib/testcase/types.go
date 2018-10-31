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
    TestCase *TestCase
    JsonFilePath string
    CsvFile string
    CsvRow string
    MutationArea string
    MutationCategory string
    MutationRule string
    MutationInfoStr interface{}
    MutationInfo MutationInfo
}

// test case execution type, includes testdata
type TestCaseExecutionInfo struct {
    TestCaseDataInfo *TestCaseDataInfo
    SetUpResult string
    HttpResult string
    ActualStatusCode int
    StartTime string
    EndTime string
    TestMessages []*TestMessage
    StartTimeUnixNano int64
    EndTimeUnixNano int64
    DurationUnixNano int64
    ActualBody []byte
    TearDownResult string
    TestResult string  // Ready, Running, Success, Fail, ParentReady, ParentRunning, ParentFailed
}

type TestMessage struct {  
    AssertionResults string
    ReponsePart string // Status, Headers, Body
    FieldName interface{}
    AssertionKey  interface{}
    ActualValue  interface{}
    ExpValue interface{}
}

//
type TestCases []TestCase

// test case type,
type TestCase map[string]*TestCaseBasics

type TestCaseBasics struct {
    Priority string         `json:"priority"`
    ParentTestCase string   `json:"parentTestCase"`
    IfGlobalSetUpTestCase bool    `json:"ifGlobalSetUpTestCase"`
    IfGlobalTearDownTestCase bool `json:"ifGlobalTearDownTestCase"`
    SetUp map[string]interface{}  `json:"setUp"`
    Inputs []interface{}     `json:"inputs"`
    Request *Request         `json:"request"`
    Response *Response       `json:"response"`
    Outputs []*OutputsDetails   `json:"outputs"`
    OutEnvVariables map[string]interface{}   `json:"outEnvVariables"`
    Session map[string]interface{}           `json:"session"`
    TearDown map[string]interface{}          `json:"tearDown"`
}

type Request struct {  
    Method string                       `json:"method"`
    Path string                         `json:"path"`
    Headers map[string]interface{}      `json:"headers"`
    QueryString map[string]interface{}  `json:"queryString"`
    Payload map[string]interface{}      `json:"payload"`
}

type Response struct {  
    Status map[string]interface{}   `json:"status"`
    Headers map[string]interface{}  `json:"headers"`
    Body map[string]interface{}     `json:"body"`
}

type OutputsDetails struct {
    FileName string
    Format string
    Data map[string][]interface{}
}

type MutationInfo struct {
    FieldPath []string
    CurrValue interface{}
    FieldType string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
    MutatedValue interface{}
}

// for report format 
type TcReportResults struct { 
    TcName string 
    IfGlobalSetUpTearDown string // SetUp, TearDown
    // CaseType string // Normal, Scenario, Mutation, Fuzz
    Priority string
    ParentTestCase string
    SetUpResult string // Success, Fail
    Path string
    Method string
    JsonFilePath string
    CsvFile string
    CsvRow string
    MutationArea string
    MutationCategory string
    MutationRule string
    MutationInfo interface{}
    HttpResult string // Success, Fail
    ActualStatusCode int
    StartTime string
    EndTime string
    TestMessages []*TestMessage
    StartTimeUnixNano int64
    EndTimeUnixNano int64
    DurationUnixNano int64
    TearDownResult string // Success, Fail
    TestResult string  // Success, Fail, ParentFailed
    CaseOrigin interface{}
}


type TcConsoleResults struct { 
    TcName string 
    Priority string
    ParentTestCase string
    JsonFilePath string
    CsvFile string
    CsvRow string
    MutationInfoStr interface{}
    TestResult string  // Ready, Running, Success, Fail, ParentReady, ParentRunning, ParentFailed
    ActualStatusCode int
    TestMessages []*TestMessage
}

