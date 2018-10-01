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
    "strconv"
    "fmt"
    "time"
    "log"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)


var db = &sql.DB{}

func init () {
    db, _ = sql.Open("mysql", "admin:password@/book")
} 

func method () {
    insert()
    query()
    update()
    delete()
}

func update () {
    tx, _ := db.Begin()
    
    tx.Exec("Update user set age = ? where uid = ?", 1, 1)

    tx.Commit()
}

func delete () {
    tx, _ := db.Begin()

    tx.Exec("DELETE FROM USER WHERE uid = ?", 1)

    tx.Commit()
}

func query () {
    stm, _ := db.Prepare("SELECT uid, username FROM USER")
    defer stm.Close()
    rows, _ = stm.Query()
    defer rows.Close()

    for rows.Next(){
         var name string
         var id int
        if err := rows.Scan(&id, &name); err != nil {
            log.Fatal(err)
        }
    }
}

func insert () {
    tx,_ := db.Begin()
    
    tx.Exec("INSERT INTO user(uid, username, age) values(?, ?, ?)", 1, "user" + strconv.Itoa(1), 1)

    tx.Commit()
}


