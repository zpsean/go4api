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
 	// "strconv"
 	"encoding/json"

    "go4api/ui/js" 
    "go4api/texttmpl"
)

func (tcReportSlice TcReportSlice) GenerateGraphicJs (resultsDir string) {
    statsFile := resultsDir + "/js/graphic.js"

    reJsons := tcReportSlice.GetGraphicJson()
    // fmt.Println(reJsons)
    texttmpl.GenerateGraphicJs(js.Graphic, statsFile, reJsons)
}

func (tcReportSlice TcReportSlice) GetGraphicJson () []string {
    var reJsons []string

    orderedByStartTime := tcReportSlice.SortByStartTime()
    //
    circlePositions, priorityLines := orderedByStartTime.GetCirclePositions()
    circlesJson, _ := json.MarshalIndent(circlePositions, "", "\t")
    reJsons = append(reJsons, string(circlesJson))

    priorityLinesJson, _ := json.MarshalIndent(priorityLines, "", "\t")
    reJsons = append(reJsons, string(priorityLinesJson))
    //
    parentChildrenlinePositions := orderedByStartTime.GetParentChildrenLinePositions(circlePositions)
    parentChildrenLineJson, _ := json.MarshalIndent(parentChildrenlinePositions, "", "\t")
    reJsons = append(reJsons, string(parentChildrenLineJson))

    return reJsons
}


func (tcReportSlice TcReportSlice) GetCirclePositions () (map[string][]interface{}, map[string][]interface{}) {
    var circlePositions = make(map[string][]interface{})
    var priorityLines = make(map[string][]interface{})
    var phase = "GlobalSetUp"
    var priority = ""

    var offset = 0

    for i, _ := range tcReportSlice {
        cx := (i % 20 + 1) * 50
        cy := (i / 20 + 2) * 50 + offset

        cradius := 0
        du := tcReportSlice[i].DurationUnixMillis
        switch {
            case du >= 1500:
                cradius = 14
            case du < 1500 && du >= 1000:
                cradius = 11
            case du < 1000 && du >= 500:
                cradius = 9
            case du < 500 && du >= 200:
                cradius = 7
            case du < 200:
                cradius = 5
        }

        ccolor := ""
        if tcReportSlice[i].TestResult == "Success" {
          ccolor = "green"
        } else if tcReportSlice[i].TestResult == "Fail" {
          // ccolor = "red"
        } else {
          ccolor = "gray"
        }
        circlePositions[tcReportSlice[i].TcName] = []interface{}{cx, cy, cradius, ccolor}
        
        //
        if tcReportSlice[i].IfGlobalSetUpTearDown != phase || tcReportSlice[i].Priority != priority {
            plstartx := 0
            plstarty := (i / 20 + 2) * 50
            plendx := 1080
            plendy := (i / 20 + 2) * 50

            priorityLines[priority] = []interface{}{plstartx, plstarty, plendx, plendy}

            priority = tcReportSlice[i].Priority

            offset = offset + 1
        }  
    }

    return circlePositions, priorityLines
}

func (tcReportSlice TcReportSlice) GetParentChildrenLinePositions (circlePositions map[string][]interface{}) [][]int {
    var linePositions [][]int

    lstartx := 0
    lstarty := 0
    lendx := 0
    lendy := 0

    for parent, _ := range tcReportSlice {
        for child, _ := range tcReportSlice {
            if tcReportSlice[parent].TcName == tcReportSlice[child].ParentTestCase {
                lstartx = circlePositions[tcReportSlice[parent].TcName][0].(int)
                lstarty = circlePositions[tcReportSlice[parent].TcName][1].(int)

                lendx = circlePositions[tcReportSlice[child].TcName][0].(int)
                lendy = circlePositions[tcReportSlice[child].TcName][1].(int)

                linePositions = append(linePositions, []int{lstartx, lstarty, lendx, lendy})
            }
        }
    }

    return linePositions
}
    