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
    "time"
    "strings"
    "database/sql"
    "strconv"
)

// var db = &sql.DB{}
// SqlCons, example,
// {
//     "mysql": {
//         "master": {}
//         "slave" : {}
//     }
//     "postgres": {
//         "master": {}
//         "slave" : {}
//     }
// }

var SqlCons map[string]map[string]*sql.DB

type SqlExec struct {
    DriverName       string
    TargetDb         string
    Stmt             string
    CmdAffectedCount int
    RowsHeaders      []string
    CmdResults       []map[string]interface{}
}

func InitConnection (driverName string) {
    dri := strings.ToLower(driverName)
    SqlCons = make(map[string]map[string]*sql.DB)

    switch dri {
    case "sql", "mysql":
        cons := InitMySqlConnection()
        SqlCons["mysql"] = cons
    case "postgres", "postgresql":
        cons := InitPgConnection()
        SqlCons["postgres"] = cons
    }
} 

//
func Run (driverName string, tgtDb string, stmt string) (int, []string, []map[string]interface{}, string) {
    // update, delete, select, insert
    s := strings.TrimSpace(stmt)
    s = strings.ToUpper(s)
    s = string([]rune(stmt)[:6])

    var err error
    cmdExecStatus := ""

    // tDb := "master"
    tDb := tgtDb
    sqlExec := &SqlExec {
        DriverName:       driverName,
        TargetDb:         tDb,
        Stmt:             stmt,
        CmdAffectedCount: -1,
        RowsHeaders:      []string{},
        CmdResults:       []map[string]interface{}{},
    }

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
    db := SqlCons[sqlExec.DriverName][sqlExec.TargetDb]

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
    db := SqlCons[sqlExec.DriverName][sqlExec.TargetDb]

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
    db := SqlCons[sqlExec.DriverName][sqlExec.TargetDb]

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
            // note, try best to get the type information to interface{}
            // int64
            // float64
            // bool
            // []byte
            // string
            // time.Time
            // nil
            switch col.(type) {
            case int:
                record[rowsHeaders[i]] = col.(int)
            case int64:
                record[rowsHeaders[i]] = col.(int64)
            case float32:
                record[rowsHeaders[i]] = col.(float32)
            case float64:
                record[rowsHeaders[i]] = col.(float64)
            case bool:
                record[rowsHeaders[i]] = col.(bool)
            case []byte:
                // postgresql's field type numeric(m, n) is recognized as []byte (i.e. []uint8)
                // for example: 99999999.9900 represeted as => []unit8 => [57 57 57 57 57 57 57 57 46 57 57 48 48] => 
                //
                // !!Note: here need more code to enhance
                var v interface{}

                s :=  string(col.([]byte))
                v, err := strconv.ParseFloat(s, 64)
                if err != nil {
                    fmt.Println("!! Err, the string can not be parsed int float64:", col.([]byte), err)
                    // panic(err)
                    v = col.([]byte)
                }
                
                record[rowsHeaders[i]] = v
            case string:
                record[rowsHeaders[i]] = col.(string)
            case time.Time:
                record[rowsHeaders[i]] = col.(time.Time)
            case nil:
                record[rowsHeaders[i]] = nil
            default:
                record[rowsHeaders[i]] = fmt.Sprint(col)
            }
        }
        rowsCount = rowsCount + 1
        rowsData = append(rowsData, record)
    }
    // fmt.Println("---rowsCount, rowsHeaders, rowsData: ", rowsCount, rowsHeaders, rowsData)
    return rowsCount, rowsHeaders, rowsData
}

func (sqlExec *SqlExec) Insert () error {
    db := SqlCons[sqlExec.DriverName][sqlExec.TargetDb]

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


