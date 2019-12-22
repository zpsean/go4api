/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package executor

import (
    "fmt"
    "time"
    "os"
    "strings"
    // "encoding/json"

    "go4api/cmd"
    // "go4api/fuzz"
    "go4api/mutation"
    "go4api/db/mysql"
    "go4api/db/postgres"
    "go4api/db/redis"
    "go4api/db/mongodb"
)

func Dispatch(ch chan int, gStart_time time.Time, gStart_str string) { 
    //
    baseUrl := GetBaseUrl(cmd.Opt)
    // make results dir
    resultsDir := MkResultsDir(gStart_str, cmd.Opt)
    resultsLogFile := resultsDir + gStart_str + ".log"
    //
    // <!!--> Note: there are two kinds of test cases dependency:
    // type 1. the parent and child has only execution dependency, no data exchange
    // type 2. the parent and child has execution dependency and data exchange dynamically
    // for type 1, the json is rendered by data tables first, then build the tcTree
    // for type 2, build the cases hierarchy first, then render the child cases using the parent's outputs
    //
    WarmUpDBConnection()
    WarmUpPgDbConnection()
    WarmUpRedisConnection()
    WarmUpMongoDBConnection()
    //
    g4Store := InitG4Store()
    //
    g4Store.GlobalSetUpRunStore.InitRun()
    g4Store.GlobalSetUpRunStore.RunPriorities(baseUrl, resultsLogFile)
    g4Store.GlobalSetUpRunStore.RunConsoleOverallReport()
    //
    g4Store.NormalRunStore.InitRun()
    g4Store.NormalRunStore.RunPriorities(baseUrl, resultsLogFile)
    g4Store.NormalRunStore.RunConsoleOverallReport()
    //
    if cmd.Opt.IfMutation {
        mutatedTcArray := mutation.MutateTcArray(g4Store.NormalRunStore.TcDs)
        g4Store.MutationRunStore.TcSlice = mutatedTcArray

        g4Store.MutationRunStore.InitRun()
        g4Store.MutationRunStore.RunPriorities(baseUrl, resultsLogFile)
        g4Store.MutationRunStore.RunConsoleOverallReport()
    }
    //
    g4Store.GlobalTeardownRunStore.InitRun()
    g4Store.GlobalTeardownRunStore.RunPriorities(baseUrl, resultsLogFile)
    g4Store.GlobalTeardownRunStore.RunConsoleOverallReport()
    //
    g4Store.RunFinalConsoleReport()
    g4Store.RunFinalReport(ch, gStart_str, resultsDir, resultsLogFile)
}


func GetBaseUrl(opt cmd.Options) string {
    baseUrl := ""
    if cmd.Opt.BaseUrl != "" {
        baseUrl = cmd.Opt.BaseUrl
    } else {
        baseUrl = cmd.GetBaseUrlFromConfig() 
    }
    if baseUrl == "" {
        fmt.Println("Warning: baseUrl is not set")
    } else {
        fmt.Println("baseUrl set to: " + baseUrl)
    }

    return baseUrl
}


func MkResultsDir(gStart_str string, opt cmd.Options) string {
    var resultsDir string

    if strings.HasSuffix(strings.TrimSpace(cmd.Opt.Testresults), "/") {
        resultsDir = cmd.Opt.Testresults + gStart_str + "/"
    } else {
        resultsDir = cmd.Opt.Testresults + "/" + gStart_str + "/"
    }

    err := os.MkdirAll(resultsDir, 0777)
    if err != nil {
      panic(err) 
    } 

    return resultsDir
}

func WarmUpDBConnection () {
    if cmd.Opt.IfSqlDb == true {
        gsql.InitConnection()
    }
}

func WarmUpPgDbConnection () {
    if cmd.Opt.IfPgDb == true {
        gpg.InitConnection()
    }
}

func WarmUpRedisConnection () {
    if cmd.Opt.IfRedis == true {
        gredis.InitRedisConnection()
    }
}

func WarmUpMongoDBConnection () {
    if cmd.Opt.IfMongoDB == true {
        gmongodb.InitMongoDBConnection()
    }
}


