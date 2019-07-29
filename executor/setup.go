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
    "go4api/lib/testcase"
)

func InitGlobalSetUpTcSlice (fullTcSlice []*testcase.TestCaseDataInfo) []*testcase.TestCaseDataInfo {
    var tcSlice []*testcase.TestCaseDataInfo
    for i, _ := range fullTcSlice {
        if fullTcSlice[i].TestCase.IfGlobalSetUpTestCase() == true {
            tcSlice = append(tcSlice, fullTcSlice[i])
        }
    }
    
    return tcSlice
}
