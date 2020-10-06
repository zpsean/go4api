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
    "fmt"
    // "time"
    // "os"
    // "sort"
    "strings"
    "net/url" 
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

func (tc *TestCase) TestCaseBasics() *TestCaseBasics {
    return (*tc)[tc.TcName()]
}

func (tc *TestCase) Priority() string {
    return (*tc)[tc.TcName()].Priority
}

func (tc *TestCase) ParentTestCase() string {
    return (*tc)[tc.TcName()].ParentTestCase
}

func (tc *TestCase) FunctionAreas() []string {
    return (*tc)[tc.TcName()].FunctionAreas
}

func (tc *TestCase) TestSuite() string {
    return (*tc)[tc.TcName()].TestSuite
}

func (tc *TestCase) IfGlobalSetUpTestCase() bool {
    return (*tc)[tc.TcName()].IfGlobalSetUpTestCase
}

func (tc *TestCase) IfGlobalTearDownTestCase() bool {
    return (*tc)[tc.TcName()].IfGlobalTearDownTestCase
}

func (tc *TestCase) SetUp() []*CommandDetails {
    return (*tc)[tc.TcName()].SetUp
}

func (tc *TestCase) Response() []map[string]interface{} {
    return (*tc)[tc.TcName()].Response
}

func (tc *TestCase) Inputs() []interface{} {
    return (*tc)[tc.TcName()].Inputs
}

func (tc *TestCase) Outputs() []*OutputsDetails {
    return (*tc)[tc.TcName()].Outputs
}

func (tc *TestCase) OutFiles() []*OutFilesDetails {
    return (*tc)[tc.TcName()].OutFiles
}

func (tc *TestCase) OutGlobalVariables() map[string]interface{} {
    return (*tc)[tc.TcName()].OutGlobalVariables
}

func (tc *TestCase) OutLocalVariables() map[string]interface{} {
    return (*tc)[tc.TcName()].OutLocalVariables
}

func (tc *TestCase) Session() map[string]interface{} {
    return (*tc)[tc.TcName()].Session
}

func (tc *TestCase) TearDown() []*CommandDetails {
    return (*tc)[tc.TcName()].TearDown
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

func (tc *TestCase) SetTestSuite(newValue string) {
    (*tc)[tc.TcName()].TestSuite = newValue
}

func (tc *TestCase) SetInputs(newValue string) {
    (*tc)[tc.TcName()].Inputs = append((*tc)[tc.TcName()].Inputs, newValue)
}

// func (tc *TestCase) SetOutputs(newValue []interface{}) {
//     (*tc)[tc.TcName()].Outputs = newValue
// }

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
    reqH := tc.ReqHeaders()
    // if has key, value already
    if len(reqH) > 0 {
        (*tc)[tc.TcName()].Request.Headers[key] = newValue
    } else {
        // need to init the headers make(map[]), otherwise, assign to nil map
        h := make(map[string]interface{})
        (*tc)[tc.TcName()].Request.Headers = h
        (*tc)[tc.TcName()].Request.Headers[key] = newValue
    }
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


// request query Payload
func (tc *TestCase) SetRequestPayload (key string, newValue interface{}) {
    (*tc)[tc.TcName()].Request.Payload[key] = newValue
}


func (tc *TestCase) UpdateTcName (newKey string) {
    mTc := TestCase{}
    mTc[newKey] = (*tc)[tc.TcName()]

    delete(*tc, tc.TcName())
    (*tc)[newKey] = mTc[newKey]
}

// type Request struct
func (tc *TestCase) Request() *Request {
    return (*tc)[tc.TcName()].Request
}

func (tc *TestCase) ReqMethod() string {
    if (*tc)[tc.TcName()].Request == nil {
        return ""
    } else {
        return (*tc)[tc.TcName()].Request.Method
    }
}

func (tc *TestCase) ReqPath() string {
    if (*tc)[tc.TcName()].Request == nil {
        return ""
    } else {
        return (*tc)[tc.TcName()].Request.Path
    }
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

    // net/url's base url format is:
    // scheme://[userinfo@]host/path[?query][#fragment]
    // but if the ? is preceded by #, like: .../str1/#/str2?q=str3...
    // then all the query string will be truncated, then one issue happes
    // to avoid this issue, here set a replacer for #
    urlStr = strings.Replace(urlStr, "#", "go4Api_wregwvasdst_ReplacerKey", -1)

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

    urlStr = strings.Replace(urlStr, "go4Api_wregwvasdst_ReplacerKey", "#", -1)

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

    //
    urlStr = strings.Replace(urlStr, "#", "go4Api_wregwvasdst_ReplacerKey", -1)

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
    
    urlStr = strings.Replace(urlStr, "go4Api_wregwvasdst_ReplacerKey", "#", -1)

    return urlStr
}


func (tc *TestCase) ReqPayload() map[string]interface{} {
    return (*tc)[tc.TcName()].Request.Payload
}


func (tc *TestCase) DelReqPayload (key string) {
    delete((*tc)[tc.TcName()].Request.Payload, key)
}


// type Response struct
// func (tc *TestCase) RespStatus() map[string]interface{} {
//     if (*tc)[tc.TcName()].Response == nil {
//         return nil
//     } else {
//         return (*tc)[tc.TcName()].Response.Status
//     }
    
// }

// func (tc *TestCase) RespHeaders() map[string]interface{} {
//     if (*tc)[tc.TcName()].Response == nil {
//         return nil
//     } else {
//         return (*tc)[tc.TcName()].Response.Headers
//     }
// }

// func (tc *TestCase) RespBody() map[string]interface{} {
//     if (*tc)[tc.TcName()].Response == nil {
//         return nil
//     } else {
//         return (*tc)[tc.TcName()].Response.Body
//     }
// }

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

func (tcData *TestCaseDataInfo) FunctionAreas() []string {
    return tcData.TestCase.FunctionAreas()
}

func (tcData *TestCaseDataInfo) TestSuite() string {
    return tcData.TestCase.TestSuite()
}

func (tcData *TestCaseDataInfo) ReqMethod() string {
    return tcData.TestCase.ReqMethod()
}

func (tcData *TestCaseDataInfo) ReqPath() string {
    return tcData.TestCase.ReqPath()
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

func (tcExecution *TestCaseExecutionInfo) ReqMethod() string {
    return tcExecution.TestCaseDataInfo.ReqMethod()
}

func (tcExecution *TestCaseExecutionInfo) ReqPath() string {
    return tcExecution.TestCaseDataInfo.ReqPath()
}

func (tcExecution *TestCaseExecutionInfo) TestCase() *TestCase {
    return tcExecution.TestCaseDataInfo.TestCase
}

func (tcExecution *TestCaseExecutionInfo) SetTestResult(value string) {
    tcExecution.TestResult = value
}

//outputs
func (tcOutDetails *OutputsDetails) GetOutputsDetailsFileName() string {
    return (*tcOutDetails).FileName
}

func (tcOutDetails *OutputsDetails) GetOutputsDetailsFormat() string {
    return (*tcOutDetails).Format
}

func (tcOutDetails *OutputsDetails) GetOutputsDetailsData() map[string][]interface{} {
    return (*tcOutDetails).Data
}


//outFiles
func (tcOutFiles *OutFilesDetails) GetTargetFileName() string {
    return (*tcOutFiles).TargetFile
}

func (tcOutFiles *OutFilesDetails) GetTargetHeader() []string {
    return (*tcOutFiles).TargetHeader
}

func (tcOutFiles *OutFilesDetails) GetSources() []string {
    return (*tcOutFiles).Sources
}

func (tcOutFiles *OutFilesDetails) GetSourcesFields() []string {
    return (*tcOutFiles).SourcesFields
}

func (tcOutFiles *OutFilesDetails) GetOperation() string {
    return (*tcOutFiles).Operation
}

func (tcOutFiles *OutFilesDetails) GetData() map[string][]interface{} {
    return (*tcOutFiles).Data
}


