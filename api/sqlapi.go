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
    "strings" 
    "encoding/json"

    "go4api/sql"
    "go4api/lib/testcase"

    gjson "github.com/tidwall/gjson"
)
// Note: for each SetUp, TesrDown, it may have more than one Command (including sql)
// for each Command, it may have more than one assertion
func RunTcSetUp (tcData testcase.TestCaseDataInfo, actualStatusCode int, actualHeader map[string][]string, actualBody []byte) string {
    cmdGroup := tcData.TestCase.SetUp()

    finalResults := Command(cmdGroup, actualStatusCode, actualHeader, actualBody)

    return finalResults
}

func RunTcTearDown (tcData testcase.TestCaseDataInfo, actualStatusCode int, actualHeader map[string][]string, actualBody []byte) string {
    cmdGroup := tcData.TestCase.TearDown()

    finalResults := Command(cmdGroup, actualStatusCode, actualHeader, actualBody)

    return finalResults
}

func Command (cmdGroup []*testcase.CommandDetails, actualStatusCode int, actualHeader map[string][]string, actualBody []byte) string {
    finalResults := "Success"
    var cmdsResults []bool
    //
    cmdGroupJsonB, _ := json.Marshal(cmdGroup)
    cmdGroupJson := string(cmdGroupJsonB)

    for i, _ := range cmdGroup {
        cmdType := gjson.Get(cmdGroupJson, fmt.Sprint(i) + "." + "cmdType")

        switch strings.ToLower(cmdType.String()) {
            case "sql":
                cmdStr := gjson.Get(cmdGroupJson, fmt.Sprint(i) + "." + "cmd")
                rowsCount, _, rowsData, sqlExecStatus := RunSql(cmdStr.String())

                if sqlExecStatus == "SqlSuccess" {
                    cmdExpResp := gjson.Get(cmdGroupJson, fmt.Sprint(i) + "." + "cmdResponse").Map()

                    singleCmdResults, _ := compareRespGroup(cmdExpResp, rowsCount, rowsData, actualStatusCode, actualHeader, actualBody)
                    cmdsResults = append(cmdsResults, singleCmdResults)
                    //
                    
                } else {
                    cmdsResults = append(cmdsResults, false)
                }
        }
    }

    for key := range cmdsResults {
        if cmdsResults[key] == false {
            finalResults = "Fail"
            break
        }
    }

    return finalResults
}

func compareRespGroup (cmdExpResp map[string]gjson.Result, rowsCount int, rowsData []map[string]interface{},
    actualStatusCode int, actualHeader map[string][]string, actualBody []byte) (bool, []*testcase.TestMessage) {
    //------
    singleCmdResults := true
    var testResults []bool
    var testMessages []*testcase.TestMessage

    for key, value := range cmdExpResp {
        cmdExpResp_sub := value.Value().(map[string]interface{})
        for assertionKey, expValueOrigin := range cmdExpResp_sub {
            
            actualValue := GetSqlActualRespValue(key, rowsCount, rowsData)

            var expValue interface{}
            switch expValueOrigin.(type) {
                case float64, int64: 
                    expValue = expValueOrigin
                default:
                    expValue = GetResponseValue(expValueOrigin.(string), actualStatusCode, actualHeader, actualBody)
            }
            
            fmt.Println("sql", key, assertionKey, actualValue, expValue)
            testRes, msg := compareCommon("sql", key, assertionKey, actualValue, expValue)

            fmt.Println("testRes, msg: ", testRes, msg)

            testMessages = append(testMessages, msg)
            testResults = append(testResults, testRes)
        }
    }

    for key := range testResults {
        if testResults[key] == false {
            singleCmdResults = false
            break
        }
    }

    return singleCmdResults, testMessages
}

func RunSql (stmt string) (int, []string, []map[string]interface{}, string) {
    // gsql.Run will return: <impacted rows : int>, <rows for select : [][]interface{}{}>, <sql status : string>
    // status: SqlSuccess, SqlFailed
    rowsCount, rowsHeaders, rowsData, sqlExecStatus := gsql.Run(stmt)

    return rowsCount, rowsHeaders, rowsData, sqlExecStatus
}


