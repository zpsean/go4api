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



