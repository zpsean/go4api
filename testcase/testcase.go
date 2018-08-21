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
    // "sync"
    "net/url" 
    // "go4api/types" 
)

// test case type
func (tc TestCase) TcName() string {
    var tcName string
    for key, _ := range tc {
        tcName = key
    }
    return tcName
}

func (tc TestCase) Priority() string {
    var priority string
    for _, value := range tc {
        priority = value.Priority
    }
    return priority
}

func (tc TestCase) ParentTestCase() string {
    var parentTestCase string
    for _, value := range tc {
        parentTestCase = value.ParentTestCase
    }
    return parentTestCase
}


// test case data type
func (tcData TestCaseDataInfo) TcName() string {
    return tcData.TestCase.TcName()
}

func (tcData TestCaseDataInfo) Priority() string {
    return tcData.TestCase.Priority()
}

func (tcData TestCaseDataInfo) ParentTestCase() string {
    return tcData.TestCase.ParentTestCase()
}





// test case execution type
func (tcExecution TestCaseExecutionInfo) TcName() string {
    return tcExecution.TestCaseDataInfo.TcName()
}

func (tcExecution TestCaseExecutionInfo) Priority() string {
    return tcExecution.TestCaseDataInfo.Priority()
}

func (tcExecution TestCaseExecutionInfo) ParentTestCase() string {
    return tcExecution.TestCaseDataInfo.ParentTestCase()
}


func (tcExecution TestCaseExecutionInfo) TestCase() TestCase {
    return tcExecution.TestCaseDataInfo.TestCase
}


func (tcExecution TestCaseExecutionInfo) SetTestResult(value string) {
    tcExecution.TestResult = value
}


// type Request struct {  
//     Method string
//     Path string
//     Inputs string
//     Headers map[string]interface{}
//     QueryString map[string]interface{}
//     Payload map[string]interface{}
// }

func (tc TestCase) ReqMethod() string {
    var reqMethod string
    for _, value := range tc {
        reqMethod = value.Request.Method
    }
    return reqMethod
}

func (tc TestCase) ReqPath() string {
    var reqPath string
    for _, value := range tc {
        reqPath = value.Request.Path
    }
    return reqPath
}

func (tc TestCase) ReqInputs() string {
    var reqInputs string
    for _, value := range tc {
        reqInputs = value.Request.Inputs
    }
    return reqInputs
}

func (tc TestCase) ReqHeaders() map[string]interface{} {
    var reqHeaders map[string]interface{}
    for _, value := range tc {
        reqHeaders = value.Request.Headers
    }
    return reqHeaders
}

func (tc TestCase) ReqQueryString() string {
    var reqQueryString string
    i := 0
    for _, value := range tc {
        for qryKey, qryValue := range value.Request.QueryString {
            if i == 0 {
                reqQueryString = fmt.Sprint(qryKey) + "=" + fmt.Sprint(qryValue)
            } else  {
                reqQueryString = reqQueryString + "&" + fmt.Sprint(qryKey) + "=" + fmt.Sprint(qryValue)
            }
            i = i + 1
        }
    }
    return reqQueryString
}

// to encode the query, also avoid the impact if string itself contains char '&'
func (tc TestCase) ReqQueryStringEncode() string {
    var reqQueryStringEncoded string
    i := 0
    for _, value := range tc {
        for qryKey, qryValue := range value.Request.QueryString {
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
    }
    return reqQueryStringEncoded
}


func (tc TestCase) ReqPayload() map[string]interface{} {
    var reqPayload map[string]interface{}
    for _, value := range tc {
        reqPayload = value.Request.Payload
    }
    return reqPayload
}


// type Response struct {  
//     Status map[string]interface{}
//     Headers map[string]interface{}
//     Body map[string]interface{}
//     Outputs []interface{}
// }


func (tc TestCase) RespStatus() map[string]interface{} {
    var reqStatus map[string]interface{}
    for _, value := range tc {
        reqStatus = value.Response.Status
    }
    return reqStatus
}


func (tc TestCase) RespHeaders() map[string]interface{} {
    var reqHeaders map[string]interface{}
    for _, value := range tc {
        reqHeaders = value.Response.Headers
    }
    return reqHeaders
}

func (tc TestCase) RespBody() map[string]interface{} {
    var reqBody map[string]interface{}
    for _, value := range tc {
        reqBody = value.Response.Body
    }
    return reqBody
}

func (tc TestCase) RespOutputs() []interface{} {
    var reqOutputs []interface{}
    for _, value := range tc {
        reqOutputs = value.Response.Outputs
    }
    return reqOutputs
}




