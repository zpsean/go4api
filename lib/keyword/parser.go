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
    // "path/filepath"

    "go4api/utils"
    "go4api/lib/testcase"
    "go4api/lib/testsuite"
)

func InitFullKwTcSlice (kwfilePathSlice []string) ([]*testcase.TestCaseDataInfo, []string) {
    var fullKwTcSlice []*testcase.TestCaseDataInfo
    var fullKwJsPathSlice []string

    kwSlice := InitKeyWordSlice(kwfilePathSlice)

    // 
    tsNameFileMap := GetTsNameFileMap(kwSlice)
    fullBasicTcSlice := GetBasicTcsInfo(kwSlice)
    
    // set kwSlice - kwSlice[i].TestCases.ParsedTestCases[k].KWTestCaseName
    for i, _ := range kwSlice {
        for j, _ := range kwSlice[i].TestCases.ParsedTestCases {
            matched := SetKwTestSuiteInfo(kwSlice[i].TestCases.ParsedTestCases[j], tsNameFileMap)

            if matched == true {
                var files []string
                tsFile := kwSlice[i].TestCases.ParsedTestCases[j].MappingToTestSuiteFile
                files = append(files, tsFile)

                kwTcsFromTs := testsuite.InitFullTsTcSlice(files)

                fullKwTcSlice = append(fullKwTcSlice, kwTcsFromTs[0:]...)

            } else {
                kwTcsFromBasicTcs := LookupKwTestCase(kwSlice[i].TestCases.ParsedTestCases[j], fullBasicTcSlice)

                fullKwTcSlice = append(fullKwTcSlice, kwTcsFromBasicTcs[0:]...)
            }  
        }

        // for js
        for _, p := range kwSlice[0].Settings.JsFuncPaths {
            absPath := utils.GetAbsPath(p)
            fullKwJsPathSlice = append(fullKwJsPathSlice, absPath)
        }
    }

    // for i, _ := range kwSlice {
    //     gKw := AnalyzeKeyWordTestCases(kwSlice[i])

    //     tc1 := LookupTestSuites(gKw.TestCases, fullTsNameSlice)
    //     tc2 := LookupTestSuites(gKw.TestCases, fullTcSlice)

    //     analyzedTestCases := Consolidate(tc1, tc2)

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

    return fullKwTcSlice, fullKwJsPathSlice
}



func SetKwTestSuiteInfo(ktc *KWTestCase, tsNameFileMap map[string]string) bool {
    var matched = false
    for k, v := range tsNameFileMap {
        if ktc.KWTestCaseName == k {
            ktc.MappingToTestSuiteId = k
            ktc.MappingToTestSuiteFile = v

            matched = true
            break
        }
    }

    return matched
}

func LookupKwTestCase(ktc *KWTestCase, fullBasicTcSlice []*testcase.TestCaseDataInfo) []*testcase.TestCaseDataInfo {
    var fullKwTcSlice []*testcase.TestCaseDataInfo

    for i, _ := range fullBasicTcSlice {
        if ktc.KWTestCaseName == fullBasicTcSlice[i].TcName() {
            fullKwTcSlice = append(fullKwTcSlice, fullBasicTcSlice[i])
        }
    }

    return fullKwTcSlice
}


func GetTsNameFileMap (kwSlice []*GKeyWord) map[string]string {
    tsNameFileMap := make(map[string]string)

    for i, _ := range kwSlice {
        paths := utils.GetAbsPaths(kwSlice[i].Settings.TestSuitePaths)

        for _, p := range paths {
            var ps []string
            ps = append(ps, p)

            tsNames := testsuite.GetTsNames(ps)
            tsNameFileMap[tsNames[0]] = p
        }
    }

    return tsNameFileMap
}

func GetBasicTcsInfo (kwSlice []*GKeyWord) ([]*testcase.TestCaseDataInfo) {
    var TcPaths []string

    for i, _ := range kwSlice {
        paths := utils.GetAbsPaths(kwSlice[i].Settings.BasicTestCasePaths)
        TcPaths = append(TcPaths, paths[0:]...)
    }

    fullTcSlice := testcase.InitFullTcSlice(TcPaths)

    return fullTcSlice
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
    gKw = InitGKeyWord(lines)
  
    return gKw
}

func InitGKeyWord (lines []string) GKeyWord {
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
        gKeyWord.PopulateSettingsOriginalContent(startLine, endLine, lines)
        gKeyWord.ParseSettingsOriginalContent()
    case "TestCases":
        gKeyWord.PopulateTestCasesOriginalContent(startLine, endLine, lines)
        gKeyWord.ParseTestCasesOriginalContent()
    // case "Keywords":
        //
    case "Variables":
        gKeyWord.PopulateVariablesOriginalContent(startLine, endLine, lines)
    default:
        fmt.Println("Warning, can not recognize the block type")
    }
}

func GetBlockType (headerLine string) string {
    var blockType string

    blockTypes := []string{"TestCases", "Settings", "Keywords", "Variables"}

    for i, _ := range blockTypes {
        if strings.Count(headerLine, blockTypes[i]) > 0 {
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
