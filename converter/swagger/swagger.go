/*
 * go4api - an api testing tool written in Go
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
    "io/ioutil"
    "strconv"
    "strings"
    // "reflect"
    "encoding/json"
    "go4api/cmd"
    "go4api/utils"
    "go4api/lib/testcase" 
    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)


var swagger Swagger2
var defJson = `{}`
var defReBuilt map[string]interface{}


func Convert () {
    // load the swagger file to json
    filePath := cmd.Opt.Swaggerfile
    resJson := utils.GetJsonFromFile(filePath)
    json.Unmarshal([]byte(resJson), &swagger)

    // build the Definitions, resolve the nest reference
    for defKey, defValue := range swagger.Definitions {
        // fmt.Println("\ndefKey: ", defKey)
        // call buildDefinitions, value stored in defJson
        buildDefinitions(defKey, defValue)
    }
    
    // Unmarshal the json to variable defReBuilt
    json.Unmarshal([]byte(defJson), &defReBuilt)
    
    // build the target testcases
    var testCases []testcase.TestCase
    //
    i := 0
    for urlPath, path := range swagger.Paths {
    	// tcname
    	for method, pathDetails := range path {
    		i = i + 1
    		g4ATcName := "Swagger2-" + method + "-" + urlPath + "-" + strconv.Itoa(i)

    		priority, parentTestCase, _, _ := pathDetails.buildG4ATcGeneralInfo()

    		g4ARequest := pathDetails.buildG4ARequest()
    		g4ARequest.Method = strings.ToUpper(method)
    		g4ARequest.Path = urlPath

    		g4AResponse := pathDetails.buildG4AResponse()

    		//
		    tCase := make(map[string]*(testcase.TestCaseBasics))
		    
		    tCaseBasics := testcase.TestCaseBasics{
                Priority:                 priority, 
                ParentTestCase:           parentTestCase,
                IfGlobalSetUpTestCase:    false,
                IfGlobalTearDownTestCase: false,
                SetUp:                    []*testcase.CommandDetails{},
                Inputs:                   []interface{}{},
                Request:                  &g4ARequest,
                Response:                 &g4AResponse,
                Outputs:                  []*testcase.OutputsDetails{},
                OutFiles:                 []*testcase.OutFilesDetails{},
                OutGlobalVariables:       map[string]interface{}{},  // OutGlobalVariables
                OutLocalVariables:        map[string]interface{}{},  // OutLocalVariables
                Session:                  map[string]interface{}{},  // Session
                TearDown:                 []*testcase.CommandDetails{},
		    }
		    //
		    tCase[g4ATcName] = &tCaseBasics

		    // tcJson, _ := json.Marshal(tCase)
		    // fmt.Print(string(tcJson))

		    testCases = append(testCases, tCase)
    	}
    }

    // marshal the testcases to json
    tcsJson, _ := json.MarshalIndent(testCases, "", "\t")
    // fmt.Print(string(tcsJson))

    // json write to file
    outPath := cmd.Opt.Swaggerfile + ".out.json"
    ioutil.WriteFile(outPath, tcsJson, 0644) 

    fmt.Println("\n! Convert Swagger API file finished !")
    fmt.Println("")
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

    // query string
    reqQS := make(map[string]interface{})
    
    // post body
    reqPL := make(map[string]interface{})

    for _, param := range pathDetails.Parameters {
        // for post, put
        if param.In == "body" {
            if ifMapHasKey(param.Schema, "$ref") {
                definitionPath := strings.Split(param.Schema["$ref"].(string), "/")

                reqPL["text"] = defReBuilt[definitionPath[len(definitionPath) - 1]]
            }
        // for query, in get, post, etc.
        } else if param.In == "query" {
            if param.Type != "" {
                switch param.Type {
                    case "string":
                        reqQS[param.Name] = ""
                    case "integer":
                        reqQS[param.Name] = 0
                    case "array": {
                        switch param.Items["type"] {
                            case "string":
                                reqQS[param.Name] = ""
                            case "integer":
                                reqQS[param.Name] = 0
                        }
                    }
                }
            }
        }
    }
    //
    g4ARequest.QueryString = reqQS
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


func buildDefinitions (defKey string, defValue interface{})  {
    if getFieldType(defValue.(map[string]interface {})) == "object" {
        defProperties := defValue.(map[string]interface {})["properties"].(map[string]interface {})

        for propKey, propValue := range defProperties {
            switch getFieldType(propValue.(map[string]interface{})) {
                case "integer":
                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey, 0)
                case "boolean":
                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey, true)
                case "string":
                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey, "")
                case "": {
                    if ifMapHasKey(propValue.(map[string]interface{}), "$ref") {
                        definitionPath := strings.Split(propValue.(map[string]interface{})["$ref"].(string), "/")
                        buildDefinitions(defKey + "." + propKey, swagger.Definitions[definitionPath[len(definitionPath) - 1]])
                    }
                }
                case "array": {
                    if ifMapHasKey(propValue.(map[string]interface{}), "items") {
                        arrayItems := propValue.(map[string]interface{})["items"].(map[string]interface{})

                        if ifMapHasKey(arrayItems, "$ref") {
                            // add one array, and then add at least one element
                            defJson, _ = sjson.Set(defJson, defKey + "." + propKey, []map[string]interface{}{})
                            
                            definitionPath := strings.Split(arrayItems["$ref"].(string), "/")
                            buildDefinitions(defKey + "." + propKey + ".0", swagger.Definitions[definitionPath[len(definitionPath) - 1]])
                        } else if ifMapHasKey(arrayItems, "type") {
                            // add one array, and then add at least one element
                            switch getFieldType(arrayItems) {
                                case "integer":
                                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey, []int{})
                                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey + ".-1", 0)
                                case "boolean":
                                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey, []bool{})
                                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey + ".-1", true)
                                case "string":
                                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey, []string{})
                                    defJson, _ = sjson.Set(defJson, defKey + "." + propKey + ".-1", "")
                            }
                        }
                    }  
                }
            }
        }
    }
}


func getMapKeys (value map[string]interface{}) []string {
    var keySlice []string
    for key, _ := range value {
        keySlice = append(keySlice, key)
    }

    return keySlice
}

func getFieldType (value map[string]interface{}) string {
    var typeValue string
    for key, v := range value {
        if key == "type" {
            typeValue = v.(string)
            break
        } else {
            typeValue = ""
        }
    }
    return typeValue
}

func ifMapHasKey (value map[string]interface{}, searchKey string) bool {
    var ifRef bool
    ifRef = false
    for key, _ := range value {
        if key == searchKey {
            ifRef = true
            break
        }
    }
    return ifRef
}



