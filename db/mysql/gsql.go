/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gsql

import (
    "fmt"
    // "os"
    // "strconv"
    "strings"
    "database/sql"
    // "reflect"
    // "encoding/json"

    "go4api/cmd"
    "go4api/utils"

    _ "github.com/go-sql-driver/mysql"
)

// var db = &sql.DB{}
var SqlCons map[string]*sql.DB

type SqlExec struct {
    TargetDb string
    Stmt string
    CmdAffectedCount int
    RowsHeaders []string
    CmdResults []map[string]interface{}
}

func InitConnection () {
    SqlCons = make(map[string]*sql.DB)

    dbs := cmd.GetDbConfig()

    for k, v := range dbs {
        envMap := utils.GetOsEnviron()
  
        ip := renderValue(v.Ip, envMap)
        port := renderValue(fmt.Sprint(v.Port), envMap)
        user := renderValue(v.UserName, envMap)
        password := renderValue(v.Password, envMap)
        // dbname := renderValue(v.Dbname, envMap)
        
        defaultSchema := ""

        conInfo := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + defaultSchema
        db, _ := sql.Open("mysql", conInfo)
        db.SetMaxOpenConns(2000)
        db.SetMaxIdleConns(1000)

        err := db.Ping()
        if err != nil {
            fmt.Println(err)
            panic(err)
        }

        dbIndicator := strings.ToLower(k)
        SqlCons[dbIndicator] = db
    }
} 

//
func renderValue (jsonStr string, feeder map[string]string) string {
    s := jsonStr
 
    if strings.HasPrefix(s, "${go4_") {
        // key, value are both string
        su := strings.Replace(s, "${go4_", "${", -1)
        for key, value := range feeder {
            k := "${" + key + "}"
            if k == su {
                s = strings.Replace(su, k, value, -1)
            }
        }
    }

    return s
}

//
func Run (tgtDb string, stmt string) (int, []string, []map[string]interface{}, string) {
    // update, delete, select, insert
    s := strings.TrimSpace(stmt)
    s = strings.ToUpper(s)
    s = string([]rune(stmt)[:6])

    var err error
    cmdExecStatus := ""

    // tDb := "master"
    tDb := tgtDb
    sqlExec := &SqlExec{tDb, stmt, -1, []string{}, []map[string]interface{}{}}

    switch strings.ToUpper(s) {
        case "UPDATE":
            err = sqlExec.Update()
        case "DELETE":
            err = sqlExec.Delete()
        case "SELECT":
            err = sqlExec.QueryWithoutParams()
        case "INSERT":
            err = sqlExec.Insert() 
    }

    if err == nil {
        cmdExecStatus = "cmdSuccess"
    } else {
        cmdExecStatus = "cmdFailed"
    }

    return sqlExec.CmdAffectedCount, sqlExec.RowsHeaders, sqlExec.CmdResults, cmdExecStatus
}

func (sqlExec *SqlExec) Update () error {
    db := SqlCons[sqlExec.TargetDb]

    // defer func() {  
    //     if r := recover(); r != nil {  
    //         fmt.Println("!! Err, Catch gsql err:", r)  
    //     }  
    // }() 

    sqlStmt, err := db.Prepare(sqlExec.Stmt)
    if err != nil {
        fmt.Println("!! Err, Catch gsql err:", err)
        panic(err)
    }
    defer sqlStmt.Close()

    res, err := sqlStmt.Exec()
    if err != nil {
        fmt.Println("!! Err, Catch gsql err:", err)
        panic(err)
    }

    if err == nil {
        rowsAffected, _ := res.RowsAffected()
        sqlExec.CmdAffectedCount = int(rowsAffected)
    }

    return err
}

func (sqlExec *SqlExec) Delete () error {
    db := SqlCons[sqlExec.TargetDb]

    // defer func() {  
    //     if r := recover(); r != nil {  
    //         fmt.Println("!! Err, Catch gsql err:", r) 
    //         sqlExec.CmdAffectedCount = -1
    //     }  
    // }()

    sqlStmt, err := db.Prepare(sqlExec.Stmt)
    if err != nil {
        fmt.Println("!! Err, Delete Prepare, Catch gsql err:", err) 
        panic(err)
    }
    defer sqlStmt.Close()

    res, err := sqlStmt.Exec()
    if err != nil {
        fmt.Println("!! Err, Delete Exec, Catch gsql err:", err)
        panic(err)
    }

    if err == nil {
        rowsAffected, _ := res.RowsAffected()
        sqlExec.CmdAffectedCount = int(rowsAffected)
    }

    return err
}

func (sqlExec *SqlExec) QueryWithoutParams () error {
    db := SqlCons[sqlExec.TargetDb]

    // defer func() {  
    //     if r := recover(); r != nil {  
    //         fmt.Println("!! Err, Catch gsql err:", r)   
    //     }  
    // }()  

    sqlStmt, err := db.Prepare(sqlExec.Stmt)
    if err != nil {
        fmt.Println("!! Err, SELECT Prepare, Catch gsql err:", err)
        panic(err)
    }
    defer sqlStmt.Close()

    rows, err := sqlStmt.Query()
    if err != nil {
        fmt.Println("!! Err, SELECT Query, Catch gsql err:", err)
        panic(err)
    }

    if err == nil {
        rowsCount, rowsHeaders, rowsData := ScanRows(rows)

        sqlExec.CmdAffectedCount = rowsCount
        sqlExec.RowsHeaders = rowsHeaders
        sqlExec.CmdResults = rowsData
    }

    return err
}

func (sqlExec *SqlExec) QueryWithParams () {

}

func ScanRows (rows *sql.Rows) (int, []string, []map[string]interface{}) {
    rowsHeaders, _ := rows.Columns()
    var rowsData []map[string]interface{}

    scanArgs := make([]interface{}, len(rowsHeaders))
    values := make([]interface{}, len(rowsHeaders))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    rowsCount := 0
    for rows.Next() {
        rows.Scan(scanArgs...)
        record := make(map[string]interface{})

        for i, col := range values {
            if col != nil {
                // note, try best to get the type information to interface{}
                switch col.(type) {
                case int:
                    record[rowsHeaders[i]] = col.(int)
                case int64:
                    record[rowsHeaders[i]] = col.(int64)
                case float32:
                    record[rowsHeaders[i]] = col.(float32)
                case float64:
                    record[rowsHeaders[i]] = col.(float64)
                default:
                    record[rowsHeaders[i]] = string(col.([]byte))
                }
            }
        }
        rowsCount = rowsCount + 1
        rowsData = append(rowsData, record)
    }

    return rowsCount, rowsHeaders, rowsData
}

func (sqlExec *SqlExec) Insert () error {
    db := SqlCons[sqlExec.TargetDb]

    // defer func() {  
    //     if r := recover(); r != nil {  
    //         fmt.Println("!! Err, Catch gsql err:", r)   
    //     }  
    // }() 

    sqlStmt, err := db.Prepare(sqlExec.Stmt)
    if err != nil {
        fmt.Println("!! Err, Insert Prepare, Catch gsql err:", err)
        panic(err)
    }
    defer sqlStmt.Close()

    res, err := sqlStmt.Exec()
    if err != nil {
        fmt.Println("!! Err, Insert Exec, Catch gsql err:", err)
        panic(err)
    }

    if err == nil {
        rowsAffected, _ := res.RowsAffected()
        sqlExec.CmdAffectedCount = int(rowsAffected)
    }

    return err
}


