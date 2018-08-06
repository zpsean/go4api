/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package logger

import (
    // "fmt"
	"os"
	// "log"
	// "io/ioutil"
	// "go4api/utils"
)

func WriteExecutionResults(resultString string, pStart string, resultsDir string) {
    // Note: potential bug, as maybe to much write happens at a time
    file, err := os.OpenFile(resultsDir + pStart + ".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
      panic(err) 
    }
    defer file.Close()
 
    file.WriteString(resultString + "\n")
}
