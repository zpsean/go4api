/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package utils

import (
    "fmt"                                                                                                                                             
    "os"
    "io/ioutil" 
    "reflect"
    "bytes"
    gjson "github.com/tidwall/gjson"
)



func GetJsonFromFile(filePath string) string {
    fi, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    
    fd, err := ioutil.ReadAll(fi)
    if err != nil {
        panic(err)
    }

    return string(fd)
}


func GetBaseUrlFromConfig(filePath string, jsonKey string) string {
    resJson := GetJsonFromFile(filePath)

    value := gjson.Parse(resJson).Get(jsonKey).Get("baseUrl")

    if value.Value() == nil {
        panic(value)
    }

    return value.String()
}

func GetTestCaseJsonsFromTestDataFile(filePath string) []interface{} {
    resJson := GetJsonFromFile(filePath)

    value := gjson.Parse(resJson).Get("TestCases")

    if value.Value() == nil {
        panic(value)
    }

    testcases := value.Array()

    var tcJsons []interface{}

    for _, tc := range testcases {
        // fmt.Println("-------> tc: ", tc)
        tcJsons = append(tcJsons, tc.Value())
    }

    return tcJsons
}

func GetTestCaseJsonsFromTestData(fjson *bytes.Buffer) []interface{} {
    resJson, _ := ioutil.ReadAll(fjson)

    testcases := gjson.Parse(string(resJson)).Get("TestCases").Array()

    var tcJsons []interface{}

    for _, tc := range testcases {
        // fmt.Println("-------> tc: ", reflect.TypeOf(tc), reflect.TypeOf(tc.String()))
        tcJsons = append(tcJsons, tc.String())
    }

    return tcJsons
}

func GetTestCaseBasicInfoFromTestData(testcase interface{}) []interface{} {
    tc_map := gjson.Parse(testcase.(string)).Map()
    tc_str := gjson.Parse(testcase.(string)).String()

    var tcInfo []interface{}
    for key, _ := range tc_map {
      tcInfo = append(tcInfo, key)
      tcInfo = append(tcInfo, gjson.Parse(tc_str).Get(key).Get("priority").String())
      tcInfo = append(tcInfo, gjson.Parse(tc_str).Get(key).Get("parentTestCase").String())
    }

    return tcInfo
}

// for request
func GetRequestForTC(tc string, tcName string) (string, string) {
    request := gjson.Parse(tc).Get(tcName).Get("request")

    apiPath := request.Get("path").String()
    apiMethod := request.Get("method").String()

    return apiPath, apiMethod
}

func GetRequestHeadersForTC(tc string, tcName string) map[string]interface{} {
    reqHeaders := gjson.Parse(tc).Get(tcName).Get("request").Get("headers").Map()

    requestHeaders := map[string]interface{}{}
    for key, value := range reqHeaders {
        requestHeaders[key] = value.String()
    }
   
    return requestHeaders
}

func GetRequestPayloadForTC(tc string, tcName string) map[string]interface{} {
    // var requestPayload map[string]interface{}
    reqPayload := gjson.Parse(tc).Get(tcName).Get("request").Get("payload").Map()

    requestPayload := map[string]interface{}{}
    for key, value := range reqPayload {
        requestPayload[key] = value.String()
    }

    return requestPayload
}



// for expect
func GetExpectedResponseForTC(tc interface{}, tcName string) (gjson.Result, gjson.Result, gjson.Result) {
    var expStatusCode, expHeader, expBody gjson.Result

    expResponse := gjson.Parse(tc.(string)).Get(tcName).Get("response")
    expStatusCode = expResponse.Get("status")
    expHeader = expResponse.Get("headers")
    expBody = expResponse.Get("body")

    return expStatusCode, expHeader, expBody
}

func GetInputsFileNameForTC(tc interface{}, tcName string) []gjson.Result {
    var expOutputs []gjson.Result

    expOutputs = gjson.Parse(tc.(string)).Get(tcName).Get("inputs").Array()

    return expOutputs
}

func GetExpectedOutputsFieldsForTC(tc interface{}, tcName string) []gjson.Result {
    var expOutputs []gjson.Result

    expOutputs = gjson.Parse(tc.(string)).Get(tcName).Get("outputs").Array()

    return expOutputs
}



// backup for use
func GetJsonKeys(expBody interface{}) {
  
  mv := reflect.ValueOf(expBody)
  for _, k := range mv.MapKeys() {
      v := mv.MapIndex(k)
      fmt.Println("aaaaaaaaa")
      fmt.Println("aaaaaaaaa", k, v)
  }
}



    
