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
    "fmt"
    // "io/ioutil"                                                                                                                                              
    "os"
    "flag"
    "go4api/utils"
)

var (
    h bool
    run bool
    convert bool
)

type Options struct {
    Testconfig string
    Testcase string 
    Testresource string
    Testresults string
    TestEnv string
    BaseUrl string
    IfScenario  bool
    IfFuzzTest  bool
    IfMutation  bool
    ConcurrencyLimit int

    Harfile string
    Swaggerfile string
}

var Opt Options

// Note: as refer to https://golang.org/doc/effective_go.html#init
// each file can have one or more init(), the init() will be run after all var evaluated
// import --> const --> var --> init()
func init() {
    defaultTestDir := utils.GetCurrentDir()
    //
    flag.BoolVar(&h, "h", false, "this help")
    flag.BoolVar(&run, "run", false, "")
    flag.BoolVar(&convert, "convert", false, "")
    //
    testconfig := flag.String("c", defaultTestDir + "/testconfig", "the path which test config in")
    testcase := flag.String("tc", defaultTestDir + "/testcase", "the path which test json in")
    testresource := flag.String("tr", defaultTestDir + "/testresource", "the path which test resource in")
    testresults := flag.String("r", defaultTestDir + "/testresults", "the path which test results in")
    testEnv := flag.String("testEnv", "QA", "the testEnv, i.e. dev, qa, uat, etc.")
    baseUrl := flag.String("baseUrl", "", "the baseUrl")
    ifScenario := flag.Bool("S", false, "if the target cases are for scenarios, which have data dependency")
    ifFuzzTest := flag.Bool("F", false, "if to run the Fuzz test")
    ifMutation := flag.Bool("M", false, "if to run the Mutation test")
    concurrency := flag.Int("cl", 100, "concurrency limitation")

    har := flag.String("harfile", "", "har file name to be converted")
    swagger := flag.String("swaggerfile", "", "har file name to be converted")

    //
    flag.Parse()
    //
    Opt.Testconfig = *testconfig
    Opt.Testcase = *testcase
    Opt.Testresource = *testresource
    Opt.Testresults = *testresults
    Opt.TestEnv = *testEnv
    Opt.BaseUrl = *baseUrl
    Opt.IfScenario = *ifScenario
    Opt.IfFuzzTest = *ifFuzzTest
    Opt.IfMutation = *ifMutation
    Opt.ConcurrencyLimit = *concurrency

    Opt.Harfile = *har
    Opt.Swaggerfile = *swagger

    if h {
        usage()
    }

    // flag.Usage = usage
}

func usage() {
    fmt.Fprintf(os.Stderr, `
go4api version: 0.12.0

Usage:
  go4api [command] [options]

Available Commands:
  run         Start a load test
  convert     Convert a HAR file / Swagger API file to a go4api Json test case

Command: run
Usage: go4api -run [-?hFMS] [-c config filename] [-t testcase path] [-d test resource path] [-r test results path] 

Options:

Command: convert
Usage: go4api -convert [-harfile har filename] [-swaggerfile swagger api filename]

Options:
`)
    flag.PrintDefaults()

    os.Exit(0)
}
