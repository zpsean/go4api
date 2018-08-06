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
    // "io"
    "bytes"
    "go4api/utils"
    "text/template"
    "path/filepath"
    "strconv"
    // "bufio"
    // simplejson "github.com/bitly/go-simplejson"
    
)

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

func GetHtmlTemplateFromFiles(file string, resultsDir string) {
    type Results struct {
        Items   []string
    }
    //
    outFile, err := os.OpenFile(resultsDir + "index.html", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()

    results := Results{
        Items:  []string{"Case 1", "Case 1"}}
 
    tmpl := template.Must(template.New(filepath.Base(file)).ParseFiles(file))

    err = tmpl.Execute(outFile, results)
    if err != nil {
      panic(err) 
    }
}

func GenerateHtmlReportFromTemplateAndVar(strVar string, resultsDir string, logResultsFile string) {
    type tcResults struct {
        Seq, CaseID, Status, CasePath string
    }
    TcResults := tcResults{}
    tcResultsList := []tcResults{}
    //
    outFile, err := os.OpenFile(resultsDir + "index.html", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()
    //
    // get the data from the log results file
    csvRows := utils.GetCsvFromFile(logResultsFile)
    // fmt.Println("csvRows: ", logResultsFile, csvRows)
    for k, csvrow := range csvRows {
        TcResults.Seq = strconv.Itoa(k + 1)
        TcResults.CaseID = csvrow[0]
        TcResults.Status = csvrow[7]
        TcResults.CasePath = csvrow[2] + " / " + csvrow[3] + " / " + csvrow[4]

        tcResultsList = append(tcResultsList, TcResults)
    }
    // fmt.Println("tcResultsList: ", tcResultsList)
    //
    tmpl := template.Must(template.New("Html").Parse(strVar))

    err = tmpl.Execute(outFile, tcResultsList)
    if err != nil {
      panic(err) 
    }
}



func GenerateJsonFileBasedOnTemplateAndCsv(jsonFilePath string, csvHeader []string, csvRow []string, tmpJsonDir string) string {
    csvMap := map[string]string{}
    //
    tmpl := template.Must(template.New(filepath.Base(jsonFilePath)).ParseFiles(jsonFilePath))

    for i, item := range csvRow {
        csvMap[csvHeader[i]] = item
    }

    // here also needs to consider add the env variables with prefix "go4_*" for username/password/athentication, etc.
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
    // os.Getenv("JAVA_HOME")

    //
    outFile, err := os.OpenFile(tmpJsonDir + filepath.Base(jsonFilePath) + ".td.json", os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()

    // Execute the template
    err = tmpl.Execute(outFile, csvMap)
    if err != nil {
      panic(err) 
    }

    return tmpJsonDir + filepath.Base(jsonFilePath) + ".td.json"
}


func GenerateJsonBasedOnTemplateAndCsv(jsonFilePath string, csvHeader []string, csvRow []string) *bytes.Buffer {
    csvMap := map[string]string{}
    //
    tmpl := template.Must(template.New(filepath.Base(jsonFilePath)).ParseFiles(jsonFilePath))

    // here also needs to consider add the env variables with prefix "go4_*" for username/password/athentication, etc.
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

    // override the env variables, using the csv data
    for i, item := range csvRow {
        csvMap[csvHeader[i]] = item
    }

    //
    fjson := &bytes.Buffer{}

    // Execute the template
    err := tmpl.Execute(fjson, csvMap)
    if err != nil {
      panic(err) 
    }

    return fjson
}


