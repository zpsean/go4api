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
    "encoding/json"

    "go4api/utils"
)

type Config map[string]*Environment

type Environment struct {
    BaseUrl string
}

var config Config

func GetBaseUrlFromConfig() string {
    if len(Opt.Testconfig) > 0 && len(Opt.TestEnv) > 0 {
        configJson := utils.GetJsonFromFile(Opt.Testconfig)
        json.Unmarshal([]byte(configJson), &config)

        return config[Opt.TestEnv].BaseUrl
    } else {
        return ""
    }
}


func GetConfig () Config {
    if len(Opt.Testconfig) > 0 && len(Opt.TestEnv) > 0 {
        configJson := utils.GetJsonFromFile(Opt.Testconfig)
        json.Unmarshal([]byte(configJson), &config)

        return config
    } else {
        return map[string]*Environment{}
    }
}