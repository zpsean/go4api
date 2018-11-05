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
    
    "go4api/sql"
    "go4api/lib/testcase" 
)

func RunTcSetUp (tcData testcase.TestCaseDataInfo) string {
    var sqlSlice []string

    for k, v := range tcData.TestCase.SetUp() {
        if k == "sql" {
            sqlSlice = append(sqlSlice, fmt.Sprint(v))
        }
    }

    tcSetUpResult := CallSql(sqlSlice)

    return tcSetUpResult
}

func RunTcTearDown (tcData testcase.TestCaseDataInfo) string {
    var sqlSlice []string
  
    for k, v := range tcData.TestCase.TearDown() {
        if k == "sql" {
            sqlSlice = append(sqlSlice, fmt.Sprint(v))
        }
    }

    tcTearDownResult := CallSql(sqlSlice)

    return tcTearDownResult
}

func CallSql (sqlSlice []string) string {
    var sqlRessult = make([]string, len(sqlSlice))  //value: SqlSuccess, SqlFailed
    tcSqlResult := "SqlSuccess"

    for i, _ := range sqlSlice {
        _, sqlRessult[i] = gsql.Run(sqlSlice[i])
        if sqlRessult[i] == "SqlFailed" {
            tcSqlResult = "SqlFailed"
            // break
        }
    }

    return tcSqlResult
}

