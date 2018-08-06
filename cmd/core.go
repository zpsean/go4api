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
    // "io/ioutil"                                                                                                                                              
    // "os"
    "flag"
    "go4api/utils"
)

func GetOptions() map[string]string {
    defaultTestDir := utils.GetCurrentDir()
    options := make(map[string]string)
    //
    testhome := flag.String("testhome", defaultTestDir + "/testhome", "the path which test scripts in")
    testresults := flag.String("testresults", defaultTestDir + "/testresults", "the path which test results in")
    testEnv := flag.String("testEnv", "QA", "the testEnv, i.e. dev, qa, uat, etc.")
    baseUrl := flag.String("baseUrl", "", "the baseUrl")

    //
    flag.Parse()
    //
    options["testhome"] = *testhome
    options["testresults"] = *testresults
    options["testEnv"] = *testEnv
    options["baseUrl"] = *baseUrl


    //
    return options
}


