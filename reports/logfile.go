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
    "fmt"
	"os"
    "bufio"
    "strings"
    "path/filepath"
    "encoding/json"
    "io"

    // "go4api/utils"
    // "go4api/ui/js"
    // "go4api/lib/testcase"
    // "go4api/texttmpl"

)

// this function is called by cmd -report, to generate report from log file
func GenerateReportsFromLogFile(resultsLogFile string) {
    resultsDir := filepath.Dir(resultsLogFile) + "/"
    tcReportSlice := ParseLogFile(resultsLogFile)

    fmt.Println(resultsDir, tcReportSlice)
    // html
    GenerateHtml(resultsDir)
    // style
    GenerateStyle(resultsDir)
    // js
    setUpResultSlice, normalResultSlice, tearDownResultSlice := tcReportSlice.ClassifyResults()
    fmt.Println(setUpResultSlice, normalResultSlice, tearDownResultSlice)
    // GenerateJs(resultsDir, resultsLogFile)

}

func ParseLogFile (resultsLogFile string) TcReportSlice {
    var tcReportSlice TcReportSlice
    //
    jsonLinesStr := "["

    f, err := os.Open(resultsLogFile)
    defer f.Close()

    if nil == err {
        buff := bufio.NewReader(f)
        for {
            line, err := buff.ReadString('\n')
            if err != nil || io.EOF == err {
                break
            }

            if strings.TrimSpace(line) != "" {
                l := strings.Replace(line, "\n", ",", -1)
                jsonLinesStr = jsonLinesStr + l
            }
        }
    }

    jsonLinesStr = strings.TrimSuffix(jsonLinesStr, ",")
    jsonLinesStr = jsonLinesStr + "]"

    json.Unmarshal([]byte(jsonLinesStr), &tcReportSlice)

    return tcReportSlice
}


