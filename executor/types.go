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
    "fmt"

    "go4api/cmd"
    "go4api/utils"
    "go4api/lib/testcase"
    "go4api/lib/tree"

    goja "github.com/dop251/goja"
)

type G4Store struct {
    OverallFail             int
    FullTcSlice             []*testcase.TestCaseDataInfo
    GlobalSetUpRunStore     *TcsRunStore
    NormalRunStore          *TcsRunStore
    GlobalTeardownRunStore  *TcsRunStore
}

type TcsRunStore struct {
    TcSlice     []*testcase.TestCaseDataInfo

    PrioritySet []string
    Root        *tree.TcNode
    TcTree      tree.TcTree
    TcTreeStats tree.TcTreeStats
    OverallFail int
}

func InitG4Store () *G4Store {
    fullTcSlice := InitFullTcSlice()
    // js files
    GetJsFiles()

    globalSetUpTcSlice := InitGlobalSetUpTcSlice(fullTcSlice)
    globalSetUpRunStore := &TcsRunStore {
        TcSlice: globalSetUpTcSlice,
        PrioritySet: []string{},
        Root: &tree.TcNode{},
        TcTree: tree.TcTree{},
        TcTreeStats: tree.TcTreeStats{},
        OverallFail: 0,
    }

    normalTcSlice := InitNormalTcSlice(fullTcSlice)
    normalRunStore := &TcsRunStore {
        TcSlice: normalTcSlice,
        PrioritySet: []string{},
        Root: &tree.TcNode{},
        TcTree: tree.TcTree{},
        TcTreeStats: tree.TcTreeStats{},
        OverallFail: 0,
    }

    globalTeardownTcSlice := InitGlobalTeardownTcSlice(fullTcSlice)
    globalTeardownRunStore := &TcsRunStore {
        TcSlice: globalTeardownTcSlice,
        PrioritySet: []string{},
        Root: &tree.TcNode{},
        TcTree: tree.TcTree{},
        TcTreeStats: tree.TcTreeStats{},
        OverallFail: 0,
    }

    g4Store := &G4Store {
        OverallFail: 0,
        FullTcSlice: fullTcSlice,
        GlobalSetUpRunStore: globalSetUpRunStore,
        NormalRunStore: normalRunStore,
        GlobalTeardownRunStore: globalTeardownRunStore,
    }

    return g4Store
}

// trial code for js
func GetJsFiles () {
    jsFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".js")
    fmt.Println(jsFileList)

    for i, _ := range jsFileList {
        srcBytes := utils.GetContentFromFile(jsFileList[i])

        src := string(srcBytes)

        fmt.Println(src)

        p, err := goja.Compile("", src, false)

        if err != nil {
            panic(err)
        }

        fmt.Println(p)

        for i = 1; i < 10; i++ {
            go GRunProgram(p)
            // vm := goja.New()
            // v, err := vm.RunProgram(p)

            // if err != nil {
            //     panic(err)
            // }

            // fmt.Println("the sum results is:", v.Export().(int64))
        }
    }
    

    // vm := goja.New()
    // v, err := vm.RunString("2 + 2")
    // if err != nil {
    //     panic(err)
    // }
    // if num := v.Export().(int64); num != 4 {
    //     panic(num)
    // }

}

func GRunProgram(p *goja.Program) {
    vm := goja.New()
    v, err := vm.RunProgram(p)

    if err != nil {
        panic(err)
    }

    fmt.Println("the sum results is:", v.Export().(int64))
}

