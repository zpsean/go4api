/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019.07
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
    "go4api/lib/keyword"
)

func GetKwFilePaths () []string {
    filePathSlice := strings.Split(cmd.Opt.KeyWord, ",")

    return filePathSlice
}

func InitFullKwTcSlice (filePaths []string) []*testcase.TestCaseDataInfo {
    // filePathSlice := GetTsFilePaths()

    fullKwTcSlice := keyword.InitFullKwTcSlice(filePaths)

    return fullKwTcSlice
}