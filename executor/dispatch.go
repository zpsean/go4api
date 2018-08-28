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
    "go4api/utils"
    "go4api/executor/fuzz"
)

func Dispatch(ch chan int, pStart_time time.Time, options map[string]string) { 
    baseUrl := GetBaseUrl(options)
    pStart := pStart_time.String()
    // get results dir
    resultsDir := GetResultsDir(pStart, options)
    //
    // <!!--> Note: there are two kinds of test cases dependency:
    // type 1. the parent and child has only execution dependency, no data exchange
    // type 2. the parent and child has execution dependency and data exchange dynamically
    // for type 1, the json is rendered by data tables first, then build the tcTree
    // for type 2, build the cases hierarchy first, then render the child cases using the parent's outputs
    //
    if options["ifScenario"] == "" {
        if options["ifMutation"] != "" {
            originMutationTcArray := GetOriginMutationTcArray(options)
            // Run(ch, pStart_time, options, pStart, baseUrl, resultsDir, originMutationTcArray)

            // fmt.Println("\noriginMutationTcArray: ", originMutationTcArray)
            // to mutate 
            mutatedTcArray := fuzz.MutateTcArray(originMutationTcArray)
            // fmt.Println("\nmutatedTcArray: ", mutatedTcArray)
            Run(ch, pStart_time, options, pStart, baseUrl, resultsDir, mutatedTcArray)
        } else if options["ifFuzzTestFirst"] != "" {
            fuzz.PrepFuzzTest(pStart_time, options)

            // GetFuzzTcArray(options)
            fuzzTcArray := GetFuzzTcArray(options)
            Run(ch, pStart_time, options, pStart, baseUrl, resultsDir, fuzzTcArray)
        } else {
            tcArray := GetTcArray(options)
            Run(ch, pStart_time, options, pStart, baseUrl, resultsDir, tcArray)
        }
        
    } else {
        RunScenario(ch, pStart_time, options, pStart, baseUrl, resultsDir)
        fmt.Println("--")
    }
}


func GetBaseUrl(options map[string]string) string {
    testenv := options["testEnv"]
    baseUrl := ""
    if options["baseUrl"] != "" {
        baseUrl = options["baseUrl"]
    } else {
        _, err := os.Stat(options["testhome"] + "/testconfig/testconfig.json")
        // fmt.Println("err: ", err)
        if err == nil {
            baseUrl = utils.GetBaseUrlFromConfig(options["testhome"] + "/testconfig/testconfig.json", testenv) 
        }
    }
    if baseUrl == "" {
        fmt.Println("Warning: baseUrl is not set")
    } else {
        fmt.Println("baseUrl set to: " + baseUrl)
    }

    return baseUrl
}


func GetResultsDir(pStart string, options map[string]string) string {
    var resultsDir string
    err := os.MkdirAll(options["testresults"] + "/" + pStart + "/", 0777)
    if err != nil {
      panic(err) 
    } else {
        resultsDir = options["testresults"] + "/" + pStart + "/"
    }

    return resultsDir
}

