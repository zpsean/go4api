/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package texttmpl

import (
    // "fmt"
    // "io/ioutil"                                                                                                                                              
    "os"
    // "strings"
    "bytes"
    "text/template"
    // "time"

    "go4api/utils"
)

type ResultsJs struct {
    SetUpStartUnixNano int64
    SetUpStart   string
    SetUpEndUnixNano int64
    SetUpEnd  string

    NormalStartUnixNano int64
    NormalStart   string
    NormalEndUnixNano int64
    NormalEnd  string

    TearDownStartUnixNano int64
    TearDownStart   string
    TearDownEndUnixNano int64
    TearDownEnd  string

    GStartUnixNano int64
    GStart   string
    GEndUnixNano int64
    GEnd  string

    TcResults string
}

type DetailsJs struct {
    StatsStr string
}

type StatsJs struct {
    StatsStr_1 string
    StatsStr_2 string
    StatsStr_Success string
    StatsStr_Fail string
    StatsStr_Status string
}

type MStatsJs struct {
    StatsStr_1 string
    StatsStr_2 string
    StatsStr_3 string
}

type GraphicJs struct {
    Circles string
    PriorityLines string
    ParentChildrenLines string
}


func GetTemplateFromString() {
    type Inventory struct {
        Material string
        Count    uint
    }
    sweaters := Inventory{"wool", 17}
    tmpl := template.Must(template.New("test").Parse("{{.Count}} of {{.Material}} \n"))

    err := tmpl.Execute(os.Stdout, sweaters)
    if err != nil {
      panic(err) 
    }
}


func GenerateDetailsJs(strVar string, targetFile string, detailsJs *DetailsJs) {
    outFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()
    //
    tmpl := template.Must(template.New("HtmlJsCss").Parse(strVar))

    err = tmpl.Execute(outFile, *detailsJs)
    if err != nil {
      panic(err) 
    }
}

func GenerateResultsJs(strVar string, targetFile string, resultsJs *ResultsJs) {
    outFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()
    //
    tmpl := template.Must(template.New("HtmlJsCss").Parse(strVar))

    err = tmpl.Execute(outFile, *resultsJs)
    if err != nil {
      panic(err) 
    }
}


func GenerateStatsJs(strVar string, targetFile string, resultsJs []string) {
    outFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()
    //
    tmpl := template.Must(template.New("HtmlJsCss").Parse(strVar))

    statsJs := StatsJs {
        StatsStr_1: resultsJs[0],
        StatsStr_2: resultsJs[1],
        StatsStr_Success: resultsJs[2],
        StatsStr_Fail: resultsJs[3],
        StatsStr_Status: resultsJs[4],
    }

    err = tmpl.Execute(outFile, statsJs)
    if err != nil {
      panic(err) 
    }
}


func GenerateMutationResultsJs(strVar string, targetFile string, resultsJs []string) {
    outFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()
    //
    tmpl := template.Must(template.New("HtmlJsCss").Parse(strVar))

    mStatsJs := MStatsJs {
        StatsStr_1: resultsJs[0],
        StatsStr_2: resultsJs[1],
        StatsStr_3: resultsJs[2],
    }

    err = tmpl.Execute(outFile, mStatsJs)
    if err != nil {
      panic(err) 
    }
}

func GenerateGraphicJs(strVar string, targetFile string, resultsJs []string) {
    outFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()
    //
    tmpl := template.Must(template.New("HtmlJsCss").Parse(strVar))

    graphicJs := GraphicJs {
        Circles: resultsJs[0],
        PriorityLines: resultsJs[1],
        ParentChildrenLines: resultsJs[2],
    }

    err = tmpl.Execute(outFile, graphicJs)
    if err != nil {
      panic(err) 
    }
}


func GenerateJsonBasedOnTemplateAndCsv(jsonFilePath string, testData map[string]interface{}) *bytes.Buffer {
    jsonTemplateBytes := utils.GetContentFromFile(jsonFilePath)
    //
    tcJson := GetTcJson(string(jsonTemplateBytes), testData)

    return tcJson
}

func GetTcJson (jsonTemplate string, testData map[string]interface{}) *bytes.Buffer {
    tmpl := template.Must(template.New("tcTemp").Parse(jsonTemplate))
    
    tcJson := &bytes.Buffer{}
    // Execute the template
    err := tmpl.Execute(tcJson, testData)
    if err != nil {
      panic(err) 
    }

    return tcJson
}


