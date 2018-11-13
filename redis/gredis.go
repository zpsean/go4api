/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gredis

import (
    "fmt"
    "os"
    "strings"

    "go4api/cmd"

    redigo "github.com/gomodule/redigo/redis"
)

var RedisCons map[string]redigo.Conn

type RedisExec struct {
    TargetRedis string
    Stmt string
    KeysCount int
    CmdResults interface{}
}

func InitRedisConnection () {
    RedisCons = make(map[string]redigo.Conn)

    reds := cmd.GetRedisConfig()

    for k, v := range reds {
        ip := v.Ip
        port := v.Port
    
        pw := ""
        pwV := v.Password
        pwV = strings.Replace(pwV, "${", "", -1)
        pwV = strings.Replace(pwV, "}", "", -1)
        if len(pwV) > 0 {
            pw = os.Getenv(pwV)
        }
        
        // defaultRedisDb := os.Getenv("go4_qa_db_defaultRedisDb")

        c, err := redigo.Dial("tcp", ip + ":" + fmt.Sprint(port))
        if err != nil {
            panic(err)
        }

        // defer c.Close()

        if len(pw) > 0 {
            if _, err = c.Do("AUTH", pw); err != nil {  
                c.Close()  
                panic(err)
            }
        }

        _, err = c.Do("PING")  
        if err != nil {
            panic(err)
        }

        key := strings.ToLower(k)
        RedisCons[key] = c
    }
} 

func Run (stmt string) (int, interface{}, string) {
    var err error
    redExecStatus := ""

    tDb := "master"
    redisExec := &RedisExec{tDb, stmt, 0, ""}
    err = redisExec.Do()

    if err == nil {
        redExecStatus = "cmdSuccess"
    } else {
        redExecStatus = "cmdFailed"
    }

    return redisExec.KeysCount, redisExec.CmdResults, redExecStatus
}

func (redisExec *RedisExec) Do () error {
    c := RedisCons[redisExec.TargetRedis]

    var err error
    var res interface{}

    cmdS := strings.Split(strings.ToUpper(redisExec.Stmt), " ")[0]
    // to support del, set, get first
    switch strings.ToUpper(cmdS) {
        case "SET":
            res, err = c.Do(redisExec.Stmt) 
            // if err == nil {
            //     redisExec.CmdResults = res
            // }
        case "GET":
            res, err = c.Do(redisExec.Stmt) 
            if err == nil {
                redisExec.CmdResults = res
            }
        case "DEL":
            res, err = c.Do(redisExec.Stmt) 
            if err == nil {
                redisExec.KeysCount = res.(int)
            }
        default:
            fmt.Println("!! Warning, Command ", redisExec.Stmt, " is not supported currently, will enhance it later")
    }
     
    return err
}


