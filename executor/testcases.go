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
    // "fmt"
    "sort"
    "strings"

    "go4api/cmd"
    "go4api/utils"
    "go4api/lib/testcase"
)

func GetTcFilePaths () []string {
    filePathSlice := strings.Split(cmd.Opt.Testcase, ",")

    return filePathSlice
}

func InitFullTcSlice (filePaths []string) []*testcase.TestCaseDataInfo { 
    // filePathSlice := GetTcFilePaths()

    fullTcSlice := testcase.InitFullTcSlice(filePaths)

    return fullTcSlice
}

func InitNormalTcSlice (fullTcSlice []*testcase.TestCaseDataInfo) []*testcase.TestCaseDataInfo {
    var tcSlice []*testcase.TestCaseDataInfo
    for i, _ := range fullTcSlice {
        if fullTcSlice[i].TestCase.IfGlobalSetUpTestCase() != true && fullTcSlice[i].TestCase.IfGlobalTearDownTestCase() != true {
            tcSlice = append(tcSlice, fullTcSlice[i])
        }
    }
    
    return tcSlice
}

func GetPrioritySet (tcArray []*testcase.TestCaseDataInfo) []string {
    // get the priorities
    var priorities []string
    for _, tcaseInfo := range tcArray {
        priorities = append(priorities, tcaseInfo.Priority())
    }
    // go get the distinct key in priorities
    keys := make(map[string]bool)
    var prioritySet []string
    for _, entry := range priorities {
        // uses 'value, ok := map[key]' to determine if map's key exists, if ok, then true
        if _, ok := keys[entry]; !ok {
            keys[entry] = true
            prioritySet = append(prioritySet, entry)
        }
    }

    prioritySet = SortPrioritySet(prioritySet)

    return prioritySet
}

func SortPrioritySet (prioritySet []string) []string {
    // Note: here is a bug, if sort as string, as the sort results is 1, 10, 11, 2, 3, etc. => fixed
    prioritySet_Int := utils.ConvertStringArrayToIntArray(prioritySet)
    sort.Ints(prioritySet_Int)
    prioritySet = utils.ConvertIntArrayToStringArray(prioritySet_Int)

    return prioritySet
}

func GetFuzzTcArray () []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo

    // csvFileList := GetCsvDataFilesForJsonFile(jsonFile, "_fuzz_dt")

    return tcArray
}

func GetOriginMutationTcArray () []testcase.TestCaseDataInfo {
    var tcArray []testcase.TestCaseDataInfo
    
    return tcArray
}


