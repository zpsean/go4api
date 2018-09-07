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
    "fmt"
    // "time"
    // "os"
    // "sort"
    "strings"
    "net/url" 
    "path/filepath"
)

// test case type - get
func (tc *TestCase) TcName() string {
    var tcName string
    for key, _ := range *tc {
        tcName = key
        break
    }
    return tcName
}

func (tc *TestCase) Priority() string {
    return (*tc)[tc.TcName()].Priority
}

func (tc *TestCase) ParentTestCase() string {
    return (*tc)[tc.TcName()].ParentTestCase
}

func (tc *TestCase) Inputs() string {
    return (*tc)[tc.TcName()].Inputs
}

func (tc *TestCase) Outputs() []interface{} {
    return (*tc)[tc.TcName()].Outputs
}

// !! ---------------------------------------
// !! --- test case type - set
// !! ---------------------------------------
func (tc *TestCase) SetPriority(newValue string) {
    (*tc)[tc.TcName()].Priority = newValue
}

func (tc *TestCase) SetParentTestCase(newValue string) {
    (*tc)[tc.TcName()].ParentTestCase = newValue
}

func (tc *TestCase) SetInputs(newValue string) {
    (*tc)[tc.TcName()].Inputs = newValue
}

func (tc *TestCase) SetOutputs(newValue []interface{}) {
    (*tc)[tc.TcName()].Outputs = newValue
}

func (tc *TestCase) SetRequestMethod (newValue string) {
    (*tc)[tc.TcName()].Request.Method = newValue
}

func (tc *TestCase) SetRequestPath (newValue string) {
    (*tc)[tc.TcName()].Request.Path = newValue
}

// request header
func (tc *TestCase) SetRequestHeader (key string, newValue string) {
    // to check if the headers has the key
    if _, ok := (*tc)[tc.TcName()].Request.Headers[key]; ok {
        (*tc)[tc.TcName()].Request.Headers[key] = newValue
    }
}

func (tc *TestCase) AddRequestHeader (key string, newValue string) {
    reqH := tc.ReqQueryString()
    // if has key, value already
    if len(reqH) > 0 {
        (*tc)[tc.TcName()].Request.Headers[key] = newValue
    }
    // } else {
        // need to init the headers make(map[]), otherwise, assign to nil map
    // }
}

func (tc *TestCase) DelRequestHeader (key string) {
    delete((*tc)[tc.TcName()].Request.Headers, key)
}

// request query string
func (tc *TestCase) SetRequestQueryString (key string, newValue string) {
    // to check if the QueryString has the key
    if _, ok := (*tc)[tc.TcName()].Request.QueryString[key]; ok {
        (*tc)[tc.TcName()].Request.QueryString[key] = newValue
    }
}

func (tc *TestCase) AddRequestQueryString (key string, newValue string) {
    // check if tc has QueryString
    reqQS := tc.ReqQueryString()
    if len(reqQS) > 0 {
        (*tc)[tc.TcName()].Request.QueryString[key] = newValue
    }
}

func (tc *TestCase) DelRequestQueryString (key string) {
    delete((*tc)[tc.TcName()].Request.QueryString, key)
}


// request query Payload??
// Note: currently, if the post data is json, then the key is "text"
func (tc *TestCase) SetRequestPayload (key string, newValue string) {
    (*tc)[tc.TcName()].Request.Payload[key] = newValue
}


func (tc *TestCase) UpdateTcName (newKey string) {
    mTc := TestCase{}
    mTc[newKey] = (*tc)[tc.TcName()]

    delete(*tc, tc.TcName())
    (*tc)[newKey] = mTc[newKey]
}


// type Request struct
func (tc *TestCase) ReqMethod() string {
    return (*tc)[tc.TcName()].Request.Method
}

func (tc *TestCase) ReqPath() string {
    return (*tc)[tc.TcName()].Request.Path
}

func (tc *TestCase) ReqHeaders() map[string]interface{} {
    return (*tc)[tc.TcName()].Request.Headers
}


func (tc *TestCase) ReqQueryString() map[string]interface{} {
    return (*tc)[tc.TcName()].Request.QueryString
}


func (tc *TestCase) ComposeReqQueryString() string {
    var reqQueryString string
    i := 0
    for qryKey, qryValue := range tc.ReqQueryString() {
        if i == 0 {
            reqQueryString = fmt.Sprint(qryKey) + "=" + fmt.Sprint(qryValue)
        } else  {
            reqQueryString = reqQueryString + "&" + fmt.Sprint(qryKey) + "=" + fmt.Sprint(qryValue)
        }
        i = i + 1
    }
    return reqQueryString
}

// to encode the query, also avoid the impact if string itself contains char '&'
func (tc *TestCase) ComposeReqQueryStringEncode() string {
    var reqQueryStringEncoded string
    i := 0
    for qryKey, qryValue := range tc.ReqQueryString() {
        if i == 0 {
            reqQueryString := fmt.Sprint(qryKey) + "=" + fmt.Sprint(qryValue)
            values, _ := url.ParseQuery(reqQueryString)
            reqQueryStringEncoded = values.Encode()
        } else  {
            reqQueryString := fmt.Sprint(qryKey) + "=" + fmt.Sprint(qryValue)
            values, _ := url.ParseQuery(reqQueryString)
            reqQueryStringEncoded = reqQueryStringEncoded + "&" + values.Encode()
        }
        i = i + 1
    }
    return reqQueryStringEncoded
}


func (tc *TestCase) UrlEncode(baseUrl string) string {
    urlStr := ""
    apiPath := tc.ReqPath()

    if strings.HasPrefix(strings.ToLower(apiPath), "http") {
        urlStr = apiPath
    } else {
        urlStr = baseUrl + apiPath
    }

    reqQueryStringEncoded := tc.ComposeReqQueryStringEncode()

    u, _ := url.Parse(urlStr)
    urlEncodedQry := u.Query().Encode()
    if len(urlEncodedQry) > 0 && len(reqQueryStringEncoded) > 0 {
        urlStr = u.Scheme + "://" + u.Host + "" + u.Path + "?" + urlEncodedQry + "&" + reqQueryStringEncoded
    } else if len (urlEncodedQry) > 0 {
        urlStr = u.Scheme + "://" + u.Host + "" + u.Path + "?" + urlEncodedQry
    } else if len (reqQueryStringEncoded) > 0 {
        urlStr = u.Scheme + "://" + u.Host + "" + u.Path + "?" + reqQueryStringEncoded
    } else {
        urlStr = u.Scheme + "://" + u.Host + "" + u.Path
    }
    return urlStr
}


func (tc *TestCase) UrlRaw(baseUrl string) string {
    urlStr := ""
    apiPath := tc.ReqPath()

    if strings.HasPrefix(strings.ToLower(apiPath), "http") {
        urlStr = apiPath
    } else {
        urlStr = baseUrl + apiPath
    }

    reqQueryStringRaw := tc.ComposeReqQueryString()

    u, _ := url.Parse(urlStr)
    urlQry := u.RawQuery
    if len(urlQry) > 0 && len(reqQueryStringRaw) > 0 {
        urlStr = u.Scheme + "://" + u.Host + "" + u.Path + "?" + urlQry + "&" + reqQueryStringRaw
    } else if len (urlQry) > 0 {
        urlStr = u.Scheme + "://" + u.Host + "" + u.Path + "?" + urlQry
    } else if len (reqQueryStringRaw) > 0 {
        urlStr = u.Scheme + "://" + u.Host + "" + u.Path + "?" + reqQueryStringRaw
    } else {
        urlStr = u.Scheme + "://" + u.Host + "" + u.Path
    }
    return urlStr
}


func (tc *TestCase) ReqPayload() map[string]interface{} {
    return (*tc)[tc.TcName()].Request.Payload
}


func (tc *TestCase) DelReqPayload (key string) {
    delete((*tc)[tc.TcName()].Request.Payload, key)
}


// type Response struct
func (tc *TestCase) RespStatus() map[string]interface{} {
    return (*tc)[tc.TcName()].Response.Status
}

func (tc *TestCase) RespHeaders() map[string]interface{} {
    return (*tc)[tc.TcName()].Response.Headers
}

func (tc *TestCase) RespBody() map[string]interface{} {
    return (*tc)[tc.TcName()].Response.Body
}

// !! ---------------------------------------
// !! --- test case data type 
// !! ---------------------------------------
func (tcData *TestCaseDataInfo) TcName() string {
    return tcData.TestCase.TcName()
}

func (tcData *TestCaseDataInfo) Priority() string {
    return tcData.TestCase.Priority()
}

func (tcData *TestCaseDataInfo) ParentTestCase() string {
    return tcData.TestCase.ParentTestCase()
}


// test case execution type
func (tcExecution *TestCaseExecutionInfo) TcName() string {
    return tcExecution.TestCaseDataInfo.TcName()
}

func (tcExecution *TestCaseExecutionInfo) Priority() string {
    return tcExecution.TestCaseDataInfo.Priority()
}

func (tcExecution *TestCaseExecutionInfo) ParentTestCase() string {
    return tcExecution.TestCaseDataInfo.ParentTestCase()
}

func (tcExecution *TestCaseExecutionInfo) TestCase() *TestCase {
    return tcExecution.TestCaseDataInfo.TestCase
}

func (tcExecution *TestCaseExecutionInfo) SetTestResult(value string) {
    tcExecution.TestResult = value
}


// for report
func (tcExecution *TestCaseExecutionInfo) TcConsoleResults() *TcConsoleResults {
    tcConsoleRes := &TcConsoleResults { 
        TcName: tcExecution.TcName(),
        Priority: tcExecution.Priority(),
        ParentTestCase: tcExecution.ParentTestCase(),
        JsonFilePath: filepath.Base(tcExecution.TestCaseDataInfo.JsonFilePath),
        CsvFile: filepath.Base(tcExecution.TestCaseDataInfo.CsvFile),
        CsvRow: tcExecution.TestCaseDataInfo.CsvRow,
        MutationInfo: tcExecution.TestCaseDataInfo.MutationInfo,
        TestResult: tcExecution.TestResult,
        ActualStatusCode: tcExecution.ActualStatusCode,
        TestMessages: tcExecution.TestMessages,
    }

    return tcConsoleRes
}


func (tcExecution *TestCaseExecutionInfo) TcReportResults() *TcReportResults {
    tcReportRes := &TcReportResults { 
        TcName: tcExecution.TcName(),
        Priority: tcExecution.Priority(),
        ParentTestCase: tcExecution.ParentTestCase(),
        JsonFilePath: tcExecution.TestCaseDataInfo.JsonFilePath,
        CsvFile: tcExecution.TestCaseDataInfo.CsvFile,
        CsvRow: tcExecution.TestCaseDataInfo.CsvRow,
        MutationInfo: tcExecution.TestCaseDataInfo.MutationInfo,
        TestResult: tcExecution.TestResult,
        ActualStatusCode: tcExecution.ActualStatusCode,
        StartTime: tcExecution.StartTime,
        EndTime: tcExecution.EndTime,
        TestMessages: tcExecution.TestMessages,
        StartTimeUnixNano: tcExecution.StartTimeUnixNano,
        EndTimeUnixNano: tcExecution.EndTimeUnixNano,
        DurationUnixNano: tcExecution.DurationUnixNano,
    }

    return tcReportRes
}





