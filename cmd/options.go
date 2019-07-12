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
    report bool
)

type Options struct {
    Testconfig string
    Testsuite  string 
    Testcase   string 
    KeyWord    string 

    Testresource string
    Testresults  string
    JsFuncs      string
    TestEnv      string
    BaseUrl      string

    IfFuzzTest    bool
    IfMutation    bool
    IfTestSuite   bool
    IfKeyWord     bool

    IfConcurrency       bool
    ConcurrencyLimit    int
    IfShowOriginRequest bool

    Harfile      string
    Swaggerfile  string
    Logfile      string

    TimeZone string

    IfSqlDb bool
    IfCache bool
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
    flag.BoolVar(&report, "report", false, "Generate report only from log file")
    //
    testconfig := flag.String("c", defaultTestDir + "/testconfig", "the path which test config in")
    testsuite := flag.String("tsuite", defaultTestDir + "/testsuite", "the path which testsuite json in")
    testcase := flag.String("tc", defaultTestDir + "/testcase", "the path which test json in")
    keyword := flag.String("kw", defaultTestDir + "/keyword", "the path which keyword in")

    testresource := flag.String("tr", defaultTestDir + "/testresource", "the path which test resource in")
    testresults := flag.String("r", defaultTestDir + "/testresults", "the path which test results in")
    js := flag.String("js", defaultTestDir + "/js", "the path which functions defined with js in")
    testEnv := flag.String("testEnv", "QA", "the testEnv, i.e. dev, qa, uat, etc.")
    baseUrl := flag.String("baseUrl", "", "the baseUrl")

    ifFuzzTest := flag.Bool("F", false, "if to run the Fuzz test")
    ifMutation := flag.Bool("M", false, "if to run the Mutation test")
    ifTestSuite := flag.Bool("TS", false, "if to run with keyword driven / testsuite mode")
    ifKeyWord := flag.Bool("K", false, "if to run with keyword driven / testsuite mode")

    ifConcurrency := flag.Bool("ifCon", true, "if to run the with concurrency mode")
    concurrency := flag.Int("cl", 100, "concurrency limitation")
    ifShowOriginRequest := flag.Bool("ifOriginReq", false, "if to show origin request, be careful, it may expose confidential info")

    har := flag.String("harfile", "", "har file name to be converted")
    swagger := flag.String("swaggerfile", "", "har file name to be converted")

    logfile := flag.String("logfile", "", "log file for report generation")

    timeZone := flag.String("timeZone", "", "timezone used, GMT+/-N:00")

    ifSqlDb := flag.Bool("ifSqlDb", true, "if test has Sql, such as MySql")
    ifCache := flag.Bool("ifCache", true, "if test has cache, such as Redis")

    //
    flag.Parse()
    //
    Opt.Testconfig = *testconfig
    Opt.Testsuite = *testsuite
    Opt.Testcase = *testcase
    Opt.KeyWord = *keyword

    Opt.Testresource = *testresource
    Opt.Testresults = *testresults
    Opt.JsFuncs = *js
    Opt.TestEnv = *testEnv
    Opt.BaseUrl = *baseUrl
    
    Opt.IfFuzzTest = *ifFuzzTest
    Opt.IfMutation = *ifMutation
    Opt.IfTestSuite = *ifTestSuite
    Opt.IfKeyWord = *ifKeyWord

    Opt.IfConcurrency = *ifConcurrency
    Opt.ConcurrencyLimit = *concurrency
    Opt.IfShowOriginRequest = *ifShowOriginRequest

    Opt.Harfile = *har
    Opt.Swaggerfile = *swagger
    Opt.Logfile = *logfile

    Opt.TimeZone = *timeZone

    Opt.IfSqlDb = *ifSqlDb
    Opt.IfCache = *ifCache

    if h {
        usage()
    }

    // flag.Usage = usage
    SetTestEnv()
    GetConfig()
}

func usage() {
    fmt.Fprintf(os.Stderr, `
go4api version: 0.20.0

Usage:
  go4api [command] [options]

Available Commands:
  run         Start a test
  convert     Convert a HAR file / Swagger API file to a go4api Json test case
  report      Generate report only from log file

Command: run
Usage: go4api -run [-?hFM] [-c config filename] [-t testcase path] [-d test resource path] [-r test results path] 

Options:

Command: convert
Usage: go4api -convert [-harfile har filename] [-swaggerfile swagger api filename]

Command: report
Usage: go4api -report [-logfile log filename]

Options:
`)
    flag.PrintDefaults()

    os.Exit(0)
}
