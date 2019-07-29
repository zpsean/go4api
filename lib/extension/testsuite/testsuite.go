/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package testsuite

import (
    // "fmt"
    // "time"
    // "os"
    // "sort"

    "go4api/lib/testcase"
)

// test suite type - get
func (ts *TestSuite) TsName() string {
    var tsName string
    for key, _ := range *ts {
        tsName = key
        break
    }
    return tsName
}

func (ts *TestSuite) TestSuiteBasics() *TestSuiteBasics {
    return (*ts)[ts.TsName()]
}

func (ts *TestSuite) Priority() string {
    return (*ts)[ts.TsName()].Priority
}

func (ts *TestSuite) TestCasePaths() []string {
    return (*ts)[ts.TsName()].TestCasePaths
}

func (ts *TestSuite) OriginalTestCases() []string {
    return (*ts)[ts.TsName()].OriginalTestCases
}

func (ts *TestSuite) Parameters() map[string]interface{} {
    return (*ts)[ts.TsName()].Parameters
}

//
// set AnalyzedTestCases
func (ts *TestSuite) SetAnalyzedTestCases (tcSlice []*testcase.TestCaseDataInfo) {
    (*ts)[ts.TsName()].AnalyzedTestCases = tcSlice
}

