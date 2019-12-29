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
    "strings"
    "database/sql"

    "go4api/cmd"
    "go4api/utils"

    _ "github.com/go-sql-driver/mysql"
)

func InitMySqlConnection () map[string]*sql.DB {
    sqlCons := make(map[string]*sql.DB)

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
        sqlCons[dbIndicator] = db
    }

    return sqlCons
} 

