/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package har

import (
    "fmt"
    // "time"
    // "os"
    // "sort"
    // "sync"
    "io/ioutil"
    "strconv"
    "strings"
    "encoding/json"
    "go4api/cmd"
    "go4api/utils"
    "go4api/testcase" 
)

func Convert () {
    var har Har

    filePath := cmd.Opt.Harfile
    resJson := utils.GetJsonFromFile(filePath)
    json.Unmarshal([]byte(resJson), &har)

    //
    var testCases []testcase.TestCase
    //
    // Note, if cnovert a bunch of files, the sequence need to use global sequence to avoid duplicate tcnames
    for i, entry := range har["log"].Entries {
        g4ATcName := entry.buildG4ATcName(i)
        priority, parentTestCase, _, _ := entry.buildG4ATcGeneralInfo()

        g4ARequest := entry.Request.buildG4ARequest()
        g4AResponse := entry.Response.buildG4AResponse()

        //
        tCase := make(map[string]*(testcase.TestCaseBasics))
        
        tCaseBasics := testcase.TestCaseBasics {
            priority, 
            parentTestCase,
            []interface{}{},
            &g4ARequest,
            &g4AResponse,
            []*testcase.OutputsDetails{},
        }

        tCase[g4ATcName] = &tCaseBasics

        // tcJson, _ := json.Marshal(tCase)
        // fmt.Print(string(tcJson))

        testCases = append(testCases, tCase)
    }
    fmt.Println("") 

    tcsJson, _ := json.MarshalIndent(testCases, "", "\t")
    // fmt.Print(string(tcsJson))

    // json write to file
    outPath := cmd.Opt.Harfile + ".out.json"
    ioutil.WriteFile(outPath, tcsJson, 0644)

    fmt.Println("\n! Convert Har file finished !")
    fmt.Println("")
}


func (harEntry Entry) buildG4ATcName (sequence int) string {
    g4ATcName := ""

    // the format of tcNmae would be: Get - url(end point) - sequenrce
    methodP := harEntry.Request.Method

    urlPath := strings.Split(harEntry.Request.Url, "?")[0]
    urlSlice := strings.Split(urlPath, "/")

    urlP := ""
    if len(urlSlice[len(urlSlice) - 1]) > 0 {
        urlP = urlSlice[len(urlSlice) - 1]
    } else {
        urlP = urlSlice[len(urlSlice) - 2]
    }

    g4ATcName = "Har-" + methodP + "-" + urlP + "-" + strconv.Itoa(sequence)

    return g4ATcName
}

func (harEntry Entry) buildG4ATcGeneralInfo () (string, string, string, string) {
    // priority, parentTestCase, inputs, outputs
    return "1", "root", "", ""
}


func (harRequest Request) buildG4ARequest () testcase.Request {
    var g4ARequest testcase.Request

    g4ARequest.Method = harRequest.Method
    g4ARequest.Path = harRequest.Url

    //
    reqHeaders := make(map[string]interface{})
    for _, header := range harRequest.Headers {
        reqHeaders[header["name"].(string)] = header["value"]
    }
    g4ARequest.Headers = reqHeaders

    //
    reqQS := make(map[string]interface{})
    for _, qs := range harRequest.QueryString {
        reqQS[qs["name"].(string)] = qs["value"]
    }
    g4ARequest.QueryString = reqQS

    //
    reqPL := make(map[string]interface{})
    for key, value := range harRequest.PostData {
        if key == "text" {
            reqPL["text"] = value
        }
        
    }
    g4ARequest.Payload = reqPL

    return g4ARequest
}

func (harResponse Response) buildG4AResponse () testcase.Response {
    var g4AResponse testcase.Response

    respStatus := make(map[string]interface{})
    respStatus["Equals"] = harResponse.Status

    g4AResponse.Status = respStatus

    return g4AResponse
}
