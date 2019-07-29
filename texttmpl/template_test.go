/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package texttmpl

import (
    "fmt"
    "testing"
    // "encoding/csv"
)

func Test_GetTcJson(t *testing.T) {
    fmt.Println("\n--> test started")

    jsonTemplate := `[
                      {
                        "FirstTestCase-001": {
                          "priority": "9",
                          "parentTestCase": "root",
                          "inputs": [],
                          "request": {
                            "method": "GET",
                            "path": "https://api.douban.com/v2/movie/top250",
                            "headers": {
                              "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
                              "h1": {{.h1}},
                              "h2": {{.h2}},
                              "h3": {{.h3}},
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
                              "start": {
                                "GreaterOrEquals": "0"
                              }
                          },
                          "outputs": []
                        }
                      }
                    ]`

    csvHeader := []string{"h1", "h2", "h3"}
    csvRow := []string{"d1", "d2", `["file1", "union", "file2", "join", "file3"]`}

    tcJson := GetTcJson(jsonTemplate, csvHeader, csvRow)

    fmt.Println("tcJson: ", tcJson)

    fmt.Println("\n--> test finished")
}




