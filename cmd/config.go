/*
 * go4api - a api testing tool written in Go
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
    Dbs map[string]*DbDetails
    Redis interface{}
}

type DbDetails struct {
    SqlCon interface{}
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

func GetDbConfig () map[string]*DbDetails {
        return config[tEnv].Dbs
}


