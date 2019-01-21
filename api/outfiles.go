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
    // "fmt"
    "strings"
    // "reflect"
    "path/filepath"
    // "encoding/json"

    "go4api/lib/testcase"
    "go4api/utils"
)


func (tcDataStore *TcDataStore) HandleOutFiles () {
    var expOutFiles []*testcase.OutFilesDetails

    tcData := tcDataStore.TcData

    // actualStatusCode := tcDataStore.HttpActualStatusCode
    // actualHeader := tcDataStore.HttpActualHeader
    actualBody := tcDataStore.HttpActualBody

    expOutFiles = tcData.TestCase.OutFiles()
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
            case "JsonToCsv":
                SaveHttpRespFile(actualBody, outputsFile)
            case "JsonToExcel":
                SaveHttpRespFile(actualBody, outputsFile)
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



  //     "outFiles": [
  //       {
  //         "filename": "storeManagementInfo.xlsx",
  //         "format": "xlsx",
  //         "data": {
  //           "content": ["$(body)"]
  //         }
  //       },
  //       {
  //         "filename": "storeManagementInfo.csv",
  //         "format": "DataToCsv",
  //         "data": {
  //           "title": ["titlevalue"],
  //           "count2": [20]
  //         }
  //       },
  //       {
  //         "targetFile": "storeManagementInfo.xlsx",
  //         "targetHeader": [],
  //         "sources": [],
  //         "operation": "DataToExcel",
  //         "data": {
  //           "content": ["$(body)"]
  //         }
  //       },
  //       {
  //         "targetFile": "storeManagementInfo.xlsx",
  //         "targetHeader": [],
  //         "sources": ["$(body)"],
  //         "operation": "JsonToCsv"
  //       },
  //       {
  //         "targetFile": "storeManagementInfo.xlsx",
  //         "targetHeader": [],
  //         "sources": ["$(sql).*"],
  //         "operation": "JsonToExcel"
  //       },
  //       {
  //         "targetFile": "storeManagementInfo.xlsx",
  //         "targetHeader": [],
  //         "sources": ["file1", "file2"],
  //         "operation": "Union"
  //       },
  //       {
  //         "targetFile": "storeManagementInfo.xlsx",
  //         "targetHeader": [],
  //         "sources": ["file1"],
  //         "operation": "CsvToExcel"
  //       },
  //       {
  //         "targetFile": "storeManagementInfo.xlsx",
  //         "targetHeader": [],
  //         "sources": ["file1"],
  //         "operation": "ExcelToCsv"
  //       }
  //     ]
  //   }
  // }
