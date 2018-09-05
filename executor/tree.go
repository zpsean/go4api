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
    // "strconv"
    "encoding/json"
    "go4api/cmd"
    "go4api/testcase"
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
    statusCountList [][]int
    tcNotExecutedList []testcase.TestCaseExecutionInfo
)

//
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
        TestMessages: "",
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
            TestMessages: "",
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
                TestMessages: "",
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


func RefreshNodeAndDirectChilrenTcResult(node *tcNode, tcRunResult string, tcStart string, tcEnd string, tcRunMessage string, 
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

func CollectNodeReadyStatus(node *tcNode, priority string) {
    for i, _ := range node.children {
        if node.children[i].TestCaseExecutionInfo.Priority() == priority {
            switch node.children[i].TestCaseExecutionInfo.TestResult { 
                case "Ready": 
                    statusReadyCount = statusReadyCount + 1
            }
        }
        CollectNodeReadyStatus(node.children[i], priority)
    }
}

func CollectNodeStatusByPriority(node *tcNode, p_index int, priority string) {
    for i, _ := range node.children {
        if node.children[i].TestCaseExecutionInfo.Priority() == priority {
            statusCountList[p_index][0] = statusCountList[p_index][0] + 1

            switch node.children[i].TestCaseExecutionInfo.TestResult { 
                case "Ready": 
                    statusCountList[p_index][1] = statusCountList[p_index][1] + 1
                case "Success": 
                    statusCountList[p_index][2] = statusCountList[p_index][2] + 1  
                case "Fail":
                    statusCountList[p_index][3] = statusCountList[p_index][3] + 1 
                default: 
                    statusCountList[p_index][4] = statusCountList[p_index][4] + 1
                    // write the cases to tcNotExecutedList
                    tcNotExecutedList = append(tcNotExecutedList, node.children[i].TestCaseExecutionInfo)
                    
            }
        }
        CollectNodeStatusByPriority(node.children[i], p_index, priority)
    }
}

func CollectOverallNodeStatus(node *tcNode, p_index int) {
    for i, _ := range node.children {
        statusCountList[p_index][0] = statusCountList[p_index][0] + 1

        switch node.children[i].TestCaseExecutionInfo.TestResult { 
            case "Ready": 
                statusCountList[p_index][1] = statusCountList[p_index][1] + 1
            case "Success": 
                statusCountList[p_index][2] = statusCountList[p_index][2] + 1  
            case "Fail":
                statusCountList[p_index][3] = statusCountList[p_index][3] + 1 
            default: 
                statusCountList[p_index][4] = statusCountList[p_index][4] + 1
        }
        CollectOverallNodeStatus(node.children[i], p_index)
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
