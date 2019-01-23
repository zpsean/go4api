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
    "strings"
    // "encoding/json"

    "go4api/lib/testcase" 
    g4http "go4api/protocal/http"
)

func (tcDataStore *TcDataStore) RunHttp (baseUrl string) (string, []*testcase.TestMessage) {
    tcDataStore.CallHttp(baseUrl)

    httpResult, httpTestMessages := tcDataStore.Compare()

    return httpResult, httpTestMessages
}

func (tcDataStore *TcDataStore) CallHttp (baseUrl string) {
    path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "request"
    tcDataStore.PrepVariablesBuiltins(path)

    tcData := tcDataStore.TcData

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

    tcDataStore.HttpActualStatusCode = actualStatusCode
    tcDataStore.HttpActualHeader = actualHeader
    tcDataStore.HttpActualBody = actualBody
}


func (tcDataStore *TcDataStore) Compare () (string, []*testcase.TestMessage) {
    var testResults []bool
    var testMessages []*testcase.TestMessage

    path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "response"
    tcDataStore.PrepVariablesBuiltins(path)
    
    // status
    testResultsS, testMessagesS := tcDataStore.CompareStatus()
    testResults = append(testResults, testResultsS[0:]...)
    testMessages = append(testMessages, testMessagesS[0:]...)
    // headers
    testResultsH, testMessagesH := tcDataStore.CompareHeaders()
    testResults = append(testResults, testResultsH[0:]...)
    testMessages = append(testMessages, testMessagesH[0:]...)
    // body
    testResultsB, testMessagesB := tcDataStore.CompareBody()
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


func (tcDataStore *TcDataStore) CompareStatus() ([]bool, []*testcase.TestMessage) {
    var testResults []bool
    var testMessages []*testcase.TestMessage

    tcData := tcDataStore.TcData
    expStatus := tcData.TestCase.RespStatus()
    actualStatusCode := tcDataStore.HttpActualStatusCode
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

func (tcDataStore *TcDataStore) CompareHeaders() ([]bool, []*testcase.TestMessage) {
    var testResults []bool
    var testMessages []*testcase.TestMessage

    tcData := tcDataStore.TcData
    expHeader := tcData.TestCase.RespHeaders()
    actualHeader := tcDataStore.HttpActualHeader
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

func (tcDataStore *TcDataStore) CompareBody() ([]bool, []*testcase.TestMessage) {
    var testResults []bool
    var testMessages []*testcase.TestMessage

    tcData := tcDataStore.TcData
    expBody := tcData.TestCase.RespBody()
    // body
    for key, value := range expBody {
        expBody_sub := value.(map[string]interface{})

        for assertionKey, expValue := range expBody_sub {
            // if path, then value - value, otherwise, key - value
            actualValue := tcDataStore.GetBodyActualValueByPath(key)
            
            testRes, msg := compareCommon("Body", key, assertionKey, actualValue, expValue)

            testMessages = append(testMessages, msg)
            testResults = append(testResults, testRes)
        }
    }

    return testResults, testMessages
} 

func (tcDataStore *TcDataStore) HandleHttpResultsForOut () {
    tcData := tcDataStore.TcData
    // write out session if has
    path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "session"
    tcDataStore.PrepVariablesBuiltins(path)

    expTcSession := tcData.TestCase.Session()
    tcDataStore.WriteSession(expTcSession)

    // write out global variables if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "outGlobalVariables"
    tcDataStore.PrepVariablesBuiltins(path)

    expOutGlobalVariables := tcData.TestCase.OutGlobalVariables()
    tcDataStore.WriteOutGlobalVariables(expOutGlobalVariables)

    // write out tc loca variables if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "outLocalVariables"
    tcDataStore.PrepVariablesBuiltins(path)

    expOutLocalVariables := tcData.TestCase.OutLocalVariables()
    tcDataStore.WriteOutTcLocalVariables(expOutLocalVariables)

    // write out files if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "outFiles"
    tcDataStore.PrepVariablesBuiltins(path)

    expOutFiles := tcData.TestCase.OutFiles()
    tcDataStore.HandleOutFiles(expOutFiles)
}

