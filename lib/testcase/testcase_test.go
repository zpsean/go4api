/*
 * go4api - an api testing tool written in Go
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
    "testing"
    // "encoding/csv"
    "encoding/json"
)

var tcSampleStr string
var tcases TestCases
var tcSample TestCase

func init() {
    tcSampleStr := `
        [
          {
            "s2ChildTestCase-001": {
              "priority": "1",
              "parentTestCase": "s2ParentTestCase-001",
              "inputs": ["s2ParentTestCase_out.csv"],
              "request": {
                "method": "GET",
                "path": "https://api.douban.com/v2/movie/subject/1292052",
                "headers": {
                  "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"
                },
                "queryString": {
                  "pageIndex": "1",
                  "pageSize": "12"
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
          }
        ]`
    
    json.Unmarshal([]byte(tcSampleStr), &tcases)
    tcSample = tcases[0]
    
}

func Test_TcName(t *testing.T) {
    res := tcSample.TcName()

    fmt.Println(res)
    if res != `s2ChildTestCase-001` {
        t.Fatalf("TcName() test failed")
    } else {
        t.Log("TcName() test passed")
    }
}

func Test_Priority(t *testing.T) {
    res := tcSample.Priority()

    fmt.Println(res)
    if res != `1` {
        t.Fatalf("Priority() test failed")
    } else {
        t.Log("Priority() test passed")
    }
}



