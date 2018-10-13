/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package reports

import (
    // "fmt"
    "os"
    // "time"
    // "strings"

    "go4api/lib/testcase"
    "go4api/ui"     
    "go4api/ui/js"  
    "go4api/ui/style"                                                                                                                                
    "go4api/utils"
    // "go4api/texttmpl"
)

// phaseEnd_time := time.Now()
// phaseEnd_str := phaseEnd_time.Format("2006-01-02 15:04:05.000000000 +0800 CST")

var (
    ExecutionResultSlice []*testcase.TcReportResults
    resultsJss = map[string]interface{}{} 
)


func GenerateTestReport(gStart_str string, gEnd_str string, resultsDir string, resultsLogFile string) {
    // html
    GenerateHtml(resultsDir)
    // style
    GenerateStyle(resultsDir)
    // js
    GenerateJs(gStart_str, gEnd_str, resultsDir, resultsLogFile)
    //
    // statsJsonBytes, _ := json.MarshalIndent(ExecutionResultSlice, "", "\t")
    // fmt.Println("ExecutionResultSlice: ", string(statsJsonBytes))
    // Get_Stats_3()
}

func GenerateHtml (resultsDir string) {
    utils.GenerateFileBasedOnVarOverride(ui.Index, resultsDir + "index.html")
    utils.GenerateFileBasedOnVarOverride(ui.Graphic, resultsDir + "graphic.html")
    utils.GenerateFileBasedOnVarOverride(ui.Details, resultsDir + "details.html")
    utils.GenerateFileBasedOnVarOverride(ui.Fuzz, resultsDir + "fuzz.html")
    utils.GenerateFileBasedOnVarOverride(ui.Mutation, resultsDir + "mutation.html")
    utils.GenerateFileBasedOnVarOverride(ui.MIndex, resultsDir + "mindex.html")
}

func GenerateStyle (resultsDir string) {
    err := os.MkdirAll(resultsDir + "style", 0777)
    if err != nil {
      panic(err) 
    }
    utils.GenerateFileBasedOnVarOverride(style.Style, resultsDir + "style/go4api.css")

    bytes := utils.DecodeBase64(style.LogoSmall)
    utils.GeneratePicture(bytes, resultsDir + "style/logosmall.png")

    bytes = utils.DecodeBase64(style.Logo)
    utils.GeneratePicture(bytes, resultsDir + "style/logo.png")

    bytes = utils.DecodeBase64(style.ArrowRight)
    utils.GeneratePicture(bytes, resultsDir + "style/arrow_right.png")

    bytes = utils.DecodeBase64(style.ArrowDown)
    utils.GeneratePicture(bytes, resultsDir + "style/arrow_down.png")
}

func GenerateJs (gStart_str string, gEnd_str string, resultsDir string, resultsLogFile string) {
    // --------
    // (0)
    err := os.MkdirAll(resultsDir + "js", 0777)
    if err != nil {
      panic(err) 
    }

    // statsFile := resultsDir + "/js/stats.js"
    // reJsons := GetStatsJson(statusCountByPriority)
    // texttmpl.GenerateStatsJs(js.Stats, statsFile, reJsons, resultsLogFile)

    // statsFile = resultsDir + "/js/mutationstats.js"
    // reJsons = GetMutationStatsJson()
    // texttmpl.GenerateMutationResultsJs(js.MutationStats, statsFile, reJsons, resultsLogFile)


    // get js/reslts.js
    GenerateReportsFromLogFile(resultsLogFile)
    
    // get js/go4api.js
    utils.GenerateFileBasedOnVarOverride(js.Js, resultsDir + "js/go4api.js")

    // 3rd js
    utils.GenerateFileBasedOnVarOverride(js.Chart, resultsDir + "js/Chart.bundle.min.js")

}




