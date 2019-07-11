/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package executor

import (
    // "fmt"
    "strings"

    "go4api/cmd"
    "go4api/lib/testcase"
    "go4api/lib/testsuite"
)

func GetTsFilePaths () []string {
    filePathSlice := strings.Split(cmd.Opt.Testsuite, ",")

    return filePathSlice
}

func InitFullTsTcSlice (filePaths []string) []*testcase.TestCaseDataInfo {
    // filePathSlice := GetTsFilePaths()

    fullTsTcSlice := testsuite.InitFullTsTcSlice(filePaths)

    return fullTsTcSlice
}