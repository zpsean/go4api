/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package statechart

import (
    // "fmt"
    "os"
    // "sync"
    // "strings"
    // "bufio"
    // "io"
    // "path/filepath"
)


func GenerateKeywordFile (strVar string, filePath string) {
    outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()

    outFile.WriteString(strVar)
}

func KWSettingsStr () {

}

func KWTestCasesStr () {
    
}