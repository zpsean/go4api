/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package main

import (
	"fmt"
  "os"
	"time"

  "go4api/cmd"
  "go4api/reports"
  "go4api/executor"
  "go4api/converter/har"
  "go4api/converter/swagger"
)

func main(){

    var Version = "0.69.2"
    var Banner = `
     ________                 ____          ___                   _
    /  ____  \   _______     / __ |        / _ \       ______    |_|
   /  /    \__| /  ___  \   / / | |       / / \ \     |  ___ \    _
   |  |   ____  | |   | |  / /__| |__    / /   \ \    | |   \ \  | |
   |  |  |_  _| | |   | | |_____ ____|  / /_____\ \   | |___/ |  | |
   |  \____| |  | |___| |       | |    / ________\ \  |  ____/   | |
    \________/  \_______/       |_|   /_/         \_\ | |        |_| 
                                                      |_|            `

    fmt.Println("\nVersion: ", Version)
    fmt.Println(Banner)

    //get the cmd options
    ch := make(chan int, 1)

    fmt.Println("\n----- Start Main -----")

    gStart := time.Now()
    gStart_str := gStart.Format("2006-01-02 15.04.05.000000000 +0800 CST")
    //
    fmt.Println("Started at: " + gStart_str)
    // fmt.Println(os.Args)

    if os.Args[1] == "-run" {
      executor.Dispatch(ch, gStart, gStart_str)
      //
      close(ch)
      x := <-ch
      fmt.Println("----- Finish Main -----")
      // this exit code to be used for CI/CD
      os.Exit(x)
    } else if os.Args[1] == "-convert" {
      if cmd.Opt.Harfile != "" {
        har.Convert()
      } else if cmd.Opt.Swaggerfile != "" {
        swagger.Convert()
      }
    } else if os.Args[1] == "-report" {
      if cmd.Opt.Logfile != "" {
        reports.GenerateReportsFromLogFile(cmd.Opt.Logfile)
      }
    } else {
      fmt.Println("Warning: no specific commnd is provided, default is to run")
      executor.Dispatch(ch, gStart, gStart_str)
      //
      close(ch)
      x := <-ch
      fmt.Println("----- Finish Main -----")
      // this exit code to be used for CI/CD
      os.Exit(x)
    }
}


