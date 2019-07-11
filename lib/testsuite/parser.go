/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019.07
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package testsuite

import (
    "fmt"
    // "time"
    "os"
    "encoding/json"

    "go4api/utils"
    "go4api/lib/testcase"
)

func InitFullTsTcSlice (filePathSlice []string) []*testcase.TestCaseDataInfo {
    var fullTsTcSlice []*testcase.TestCaseDataInfo

    tsSlice := InitTestSuiteSlice(filePathSlice)

    for i, _ := range tsSlice {
        tsuite := AnalyzeTestSuiteTestCases(tsSlice[i])

        analyzedTestCases := (*tsuite)[tsuite.TsName()].AnalyzedTestCases

        // Note: to avoid the possibility of the case duplication, here is to put the TestSuite prefix to tcName
        // Please remember also need to update the ParentTestCase name
        for i, _ := range analyzedTestCases {
            tsName := tsuite.TsName()

            tcName := analyzedTestCases[i].TestCase.TcName()
            parentTestCaseName := analyzedTestCases[i].TestCase.ParentTestCase()

            analyzedTestCases[i].TestCase.UpdateTcName(tsName + "-" + tcName)
            if parentTestCaseName != "root" {
                analyzedTestCases[i].TestCase.SetParentTestCase(tsName + "-" + parentTestCaseName)
            }

        }

        fullTsTcSlice = append(fullTsTcSlice, analyzedTestCases[0:]...)
    }

    return fullTsTcSlice
}

func InitTestSuiteSlice (filePathSlice []string) []*TestSuite { 
    var tsSlice []*TestSuite
    var suiteFileList []string

    for i, _ := range filePathSlice {
        // to support pattern later
        // matches, _ := filepath.Glob(filePathSlice[i])

        suiteFileListTemp, _ := utils.WalkPath(filePathSlice[i], ".testsuite")
        suiteFileList = append(suiteFileList, suiteFileListTemp[0:]...)
    }

    for _, suiteFile := range suiteFileList {
        tsuite := ConstructTsInfosWithoutDt(suiteFile)

        tsSlice = append(tsSlice, &tsuite)
    }

    return tsSlice
}

// to populate AnalyzedTestCases, 
// if TestCasePaths is defined, use path to generate
// otherwise, use OriginalTestCases
func AnalyzeTestSuiteTestCases (tsuite *TestSuite) *TestSuite {
    if len(tsuite.TestCasePaths()) > 0 {
        fullTcSlice := testcase.InitFullTcSlice(tsuite.TestCasePaths())

        tsuite.SetAnalyzedTestCases(fullTcSlice)
    }

    return tsuite
}

func ConstructTsInfosWithoutDt (jsonFile string) TestSuite {
    var tsuite TestSuite

    jsonStr := utils.GetJsonFromFile(jsonFile)

    err := json.Unmarshal([]byte(jsonStr), &tsuite)
    if err != nil {
        fmt.Println("!! Error, parse Json into testsuite failed: ", jsonFile, ": ", err)
        os.Exit(1)
    }
  
    return tsuite
}

