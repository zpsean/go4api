/*
 * go4api - an api testing tool written in Go
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
    CmdStr string
    CmdKey string
    CmdValue string
    CmdAffectedCount int
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

func Run (cmdStr string, cmdKey string, cmdValue string) (int, interface{}, string) {
    var err error
    redExecStatus := ""
    
    tDb := "master"
    redisExec := &RedisExec{tDb, cmdStr, cmdKey, cmdValue, 0, ""}
    err = redisExec.Do()

    if err == nil {
        redExecStatus = "cmdSuccess"
    } else {
        redExecStatus = "cmdFailed"
    }

    return redisExec.CmdAffectedCount, redisExec.CmdResults, redExecStatus
}

func (redisExec *RedisExec) Do () error {
    c := RedisCons[redisExec.TargetRedis]

    var err error
    var res interface{}

    // to support del, set, get first
    switch strings.ToUpper(redisExec.CmdStr) {
        case "SET":
            res, err = c.Do(redisExec.CmdStr, redisExec.CmdKey, redisExec.CmdValue) 
            if err == nil {
                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = res
            }
        case "GET":
            res, err = c.Do(redisExec.CmdStr, redisExec.CmdKey) 
            if err == nil {
                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = res
            }
        case "DEL":
            res, err = c.Do(redisExec.CmdStr, redisExec.CmdKey) 
            if err == nil {
                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = res
            }
        case "EXISTS":
            res, err = c.Do(redisExec.CmdStr, redisExec.CmdKey) 
            if err == nil {
                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = res
            }
        default:
            redisExec.CmdAffectedCount = -1
            fmt.Println("!! Warning, Command ", redisExec.CmdStr, " is not supported currently, will enhance it later")
    }

    return err
}


