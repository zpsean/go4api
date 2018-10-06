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
    "go4api/fuzz/fuzz"
    "go4api/fuzz/mutation"
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
            //
            mutatedTcArray := mutation.MutateTcArray(originMutationTcArray)
            setUpTcSlice := GetSetupTcSlice(mutatedTcArray)
            RunSetup(ch, pStart_time, pStart, baseUrl, resultsDir, setUpTcSlice)
            //
            mutatedTcArray = mutation.MutateTcArray(originMutationTcArray)
            normalTcSlice := GetNormalTcSlice(mutatedTcArray)
            Run(ch, pStart_time, pStart, baseUrl, resultsDir, normalTcSlice)
            //
            mutatedTcArray = mutation.MutateTcArray(originMutationTcArray)
            teardownTcSlice := GetTeardownTcSlice(mutatedTcArray)
            RunTeardown(ch, pStart_time, pStart, baseUrl, resultsDir, teardownTcSlice)
        } else if cmd.Opt.IfFuzzTest {
            fuzz.PrepFuzzTest(pStart_time)
            //
            fuzzTcArray := GetFuzzTcArray()
            Run(ch, pStart_time, pStart, baseUrl, resultsDir, fuzzTcArray)
        } else {
            tcArray := GetTcArray()
            setUpTcSlice := GetSetupTcSlice(tcArray)
            RunSetup(ch, pStart_time, pStart, baseUrl, resultsDir, setUpTcSlice)
            //
            tcArray = GetTcArray()
            normalTcSlice := GetNormalTcSlice(tcArray)
            Run(ch, pStart_time, pStart, baseUrl, resultsDir, normalTcSlice)
            //
            tcArray = GetTcArray()
            teardownTcSlice := GetTeardownTcSlice(tcArray)
            RunTeardown(ch, pStart_time, pStart, baseUrl, resultsDir, teardownTcSlice)
        }
    } else {
        RunScenario(ch, pStart_time, pStart, baseUrl, resultsDir)
    }
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

