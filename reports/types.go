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
 	// "encoding/json"

	"go4api/lib/testcase"
)

type TcReportSlice []*testcase.TcReportResults

func (tcReportSlice TcReportSlice) ClassifyResults () (TcReportSlice, TcReportSlice, TcReportSlice) {
    var setUpResultSlice TcReportSlice
    var normalResultSlice TcReportSlice
    var tearDownResultSlice TcReportSlice

    for i, _ := range tcReportSlice {
        switch tcReportSlice[i].IfGlobalSetUpTearDown {
            case "SetUp":
                setUpResultSlice = append(setUpResultSlice, tcReportSlice[i])
            case "TearDown":
                tearDownResultSlice = append(tearDownResultSlice, tcReportSlice[i])
            default:
                normalResultSlice = append(normalResultSlice, tcReportSlice[i])
        }
    }

    return setUpResultSlice, normalResultSlice, tearDownResultSlice
}

