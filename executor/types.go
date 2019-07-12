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
    "go4api/lib/testcase"
    "go4api/lib/tree"
    "go4api/js"
    "go4api/cmd"
    "go4api/utils"
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
    var fullTcSlice []*testcase.TestCaseDataInfo
    var jsFileList, fullKwJsPathSlice []string

    if cmd.Opt.IfTestSuite {
        filePaths := GetTsFilePaths()
        fullTcSlice = InitFullTsTcSlice(filePaths)

        jsFileList, _ = utils.WalkPath(cmd.Opt.JsFuncs, ".js")
    } else if cmd.Opt.IfKeyWord {
        filePaths := GetKwFilePaths()
        fullTcSlice, fullKwJsPathSlice = InitFullKwTcSlice(filePaths)

        jsFileList, _ = utils.WalkPath(fullKwJsPathSlice[0], ".js")
    } else {
        filePaths := GetTcFilePaths()
        fullTcSlice = InitFullTcSlice(filePaths)

        jsFileList, _ = utils.WalkPath(cmd.Opt.JsFuncs, ".js")
    }

    // js files
    gjs.InitJsFunctions(jsFileList)

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

