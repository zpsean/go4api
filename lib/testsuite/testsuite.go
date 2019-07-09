/*
 * go4api - a api testing tool written in Go
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

func (ts *TestSuite) TcRootPath() string {
    return (*ts)[ts.TsName()].TcRootPath
}

func (ts *TestSuite) TestCases() []string {
    return (*ts)[tc.TsName()].TestCases
}

func (ts *TestSuite) Parameters() map[string]interface{} {
    return (*ts)[tc.TsName()].Parameters
}