/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package keyword

import (
    "fmt"
    "os"
    "strings"
    "bufio"
    "io"

    "go4api/utils"
    "go4api/lib/testcase"
    // "go4api/lib/testsuite"
)

func InitFullKwTcSlice (filePathSlice []string) []*testcase.TestCaseDataInfo {
    var fullKwTcSlice []*testcase.TestCaseDataInfo

    // fullTsTcSlice := testsuite.InitFullTsTcSlice(filePaths)
    // fullTcSlice := testcase.InitFullTcSlice(filePaths)

    // kwSlice := InitKeyWordSlice(filePathSlice)

    // for i, _ := range kwSlice {
    //     gKw := AnalyzeKeyWordTestCases(kwSlice[i])

    //     gKw := LookupTestSuiteOrTestCases(gKw.TestCases, fullTsTcSlice, fullTcSlice)

    //     analyzedTestCases := (*gKw).TestCases

    //     // Note: to avoid the possibility of the case duplication, here is to put the keyword prefix to tsName or tcName
    //     // Please remember also need to update the ParentTestCase name, and except for root
    //     for i, _ := range analyzedTestCases {
    //         tsName := tsuite.TsName()

    //         tcName := analyzedTestCases[i].TestCase.TcName()
    //         parentTestCaseName := analyzedTestCases[i].TestCase.ParentTestCase()

    //         analyzedTestCases[i].TestCase.UpdateTcName(tsName + "-" + tcName)
    //         if parentTestCaseName != "root" {
    //             analyzedTestCases[i].TestCase.SetParentTestCase(tsName + "-" + parentTestCaseName)
    //         }

    //     }

    //     fullTsTcSlice = append(fullTsTcSlice, analyzedTestCases[0:]...)
    // }

    return fullKwTcSlice
}

func AnalyzeKeyWordTestCases () {

}

func LookupTestCases () {

}

func InitKeyWordSlice (filePathSlice []string) []*GKeyWord { 
    var kwSlice []*GKeyWord
    var kwFileList []string

    for i, _ := range filePathSlice {
        // to support pattern later
        // matches, _ := filepath.Glob(filePathSlice[i])

        kwFileListTemp, _ := utils.WalkPath(filePathSlice[i], ".keyword")
        kwFileList = append(kwFileList, kwFileListTemp[0:]...)
    }

    for _, kwFile := range kwFileList {
        gKw := ConstructKwInfosWithoutDt(kwFile)

        kwSlice = append(kwSlice, &gKw)
    }

    return kwSlice
}


func ConstructKwInfosWithoutDt (kwFile string) GKeyWord {
    var gKw GKeyWord
    var lines []string

    lines, _ = readLines(kwFile)
    gKw = GetGKeyWord(lines)
  
    return gKw
}

func GetGKeyWord (lines []string) GKeyWord {
    // Note: each block has the leading line with prefix '*** TestCases / Settings / Keywords / Variables /...''
    var blockHeaderLines []int
    gKeyWord := GKeyWord{}

    linesCount := len(lines)
    // get the block header line numbers, starting from line 0
    for i, line := range lines {
        if strings.HasPrefix(strings.TrimSpace(line), "***") {
            blockHeaderLines = append(blockHeaderLines, i)
        }
    }

    headerCount := len(blockHeaderLines)

    for i, _ := range blockHeaderLines {
        if i != headerCount - 1 {
            // passing starting line, ending line, line for each block
            FullfillBlock(&gKeyWord, blockHeaderLines[i], blockHeaderLines[i + 1] - 1, lines)
        } else {
            FullfillBlock(&gKeyWord, blockHeaderLines[i], linesCount - 1, lines)
        }
    }

    return gKeyWord
}

func FullfillBlock (gKeyWord *GKeyWord, startLine int, endLine int, lines []string) {
    blockType := GetBlockType(lines[startLine])

    switch blockType {
    case "Settings":
        settings := &Settings {
            StartLine: startLine,
            EndLine: endLine,
        }

        gKeyWord.Settings = settings
    case "TestCases":
        testCases := &TestCases {
            StartLine: startLine,
            EndLine: endLine,
        }

        gKeyWord.TestCases = testCases
    // case "Keywords":
        //
    case "Variables":
        variables := &Variables {
            StartLine: startLine,
            EndLine: endLine,
        }

        gKeyWord.Variables = variables
    default:
        fmt.Println("Warning, can not recognize the block type")
    }
}

func GetBlockType (headerLine string) string {
    var blockType string

    blockTypes := []string{"TestCases", "Settings", "Keywords", "Variables"}

    for i, _ := range blockTypes {
        if strings.Count(headerLine, blockTypes[i]) > 1 {
            blockType = blockTypes[i]
            break
        }
    }

    return blockType
}

func readLines (path string) (lines []string, err error){  
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()
 
    rd := bufio.NewReader(f)
    for {
            line, err := rd.ReadString('\n')

            line = strings.Replace(line, "\n", "", -1)
            lines = append(lines, line)

            // fmt.Println(line)
          
            if err != nil || io.EOF == err {
                break
            }  
        }

    return
}  
