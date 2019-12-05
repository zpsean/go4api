/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gpg

import (
    "fmt"
    "os"
    // "strconv"
    "strings"
    "database/sql"
    // "reflect"
    // "encoding/json"

    "go4api/cmd"

    _ "github.com/lib/pq"
)

// var db = &sql.DB{}
var PgCons map[string]*sql.DB

type SqlExec struct {
    TargetDb string
    Stmt string
    CmdAffectedCount int
    RowsHeaders []string
    CmdResults []map[string]interface{}
}

func InitConnection () {
    PgCons = make(map[string]*sql.DB)

    //
    dbs := cmd.GetPgDbConfig()

    for k, v := range dbs {
        ip := v.Ip
        port := v.Port
        user := v.UserName
    
        pw := ""
        pwV := v.Password
        pwV = strings.Replace(pwV, "${", "", -1)
        pwV = strings.Replace(pwV, "}", "", -1)
        if len(pwV) > 0 {
            if len(os.Getenv(pwV)) > 0 {
                pw = os.Getenv(pwV)
            } else {
                pw = pwV
            }
        }
        // conInfo := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
        // conInfo := "postgres://" +  user + ":" + pw + "@" + ip + ":" + fmt.Sprint(port)
        u := "user=" + user 
        pa := " password=" + pw
        h := " host=" + ip
        po := " port=" + fmt.Sprint(port)
        d := " dbname=" + v.Dbname
        ssl := " sslmode=" + v.Sslmode
        conInfo := u + pa + h + po + d + ssl

        db, _ := sql.Open("postgres", conInfo)
        db.SetMaxOpenConns(2000)
        db.SetMaxIdleConns(1000)

        err := db.Ping()
        if err != nil {
            fmt.Println("Err, pg connection is not established. ", err)
            panic(err)
        }

        dbIndicator := strings.ToLower(k)
        PgCons[dbIndicator] = db
    }
    // fmt.Println("PgCons: ", PgCons)
} 

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
    db := PgCons[sqlExec.TargetDb]

    // defer func() {  
    //     if r := recover(); r != nil {  
    //         fmt.Println("!! Err, Catch gsql err:", r)  
    //     }  
    // }() 

    sqlStmt, err := db.Prepare(sqlExec.Stmt)
    if err != nil {
        fmt.Println("!! Err, Catch pg err:", err)
        panic(err)
    }
    defer sqlStmt.Close()

    res, err := sqlStmt.Exec()
    if err != nil {
        fmt.Println("!! Err, Catch pg err:", err)
        panic(err)
    }

    if err == nil {
        rowsAffected, _ := res.RowsAffected()
        sqlExec.CmdAffectedCount = int(rowsAffected)
    }

    return err
}

func (sqlExec *SqlExec) Delete () error {
    db := PgCons[sqlExec.TargetDb]

    // defer func() {  
    //     if r := recover(); r != nil {  
    //         fmt.Println("!! Err, Catch gsql err:", r) 
    //         sqlExec.CmdAffectedCount = -1
    //     }  
    // }()

    sqlStmt, err := db.Prepare(sqlExec.Stmt)
    if err != nil {
        fmt.Println("!! Err, Delete Prepare, Catch pg err:", err) 
        panic(err)
    }
    defer sqlStmt.Close()

    res, err := sqlStmt.Exec()
    if err != nil {
        fmt.Println("!! Err, Delete Exec, Catch pg err:", err)
        panic(err)
    }

    if err == nil {
        rowsAffected, _ := res.RowsAffected()
        sqlExec.CmdAffectedCount = int(rowsAffected)
    }

    return err
}

func (sqlExec *SqlExec) QueryWithoutParams () error {
    db := PgCons[sqlExec.TargetDb]

    // defer func() {  
    //     if r := recover(); r != nil {  
    //         fmt.Println("!! Err, Catch gsql err:", r)   
    //     }  
    // }()  
    sqlStmt, err := db.Prepare(sqlExec.Stmt)
    if err != nil {
        fmt.Println("!! Err, SELECT Prepare, Catch pg err:", err)
        panic(err)
    }
    defer sqlStmt.Close()

    rows, err := sqlStmt.Query()
    if err != nil {
        fmt.Println("!! Err, SELECT Query, Catch pg err:", err)
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
    db := PgCons[sqlExec.TargetDb]

    // defer func() {  
    //     if r := recover(); r != nil {  
    //         fmt.Println("!! Err, Catch gsql err:", r)   
    //     }  
    // }() 

    sqlStmt, err := db.Prepare(sqlExec.Stmt)
    if err != nil {
        fmt.Println("!! Err, Insert Prepare, Catch pg err:", err)
        panic(err)
    }
    defer sqlStmt.Close()

    res, err := sqlStmt.Exec()
    if err != nil {
        fmt.Println("!! Err, Insert Exec, Catch pg err:", err)
        panic(err)
    }

    if err == nil {
        rowsAffected, _ := res.RowsAffected()
        sqlExec.CmdAffectedCount = int(rowsAffected)
    }

    return err
}


