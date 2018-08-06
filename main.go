/*
 * go4api - a api testing tool written in Go
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
  "go4api/executor" 
  "os"
	"time"
  "go4api/cmd"
  // "go4api/utils"
)

func main(){

    var Version = "0.10.1"
    var Banner = `
     _______                 ____          ___                   _
    /  ____  \   ______     / __ |        / _ \       ______    |_|
   /  /    \__| /  __  \   / / | |       / / \ \     |  ___ \    _
   |  |   ____  | |  | |  / /__| |__    / /   \ \    | |   \ \  | |
   |  |  |_  _| | |  | | |_____ ____|  / /_____\ \   | |___/ |  | |
   |  \____| |  | |__| |       | |    / ________\ \  |  ____/   | |
    \________/   \____/        |_|   /_/         \_\ | |        |_| 
                                                     |_|            `

    fmt.Println("\nVersion: ", Version)
    fmt.Println(Banner)

    //get the cmd options
    options := cmd.GetOptions()

    // fmt.Println(os.Args)
    ch := make(chan int, 1)

    fmt.Println("\n----- Start Main -----")

    pStart := string(time.Now().String())
    // fmt.Println("pStart: ", pStart)
    //
    executor.Run(ch, pStart, options)

    x := <-ch
    fmt.Println("----- Finish Main -----\n")

    close(ch)
    // this exit code to be used for CI/CD
    os.Exit(x)
}
