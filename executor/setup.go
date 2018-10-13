/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package executor

import (
    "fmt"
    // "time"

    "go4api/lib/testcase"
    "go4api/sql"
)

func GetSetupTcSlice (tcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var setUpTcSlice []testcase.TestCaseDataInfo
    for i, _ := range tcArray {
        if tcArray[i].TestCase.IfGlobalSetUpTestCase() == true {
            setUpTcSlice = append(setUpTcSlice, tcArray[i])
        }
    }

    return setUpTcSlice
}

func RunSetup(ch chan int, baseUrl string, resultsDir string, resultsLogFile string, tcArray []testcase.TestCaseDataInfo) { 
    // sqlSetUpTcSlice, notSqlSetUpTcSlice := ClassifySetUp(tcArray)

    // prioritySet, root, tcTree, tcTreeStats := RunBefore(notSqlSetUpTcSlice)
    // fmt.Println("\n====> setup test cases execution starts!") 
    // RunPriorities(ch, gStart, baseUrl, resultsDir, notSqlSetUpTcSlice, prioritySet, root, tcTree, tcTreeStats)
    // RunConsoleOverallReport(ch, gStart_time, gStart, resultsDir, notSqlSetUpTcSlice, root, tcTree, tcTreeStats)

    // // -- for sql execution
    // RunSqlSetUpTc(sqlSetUpTcSlice)
}

func ClassifySetUp (tcArray []testcase.TestCaseDataInfo) ([]testcase.TestCaseDataInfo, []testcase.TestCaseDataInfo) {
    var sqlSetUpTcSlice []testcase.TestCaseDataInfo
    var notSqlSetUpTcSlice []testcase.TestCaseDataInfo

    for i, _ := range tcArray {
        ifSql := false
        for k, _ := range tcArray[i].TestCase.SetUp() {
            if k == "sql" {
                sqlSetUpTcSlice = append(sqlSetUpTcSlice, tcArray[i])
                ifSql = true
            }
        }
        if ifSql == false {
            notSqlSetUpTcSlice = append(notSqlSetUpTcSlice, tcArray[i])
        }
    }

    return sqlSetUpTcSlice, notSqlSetUpTcSlice
}

func RunSqlSetUpTc (sqlTcSlice []testcase.TestCaseDataInfo) {
    var sqlSlice []string

    for i, _ := range sqlTcSlice {
        for k, v := range sqlTcSlice[i].TestCase.SetUp() {
            if k == "sql" {
                sqlSlice = append(sqlSlice, fmt.Sprint(v))
            }
        }
    }

    if len(sqlSlice) > 0 {
        ip, port, user, pw, defaultDB := gsql.GetDBConnInfo()
        gsql.InitConnection(ip, port, user, pw, defaultDB)
    }

    for i, _ := range sqlSlice {
        gsql.Run(sqlSlice[i])
    }
}

