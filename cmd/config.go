/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package cmd

import (
    "fmt"
    "encoding/json"

    "go4api/utils"
)

var config Config
var tEnv string

type Config map[string]*Environment

type Environment struct {
    BaseUrl string
    TimeZone string
    Dbs map[string]*DbDetails
    PostgreSql map[string]*DbDetails
    Redis map[string]*RedisDetails
    MongoDB map[string]*MongDBDetails
}

//
type DbDetails struct {
    SqlCon interface{}
    Ip string
    Port interface{}
    UserName string
    Password string
    Dbname   string
    Sslmode  string
}

type RedisDetails struct {
    RedisCon interface{}
    Ip string
    Port interface{}
    UserName string
    Password string
}

type MongDBDetails struct {
    MongDBCon interface{}
    Ip string
    Port interface{}
    UserName string
    Password string
}

//
func SetTestEnv () {
    if Opt.TestEnv != "" {
        tEnv = Opt.TestEnv
    } else {
        tEnv = "QA"
    }
}
    
func GetConfig () Config {
    if len(Opt.Testconfig) > 0 {
        configJson := utils.GetJsonFromFile(Opt.Testconfig)

        e := json.Unmarshal([]byte(configJson), &config)
        if e != nil {
            fmt.Println("Unmarshal error: ")
            panic(e)
        }

        return config
    } else {
        return map[string]*Environment{}
    }
}

func GetBaseUrlFromConfig () string {
        return config[tEnv].BaseUrl
}

func GetTimeZoneConfig () string {
        return config[tEnv].TimeZone
}

// mysql
func GetDbConfig () map[string]*DbDetails {
        return config[tEnv].Dbs
}

// postgresql
func GetPgDbConfig () map[string]*DbDetails {
        return config[tEnv].PostgreSql
}

func GetRedisConfig () map[string]*RedisDetails {
        return config[tEnv].Redis
}

func GetMongoDBConfig () map[string]*MongDBDetails {
        return config[tEnv].MongoDB
}


