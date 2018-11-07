/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
    // "fmt"
    // "time"
    // "sync"
    // "reflect"
    // "net/http"  
    "strings"
    // "encoding/json"

    // "go4api/cmd"
    "go4api/lib/testcase"                                                                                                                             
    "go4api/assertion"
    g4http "go4api/protocal/http"
    // "go4api/sql"
)


func CallHttp(baseUrl string, tcData testcase.TestCaseDataInfo) (int, map[string][]string, []byte) {
    // urlStr := tcData.TestCase.UrlRaw(baseUrl)
    urlStr := tcData.TestCase.UrlEncode(baseUrl)
    //
    apiMethodSelector, apiMethod, bodyText, bodyMultipart, boundary := GetPayloadInfo(tcData)
    //
    reqHeaders := make(map[string]interface{})
    reqHeaders = tcData.TestCase.ReqHeaders()
    // set the boundary to headers, if multipart
    if boundary != "" {
        reqHeaders["Content-Type"] = boundary
    }

    // < !! ----------- !! >
    // (1). Actual response
    var actualStatusCode int
    var actualHeader map[string][]string
    var actualBody []byte
    // 
    httpRequest := g4http.HttpRestful{}
    if apiMethodSelector == "POSTMultipart" {
        actualStatusCode, actualHeader, actualBody = httpRequest.Request(urlStr, apiMethod, reqHeaders, bodyMultipart)    
    } else {
        actualStatusCode, actualHeader, actualBody = httpRequest.Request(urlStr, apiMethod, reqHeaders, bodyText)
        }

    return actualStatusCode, actualHeader, actualBody
}


func Compare(tcName string, actualStatusCode int, actualHeader map[string][]string, actualBody []byte, 
        expStatus map[string]interface{}, expHeader map[string]interface{}, expBody map[string]interface{}) (string, []*testcase.TestMessage) {
    //
    var testResults []bool
    var testMessages []*testcase.TestMessage
    // status
    testResultsS, testMessagesS := CompareStatus(actualStatusCode, expStatus)
    testResults = append(testResults, testResultsS[0:]...)
    testMessages = append(testMessages, testMessagesS[0:]...)
    // headers
    testResultsH, testMessagesH := CompareHeaders(actualHeader, expHeader)
    testResults = append(testResults, testResultsH[0:]...)
    testMessages = append(testMessages, testMessagesH[0:]...)
    // body
    testResultsB, testMessagesB := CompareBody(actualBody, expBody)
    testResults = append(testResults, testResultsB[0:]...)
    testMessages = append(testMessages, testMessagesB[0:]...)

    // default finalResults
    finalResults := "Success"

    for key := range testResults {
        if testResults[key] == false {
            finalResults = "Fail"
            break
        }
    }
    // testMessagesJson, _ := json.Marshal(testMessages)
    // testMessagesJsonStr := string(testMessagesJson)
    
    return finalResults, testMessages
} 


func CompareStatus(actualStatusCode int, expStatus map[string]interface{}) ([]bool, []*testcase.TestMessage) {
    var testResults []bool
    var testMessages []*testcase.TestMessage
    // status
    for assertionKey, expValue := range expStatus {
        actualValue := actualStatusCode
        key := "status"

        testRes, msg := compareCommon("Status", key, assertionKey, actualValue, expValue)
        
        testMessages = append(testMessages, msg)
        testResults = append(testResults, testRes)
    }

    return testResults, testMessages
} 

func CompareHeaders(actualHeader map[string][]string, expHeader map[string]interface{}) ([]bool, []*testcase.TestMessage) {
    var testResults []bool
    var testMessages []*testcase.TestMessage
    // headers
    for key, value := range expHeader {
        expHeader_sub := value.(map[string]interface{})
        //
        for assertionKey, expValue := range expHeader_sub {
            actualValue := strings.Join(actualHeader[key], ",")

            testRes, msg := compareCommon("Headers", key, assertionKey, actualValue, expValue)

            testMessages = append(testMessages, msg)
            testResults = append(testResults, testRes)
        } 
    }

    return testResults, testMessages
} 

func CompareBody(actualBody []byte, expBody map[string]interface{}) ([]bool, []*testcase.TestMessage) {
    var testResults []bool
    var testMessages []*testcase.TestMessage
    // body
    for key, value := range expBody {
        expBody_sub := value.(map[string]interface{})

        for assertionKey, expValue := range expBody_sub {
            // if path, then value - value, otherwise, key - value
            actualValue := GetActualValueByJsonPath(key, actualBody)
            
            testRes, msg := compareCommon("Body", key, assertionKey, actualValue, expValue)

            testMessages = append(testMessages, msg)
            testResults = append(testResults, testRes)
        }
    }

    return testResults, testMessages
} 

func compareCommon (reponsePart string, key string, assertionKey string, actualValue interface{}, expValue interface{}) (bool, *testcase.TestMessage) {
    // Note: As get Go nil, for JSON null, need special care, two possibilities:
    // p1: expResult -> null, but can not find out actualValue, go set it to nil, i.e. null (assertion -> false)
    // p2: expResult -> null, actualValue can be founc, and its value --> null (assertion -> true)
    // but here can not distinguish them
    assertionResults := ""
    var testRes bool

    if actualValue == nil || expValue == nil {
        // if only one nil
        if actualValue != nil || expValue != nil {
            assertionResults = "Failed"
            testRes = false
        // both nil
        } else {
            assertionResults = "Success"
            testRes = true
        }
    // no nil
    } else {
        // call the assertion function
        testResult := assertion.CallAssertion(assertionKey, actualValue, expValue)
        // fmt.Println("--->", key, assertionKey, actualValue, expValue, reflect.TypeOf(actualValue), reflect.TypeOf(expValue), testResult)
        if testResult == false {
            assertionResults = "Failed"
            testRes = false
        } else {
            assertionResults = "Success"
            testRes = true
        }
    }
    //
    msg := testcase.TestMessage {
        AssertionResults: assertionResults,
        ReponsePart: reponsePart,
        FieldName: key,
        AssertionKey:  assertionKey,
        ActualValue: actualValue,
        ExpValue: expValue,   
    }

    return testRes, &msg
}





