/*
 * go4api - a api testing tool written in Go
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
    "strings"
    "bytes"
    "text/template"
    // "time"

    "go4api/utils"
)

type ResultsJs struct {
    PStart_time int64
    PStart   string
    PEnd_time int64
    PEnd  string
    TcReportStr string
}

type StatsJs struct {
    StatsStr string
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


func GenerateStatsJs(strVar string, targetFile string, statsJs *StatsJs, logResultsFile string) {
    outFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()
    //
    tmpl := template.Must(template.New("HtmlJsCss").Parse(strVar))

    err = tmpl.Execute(outFile, *statsJs)
    if err != nil {
      panic(err) 
    }
}

func GenerateResultsJs(strVar string, targetFile string, resultsJs *ResultsJs, logResultsFile string) {
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


func GenerateJsonBasedOnTemplateAndCsv(jsonFilePath string, csvHeader []string, csvRow []string) *bytes.Buffer {
    jsonTemplateBytes := utils.GetContentFromFile(jsonFilePath)
    //
    tcJson := GetTcJson(string(jsonTemplateBytes), csvHeader, csvRow)

    return tcJson
}

func GetTcJson (jsonTemplate string, csvHeader []string, csvRow []string) *bytes.Buffer {
    csvMap := map[string]string{}

    tmpl := template.Must(template.New("tcTemp").Parse(jsonTemplate))
    
    // consider add the env variables with prefix "go4_*" for username/password/athentication, etc.
    csvMap = GetOsEnviron()
    // override the env variables, using the csv data
    for i, item := range csvRow {
        csvMap[csvHeader[i]] = item
    }
    //
    tcJson := &bytes.Buffer{}
    // Execute the template
    err := tmpl.Execute(tcJson, csvMap)
    if err != nil {
      panic(err) 
    }

    return tcJson
}

func GetOsEnviron () map[string]string {
    csvMap := map[string]string{}
    // consider add the env variables with prefix "go4_*" for username/password/athentication, etc.
    var envArray []string
    envArray = os.Environ()
    for _, env := range envArray {
        // find out the first = position, to get the key
        env_k := strings.Split(env, "=")[0]
        if strings.HasPrefix(env_k, "go4_") {
            if strings.TrimLeft(env_k, "go4_") != "" {
                csvMap[strings.TrimLeft(env_k, "go4_")] = os.Getenv(env_k)
            }
        } 
    }

    return csvMap
}


