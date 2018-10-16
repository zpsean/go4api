/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package reports

import (
 	// "fmt"
 	"encoding/json"

 	"go4api/ui/js"
 	"go4api/texttmpl"
)

func GenerateResultsJs (tcReportSlice TcReportSlice, resultsDir string) {
	normalResultsJs := tcReportSlice.GetResultsJs()

    resultsFile := resultsDir + "/js/results.js"
    texttmpl.GenerateResultsJs(js.Results, resultsFile, normalResultsJs)
}

func (tcReportSlice TcReportSlice) GetResultsJs () *texttmpl.ResultsJs {
    resultsJs := tcReportSlice.InitResultsJs()

    setUpResultSlice, normalResultSlice, tearDownResultSlice := tcReportSlice.ClassifyResults()
    
    setUpResultSlice.UpdateResultsJsForSetUp(resultsJs)
    normalResultSlice.UpdateResultsJsForNormal(resultsJs)
    tearDownResultSlice.UpdateResultsJsForTearDown(resultsJs)

    tcReportSlice.UpdateResultsJsForGlobal(resultsJs)
    tcReportSlice.UpdateResultsJsForResultsArray(resultsJs)

    return resultsJs
}

func (tcReportSlice TcReportSlice) InitResultsJs () *texttmpl.ResultsJs {
    resultsJs := texttmpl.ResultsJs {
        SetUpStartUnixNano: 0,
        SetUpStart: `""`,
        SetUpEndUnixNano: 0,
        SetUpEnd: `""`,

        NormalStartUnixNano: 0,
        NormalStart: `""`,
        NormalEndUnixNano: 0,
        NormalEnd: `""`,

        TearDownStartUnixNano: 0,
        TearDownStart: `""`,
        TearDownEndUnixNano: 0,
        TearDownEnd: `""`,

        GStartUnixNano: 0,
        GStart: `""`,
        GEndUnixNano: 0,
        GEnd: `""`,

        TcResults: "[]",
    }

    return &resultsJs
}

func (tcReportSlice TcReportSlice) UpdateResultsJsForSetUp (resultsJs *texttmpl.ResultsJs) {
    if len(tcReportSlice) > 0 {
        orderedByStartTime := tcReportSlice.SortByStartTime()
        orderedByEndTime := tcReportSlice.SortByEndTime()

        totalTc := len(tcReportSlice)

        resultsJs.SetUpStartUnixNano = orderedByStartTime[0].StartTimeUnixNano
        resultsJs.SetUpStart = `"` + orderedByStartTime[0].StartTime + `"`
        resultsJs.SetUpEndUnixNano = orderedByEndTime[totalTc - 1].EndTimeUnixNano
        resultsJs.SetUpEnd = `"` + orderedByEndTime[totalTc - 1].EndTime + `"`
    }
}

func (tcReportSlice TcReportSlice) UpdateResultsJsForNormal (resultsJs *texttmpl.ResultsJs) {
    if len(tcReportSlice) > 0 {
        orderedByStartTime := tcReportSlice.SortByStartTime()
        orderedByEndTime := tcReportSlice.SortByEndTime()

        totalTc := len(tcReportSlice)

        resultsJs.NormalStartUnixNano = orderedByStartTime[0].StartTimeUnixNano
        resultsJs.NormalStart = `"` + orderedByStartTime[0].StartTime + `"`
        resultsJs.NormalEndUnixNano = orderedByEndTime[totalTc - 1].EndTimeUnixNano
        resultsJs.NormalEnd = `"` + orderedByEndTime[totalTc - 1].EndTime + `"`
    }
}

func (tcReportSlice TcReportSlice) UpdateResultsJsForTearDown (resultsJs *texttmpl.ResultsJs) {
    if len(tcReportSlice) > 0 {
        orderedByStartTime := tcReportSlice.SortByStartTime()
        orderedByEndTime := tcReportSlice.SortByEndTime()

        totalTc := len(tcReportSlice)

        resultsJs.TearDownStartUnixNano = orderedByStartTime[0].StartTimeUnixNano
        resultsJs.TearDownStart = `"` + orderedByStartTime[0].StartTime + `"`
        resultsJs.TearDownEndUnixNano = orderedByEndTime[totalTc - 1].EndTimeUnixNano
        resultsJs.TearDownEnd = `"` + orderedByEndTime[totalTc - 1].EndTime + `"`
    }
}

func (tcReportSlice TcReportSlice) UpdateResultsJsForGlobal (resultsJs *texttmpl.ResultsJs) {
    if len(tcReportSlice) > 0 {
        orderedByStartTime := tcReportSlice.SortByStartTime()
        orderedByEndTime := tcReportSlice.SortByEndTime()

        totalTc := len(tcReportSlice)

        resultsJs.GStartUnixNano = orderedByStartTime[0].StartTimeUnixNano
        resultsJs.GStart = `"` + orderedByStartTime[0].StartTime + `"`
        resultsJs.GEndUnixNano = orderedByEndTime[totalTc - 1].EndTimeUnixNano
        resultsJs.GEnd = `"` + orderedByEndTime[totalTc - 1].EndTime + `"`
    }
}

func (tcReportSlice TcReportSlice) UpdateResultsJsForResultsArray (resultsJs *texttmpl.ResultsJs) {
    tcReportBytes, _ := json.Marshal(tcReportSlice)
    tcReportStr := string(tcReportBytes)

    resultsJs.TcResults = tcReportStr
}

