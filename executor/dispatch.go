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
    "strings"

    "go4api/cmd"
    "go4api/fuzz/fuzz"
    "go4api/fuzz/mutation"
)

func Dispatch(ch chan int, gStart_time time.Time) { 
    //
    baseUrl := GetBaseUrl(cmd.Opt)
    gStart_str := gStart_time.Format("2006-01-02 15.04.05.000000000 +0800 CST")
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
    if !cmd.Opt.IfScenario {
        if cmd.Opt.IfMutation {
            originMutationTcArray := GetOriginMutationTcArray()
            setUpTcSlice := GetSetupTcSlice(originMutationTcArray)
            RunSetup(ch, baseUrl, resultsDir, resultsLogFile, setUpTcSlice)
            //
            originMutationTcArray = GetOriginMutationTcArray()
            originNormalTc := GetNormalTcSlice(originMutationTcArray)
            mutatedTcArray := mutation.MutateTcArray(originNormalTc)
            Run(ch, baseUrl, resultsDir, resultsLogFile, mutatedTcArray)
            //
            originMutationTcArray = GetOriginMutationTcArray()
            teardownTcSlice := GetTeardownTcSlice(originMutationTcArray)
            RunTeardown(ch, baseUrl, resultsDir, resultsLogFile, teardownTcSlice)
        } else if cmd.Opt.IfFuzzTest {
            fuzz.PrepFuzzTest(gStart_time)
            //
            fuzzTcArray := GetFuzzTcArray()
            Run(ch, baseUrl, resultsDir, resultsLogFile, fuzzTcArray)
        } else {
            tcArray := GetTcArray()
            setUpTcSlice := GetSetupTcSlice(tcArray)
            RunSetup(ch, baseUrl, resultsDir, resultsLogFile, setUpTcSlice)
            //
            tcArray = GetTcArray()
            normalTcSlice := GetNormalTcSlice(tcArray)
            Run(ch, baseUrl, resultsDir, resultsLogFile, normalTcSlice)
            //
            tcArray = GetTcArray()
            teardownTcSlice := GetTeardownTcSlice(tcArray)
            RunTeardown(ch, baseUrl, resultsDir, resultsLogFile, teardownTcSlice)
        }
    } else {
        jsonFileList := GetJsonFiles()
        //
        tcArray := ConstructChildTcInfosBasedOnParentRoot(jsonFileList, "root" , "_dt")
        setUpTcSlice := GetSetupTcSlice(tcArray)
        RunSetup(ch, baseUrl, resultsDir, resultsLogFile, setUpTcSlice)
        //
        tcArray = ConstructChildTcInfosBasedOnParentRoot(jsonFileList, "root" , "_dt")
        RunScenario(ch, baseUrl, resultsDir, resultsLogFile, jsonFileList, tcArray)
        //
        tcArray = ConstructChildTcInfosBasedOnParentRoot(jsonFileList, "root" , "_dt")
        teardownTcSlice := GetTeardownTcSlice(tcArray)
        RunTeardown(ch, baseUrl, resultsDir, resultsLogFile, teardownTcSlice)
        //
        RunFinalReport(ch, gStart_str, resultsDir, resultsLogFile)
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

