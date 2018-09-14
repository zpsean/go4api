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
    "os"
    "time"
    "fmt"
    "sync"
    "encoding/json"

    "go4api/cmd"
    "go4api/lib/testcase"
    "go4api/api"
)

type tcNode struct{
    TestCaseExecutionInfo testcase.TestCaseExecutionInfo
    children []*tcNode // for child
}


var (
    tcTree = map[string]*tcNode{}
    root *tcNode

    statusReadyCount int

    statusCountByPriority = map[string]map[string]int{} 
    tcExecutedByPriority = map[string]map[string][]*testcase.TestCaseExecutionInfo{}
    tcNotExecutedByPriority = map[string]map[string][]*testcase.TestCaseExecutionInfo{}
)

func InitVariables(prioritySet []string) {
    statusReadyCount = 0

    for _, priority := range prioritySet {
        // Ready, Success, Fail, ParentFailed
        statusCountByPriority[priority] = map[string]int{}
        tcExecutedByPriority[priority] = map[string][]*testcase.TestCaseExecutionInfo{}
        tcNotExecutedByPriority[priority] = map[string][]*testcase.TestCaseExecutionInfo{}
    }

    statusCountByPriority["Overall"] = map[string]int{}
    tcExecutedByPriority["Overall"] = map[string][]*testcase.TestCaseExecutionInfo{}
    tcNotExecutedByPriority["Overall"] = map[string][]*testcase.TestCaseExecutionInfo{}
}

func GetDummyRootTc() testcase.TestCase {
    var rootTc testcase.TestCase
    str := `{
              "rootTcNmae": {
                "priority": "0",
                "parentTestCase": "0"
              }
            }`
    json.Unmarshal([]byte(str), &rootTc)

    return rootTc
}

func AddNode(TcaseExecution testcase.TestCaseExecutionInfo) bool {
    node := &tcNode{
        TestCaseExecutionInfo: TcaseExecution, 
        children: []*tcNode{},
    }
    var ifAdded bool
    
    parentTcName := node.TestCaseExecutionInfo.ParentTestCase()
    // the dummy root tese case's parentTestCase == "0"
    if parentTcName == "0" {
        root = node
        tcTree["root"] = root

        ifAdded = true
    } else if parentTcName == "root" {
        parent := tcTree[parentTcName]
        parent.children = append(parent.children, node)

        ifAdded = true
    } else {
        // below is to find out the right parent node, then add the child node for it if found
        c := make(chan *tcNode)
        go func(c chan *tcNode) {
            defer close(c)
            SearchNode(c, root, parentTcName)
        }(c)

        var resNodes []*tcNode
        for n := range c {
            resNodes = append(resNodes, n)
        }

        if len(resNodes) > 1 {
            fmt.Println("\n!! Error, more than one parent node found, please verify the test data")
            os.Exit(1)
        } else if len(resNodes) == 1 {
            parent := resNodes[0]
            parent.children = append(parent.children, node)   

            ifAdded = true
        } else {
            ifAdded = false 
        }
    }

    return ifAdded
}

func BuildTree(tcArray []testcase.TestCaseDataInfo) (*tcNode, map[string]*tcNode) {
    fmt.Println("\n---- Build the tcTree - Start ----")
    var tcArrayTree []testcase.TestCaseDataInfo
    var tcArrayNotTree []testcase.TestCaseDataInfo
    // Step 1: add the root node
    rootTc := GetDummyRootTc()
    rootTcaseData := testcase.TestCaseDataInfo {
        TestCase: &rootTc,
        JsonFilePath: "",
        CsvFile: "",
        CsvRow: "",
    }

    rootTcaseExecution := testcase.TestCaseExecutionInfo {
        TestCaseDataInfo: &rootTcaseData,
        TestResult: "",
        ActualStatusCode: 0,
        StartTime: "",
        EndTime: "",
        TestMessages: []*testcase.TestMessage{},
        StartTimeUnixNano: 0,
        EndTimeUnixNano: 0,
        DurationUnixNano: 0,
    }

    AddNode(rootTcaseExecution)

    // Step 2: add the node, init the tcArrayTree, tcArrayNotTree
    for i, _ := range tcArray {
        tcaseExecution := testcase.TestCaseExecutionInfo {
            TestCaseDataInfo: &tcArray[i],
            TestResult: "",
            ActualStatusCode: 0,
            StartTime: "",
            EndTime: "",
            TestMessages: []*testcase.TestMessage{},
            StartTimeUnixNano: 0,
            EndTimeUnixNano: 0,
            DurationUnixNano: 0,
        }

        ifAdded := AddNode(tcaseExecution)
        if ifAdded && true {
            tcArrayTree = append(tcArrayTree, tcArray[i])
        } else {
            tcArrayNotTree = append(tcArrayNotTree, tcArray[i])
        }
    }
    fmt.Println("tcArrayNotTree: ", len(tcArrayNotTree))

    // Step 3: loop the tcArrayNotTree, until none can be added to tree anymore
    // if the parentTestCase does not exist in full set of tc (i.e. tcArray), just skip the test case
    for {
        var tcArrayNotTreeTemp []testcase.TestCaseDataInfo
        for i, _ := range tcArrayNotTree {
            tcArrayNotTreeTemp = append(tcArrayNotTreeTemp, tcArrayNotTree[i])
        }
        //
        len1 := len(tcArrayNotTreeTemp)
        for i, _ := range tcArrayNotTreeTemp {
            tcaseExecution := testcase.TestCaseExecutionInfo {
                TestCaseDataInfo: &tcArrayNotTreeTemp[i],
                TestResult: "",
                ActualStatusCode: 0,
                StartTime: "",
                EndTime: "",
                TestMessages: []*testcase.TestMessage{},
                StartTimeUnixNano: 0,
                EndTimeUnixNano: 0,
                DurationUnixNano: 0,
            }

            ifAdded := AddNode(tcaseExecution)
            
            if ifAdded && true {
                tcArrayTree = append(tcArrayTree, tcArrayNotTreeTemp[i])
                tcArrayNotTree = RemoveArrayItem(tcArrayNotTree, tcArrayNotTreeTemp[i])
            }
        }
        len2 := len(tcArrayNotTree)
        // if can not add / remove item anymore
        if len1 == len2 {
            break
        }
    }
    // fmt.Println("----------------")
    // fmt.Println("---- build tree, step 3 END ----")

    if len(tcArrayNotTree) > 0 {
        fmt.Println("!!!Attention, there are test cases which parentTestCase does not exists\n", tcArrayNotTree)
    }

    fmt.Println("---- Build the tcTree - END ----")
    return root, tcTree
}


func SearchNode(c chan *tcNode, node *tcNode, testCaseName string) {
    for i, _ := range node.children {
        if node.children[i].TestCaseExecutionInfo.TcName() == testCaseName {
            c <- node.children[i]
        } else {
            SearchNode(c, node.children[i], testCaseName)
        }
    }
}


func InitNodesRunResult(node *tcNode, runResult string) {
    for i, _ := range node.children {
        // check if the current tese case is parent to others, if yes, then refresh
        if node.children[i].TestCaseExecutionInfo.ParentTestCase() == "root"{
            node.children[i].TestCaseExecutionInfo.TestResult = "Ready"
        } else {
            node.children[i].TestCaseExecutionInfo.TestResult = "ParentReady"
        }
        InitNodesRunResult(node.children[i], "")
    }
}


func ScheduleNodes(node *tcNode, wg *sync.WaitGroup, priority string, resultsChan chan testcase.TestCaseExecutionInfo, 
        pStart string, baseUrl string, resultsDir string) {
    //
    tick := 0
    max := cmd.Opt.ConcurrencyLimit
    //
    // note: does data copy happen if use n but not node.children[i]?
    for i, _ := range node.children {
        if priority == node.children[i].TestCaseExecutionInfo.Priority() && node.children[i].TestCaseExecutionInfo.TestResult == "Ready"{
            wg.Add(1)
            // Note: to prevent to tcp connection, here set a max, then sleep for a while
            if tick % max == 0 {
                time.Sleep(500 * time.Millisecond)
                go api.HttpApi(wg, resultsChan, pStart, baseUrl, *(node.children[i].TestCaseExecutionInfo.TestCaseDataInfo), resultsDir)
            } else {
                go api.HttpApi(wg, resultsChan, pStart, baseUrl, *(node.children[i].TestCaseExecutionInfo.TestCaseDataInfo), resultsDir)
            }

            tick = tick + 1
        }
        
        ScheduleNodes(node.children[i], wg, priority, resultsChan, pStart, baseUrl, resultsDir)
    }
}


func RefreshNodeAndDirectChilrenTcResult(node *tcNode, tcRunResult string, tcStart string, tcEnd string, tcRunMessage []*testcase.TestMessage, 
        tcStartUnixNano int64, tcEndUnixNano int64) {
    // fmt.Println("ccc the node to be refreshed: ", node, &node)
    node.TestCaseExecutionInfo.TestResult = tcRunResult
    node.TestCaseExecutionInfo.StartTime = tcStart
    node.TestCaseExecutionInfo.EndTime = tcEnd
    node.TestCaseExecutionInfo.TestMessages = tcRunMessage
    node.TestCaseExecutionInfo.StartTimeUnixNano = tcStartUnixNano
    node.TestCaseExecutionInfo.EndTimeUnixNano = tcEndUnixNano
    node.TestCaseExecutionInfo.DurationUnixNano = tcEndUnixNano - tcStartUnixNano

    for i, _ := range node.children {
        if tcRunResult == "Fail"{
            node.children[i].TestCaseExecutionInfo.TestResult = "ParentFailed"
        } else if tcRunResult == "Success"{
            node.children[i].TestCaseExecutionInfo.TestResult = "Ready"
        }
    }
}

func CollectNodeReadyStatusByPriority(node *tcNode, priority string) {
    for i, _ := range node.children {
        if node.children[i].TestCaseExecutionInfo.Priority() == priority {
            switch node.children[i].TestCaseExecutionInfo.TestResult { 
                case "Ready": 
                    statusReadyCount = statusReadyCount + 1
            }
        }
        CollectNodeReadyStatusByPriority(node.children[i], priority)
    }
}

func CollectNodeStatusByPriority(node *tcNode, priority string) {
    for i, _ := range node.children {
        if node.children[i].TestCaseExecutionInfo.Priority() == priority {
            statusCountByPriority[priority]["Total"] += 1
            tcExecutedByPriority[priority]["Total"] = append(tcExecutedByPriority[priority]["Total"], &(node.children[i].TestCaseExecutionInfo))

            switch node.children[i].TestCaseExecutionInfo.TestResult { 
                case "Ready": 
                    statusCountByPriority[priority]["Ready"] += 1
                case "Success": 
                    statusCountByPriority[priority]["Success"] += 1
                case "Fail":
                    statusCountByPriority[priority]["Fail"] += 1
                default: 
                    statusCountByPriority[priority]["ParentFailed"] += 1
                    tcNotExecutedByPriority[priority]["ParentFailed"] = append(tcNotExecutedByPriority[priority]["ParentFailed"], &(node.children[i].TestCaseExecutionInfo))
            }
        }
        CollectNodeStatusByPriority(node.children[i], priority)
    }
}

func CollectOverallNodeStatus(node *tcNode, key string) {
    for i, _ := range node.children {
        statusCountByPriority[key]["Total"] += 1
        tcExecutedByPriority[key]["Total"] = append(tcExecutedByPriority[key]["Total"], &(node.children[i].TestCaseExecutionInfo))

        switch node.children[i].TestCaseExecutionInfo.TestResult { 
            case "Ready": 
                statusCountByPriority[key]["Ready"] += 1
            case "Success": 
                statusCountByPriority[key]["Success"] += 1
            case "Fail":
                statusCountByPriority[key]["Fail"] += 1
            default: 
                statusCountByPriority[key]["ParentFailed"] += 1
                tcNotExecutedByPriority[key]["ParentFailed"] = append(tcNotExecutedByPriority[key]["ParentFailed"], &(node.children[i].TestCaseExecutionInfo))
        }
        CollectOverallNodeStatus(node.children[i], key)
    }
}


func ShowNodes(node *tcNode) {
    fmt.Println("\nShow node: ", node.TestCaseExecutionInfo.Priority(), node.TestCaseExecutionInfo.TcName(), node.TestCaseExecutionInfo.TestResult)
    for i, _ := range node.children {
        // fmt.Println("NN - C node:", &n, n)
        ShowNodes(node.children[i])
    }
}


func RemoveArrayItem(sourceArray []testcase.TestCaseDataInfo, tcData testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    // fmt.Println("RemoveArryaItem", sourceArray, tc)
    var resultArray []testcase.TestCaseDataInfo
    // resultArray := append(sourceArray[:index], sourceArray[index + 1:]...)
    for index, tc_i := range sourceArray {
        if tc_i.TcName() == tcData.TcName() {
            resultArray = append(sourceArray[:index], sourceArray[index + 1:]...)
            break
        }
    }

    return resultArray
}
