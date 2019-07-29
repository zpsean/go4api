/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package reports

import (
	"os"
)


func OpenExecutionResultsLogFile(logFile string) *os.File {
    file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
      panic(err) 
    }
    // defer file.Close()

 	return file
}

func WriteExecutionResults(resultString string, file *os.File) { 
    file.WriteString(resultString + "\n")
}



func CloseExecutionResultsLogFile(file *os.File) {
    // Note: potential bug, as maybe to much write happens at a time
    // defer file.Close()
 
    file.Close()
}
