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
    // "fmt"
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
    Redis map[string]*RedisDetails
}

type DbDetails struct {
    SqlCon interface{}
    Ip string
    Port interface{}
    UserName string
    Password string
}

type RedisDetails struct {
    RedisCon interface{}
    Ip string
    Port interface{}
    UserName string
    Password string
}

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
        json.Unmarshal([]byte(configJson), &config)

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

func GetDbConfig () map[string]*DbDetails {
        return config[tEnv].Dbs
}

func GetRedisConfig () map[string]*RedisDetails {
        return config[tEnv].Redis
}


