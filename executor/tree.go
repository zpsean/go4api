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
    // "os"
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

    findNode **tcNode

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

    // fmt.Println("\nc-node: ", node)
    var ifAdded bool
    
    parentTCname := node.TestCaseExecutionInfo.ParentTestCase()
    // the dummy root tese case's parentTestCase == "0"
    if parentTCname == "0" {
        root = node
        tcTree["root"] = root

        ifAdded = true
    } else if parentTCname == "root" {
        parent := tcTree[parentTCname]
        parent.children = append(parent.children, node)   

        ifAdded = true
    } else {
        findNode = nil

        // below is to find out the right parent node, then add the child node for it
        SearchNode(&root, parentTCname)
        if findNode != nil {
            parent := *findNode
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
    // here seperate the tcArray into two parts, (appended to tree) and (not yet appended to tree)
    // Step 1: add the root node
    rootTc := GetDummyRootTc()
    rootTcaseData := testcase.TestCaseDataInfo {
        TestCase: rootTc,
        JsonFilePath: "",
        CsvFile: "",
        CsvRow: "",
    }

    rootTcaseExecution := testcase.TestCaseExecutionInfo {
        TestCaseDataInfo: rootTcaseData,
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
    for _, tcData := range tcArray {
        // if parentTestCase name can be found in tree
        // tcRunResult is "" for init
        tcaseExecution := testcase.TestCaseExecutionInfo {
            TestCaseDataInfo: tcData,
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
        // fmt.Println("!!now try to add tc: ", tc[0].(string), ifAdded)
        if ifAdded && true {
            tcArrayTree = append(tcArrayTree, tcData)
        } else {
            tcArrayNotTree = append(tcArrayNotTree, tcData)
        }
    }

    // fmt.Println("\n---- build tree, step 3 ----")

    // Step 3: loop the tcArrayNotTree, till all added, can set a strategy that:
    // (1): if the parentTestCase does not exist in full set of tc (i.e. tcArray), then set its parent as root
    // (2): or just skip the test case
    // fmt.Println("len tcArrayTree: ", len(tcArrayTree), tcArrayTree)
    // fmt.Println("len tcArrayNotTree: ", len(tcArrayNotTree), tcArrayNotTree)
    // fmt.Println("----------------")
    // loop the tcArrayNotTree, until none can be added to tree anymore
    for {
        len1 := len(tcArrayNotTree)
        // here may be bug, as once the item removed, the index would be out of range 
        // => fixed, using the name but not the index to remove
        for _, tcData := range tcArrayNotTree {
            // fmt.Println("!!now try to add tc - 0: ", tc[0].(string), tcArrayNotTree, tc)
            // if parentTestCase name can be found in tree
            tcaseExecution := testcase.TestCaseExecutionInfo {
                TestCaseDataInfo: tcData,
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
            // fmt.Println("!!now try to add tc: ", tc[0].(string), ifAdded, tcArrayNotTree, tc)
            if ifAdded && true {
                tcArrayTree = append(tcArrayTree, tcData)
                tcArrayNotTree = RemoveArrayItem(tcArrayNotTree, tcData)
            }
        }

        // fmt.Println("\nlen tcArrayNotTree: ", len(tcArrayNotTree), tcArrayNotTree)
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


// Note: here may be bug, as even the right node has been found, but the nil returned
// seems it always return the last leaf of the tree
// the bug impact the buil tree, and the refresh nodes for tcRunResult
// to fix: guess the return clause, it needs to return the pointer (address) of the node of tree, but not the variable in function
// now used the address as the temporary solution
func SearchNode(node **tcNode, testCaseNmae string) **tcNode {
    for i, _ := range (*node).children {
        if (*node).children[i].TestCaseExecutionInfo.TcName() == testCaseNmae {
            findNode = &((*node).children[i])
            return findNode
        } else {
            SearchNode(&((*node).children[i]), testCaseNmae)
        }
    }
    return findNode
}


func InitNodesRunResult(node *tcNode, runResult string) {
    for _, n := range node.children {
        // check if the current tese case is parent to others, if yes, then refresh
        if n.TestCaseExecutionInfo.ParentTestCase() == "root"{
            n.TestCaseExecutionInfo.TestResult = "Ready"
        } else {
            n.TestCaseExecutionInfo.TestResult = "ParentReady"
        }
        InitNodesRunResult(n, "")
    }
}


func ScheduleNodes(node *tcNode, wg *sync.WaitGroup, priority string, resultsChan chan testcase.TestCaseExecutionInfo, 
        pStart string, baseUrl string, resultsDir string) {
    //
    tick := 0
    max := cmd.Opt.ConcurrencyLimit
    //
    for _, n := range node.children {
        if priority == n.TestCaseExecutionInfo.Priority() && n.TestCaseExecutionInfo.TestResult == "Ready"{
            wg.Add(1)
            // Note: to prevent to tcp connection, here set a max, then sleep for a while
            if tick % max == 0 {
                time.Sleep(500 * time.Millisecond)
                go api.HttpApi(wg, resultsChan, pStart, baseUrl, n.TestCaseExecutionInfo.TestCaseDataInfo, resultsDir)
            } else {
                go api.HttpApi(wg, resultsChan, pStart, baseUrl, n.TestCaseExecutionInfo.TestCaseDataInfo, resultsDir)
            }

            tick = tick + 1
        }
        
        ScheduleNodes(n, wg, priority, resultsChan, pStart, baseUrl, resultsDir)
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

    for _, n := range node.children {
        if tcRunResult == "Fail"{
            n.TestCaseExecutionInfo.TestResult = "ParentFailed"
        } else if tcRunResult == "Success"{
            n.TestCaseExecutionInfo.TestResult = "Ready"
        }
    }
}

func CollectNodeReadyStatus(node *tcNode, priority string) {
    for _, n := range node.children {
        if n.TestCaseExecutionInfo.Priority() == priority {
            switch n.TestCaseExecutionInfo.TestResult { 
                case "Ready": 
                    statusReadyCount = statusReadyCount + 1
            }
        }
        CollectNodeReadyStatus(n, priority)
    }
}

func CollectNodeStatusByPriority(node *tcNode, p_index int, priority string) {
    for _, n := range node.children {
        if n.TestCaseExecutionInfo.Priority() == priority {
            statusCountList[p_index][0] = statusCountList[p_index][0] + 1

            switch n.TestCaseExecutionInfo.TestResult { 
                case "Ready": 
                    statusCountList[p_index][1] = statusCountList[p_index][1] + 1
                case "Success": 
                    statusCountList[p_index][2] = statusCountList[p_index][2] + 1  
                case "Fail":
                    statusCountList[p_index][3] = statusCountList[p_index][3] + 1 
                default: 
                    statusCountList[p_index][4] = statusCountList[p_index][4] + 1
                    // write the cases to tcNotExecutedList
                    tcNotExecutedList = append(tcNotExecutedList, n.TestCaseExecutionInfo)
                    
            }
        }
        CollectNodeStatusByPriority(n, p_index, priority)
    }
}

func CollectOverallNodeStatus(node *tcNode, p_index int) {
    for _, n := range node.children {
        statusCountList[p_index][0] = statusCountList[p_index][0] + 1

        switch n.TestCaseExecutionInfo.TestResult { 
            case "Ready": 
                statusCountList[p_index][1] = statusCountList[p_index][1] + 1
            case "Success": 
                statusCountList[p_index][2] = statusCountList[p_index][2] + 1  
            case "Fail":
                statusCountList[p_index][3] = statusCountList[p_index][3] + 1 
            default: 
                statusCountList[p_index][4] = statusCountList[p_index][4] + 1
        }
        CollectOverallNodeStatus(n, p_index)
    }
}


func ShowNodes(node *tcNode) {
    fmt.Println("\nNN P node:", node.TestCaseExecutionInfo.Priority(), node.TestCaseExecutionInfo.TcName(), node.TestCaseExecutionInfo.TestResult)
    for _, n := range node.children {
        // fmt.Println("NN - C node:", &n, n)
        ShowNodes(n)
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
