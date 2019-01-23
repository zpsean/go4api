/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
    "os"
    "fmt"
    "strings"
    "reflect"
    "path/filepath"
    "encoding/json"
    "encoding/csv"

    "go4api/lib/testcase"
    "go4api/utils"

    gjson "github.com/tidwall/gjson"
    xlsx "github.com/tealeg/xlsx"
)


func (tcDataStore *TcDataStore) HandleOutFiles (expOutFiles []*testcase.OutFilesDetails) {
    tcData := tcDataStore.TcData
    actualBody := tcDataStore.HttpActualBody

    if len(expOutFiles) > 0 {
        // get the actual value from actual body based on the fields in json outputs
        for i, _ := range expOutFiles {
            var keyStrList []string
            var valueStrList []string
            // (1). get the full path of outputsfile
            tempDir := utils.CreateTempDir(tcData.JsonFilePath)
            outputsFile := filepath.Join(tempDir, expOutFiles[i].GetTargetFileName())
            os.Remove(outputsFile)

            // (2). get operation
            outFilesOperation := expOutFiles[i].GetOperation()
            outFilesData := expOutFiles[i].GetData()

            switch strings.ToLower(outFilesOperation) {
            case "DataToCsv":
                keyStrList, valueStrList = tcDataStore.GetOutputsCsvData(outFilesData)
                // write csv header
                utils.GenerateCsvFileBasedOnVarOverride(keyStrList, outputsFile)
                // write csv data
                utils.GenerateCsvFileBasedOnVarAppend(valueStrList, outputsFile)
            case "ObjectToExcel":
                SaveHttpRespFile(actualBody, outputsFile)
            case "jsontocsv":
                sources := expOutFiles[i].GetSources()
                sourcesFields := expOutFiles[i].GetSourcesFields()
                targetHeader := expOutFiles[i].GetTargetHeader()

                SaveJsonToCsvFile(tcDataStore, sources, sourcesFields, targetHeader, outputsFile)
            case "jsontoexcel":
                sources := expOutFiles[i].GetSources()
                sourcesFields := expOutFiles[i].GetSourcesFields()
                targetHeader := expOutFiles[i].GetTargetHeader()

                SaveJsonToExcelFile(tcDataStore, sources, sourcesFields, targetHeader, outputsFile)
            case "UnionExcel":
                SaveHttpRespFile(actualBody, outputsFile)
            case "UnionCsv":
                SaveHttpRespFile(actualBody, outputsFile)
            case "CsvToExcel":
                SaveHttpRespFile(actualBody, outputsFile)
            case "ExcelToCsv":
                SaveHttpRespFile(actualBody, outputsFile)
            }   
        } 
    }
}


func SaveJsonToCsvFile (tcDataStore *TcDataStore, sources []string, sourcesFields []string, targetHeader []string, outputsFile string) {
    outFile := OpenOutFileForAppend(outputsFile)
    w := csv.NewWriter(outFile)

    w.Write(targetHeader)
 
    jsonValue := tcDataStore.GetResponseValue(sources[0])
    valueType := reflect.TypeOf(jsonValue).Kind().String()

    // http returns []bytes -> string, sql returns []map[string]interface{}
    jsonStr := ""
    if valueType == "string" {
        jsonStr = jsonValue.(string)
    } else {
        jsonBytes, _ := json.Marshal(jsonValue)
        jsonStr = string(jsonBytes)
    }
    // fmt.Println("jsonStr: ", jsonStr)

    jsonLength := int(gjson.Get(jsonStr, "#").Int())

    for i := 0; i < jsonLength; i++ {
        if len(sourcesFields) == 0 {
        } else {
            var tmpSlice []string
            for ii, _ := range sourcesFields {
                fValue := gjson.Get(jsonStr, fmt.Sprint(i) + "." + sourcesFields[ii])
                tmpSlice = append(tmpSlice, fValue.String())
            }

            w.Write(tmpSlice)
        }
    }

    w.Flush()
    outFile.Close()
}

func SaveJsonToExcelFile (tcDataStore *TcDataStore, sources []string, sourcesFields []string, targetHeader []string, outputsFile string) {
    var file *xlsx.File
    var sheet *xlsx.Sheet
    var row *xlsx.Row
    var cell *xlsx.Cell
    var err error

    file = xlsx.NewFile()
    sheet, err = file.AddSheet("Sheet1")
    if err != nil {
        fmt.Printf(err.Error())
    }

    row = sheet.AddRow()
    for i, _ := range targetHeader {
        cell = row.AddCell()
        cell.Value = targetHeader[i]
    }

    jsonValue := tcDataStore.GetResponseValue(sources[0])
    valueType := reflect.TypeOf(jsonValue).Kind().String()

    // http returns []bytes -> string, sql returns []map[string]interface{}
    jsonStr := ""
    if valueType == "string" {
        jsonStr = jsonValue.(string)
    } else {
        jsonBytes, _ := json.Marshal(jsonValue)
        jsonStr = string(jsonBytes)
    }

    jsonLength := int(gjson.Get(jsonStr, "#").Int())

    for i := 0; i < jsonLength; i++ {
        row = sheet.AddRow()

        if len(sourcesFields) == 0 {
        } else {
            for ii, _ := range sourcesFields {
                fValue := gjson.Get(jsonStr, fmt.Sprint(i) + "." + sourcesFields[ii])

                cell = row.AddCell()
                cell.Value = fValue.String()
            }
        }
    }

    err = file.Save(outputsFile)
    if err != nil {
        fmt.Printf(err.Error())
    }
}


func OpenOutFileForAppend(logFile string) *os.File {
    file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
      panic(err) 
    }

    return file
}

