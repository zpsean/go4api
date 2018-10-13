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
    // "os"
    // "strings"

    "go4api/lib/testcase"
    "go4api/sql"
)

func GetTeardownTcSlice (tcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var teardownTcSlice []testcase.TestCaseDataInfo
    for i, _ := range tcArray {
        if tcArray[i].TestCase.IfGlobalTearDownTestCase() == true {
            teardownTcSlice = append(teardownTcSlice, tcArray[i])
        }
    }
    
    return teardownTcSlice
}


func RunTeardown(ch chan int, baseUrl string, resultsLogFile string, resultsDir string, tcArray []testcase.TestCaseDataInfo) { 
    // sqlTearDownTcSlice, notSqlTearDownTcSlice := ClassifyTearDown(tcArray)

    // prioritySet, root, tcTree, tcTreeStats := RunBefore(notSqlTearDownTcSlice)
    // fmt.Println("\n====> teardown test cases execution starts!") 
    // RunPriorities(ch, gStart, baseUrl, resultsDir, notSqlTearDownTcSlice, prioritySet, root, tcTree, tcTreeStats)
    // RunConsoleOverallReport(ch, gStart_time, gStart, resultsDir, notSqlTearDownTcSlice, root, tcTree, tcTreeStats)

    // // -- for sql execution
    // RunSqlTearDownTc(sqlTearDownTcSlice)
}

func ClassifyTearDown (tcArray []testcase.TestCaseDataInfo) ([]testcase.TestCaseDataInfo, []testcase.TestCaseDataInfo) {
    var sqlTearDownTcSlice []testcase.TestCaseDataInfo
    var notSqlTearDownTcSlice []testcase.TestCaseDataInfo

    for i, _ := range tcArray {
        ifSql := false
        for k, _ := range tcArray[i].TestCase.TearDown() {
            if k == "sql" {
                sqlTearDownTcSlice = append(sqlTearDownTcSlice, tcArray[i])
                ifSql = true
            }
        }
        if ifSql == false {
            notSqlTearDownTcSlice = append(notSqlTearDownTcSlice, tcArray[i])
        }
    }

    return sqlTearDownTcSlice, notSqlTearDownTcSlice
}

func RunSqlTearDownTc (sqlTcSlice []testcase.TestCaseDataInfo) {
    var sqlSlice []string

    for i, _ := range sqlTcSlice {
        for k, v := range sqlTcSlice[i].TestCase.TearDown() {
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
