/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package tree

import (                                                                                                                                             
    // "os"
    // "time"
    // "fmt"
    // "sync"
    // "encoding/json"

    "go4api/lib/testcase"
)

type TcTreeStats struct{
    StatusCountByPriority map[string]map[string]int
    TcNotExecutedByPriority map[string]map[string][]*testcase.TestCaseExecutionInfo
}


func CreateTcTreeStats (prioritySet []string) TcTreeStats {
    StatusKeys := []string{"Ready", "Success", "Fail", "ParentFailed", "ParentSkipped"}
    statusCountByPriority := map[string]map[string]int{} 
    tcNotExecutedByPriority := map[string]map[string][]*testcase.TestCaseExecutionInfo{}

    for _, priority := range prioritySet {
        statusCountByPriority[priority] = map[string]int{}
        tcNotExecutedByPriority[priority] = map[string][]*testcase.TestCaseExecutionInfo{}
        
        for _, status := range StatusKeys {
            statusCountByPriority[priority][status] = 0
        }
    }

    statusCountByPriority["Overall"] = map[string]int{}
    tcNotExecutedByPriority["Overall"] = map[string][]*testcase.TestCaseExecutionInfo{}

    for _, status := range StatusKeys {
        statusCountByPriority["Overall"][status] = 0
    }

    tcTreeStats := TcTreeStats {
        StatusCountByPriority: statusCountByPriority,
        TcNotExecutedByPriority: tcNotExecutedByPriority,
    }

    return tcTreeStats
}

func (tcTreeStats TcTreeStats) ResetTcTreeStats (priority string) {
    StatusKeys := []string{"Ready", "Success", "Fail", "ParentFailed", "ParentSkipped"}

    for _, status := range StatusKeys {
        tcTreeStats.StatusCountByPriority[priority][status] = 0
        //
        tcTreeStats.TcNotExecutedByPriority[priority] = map[string][]*testcase.TestCaseExecutionInfo{}
    }
    tcTreeStats.StatusCountByPriority[priority]["Total"] = 0
}

func (tcTreeStats TcTreeStats) DeductReadyCount(priority string) {
    tcTreeStats.StatusCountByPriority[priority]["Ready"] -= 1
}

func (tcTreeStats TcTreeStats) CollectNodeStatusByPriority(node *TcNode, priority string) {
    for i, _ := range node.Children {
        if node.Children[i].TestCaseExecutionInfo.Priority() == priority {
            tcTreeStats.collectNodeStatusCommon(node, i, priority)
        }
        tcTreeStats.CollectNodeStatusByPriority(node.Children[i], priority)
    }
}

// key can be "Overall"
func (tcTreeStats TcTreeStats) CollectOverallNodeStatus(node *TcNode, key string) {
    for i, _ := range node.Children {
        tcTreeStats.collectNodeStatusCommon(node, i, key)
        tcTreeStats.CollectOverallNodeStatus(node.Children[i], key)
    }
}

func (tcTreeStats TcTreeStats) collectNodeStatusCommon(node *TcNode, i int, key string) {
    tcTreeStats.StatusCountByPriority[key]["Total"] += 1

    switch node.Children[i].TestCaseExecutionInfo.TestResult { 
        case "Ready": 
            tcTreeStats.StatusCountByPriority[key]["Ready"] += 1
        case "Success": 
            tcTreeStats.StatusCountByPriority[key]["Success"] += 1
        case "Fail":
            tcTreeStats.StatusCountByPriority[key]["Fail"] += 1
        case "ParentFailed":
            tcTreeStats.StatusCountByPriority[key]["ParentFailed"] += 1
            tcTreeStats.TcNotExecutedByPriority[key]["ParentFailed"] = append(tcTreeStats.TcNotExecutedByPriority[key]["ParentFailed"], node.Children[i].TestCaseExecutionInfo)
        case "ParentSkipped":
            tcTreeStats.StatusCountByPriority[key]["ParentSkipped"] += 1
            tcTreeStats.TcNotExecutedByPriority[key]["ParentSkipped"] = append(tcTreeStats.TcNotExecutedByPriority[key]["ParentSkipped"], node.Children[i].TestCaseExecutionInfo)
    }
}

