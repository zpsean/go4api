/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
    "go4api/lib/testcase" 
)

type TcDataStore struct {
    TcData *testcase.TestCaseDataInfo

    TcLocalVariables map[string]interface{}

    HttpExpStatus    map[string]interface{}
    HttpExpHeader    map[string]interface{}
    HttpExpBody      map[string]interface{}
    HttpActualStatusCode int
    HttpActualHeader map[string][]string
    HttpActualBody   []byte

    HttpUrl    string

    CmdSection string // setUp, tearDown
    CmdGroupLength int
    
    CmdType       string // sql, redis, init, etc.
    CmdExecStatus string
    CmdAffectedCount int
    CmdResults    interface{}
}

func InitTcDataStore (tcData *testcase.TestCaseDataInfo) *TcDataStore {
    tcDataStore := &TcDataStore {
        TcData:               tcData,

        TcLocalVariables:     map[string]interface{}{},

        HttpExpStatus:        map[string]interface{}{},
        HttpExpHeader:        map[string]interface{}{},
        HttpExpBody:          map[string]interface{}{},
        HttpActualStatusCode: -1,
        HttpActualHeader:     map[string][]string{},
        HttpActualBody:       []byte{},

        HttpUrl:          "",

        CmdSection:       "",
        CmdGroupLength:   0,

        CmdType:          "",
        CmdExecStatus:    "",
        CmdAffectedCount: -1,
        CmdResults:       -1,
    }
    
    return tcDataStore
}




