/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gsql

import (
    // "os"
    "strconv"
    "fmt"
    // "time"
    // "log"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)

var db = &sql.DB{}

func InitConnection (ip string, port string, user string, pw string, defaultDB string) {
    conInfo := user + ":" + pw + "@tcp(" + ip + ":" + port + ")/" + defaultDB
    db, _ = sql.Open("mysql", conInfo)

    err := db.Ping()
    if err != nil {
        panic(err)
    }
} 


func Update () {
    tx, _ := db.Begin()
    
    tx.Exec("Update user set age = ? where uid = ?", 1, 1)

    tx.Commit()
}

func Delete (stmt string) {
    tx, _ := db.Begin()

    tx.Exec(stmt)

    tx.Commit()
}

func QueryWithoutParams () {
    tx, _ := db.Begin()

    tx.Exec("SELECT * FROM STORE;")

    tx.Commit()
}

func QueryWithParams () {
    stm, _ := db.Prepare("SELECT * FROM STORE;")
    defer stm.Close()
    rows, _ := stm.Query()
    defer rows.Close()

    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))
    for i := range values {
        scanArgs[i] = &values[i]
    }
     
    for rows.Next() {
        rows.Scan(scanArgs...)
        record := make(map[string]string)
        for i, col := range values {
            if col != nil {
                record[columns[i]] = string(col.([]byte))
            }
        }
        fmt.Println(record)
    }
}

func Insert () {
    tx,_ := db.Begin()
    
    tx.Exec("INSERT INTO user(uid, username, age) values(?, ?, ?)", 1, "user" + strconv.Itoa(1), 1)

    tx.Commit()
}


