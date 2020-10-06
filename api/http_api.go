/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
    "fmt" 
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
    tcDataStore.PrepEmbeddedFunctions(path)

    tcData := tcDataStore.TcData

    // urlStr := tcData.TestCase.UrlRaw(baseUrl)
    urlStr := tcData.TestCase.UrlEncode(baseUrl)

    tcDataStore.HttpUrl = urlStr
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
    tcDataStore.PrepEmbeddedFunctions(path)


    httpExpResp := tcDataStore.TcData.TestCase.Response()
    testResults, testMessages = tcDataStore.CompareHttpRespGroup(httpExpResp)

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


func (tcDataStore *TcDataStore) CompareHttpRespGroup (httpExpResp []map[string]interface{}) ([]bool, []*testcase.TestMessage){
    var testResults []bool
    var testMessages []*testcase.TestMessage

    for _, v := range httpExpResp {
        for key, value := range v {
            httpExpResp_sub := value.(map[string]interface{})
            for assertionKey, expValueOrigin := range httpExpResp_sub {
                switch assertionKey {
                case "HasMapKey", "NotHasMapKey":
                    
                case "IsNull", "IsNotNull":

                default:
                    actualValue := tcDataStore.GetResponseValue(key)

                    var expValue interface{}
                    switch expValueOrigin.(type) {
                        case float64, int64, nil: 
                            expValue = expValueOrigin
                        default:
                            expValue = tcDataStore.GetResponseValue(fmt.Sprint(expValueOrigin))
                    }
                        
                    // $(status), $(headers), $(body)
                    var part string
                    switch {
                        case strings.HasPrefix(key, "$(status)"): 
                            part = "HTTP.Status"
                        case strings.HasPrefix(key, "$(headers)"): 
                            part = "HTTP.Headers"
                        case strings.HasPrefix(key, "$(body)"): 
                            part = "HTTP.Body"
                        default:
                            part = "HTTP"
                    }
                    testRes, msg := compareCommon(part, key, assertionKey, actualValue, expValue)
                    
                    testMessages = append(testMessages, msg)
                    testResults = append(testResults, testRes)
                }
            }
        }
    }

    return testResults, testMessages
}


//
func (tcDataStore *TcDataStore) HandleHttpResultsForOut () {
    // write out session if has
    path := "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "session"
    tcDataStore.PrepEmbeddedFunctions(path)

    expTcSession := tcDataStore.TcData.TestCase.Session()
    tcDataStore.WriteSession(expTcSession)

    // write out global variables if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "outGlobalVariables"
    tcDataStore.PrepEmbeddedFunctions(path)

    expOutGlobalVariables := tcDataStore.TcData.TestCase.OutGlobalVariables()
    tcDataStore.WriteOutGlobalVariables(expOutGlobalVariables)

    // write out tc loca variables if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "outLocalVariables"
    tcDataStore.PrepEmbeddedFunctions(path)

    expOutLocalVariables := tcDataStore.TcData.TestCase.OutLocalVariables()
    tcDataStore.WriteOutTcLocalVariables(expOutLocalVariables)

    // write out files if has
    path = "TestCase." + tcDataStore.TcData.TestCase.TcName() + "." + "outFiles"
    tcDataStore.PrepEmbeddedFunctions(path)

    expOutFiles := tcDataStore.TcData.TestCase.OutFiles()
    tcDataStore.HandleOutFiles(expOutFiles)
}

