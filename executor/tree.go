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
    // "time"
    "fmt"
    "sync"
    "go4api/api"
    "go4api/utils"
    simplejson "github.com/bitly/go-simplejson"
    // "strconv"
)

type tcNode struct{
    tcName string
    parentTestCase string // for parent
    tcRunResult string // Ready, Running, Success, Fail, ParentReady, ParentRunning, ParentFailed
    tc []interface{}
    children []*tcNode // for child
}

var (
    tcTree = map[string]*tcNode{}
    root *tcNode

    findNode **tcNode

    statusReadyCount int
    statusCountList [][]int
)



//
func GetDummyRootTc() []interface{} {
    // dummy root tc => {"root", "0", "0", rooTC, "", "", ""}
    rootTC, _ := simplejson.NewJson([]byte("{}}"))
    var rootTcInfo []interface{}
    rootTcInfo = append(rootTcInfo, "root")
    rootTcInfo = append(rootTcInfo, "0")
    rootTcInfo = append(rootTcInfo, "0")
    rootTcInfo = append(rootTcInfo, rootTC)
    rootTcInfo = append(rootTcInfo, "")
    rootTcInfo = append(rootTcInfo, "")
    rootTcInfo = append(rootTcInfo, "")

    return rootTcInfo
}

func AddNode(tcName string, parentTestCase string, tcRunResult string, tc []interface{}) bool {
    node := &tcNode{
        tcName: tcName, 
        parentTestCase: parentTestCase, 
        tcRunResult: tcRunResult, 
        tc: tc, 
        children: []*tcNode{}}

    // fmt.Println("\nc-node: ", node)
    var ifAdded bool
    // the dummy root tese case's parentTestCase == "0"
    if parentTestCase == "0" {
        root = node
        tcTree["root"] = root

        ifAdded = true
    } else if parentTestCase == "root" {
        parent := tcTree[parentTestCase]
        parent.children = append(parent.children, node)   

        ifAdded = true
    } else {
        findNode = nil

        // below is to find out the right parent node, then add the child node for it
        SearchNode(&root, parentTestCase)
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

func BuildTree(tcArray [][]interface{}) (*tcNode, map[string]*tcNode) {
    fmt.Println("\n---- Build the tcTree - Start ----")
    var tcArrayTree [][]interface{}
    var tcArrayNotTree [][]interface{}
    // here seperate the tcArray into two parts, (appended to tree) and (not yet appended to tree)
    // Step 1: add the root node
    rootTcInfo := GetDummyRootTc()
    AddNode(rootTcInfo[0].(string), rootTcInfo[2].(string), "", rootTcInfo)

    // Step 2: add the node, init the tcArrayTree, tcArrayNotTree
    for _, tc := range tcArray {
        // if parentTestCase name can be found in tree
        // tcRunResult is "" for init
        ifAdded := AddNode(tc[0].(string), tc[2].(string), "", tc)
        // fmt.Println("!!now try to add tc: ", tc[0].(string), ifAdded)
        if ifAdded && true {
            tcArrayTree = append(tcArrayTree, tc)
        } else {
            tcArrayNotTree = append(tcArrayNotTree, tc)
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

        for i, tc := range tcArrayNotTree {
            // if parentTestCase name can be found in tree
            ifAdded := AddNode(tc[0].(string), tc[2].(string), "", tc)
            // fmt.Println("!!now try to add tc: ", tc[0].(string), ifAdded)
            if ifAdded && true {
                tcArrayTree = append(tcArrayTree, tc)
                tcArrayNotTree = utils.RemoveArryaItem(tcArrayNotTree, i)
            }
        }

        // fmt.Println("\nlen tcArrayNotTree: ", len(tcArrayNotTree), tcArrayNotTree)
        len2 := len(tcArrayNotTree)
        // if can not add anymore
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
        if (*node).children[i].tcName == testCaseNmae {
            findNode = &((*node).children[i])
            return findNode
        } else {
            SearchNode(&((*node).children[i]), testCaseNmae)
        }
    }
    return findNode
}


func SearchNode5(node *tcNode, testCaseNmae string) *tcNode {
    var findNode *tcNode
    for _, n := range node.children {
        if n.tcName == testCaseNmae {
            fmt.Println("ccccccccc tree found node for 22: ", testCaseNmae, n)
            findNode = node
            return findNode
        } else {
            findNode = SearchNode5(n, testCaseNmae)
        }
        fmt.Println("???c why can touch here?")
    }
    return nil
}

func InitNodesRunResult(node *tcNode, runResult string) {
    for _, n := range node.children {
        // check if the current tese case is parent to others, if yes, then refresh
        if n.parentTestCase == "root"{
            n.tcRunResult = "Ready"
        } else {
            n.tcRunResult = "ParentReady"
        }
        InitNodesRunResult(n, "")
    }
}


func ScheduleNodes(node *tcNode, wg *sync.WaitGroup, options map[string]string, priority string, resultsChan chan []interface{}, 
        pStart string, baseUrl string, resultsDir string) {
    //
    for _, n := range node.children {
        if priority == n.tc[1].(string) && n.tcRunResult == "Ready"{
            wg.Add(1)
            // Note: how to control one test case not be be run more than once???
            go api.HttpApi(wg, resultsChan, options, pStart, baseUrl, n.tc, resultsDir)
        }
        
        ScheduleNodes(n, wg, options, priority, resultsChan, pStart, baseUrl, resultsDir)
    }
}


func RefreshNodesRunResult(node *tcNode, tcRunResultsArray [][]interface{}) {
    for _, n := range node.children {
        // check if the current tese case is parent to others, if yes, then refresh
        var childrenTcRunResultsArray [][]interface{}
        for _, v := range tcRunResultsArray {
            if v[1] == "Fail"{
                var vv []interface{}
                vv[0] = ""
                vv[0] = "ParentFailed"
                childrenTcRunResultsArray = append(childrenTcRunResultsArray, vv)
                RefreshNodesRunResult(n, childrenTcRunResultsArray)
            }
        }
        
        // match the test case
        // if n.tc[0].(string) == tcRunResultsArray[0].(string){
        //     // update the tcRunResult, it is tc[7]
        //     n.tc[7] = tcRunResultsArray[1].(string)
        // }
        RefreshNodesRunResult(n, tcRunResultsArray)
    }
}


func RefreshNodeAndDirectChilrenTcResult(node *tcNode, tcRunResult string) {
    // fmt.Println("ccc the node to be refreshed: ", node, &node)
    node.tcRunResult = tcRunResult
    for _, n := range node.children {
        if tcRunResult == "Fail"{
            n.tcRunResult = "ParentFailed"
        } else if tcRunResult == "Success"{
            n.tcRunResult = "Ready"
        }
    }
}

func CollectNodeReadyStatus(node *tcNode, priority string) {
    for _, n := range node.children {
        if n.tc[1].(string) == priority {
            switch n.tcRunResult { 
                case "Ready": 
                    statusReadyCount = statusReadyCount + 1
            }
        }
        CollectNodeReadyStatus(n, priority)
    }
}

func CollectNodeStatusByPriority(node *tcNode, p_index int, priority string) {
    for _, n := range node.children {
        if n.tc[1].(string) == priority {
            statusCountList[p_index][0] = statusCountList[p_index][0] + 1

            switch n.tcRunResult { 
                case "Ready": 
                    statusCountList[p_index][1] = statusCountList[p_index][1] + 1
                case "Success": 
                    statusCountList[p_index][2] = statusCountList[p_index][2] + 1  
                case "Fail":
                    statusCountList[p_index][3] = statusCountList[p_index][3] + 1 
                default: 
                    statusCountList[p_index][4] = statusCountList[p_index][4] + 1
            }
        }
        CollectNodeStatusByPriority(n, p_index, priority)
    }
}

func CollectOverallNodeStatus(node *tcNode, p_index int) {
    for _, n := range node.children {
        statusCountList[p_index][0] = statusCountList[p_index][0] + 1

        switch n.tcRunResult { 
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
    fmt.Println("\nNN P node:", &node, node, node.children)
    for _, n := range node.children {
        // fmt.Println("NN - C node:", &n, n)
        ShowNodes(n)
    }
}

