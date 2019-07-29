/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2019.07
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package extension

import (
    // "fmt"
    "strings"

    "go4api/cmd"
    "go4api/lib/testcase"
    "go4api/lib/extension/statechart"
)

func GetScFilePaths () []string {
    filePathSlice := strings.Split(cmd.Opt.StateChart, ",")

    return filePathSlice
}

func InitFullScTcSlice (filePaths []string) ([]*testcase.TestCaseDataInfo) {
    // filePathSlice := GetTsFilePaths()

    fullKwTcSlice := statechart.InitFullScTcSlice(filePaths)

    return fullKwTcSlice
}