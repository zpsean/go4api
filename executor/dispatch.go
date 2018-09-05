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
    "fmt"
    "time"
    "os"
    "go4api/cmd"
    "go4api/utils"
    "go4api/executor/fuzz"
)

func Dispatch(ch chan int, pStart_time time.Time) { 
    //
    baseUrl := GetBaseUrl(cmd.Opt)
    pStart := pStart_time.String()
    // get results dir
    resultsDir := GetResultsDir(pStart, cmd.Opt)
    //
    // <!!--> Note: there are two kinds of test cases dependency:
    // type 1. the parent and child has only execution dependency, no data exchange
    // type 2. the parent and child has execution dependency and data exchange dynamically
    // for type 1, the json is rendered by data tables first, then build the tcTree
    // for type 2, build the cases hierarchy first, then render the child cases using the parent's outputs
    //
    if !cmd.Opt.IfScenario {
        if cmd.Opt.IfMutation {
            originMutationTcArray := GetOriginMutationTcArray()
            // Run(ch, pStart_time, options, pStart, baseUrl, resultsDir, originMutationTcArray)

            // fmt.Println("\noriginMutationTcArray: ", originMutationTcArray)
            // to mutate 
            mutatedTcArray := fuzz.MutateTcArray(originMutationTcArray)
            // fmt.Println("\nmutatedTcArray: ", mutatedTcArray)
            Run(ch, pStart_time, pStart, baseUrl, resultsDir, mutatedTcArray)
        } else if cmd.Opt.IfFuzzTest {
            fuzz.PrepFuzzTest(pStart_time)

            // GetFuzzTcArray(options)
            fuzzTcArray := GetFuzzTcArray()
            Run(ch, pStart_time, pStart, baseUrl, resultsDir, fuzzTcArray)
        } else {
            tcArray := GetTcArray()
            // fmt.Println("\n tcArray: ", tcArray)
            Run(ch, pStart_time, pStart, baseUrl, resultsDir, tcArray)
        }
    } else {
        RunScenario(ch, pStart_time, pStart, baseUrl, resultsDir)
    }
}


func GetBaseUrl(opt cmd.Options) string {
    testenv := cmd.Opt.TestEnv
    baseUrl := ""
    if cmd.Opt.BaseUrl != "" {
        baseUrl = cmd.Opt.BaseUrl
    } else {
        _, err := os.Stat(cmd.Opt.Testconfig + "/config.json")
        // fmt.Println("err: ", err)
        if err == nil {
            baseUrl = utils.GetBaseUrlFromConfig(cmd.Opt.Testconfig + "/config.json", testenv) 
        }
    }
    if baseUrl == "" {
        fmt.Println("Warning: baseUrl is not set")
    } else {
        fmt.Println("baseUrl set to: " + baseUrl)
    }

    return baseUrl
}


func GetResultsDir(pStart string, opt cmd.Options) string {
    var resultsDir string
    err := os.MkdirAll(cmd.Opt.Testresults + "/" + pStart + "/", 0777)
    if err != nil {
      panic(err) 
    } else {
        resultsDir = cmd.Opt.Testresults + "/" + pStart + "/"
    }
    return resultsDir
}

