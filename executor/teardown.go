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
    "time"
    "os"
    "strings"

    "go4api/cmd"
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


func RunTeardown(ch chan int, pStart_time time.Time, pStart string, baseUrl string, resultsDir string, tcArray []testcase.TestCaseDataInfo) { 
    sqlTearDownTcSlice, notSqlTearDownTcSlice := ClassifyTearDown(tcArray)

    prioritySet, root, tcTree, tcTreeStats := RunBefore(notSqlTearDownTcSlice)

    fmt.Println("\n====> teardown test cases execution starts!") 

    RunPriorities(ch, pStart, baseUrl, resultsDir, notSqlTearDownTcSlice, prioritySet, root, tcTree, tcTreeStats)

    RunAfter(ch, pStart_time, pStart, resultsDir, notSqlTearDownTcSlice, root, tcTree, tcTreeStats)

    // -- for sql execution
    RunSqlTearDown(sqlTearDownTcSlice)
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

func RunSqlTearDown (sqlTearDownTcSlice []testcase.TestCaseDataInfo) {
    var sqlSlice []string

    for i, _ := range sqlTearDownTcSlice {
        for k, v := range sqlTearDownTcSlice[i].TestCase.TearDown() {
            if k == "sql" {
                sqlSlice = append(sqlSlice, fmt.Sprint(v))
            }
        }
    }

    ip, port, user, pw, defaultDB := GetDBConnInfo()
    gsql.InitConnection(ip, port, user, pw, defaultDB)

    gsql.Delete(sqlSlice[0])
}

func GetDBConnInfo () (string, string, string, string, string) {
    var ip, port, user, pw, defaultDB string

    testEnv := ""
    if cmd.Opt.TestEnv != "" {
        testEnv = cmd.Opt.TestEnv
    } else {
        testEnv = "qa"
    }

    switch strings.ToLower(testEnv) {
        case "qa":
            ip = os.Getenv("go4_qa_db_ip")
            port = os.Getenv("go4_qa_db_port")
            user = os.Getenv("go4_qa_db_username")
            pw = os.Getenv("go4_qa_db_password")
            defaultDB = os.Getenv("go4_qa_db_defaultDB")
        case "dev":
            ip = os.Getenv("go4_dev_db_ip")
            port = os.Getenv("go4_dev_db_port")
            user = os.Getenv("go4_dev_db_username")
            pw = os.Getenv("go4_dev_db_password")
            defaultDB = os.Getenv("go4_dev_db_defaultDB")
    }

    return ip, port, user, pw, defaultDB
}



