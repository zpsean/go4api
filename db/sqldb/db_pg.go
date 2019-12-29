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

    _ "github.com/lib/pq"
)

// var db = &sql.DB{}
func InitPgConnection () map[string]*sql.DB {
    sqlCons := make(map[string]*sql.DB)
    //
    dbs := cmd.GetPgDbConfig()

    for k, v := range dbs {
        envMap := utils.GetOsEnviron()
  
        ip := renderValue(v.Ip, envMap)
        port := renderValue(fmt.Sprint(v.Port), envMap)
        userName := renderValue(v.UserName, envMap)
        password := renderValue(v.Password, envMap)
        dbname := renderValue(v.Dbname, envMap)
        sslmode := renderValue(v.Sslmode, envMap)

        // conInfo := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
        // conInfo := "postgres://" +  user + ":" + pw + "@" + ip + ":" + fmt.Sprint(port)
        h := " host=" + ip
        u := "user=" + userName 
        pa := " password=" + password
        po := " port=" + port
        d := " dbname=" + dbname
        ssl := " sslmode=" + sslmode
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
        sqlCons[dbIndicator] = db
    }

    return sqlCons
} 
