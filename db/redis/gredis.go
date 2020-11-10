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
    "strings"

    "go4api/cmd"
    "go4api/utils"

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
        envMap := utils.GetOsEnviron()
  
        ip := renderValue(v.Ip, envMap)
        port := renderValue(fmt.Sprint(v.Port), envMap)
        password := renderValue(v.Password, envMap)
        
        // defaultRedisDb := os.Getenv("go4_qa_db_defaultRedisDb")

        c, err := redigo.Dial("tcp", ip + ":" + fmt.Sprint(port))
        if err != nil {
            panic(err)
        }

        // defer c.Close()

        if len(password) > 0 {
            if _, err = c.Do("AUTH", password); err != nil {  
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
                r, _ := redigo.String(res, nil)

                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = r
            }
        case "GET":
            res, err = c.Do(redisExec.CmdStr, redisExec.CmdKey) 
            if err == nil {
                r, _ := redigo.String(res, nil)

                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = r
            }
        case "DEL":
            res, err = c.Do(redisExec.CmdStr, redisExec.CmdKey) 
            if err == nil {
                r, _ := redigo.Int(res, nil)

                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = r
            }
        case "EXISTS":
            res, err = c.Do(redisExec.CmdStr, redisExec.CmdKey) 
            if err == nil {
                r, _ := redigo.Int(res, nil)

                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = r
            }
        case "KEYS":
            res, err = c.Do(redisExec.CmdStr, redisExec.CmdKey) 
            if err == nil {
                r, _ := redigo.Strings(res, nil)
    
                redisExec.CmdAffectedCount = 1
                redisExec.CmdResults = r
            }
        default:
            redisExec.CmdAffectedCount = -1
            fmt.Println("!! Warning, Command ", redisExec.CmdStr, " is not supported currently, will enhance it later")
    }

    return err
}


func renderValue (jsonStr string, feeder map[string]string) string {
    s := jsonStr
    
    for key, value := range feeder {
        k := "${" + key + "}"
        if k == s {
            s = strings.Replace(s, k, value, -1)
        }
    }

    return s
}


