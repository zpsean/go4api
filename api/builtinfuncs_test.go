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
    "fmt"
    "testing"
    // "strings"
    // "encoding/csv"
    "encoding/json"

    "go4api/lib/testcase"
    "go4api/api"

    // gjson "github.com/tidwall/gjson"
    // sjson "github.com/tidwall/sjson"
)

var tcSampleStr_2 string
var value_2 interface{}
var tcData testcase.TestCaseDataInfo

func init () {
    tcSampleStr_2 = `
      {
        "TestCase": {
          "casename-0001": {
            "priority": "1",
            "parentTestCase": "s2ParentTestCase-002",
            "inputs": [],
            "request": {
              "method": "GET",
              "path": "https://api.dummysite.com/v2/movie/subject/1292052",
              "payload": {
                "text": { 
                          "builtin1": {"Fn::CurrentTimeStampString": ""},
                          "builtin2": {"Fn::CurrentTimeStampString": "2006-01-02 15:04:05.999"},
                          "builtin3": {"Fn::CurrentTimeStampUnixMilli": ""},
                          "builtin4": {"Fn::NextAlphaNumeric": 33},
                          "builtin5": {"Fn::NextInt": [644, 6447]},
                          "builtin6": {"Fn::NextStringNumeric": 31},
                          "builtin7": {"Fn::ToString": 31},
                          "builtin8": {"Fn::ToString": {"Fn::CurrentTimeStampUnixMilli": []}},
                          "builtin9": {"Fn::ToString": 1234132.9876723},
                          "builtin10": {"Fn::Join" : [":", ["a", "b", "c"]]},
                          "builtin11": {"Fn::Split" : ["|", "a|b|c"]},
                          "date": 1541153618906,
                          "nullValue": null,
                          "blankMap": {}
                        }
              }
            },
            "response": {
              "status": {
                "Equals": 200
              }
            }
          }
        }
      }
        `
    json.Unmarshal([]byte(tcSampleStr_2), &value_2)

    json.Unmarshal([]byte(tcSampleStr_2), &tcData)

    tcJson, _ := json.MarshalIndent(tcData, "", "\t")
    fmt.Println(tcData)
    fmt.Println("origin tcdata: ", string(tcJson))
}


func Test_EvaluateBuiltinFunctions (t *testing.T) {
    res := api.EvaluateBuiltinFunctions(tcData)
    
    resj, _ := json.MarshalIndent(res, "", "\t")
    fmt.Println(string(resj))
    
    a := "12"
    if len(a) != 13 {
        t.Fatalf("json parse failed")
    } else {
        t.Log("json parse passed")
    }
}



