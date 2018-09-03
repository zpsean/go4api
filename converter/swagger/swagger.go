/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package swagger

import (
    "fmt"
    // "time"
    // "os"
    // "sort"
    // "sync"
    // "io/ioutil"
    "strconv"
    "strings"
    "encoding/json"
    // "go4api/cmd"
    "go4api/utils"
    "go4api/testcase" 
)

func Convert () {
    var swagger Swagger2

    // filePath := cmd.Opt.Swagger
    filePath := "/Users/pingzhu/Downloads/go/run/testhome/testresource/swagger.json"
    resJson := utils.GetJsonFromFile(filePath)
    json.Unmarshal([]byte(resJson), &swagger)

    //
    var testCases []testcase.TestCase
    //
    i := 0
    for urlPath, path := range swagger.Paths {
    	// tcname
    	for method, pathDetails := range path {
    		i = i + 1
    		g4ATcName := "Swagger2-" + method + "-" + urlPath + "-" + strconv.Itoa(i)

    		priority, parentTestCase, inputs, _ := pathDetails.buildG4ATcGeneralInfo()

    		g4ARequest := pathDetails.buildG4ARequest()
    		g4ARequest.Method = strings.ToUpper(method)
    		g4ARequest.Path = urlPath

    		g4AResponse := pathDetails.buildG4AResponse()

    		//
		    tCase := make(map[string]testcase.TestCaseBasics)
		    
		    tCaseBasics := testcase.TestCaseBasics{
		        priority, 
		        parentTestCase,
		        inputs,
		        g4ARequest,
		        g4AResponse,
		        []interface{}{},
		    }
		    //
		    tCase[g4ATcName] = tCaseBasics

		    tcJson, _ := json.Marshal(tCase)
		    fmt.Print(string(tcJson))

		    testCases = append(testCases, tCase)
    	}
    }

    fmt.Println("")

    tcsJson, _ := json.MarshalIndent(testCases, "", "\t")
    fmt.Print(string(tcsJson))
}


func (pathDetails PathDetails) buildG4ATcGeneralInfo () (string, string, string, string) {
    // priority, parentTestCase, inputs, outputs
    return "1", "root", "", ""
}


func (pathDetails PathDetails) buildG4ARequest () testcase.Request {
    var g4ARequest testcase.Request

    //
    reqHeaders := make(map[string]interface{})
    // for _, consumes := range pathDetails.Consumes {
        reqHeaders["Content-Type"] = strings.Join(pathDetails.Consumes, ";")
    // }
    g4ARequest.Headers = reqHeaders

    //
    reqQS := make(map[string]interface{})
    g4ARequest.QueryString = reqQS

    //
    reqPL := make(map[string]interface{})
    g4ARequest.Payload = reqPL

    return g4ARequest
}

func (pathDetails PathDetails) buildG4AResponse () testcase.Response {
    var g4AResponse testcase.Response

    respHeaders := make(map[string]interface{})
    // for _, consumes := range pathDetails.Produces {
        respHeaders["Content-Type"] = strings.Join(pathDetails.Produces, ";")
    // }
    g4AResponse.Headers = respHeaders

    respStatus := make(map[string]interface{})
    respStatus["Equals"] = 200

    g4AResponse.Status = respStatus

    return g4AResponse
}


func (definitions Definitions) resolveDefinitionNest () {

}


func (definition Definition) BuildJsonExample () {
    // for key, value := range definition.Properties {

    // }
}


