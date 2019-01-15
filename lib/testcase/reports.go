/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package testcase

import (
    // "fmt"
    "path/filepath"
    "encoding/json"

    "go4api/cmd"
    gsession "go4api/lib/session"
)

// for report
func (tcExecution *TestCaseExecutionInfo) TcConsoleResults() *TcConsoleResults {
    tcConsoleRes := &TcConsoleResults { 
        TcName: tcExecution.TcName(),
        Priority: tcExecution.Priority(),
        ParentTestCase: tcExecution.ParentTestCase(),
        JsonFilePath: filepath.Base(tcExecution.TestCaseDataInfo.JsonFilePath),
        CsvFile: filepath.Base(tcExecution.TestCaseDataInfo.CsvFile),
        CsvRow: tcExecution.TestCaseDataInfo.CsvRow,
        MutationInfoStr: tcExecution.TestCaseDataInfo.MutationInfoStr,
        SetUpResult: tcExecution.SetUpResult,
        HttpResult: tcExecution.HttpResult,
        TearDownResult: tcExecution.TearDownResult,
        TestResult: tcExecution.TestResult,
        ActualStatusCode: tcExecution.ActualStatusCode,
        HttpTestMessages: tcExecution.HttpTestMessages,
    }

    return tcConsoleRes
}


func (tcExecution *TestCaseExecutionInfo) TcReportResults() *TcReportResults {
    ifGlobalSetUpTearDown := ""
    if tcExecution.TestCaseDataInfo.TestCase.IfGlobalSetUpTestCase() == true {
        ifGlobalSetUpTearDown = "GlobalSetUp"
    } else if tcExecution.TestCaseDataInfo.TestCase.IfGlobalTearDownTestCase() == true {
        ifGlobalSetUpTearDown = "GlobalTearDown"
    } else {
        ifGlobalSetUpTearDown = "RegularCases"
    }

    var caseOrigin interface{}
    var actualBody interface{}
    var actualHeader interface{}
    var globalVariables interface{}
    var tcSession interface{}
    if cmd.Opt.IfShowOriginRequest == true {
        caseOrigin = tcExecution.TestCaseDataInfo.TestCase

        json.Unmarshal(tcExecution.ActualBody, &actualBody)

        actualHeader = tcExecution.ActualHeader
        globalVariables = gsession.LoopGlobalVariables()
        tcSession = gsession.LookupTcSession(tcExecution.TcName())
    }

    tcReportRes := &TcReportResults { 
        TcName: tcExecution.TcName(),
        IfGlobalSetUpTearDown: ifGlobalSetUpTearDown,
        // CaseType: 
        Priority: tcExecution.Priority(),
        ParentTestCase: tcExecution.ParentTestCase(),
        SetUpResult: tcExecution.SetUpResult,
        SetUpTestMessages: tcExecution.SetUpTestMessages,
        Path: tcExecution.ReqPath(),
        Method: tcExecution.ReqMethod(),
        JsonFilePath: tcExecution.TestCaseDataInfo.JsonFilePath,
        CsvFile: tcExecution.TestCaseDataInfo.CsvFile,
        CsvRow: tcExecution.TestCaseDataInfo.CsvRow,
        MutationArea: tcExecution.TestCaseDataInfo.MutationArea,
        MutationCategory:tcExecution.TestCaseDataInfo.MutationCategory,
        MutationRule: tcExecution.TestCaseDataInfo.MutationRule,
        MutationInfo: tcExecution.TestCaseDataInfo.MutationInfo,
        HttpResult: tcExecution.HttpResult,
        ActualStatusCode: tcExecution.ActualStatusCode,
        StartTime: tcExecution.StartTime,
        EndTime: tcExecution.EndTime,
        HttpTestMessages: tcExecution.HttpTestMessages,
        StartTimeUnixNano: tcExecution.StartTimeUnixNano,
        EndTimeUnixNano: tcExecution.EndTimeUnixNano,
        DurationUnixNano: tcExecution.DurationUnixNano,
        DurationUnixMillis: tcExecution.DurationUnixNano / 1000000,
        TearDownResult: tcExecution.TearDownResult,
        TearDownTestMessages: tcExecution.TearDownTestMessages,
        TestResult: tcExecution.TestResult,
        CaseOrigin: caseOrigin,
        GlobalVariables: globalVariables,
        Session: tcSession,
        ActualHeader: actualHeader,
        ActualBody: actualBody,
    }

    return tcReportRes
}


