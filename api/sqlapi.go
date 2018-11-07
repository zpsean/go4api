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
func RunTcSetUp (tcData testcase.TestCaseDataInfo, actualStatusCode int, actualHeader map[string][]string, actualBody []byte) (string, [][]*testcase.TestMessage) {
    cmdGroup := tcData.TestCase.SetUp()

    finalResults, finalTestMessages := Command(cmdGroup, actualStatusCode, actualHeader, actualBody)

    return finalResults, finalTestMessages
}

func RunTcTearDown (tcData testcase.TestCaseDataInfo, actualStatusCode int, actualHeader map[string][]string, actualBody []byte) (string, [][]*testcase.TestMessage) {
    cmdGroup := tcData.TestCase.TearDown()

    finalResults, finalTestMessages := Command(cmdGroup, actualStatusCode, actualHeader, actualBody)

    return finalResults, finalTestMessages
}

func Command (cmdGroup []*testcase.CommandDetails, actualStatusCode int, actualHeader map[string][]string, actualBody []byte) (string, [][]*testcase.TestMessage) {
    finalResults := "Success"
    var finalTestMessages = [][]*testcase.TestMessage{}
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

                    singleCmdResults, testMessages := compareRespGroup(cmdExpResp, rowsCount, rowsData, actualStatusCode, actualHeader, actualBody)

                    cmdsResults = append(cmdsResults, singleCmdResults)
                    finalTestMessages = append(finalTestMessages, testMessages)
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

    return finalResults, finalTestMessages
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
            
            testRes, msg := compareCommon("sql", key, assertionKey, actualValue, expValue)
            
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


