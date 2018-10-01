/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package executor

import (                                                                                                                                             
    "os"
    // "fmt"
    "strings"
    // "regexp"
    // "path/filepath"
    
    // "go4api/utils"
    // "go4api/lib/csv"
    // "go4api/cmd"
    "go4api/lib/session"
)

// data lookup sequence, latter override former one(s)
// config (json) -> env variables (key, value) -> session (key, value) -> data file (*_dt) / data file (inputs)

func MergeTestData (csvHeader []string, csvRow []interface{}) map[string]interface{} {
    var finalMap = make(map[string]interface{})
    // check if config

    // 2
    envMap := GetOsEnviron()
    for k, v := range envMap {
        finalMap[k] = v
    }

    // 3
    sessionMap := gsession.LookupParentSession("")
    for k, v := range sessionMap {
        finalMap[k] = v
    }

    // 4
    dtMap := ConvertCsvRowToMap(csvHeader, csvRow)
    for k, v := range dtMap {
        finalMap[k] = v
    }

    return finalMap
}


// 1. from cmd, getconfig()


// 2. env variables (key, value), with prefix go4_
func GetOsEnviron () map[string]string {
    csvMap := map[string]string{}
    // consider add the env variables with prefix "go4_*" for username/password/athentication, etc.
    var envArray []string

    envArray = os.Environ()
    for _, env := range envArray {
        // find out the first = position, to get the key
        env_k := strings.Split(env, "=")[0]
        if strings.HasPrefix(env_k, "go4_") {
            if strings.TrimLeft(env_k, "go4_") != "" {
                csvMap[strings.TrimLeft(env_k, "go4_")] = os.Getenv(env_k)
            }
        } 
    }

    return csvMap
}


// 3. session, if parent has seesion, all direct child would have it (mainly for scenario)


// 4. data file (*_dt) / data file (inputs)
func ConvertCsvRowToMap (csvHeader []string, csvRow []interface{}) map[string]interface{} {
    csvMap := map[string]interface{}{}

    for i, item := range csvRow {
        csvMap[csvHeader[i]] = item
    }

    return csvMap
}

