/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package g4json

import (
    "fmt"
    "testing"
    // "strings"
    // "encoding/csv"
	"encoding/json"
)

var tcSampleStr string
var tcSampleStr_2 string
var value interface{}
var value_2 interface{}

func init () {
    tcSampleStr = `
        {
          "priority": "1",
          "parentTestCase": "s2ParentTestCase-001",
          "inputs": ["s2ParentTestCase_out.csv"],
          "request": {
            "method": "GET",
            "path": "https://api.dummysite.com/v2/movie/subject/1292052",
            "headers": {
              "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"
            },
            "queryString": {
              "pageIndex": "1",
              "pageSize": "12"
            },
            "payload": {
	          "text": {
	                    "couIds": [
		                    {
		                        "id": "id11",
		                        "status": "status11",
		                        "Name": "name11"
		                    },
		                    {
		                        "id": "id22",
		                        "status": "status22",
		                        "Name": "name22"
		                    },
		                    null,
		                    {},
		                    [],
		            		123,
		            		123.0,
		            		123.44
	                    ],
	                    "date": 1541153618906,
	                    "nullValue": null,
	                    "intValue": 12345,
	                    "floatValue1": 12345.0,
	                    "floatValue2": 12345.555,
	                    "boolValue": true
	                  }
	        }
          },
          "response": {
            "status": {
              "Equals": 200
            },
            "headers": {
              "Content-Type": {
                "Contains": "application/json"
              }
            },
            "body": {
              "$.title": {
                "Contains": "{{.title}}"
              }
            }
          }
        }
        `

    tcSampleStr_2 = `
        {
          "priority": "1",
          "parentTestCase": "s2ParentTestCase-002",
          "inputs": [],
          "request": {
            "method": "GET",
            "path": "https://api.dummysite.com/v2/movie/subject/1292052",
            "payload": {
	          "text": {
	          			"field1": {"Fn::F1": []},
	                    "date": 1541153618906,
	                    "nullValue": null,
	                    "field2": {"Fn::F2": [12, {"Fn::F3": ["aaa", "bbbb"]}]},
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
        `
    
    json.Unmarshal([]byte(tcSampleStr), &value)
    json.Unmarshal([]byte(tcSampleStr_2), &value_2)
}

func Test_GetFieldsDetails (t *testing.T) {
    res := GetFieldsDetails(value)

    fmt.Println(len(res))

    // resj,_ := json.MarshalIndent(res, "", "\t")
    // fmt.Println(string(resj))

    if len(res) != 45 {
        t.Fatalf("json parse failed")
    } else {
        t.Log("json parse passed")
    }
}

func Test_GetJsonNodesLevel (t *testing.T) {
    res := GetFieldsDetails(value)

    max := GetJsonNodesLevel(res)
    
    if max != 6 {
        t.Fatalf("json parse failed")
    } else {
        t.Log("json parse passed")
    }
}

func Test_GetJsonNodesPath (t *testing.T) {
    res := GetFieldsDetails(value)

    a := GetJsonNodesPath(res)
    
    // resj,_ := json.MarshalIndent(a, "", "\t")
    // fmt.Println(string(resj), len(a))

    if len(a) != 45 {
        t.Fatalf("json parse failed")
    } else {
        t.Log("json parse passed")
    }
}

func Test_GetJsonLeavesPath (t *testing.T) {
    res := GetFieldsDetails(value)

    a := GetJsonLeavesPath(res)
    
    // resj,_ := json.MarshalIndent(a, "", "\t")
    // fmt.Println(string(resj), len(a))

    if len(a) != 29 {
        t.Fatalf("json parse failed")
    } else {
        t.Log("json parse passed")
    }
}


func Test_GetJsonNodesLevel_2 (t *testing.T) {
    res := GetFieldsDetails(value_2)

    max := GetJsonNodesLevel(res)
    
    if max != 8 {
        t.Fatalf("json parse failed")
    } else {
        t.Log("json parse passed")
    }
}


func Test_GetJsonLeavesPath_2 (t *testing.T) {
    res := GetFieldsDetails(value_2)

    a := GetJsonLeavesPath(res)
    
    resj,_ := json.MarshalIndent(a, "", "\t")
    fmt.Println(string(resj), len(a))
	// ...
	// "request.payload.text.field1.Fn::F1",
	// "request.payload.text.field2.Fn::F2.0",
	// "request.payload.text.field2.Fn::F2.1.Fn::F3.1",
	// "request.payload.text.field2.Fn::F2.1.Fn::F3.0"

    if len(a) != 13 {
        t.Fatalf("json parse failed")
    } else {
        t.Log("json parse passed")
    }
}

// usage: go test -test.bench=".*" -count=5
// 2018-11-01: Benchmark_GetFieldsDetails-8   	   50000	     28625 ns/op
func Benchmark_GetFieldsDetails (b *testing.B) {
	// use b.N for looping 
    for i := 0; i < b.N; i++ { 
        GetFieldsDetails(value)
    }
}

