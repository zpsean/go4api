/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package executor

import (                                                                                                                                             
    "os"
    "fmt"
    "strings"
    "regexp"
    "path/filepath"
    
    "go4api/utils"
    "go4api/lib/csv"
)


func GetBasicInputsFilesPerFile(filePath string) []string {
    fileInputsInfos := GetBasicInputsInfos(filePath)
    // suppose one child tc per file currently
    inputsFiles := GenerateInputsFileWithConsolidatedData(filePath, fileInputsInfos[0])

    // fmt.Println("--> :", filePath, inputsFiles, fileInputsInfos[0])
    return inputsFiles
}

func GetBasicInputsInfos(filePath string) [][]string {
    // as the raw Jsonfile itself is template, may not be valid json fomat, before rendered by data
    contentsBytes := utils.GetContentFromFile(filePath)
    contents := string(contentsBytes)
    // add some space a make regexp works well
    contents = strings.Replace(contents, `[`, "[ ", -1)
    contents = strings.Replace(contents, `{`, "{ ", -1)

    var inputsInfos []string
    var fileInputsInfos [][]string
    var inputsPosSlice []int

    reg := regexp.MustCompile(`[\p{L}\w\pP]+`)
    wordSlice := reg.FindAllString(contents, -1)

    for i, value := range wordSlice {
        if strings.Contains(value, `"inputs"`) {
            inputsPosSlice = append(inputsPosSlice, i)
        }
    }
    for _, inputsPos := range inputsPosSlice {
        for ii := inputsPos + 1; ; ii ++ {
            v := strings.Replace(wordSlice[ii], `"`, "", -1)
            v = strings.Replace(v, `[`, "", -1)
            v = strings.Replace(v, `]`, "", -1)
            v = strings.Replace(v, `,`, "", -1)
            v = strings.Replace(v, ` `, "", -1)
            if len(v) > 0 {
                inputsInfos = append(inputsInfos, v)
                // fmt.Println(v)
            }
            if strings.Contains(wordSlice[ii], `]`) {
                break
            }
        }
        fileInputsInfos = append(fileInputsInfos, inputsInfos)
    }

    return fileInputsInfos
}

// to apply the csv operator: union, join, append for inputs files    
func GenerateInputsFileWithConsolidatedData (filePath string, inputsInfos []string) []string {
    // inputsInfos => ["s1ParentTestCase_out.csv", "join", "s1ParentTestCase_out.csv"]
    var inputsFiles []string
    tempDir := filepath.Dir(filePath) + "/temp"
    // 1. len(inputsInfos)
    switch {
        case len(inputsInfos) == 0:
            inputsFiles = []string{}
        case len(inputsInfos) == 1:
            inputsFiles = []string{tempDir + "/" + inputsInfos[0]}
        case len(inputsInfos) > 1:
            if len(inputsInfos) % 2 != 1 {
                fmt.Println("!! Error, inputs contents error, please check")
                os.Exit(1)
            }
            for i := 1; i <= len(inputsInfos) / 2; i ++ {
                operator := strings.ToLower(inputsInfos[2 * (i - 1) + 1])
                if operator != "union" && operator != "join" && operator != "append" {
                    fmt.Println("!! Error, inputs operator error, please check")
                    os.Exit(1)
                }
            }
            //
            var leftCsvPtr *gcsv.Gcsv
            // init the leftCsvPtr with first file
            leftFile := filepath.Join(tempDir, inputsInfos[0])
            lContentBytes := utils.GetContentFromFile(leftFile)
            leftCsv := gcsv.GetCsv(string(lContentBytes))
            leftCsvPtr = &leftCsv

            for i := 1; i <= len(inputsInfos) / 2; i ++ {
                operator := strings.ToLower(inputsInfos[2 * (i - 1) + 1])
                //
                rightFile := filepath.Join(tempDir, inputsInfos[2 * (i - 1) + 2])
                rContentBytes := utils.GetContentFromFile(rightFile)
                //
                switch operator {
                    case "union":
                        leftCsvPtr.Union(string(rContentBytes))
                    case "join":
                        leftCsvPtr.Join(string(rContentBytes))
                    case "append":
                        leftCsvPtr.Append(string(rContentBytes))
                }
            }
            inputsFiles = []string{tempDir + "/" + filepath.Base(filePath) + "_" + "inputs_cosolidated.csv"}
            
            writeGcsvToCsv(leftCsvPtr, inputsFiles[0])
    }
    return inputsFiles
}

func writeGcsvToCsv (gcsvPtr *gcsv.Gcsv, outFile string) {
    // header
    utils.GenerateCsvFileBasedOnVarOverride(gcsvPtr.Header, outFile)
    // data
    for i, _ := range gcsvPtr.DataRows {
        utils.GenerateCsvFileBasedOnVarAppend(gcsvPtr.DataRows[i], outFile)
    }
}

