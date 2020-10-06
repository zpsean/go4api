/*
 * go4api - an api testing tool written in Go
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
    TestCase         *TestCase
    JsonFilePath     string
    CsvFile          string
    CsvRow           string
    MutationArea     string
    MutationCategory string
    MutationRule     string
    MutationInfoStr  interface{}
    MutationInfo     MutationInfo
}

// test case execution type, includes testdata
type TestCaseExecutionInfo struct {
    TestCaseDataInfo  *TestCaseDataInfo
    SetUpResult       string
    SetUpTestMessages [][]*TestMessage
    HttpResult        string
    ActualStatusCode  int
    StartTime         string
    EndTime           string
    HttpTestMessages  []*TestMessage
    StartTimeUnixNano int64
    EndTimeUnixNano   int64
    DurationUnixNano  int64
    ActualBody        []byte
    ActualHeader      map[string][]string
    HttpUrl           string
    TearDownResult    string
    TearDownTestMessages [][]*TestMessage
    TestResult           string  // Ready, Running, Success, Fail, ParentReady, ParentRunning, ParentFailed
    LocalVariables       interface{}
}

type TestMessage struct {  
    AssertionResults string
    ReponsePart      string // Status, Headers, Body
    FieldName        interface{}
    AssertionKey     interface{}
    ActualValue      interface{}
    ExpValue         interface{}
}

//
type TestCases []TestCase

// test case type,
type TestCase map[string]*TestCaseBasics

type TestCaseBasics struct {
    Priority       string                      `json:"priority"`
    ParentTestCase string                      `json:"parentTestCase"`
    FunctionAreas  []string                    `json:"functionAreas"`
    TestSuite      string                      `json:"testSuite"`
    IfGlobalSetUpTestCase    bool              `json:"ifGlobalSetUpTestCase"`
    IfGlobalTearDownTestCase bool              `json:"ifGlobalTearDownTestCase"`
    SetUp       []*CommandDetails              `json:"setUp"`
    Inputs      []interface{}                  `json:"inputs"`
    Request     *Request                       `json:"request"`
    Response    []map[string]interface{}       `json:"response"`
    Outputs     []*OutputsDetails              `json:"outputs"`
    OutFiles    []*OutFilesDetails             `json:"outFiles"`
    OutGlobalVariables map[string]interface{}  `json:"outGlobalVariables"`
    OutLocalVariables  map[string]interface{}  `json:"outLocalVariables"`
    Session   map[string]interface{}           `json:"session"`
    TearDown  []*CommandDetails                `json:"tearDown"`
}

//
type Request struct {  
    Method      string                  `json:"method"`
    Path        string                  `json:"path"`
    Headers     map[string]interface{}  `json:"headers"`
    QueryString map[string]interface{}  `json:"queryString"`
    Payload     map[string]interface{}  `json:"payload"`
}

type Payload struct {
    TextJson       interface{}             `json:"textJson"`
    Text           interface{}             `json:"text"`
    MultipartForm  *MultipartForm          `json:"multipartForm"`
    Form           map[string]interface{}  `json:"form"`
}

type MultipartForm struct {
    Name        string                 `json:"name"`
    Value       string                 `json:"value"`
    Type        string                 `json:"type"`
    MIMEHeader  map[string]interface{} `json:"mIMEHeader"`
}


// --
// type Response struct {  
//     Status  map[string]interface{}  `json:"status"`
//     Headers map[string]interface{}  `json:"headers"`
//     Body    map[string]interface{}  `json:"body"`
// }
// type Response []map[string]interface{}

type OutputsDetails struct {
    FileName string
    Format   string
    Data     map[string][]interface{}
}

type OutFilesDetails struct {
    TargetFile    string
    TargetHeader  []string
    Sources       []string
    SourcesFields []string
    Operation     string
    Data          map[string][]interface{}
}

type CommandDetails struct {
    CmdType            string                    `json:"cmdType"`
    CmdSource          string                    `json:"cmdSource"`
    Cmd                interface{}               `json:"cmd"`
    CmdResponse        []map[string]interface{}  `json:"cmdResponse"`
    OutGlobalVariables map[string]interface{}    `json:"outGlobalVariables"`
    OutLocalVariables  map[string]interface{}    `json:"outLocalVariables"`
    Session            map[string]interface{}    `json:"session"`
    OutFiles           []*OutFilesDetails        `json:"outFiles"`
}

type MutationInfo struct {
    FieldPath    []string
    CurrValue    interface{}
    FieldType    string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
    MutatedValue interface{}
}

// for report format 
type TcReportResults struct { 
    TcName                string 
    IfGlobalSetUpTearDown string // SetUp, TearDown
    // CaseType string // Normal, Scenario, Mutation, Fuzz
    Priority       string
    ParentTestCase string
    FunctionAreas  []string
    TestSuite      string   
    SetUpResult     string // Success, Fail
    SetUpTestMessages [][]*TestMessage
    Path          string
    Method        string
    JsonFilePath  string
    CsvFile       string
    CsvRow        string
    MutationArea  string
    MutationCategory string
    MutationRule     string
    MutationInfo     interface{}
    HttpResult       string // Success, Fail
    ActualStatusCode int
    StartTime        string
    EndTime          string
    HttpTestMessages []*TestMessage
    StartTimeUnixNano int64
    EndTimeUnixNano   int64
    DurationUnixNano  int64
    DurationUnixMillis   int64
    TearDownResult       string // Success, Fail
    TearDownTestMessages [][]*TestMessage
    TestResult           string  // Success, Fail, ParentFailed
    HttpUrl              string
    CaseOrigin           interface{}
    GlobalVariables      interface{}
    Session              interface{}
    LocalVariables interface{}
    ActualHeader   interface{}
    ActualBody     interface{}
}


type TcConsoleResults struct { 
    TcName         string 
    Priority       string
    ParentTestCase string
    JsonFilePath   string
    CsvFile        string
    CsvRow         string
    MutationInfoStr interface{}
    SetUpResult     string
    HttpResult      string  
    TearDownResult  string
    TestResult      string // Ready, Running, Success, Fail, ParentReady, ParentRunning, ParentFailed
    ActualStatusCode int
    HttpTestMessages []*TestMessage
}

