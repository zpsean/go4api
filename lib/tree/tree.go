/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package tree

import (                                                                                                                                             
    "os"
    "fmt"
    "encoding/json"

    "go4api/lib/testcase"
)

type TcNode struct {
    TestCaseExecutionInfo *testcase.TestCaseExecutionInfo
    Children []*TcNode // for child
}

type TcTree map[string]*TcNode

func CreateTcTree () TcTree {
    var tcTree = TcTree{}

    return tcTree
}


func (tcTree TcTree) BuildTree (tcArray []*testcase.TestCaseDataInfo) (*TcNode) {
    fmt.Println("\n---- Build the tcTree - Start ----")

    root, ifAdded := tcTree.BuildRootNode()

    if ifAdded && true {
        tcArrayTree, tcArrayNotTree := tcTree.BuildRootDirectChildrenNodes(root, tcArray)

        tcTree.LoopAndBuildOtherNodes(root, tcArrayTree, tcArrayNotTree)
    } else {
        fmt.Println("\n!! Error, build root for TcTree Failed")
        os.Exit(1)
    }
  
    fmt.Println("---- Build the tcTree - END ----")

    return root
}

func GetDummyRootTc () testcase.TestCase {
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

func (tcTree TcTree) BuildRootNode () (*TcNode, bool) {
    rootTc := GetDummyRootTc()

    rootTcaseData := testcase.TestCaseDataInfo {
        TestCase: &rootTc,
        JsonFilePath: "",
        CsvFile: "",
        CsvRow: "",
    }

    rootTcaseExecution := &testcase.TestCaseExecutionInfo {
        TestCaseDataInfo: &rootTcaseData,
        TestResult: "",
        ActualStatusCode: 0,
        StartTime: "",
        EndTime: "",
        HttpTestMessages: []*testcase.TestMessage{},
        StartTimeUnixNano: 0,
        EndTimeUnixNano: 0,
        DurationUnixNano: 0,
    }

    root, ifAdded := tcTree.AddRootNode(rootTcaseExecution) 

    return root, ifAdded
}

func (tcTree TcTree) BuildRootDirectChildrenNodes (root *TcNode, tcArray []*testcase.TestCaseDataInfo) ([]*testcase.TestCaseDataInfo, []*testcase.TestCaseDataInfo) {
    var tcArrayTree []*testcase.TestCaseDataInfo
    var tcArrayNotTree []*testcase.TestCaseDataInfo

    for i, _ := range tcArray {
        tcaseExecution := &testcase.TestCaseExecutionInfo {
            TestCaseDataInfo: tcArray[i],
            TestResult: "",
            ActualStatusCode: 0,
            StartTime: "",
            EndTime: "",
            HttpTestMessages: []*testcase.TestMessage{},
            StartTimeUnixNano: 0,
            EndTimeUnixNano: 0,
            DurationUnixNano: 0,
        }

        ifAdded := tcTree.AddNode(root, tcaseExecution)
        if ifAdded && true {
            tcArrayTree = append(tcArrayTree, tcArray[i])
        } else {
            tcArrayNotTree = append(tcArrayNotTree, tcArray[i])
        }
    }

    return tcArrayTree, tcArrayNotTree
}

func (tcTree TcTree) LoopAndBuildOtherNodes (root *TcNode, tcArrayTree []*testcase.TestCaseDataInfo, tcArrayNotTree []*testcase.TestCaseDataInfo) {
    for {
        var tcArrayNotTreeTemp []*testcase.TestCaseDataInfo
        for i, _ := range tcArrayNotTree {
            tcArrayNotTreeTemp = append(tcArrayNotTreeTemp, tcArrayNotTree[i])
        }
        //
        len1 := len(tcArrayNotTreeTemp)
        for i, _ := range tcArrayNotTreeTemp {
            tcaseExecution := &testcase.TestCaseExecutionInfo {
                TestCaseDataInfo: tcArrayNotTreeTemp[i],
                TestResult: "",
                ActualStatusCode: 0,
                StartTime: "",
                EndTime: "",
                HttpTestMessages: []*testcase.TestMessage{},
                StartTimeUnixNano: 0,
                EndTimeUnixNano: 0,
                DurationUnixNano: 0,
            }

            ifAdded := tcTree.AddNode(root, tcaseExecution)
            
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

    if len(tcArrayNotTree) > 0 {
        notTree, _ := json.Marshal(tcArrayNotTree)
        fmt.Println("!!!Attention, there are test cases which parentTestCase does not exists\n", string(notTree))
    }
}

func (tcTree TcTree) AddRootNode (TcaseExecution *testcase.TestCaseExecutionInfo) (*TcNode, bool) {
    node := &TcNode{
        TestCaseExecutionInfo: TcaseExecution, 
        Children: []*TcNode{},
    }
    var ifAdded bool
    
    parentTcName := node.TestCaseExecutionInfo.ParentTestCase()
    // the dummy root tese case's parentTestCase == "0"
    if parentTcName == "0" {
        tcTree["root"] = node
        ifAdded = true
    } else {
        fmt.Println("\n!! Error, root node can not be added")
        os.Exit(1)
    }

    return node, ifAdded
}

func (tcTree TcTree) AddNode (root *TcNode, TcaseExecution *testcase.TestCaseExecutionInfo) bool {
    node := &TcNode{
        TestCaseExecutionInfo: TcaseExecution, 
        Children: []*TcNode{},
    }
    var ifAdded bool
    
    parentTcName := node.TestCaseExecutionInfo.ParentTestCase()
    // the dummy root tese case's parentTestCase == "0"
    if parentTcName == "root" {
        parent := tcTree[parentTcName]
        parent.Children = append(parent.Children, node)

        ifAdded = true
    } else if parentTcName != "0" {
        // below is to find out the right parent node, then add the child node for it if found
        c := make(chan *TcNode)
        go func(c chan *TcNode) {
            defer close(c)
            tcTree.SearchNode(c, root, parentTcName)
        }(c)

        var resNodes []*TcNode
        for n := range c {
            resNodes = append(resNodes, n)
        }

        if len(resNodes) > 1 {
            fmt.Println("\n!! Error, more than one parent node found, please verify the test data")
            os.Exit(1)
        } else if len(resNodes) == 1 {
            parent := resNodes[0]
            parent.Children = append(parent.Children, node)   

            ifAdded = true
        } else {
            ifAdded = false 
        }
    } else {
        fmt.Println("\n!! Error, node with parentTcName = 0, it shoud be for root only, not added")
        os.Exit(1)
    }

    return ifAdded
}


func (tcTree TcTree) SearchNode (c chan *TcNode, node *TcNode, testCaseName string) {
    for i, _ := range node.Children {
        if node.Children[i].TestCaseExecutionInfo.TcName() == testCaseName {
            c <- node.Children[i]
        } else {
            tcTree.SearchNode(c, node.Children[i], testCaseName)
        }
    }
}


func (tcTree TcTree) InitNodesRunResult (node *TcNode, runResult string) {
    for i, _ := range node.Children {
        // check if the current tese case is parent to others, if yes, then refresh
        if node.Children[i].TestCaseExecutionInfo.ParentTestCase() == "root"{
            node.Children[i].TestCaseExecutionInfo.TestResult = "Ready"
        } else {
            node.Children[i].TestCaseExecutionInfo.TestResult = "ParentReady"
        }
        tcTree.InitNodesRunResult(node.Children[i], "")
    }
}

func (tcTree TcTree) RefreshNodeAndDirectChilrenTcResult(node *TcNode, tcRunResult string, tcStart string, tcEnd string, tcRunMessage []*testcase.TestMessage, 
        tcStartUnixNano int64, tcEndUnixNano int64) {
    // --------
    node.TestCaseExecutionInfo.TestResult = tcRunResult
    node.TestCaseExecutionInfo.StartTime = tcStart
    node.TestCaseExecutionInfo.EndTime = tcEnd
    node.TestCaseExecutionInfo.HttpTestMessages = tcRunMessage
    node.TestCaseExecutionInfo.StartTimeUnixNano = tcStartUnixNano
    node.TestCaseExecutionInfo.EndTimeUnixNano = tcEndUnixNano
    node.TestCaseExecutionInfo.DurationUnixNano = tcEndUnixNano - tcStartUnixNano

    for i, _ := range node.Children {
        if tcRunResult == "Fail"{
            node.Children[i].TestCaseExecutionInfo.TestResult = "ParentFailed"
        } else if tcRunResult == "Success"{
            node.Children[i].TestCaseExecutionInfo.TestResult = "Ready"
        }
    }
}

func (tcTree TcTree) RefreshNodeAndChilrenTcResult(node *TcNode, tcRunResult string, tcStart string, tcEnd string, tcRunMessage []*testcase.TestMessage, 
        tcStartUnixNano int64, tcEndUnixNano int64) {
    // --------
    node.TestCaseExecutionInfo.TestResult = tcRunResult
    node.TestCaseExecutionInfo.StartTime = tcStart
    node.TestCaseExecutionInfo.EndTime = tcEnd
    node.TestCaseExecutionInfo.HttpTestMessages = tcRunMessage
    node.TestCaseExecutionInfo.StartTimeUnixNano = tcStartUnixNano
    node.TestCaseExecutionInfo.EndTimeUnixNano = tcEndUnixNano
    node.TestCaseExecutionInfo.DurationUnixNano = tcEndUnixNano - tcStartUnixNano

    for i, _ := range node.Children {
        switch tcRunResult { 
        case "Fail": 
            // node.Children[i].TestCaseExecutionInfo.TestResult = "ParentFailed"

            tcTree.RefreshNodeAndChilrenTcResult(node.Children[i], "ParentFailed", tcStart, tcEnd, 
                tcRunMessage, tcStartUnixNano, tcEndUnixNano)
        case "ParentFailed": 
            // node.Children[i].TestCaseExecutionInfo.TestResult = "ParentSkipped"

            tcTree.RefreshNodeAndChilrenTcResult(node.Children[i], "ParentSkipped", tcStart, tcEnd, 
                tcRunMessage, tcStartUnixNano, tcEndUnixNano)
        case "ParentSkipped": 
            // node.Children[i].TestCaseExecutionInfo.TestResult = "ParentSkipped"

            tcTree.RefreshNodeAndChilrenTcResult(node.Children[i], "ParentSkipped", tcStart, tcEnd, 
                tcRunMessage, tcStartUnixNano, tcEndUnixNano)
        case "Success":
            node.Children[i].TestCaseExecutionInfo.TestResult = "Ready"
        }
    }
}

func (tcTree TcTree) CollectNodeReadyByPriority(c chan *TcNode, node *TcNode, priority string) {
    for i, _ := range node.Children {
        if node.Children[i].TestCaseExecutionInfo.Priority() == priority {
            switch node.Children[i].TestCaseExecutionInfo.TestResult { 
                case "Ready": 
                    c <- node.Children[i]
            }
        }
        tcTree.CollectNodeReadyByPriority(c, node.Children[i], priority)
    }
}

func (tcTree TcTree) ShowNodes(node *TcNode) {
    fmt.Println("\nShow node: ", node)
    for i, _ := range node.Children {
        tcTree.ShowNodes(node.Children[i])
    }
}


func RemoveArrayItem(sourceArray []*testcase.TestCaseDataInfo, tcData *testcase.TestCaseDataInfo) []*testcase.TestCaseDataInfo {
    // fmt.Println("RemoveArryaItem", sourceArray, tc)
    var resultArray []*testcase.TestCaseDataInfo
    // resultArray := append(sourceArray[:index], sourceArray[index + 1:]...)
    for index, tc_i := range sourceArray {
        if tc_i.TcName() == tcData.TcName() {
            resultArray = append(sourceArray[:index], sourceArray[index + 1:]...)
            break
        }
    }

    return resultArray
}
