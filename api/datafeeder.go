/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (                                                                                                                                             
    "os"
    // "fmt"
    "strings"
    // "encoding/json"

    gsession "go4api/lib/session"
    // "go4api/cmd"
)

// data lookup sequence, latter override former one(s)
// config (json) -> env variables (key, value) -> session (key, value) -> data file (*_dt) / data file (inputs)

// Note: here may occur: fatal error: concurrent map iteration and map write, => need to fix
func (tcDataStore *TcDataStore) MergeTestData () map[string]interface{} {
    var finalMap = make(map[string]interface{})
    // check if config

    // 2
    envMap := GetOsEnviron()
    for k, v := range envMap {
        finalMap[k] = v
    }

    globalVariables := gsession.LoopGlobalVariables()
    for k, v := range globalVariables {
        finalMap[k] = v
    }

    // 3
    pTcName := tcDataStore.TcData.TestCase.ParentTestCase()
    pSssionMap := gsession.LookupTcSession(pTcName)
    for k, v := range pSssionMap {  // => concurrent map iteration and map write?
        finalMap[k] = v
    }

    tcName := tcDataStore.TcData.TestCase.TcName()
    sessionMap := gsession.LookupTcSession(tcName)
    for k, v := range sessionMap { // => concurrent map iteration and map write?
        finalMap[k] = v
    }

    tcLocalVariables := tcDataStore.TcLocalVariables
    for k, v := range tcLocalVariables {
        finalMap[k] = v
    }

    // fmt.Println("")
    // fmt.Println("finalMap: ", finalMap)

    // var ff interface{}
    // ss, _ :=  json.Marshal(tcDataStore.TcLocalVariables)
    // fmt.Println("TcLocalVariables: ", string(ss))
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

