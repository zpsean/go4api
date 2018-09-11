/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package pairwise

import (
    "fmt"
    "sort"
    "reflect"

    combins "go4api/lib/combination"
)

// -------------------------------------------------------------------------
// Note: for pairewise, there are two kinds of algorithm:
// --> algorithm 1. use the pairWises to build the test case data (test case)
// --> algorithm 2. get the full combination first, the filter (remove) them based on the pairWises if repeated
// -------------------------------------------------------------------------
// The algorithm 1: is rewitten the code for AllPairs (python) using Golang
// refer to: https://github.com/thombashi/allpairspy
//
// key steps of this algorithm is:
// 1. get the value value_matrix (combs), and maxUniquePairsExpected, workingItemMatrix, 
// 2. try the next item to be add, with computing and comparing existing pairs, some weights, like most new pairs, in, out, etc.
// 3. sort the items in one vetor of the workingItemMatrix (workingItemMatrix(m))
// 4. add the highly recommended item of the sort items (step 3)
// 5. add the item to existing pairs
// 6. repeat Step 2 ~ Step 5
// -------------------------------------------------------------------------

// here is the entry for pairwise algorithm 1
func GetPairWise(c chan []interface{} ,combins [][]interface{}, PwLength int) {
    // init the PairsStorage
    var pairs PairsStorage
    pairs.PwLength = PwLength

    pwNodes := map[string]Node{}
    pairs.PwNodes = pwNodes

    for i := 0; i < PwLength; i++ {
        pairs.PwCombsArr = append(pairs.PwCombsArr, [][]interface{}{})
    }

    // init the AllPairs
    var allPairs AllPairs

    allPairs.PwLength = PwLength
    allPairs.Pairs = pairs
    allPairs.MaxPairWiseCombinationNumber = GetMaxPairWiseCombinationNumber(combins, PwLength)
    allPairs.WorkingItemMatrix = GetWorkingItemMatrix(combins)
    
    // to get the data 
    // for debug loop
    loopDepth := 0
    for {
        returnedTestCaseData := allPairs.NextPairWiseTestCaseData()

        if len(returnedTestCaseData) == 0 {
            break
        } else {
            c <- returnedTestCaseData
            // fmt.Println(returnedTestCaseData)
        }

        if loopDepth == 30 {
            fmt.Println("loopDepth: ", loopDepth)
            break
        }
    }
    // fmt.Println("touch here ????? => yes")
    // allPairs.NextPairWiseTestCaseData()
    // allPairs.NextPairWiseTestCaseData()
}


// to define the item for each element
type AllPairs struct {
    PwLength int
    Pairs PairsStorage
    MaxPairWiseCombinationNumber int
    WorkingItemMatrix [][]Item  // [][]Item
    PairWiseTestCaseData [][]interface{} // [][]Item.Value
}


type Item struct {
    Id string
    Value interface{}
    Weights []int
}

//
func (item Item) Append(weight int) {
    item.Weights = append(item.Weights, weight)
}


// Note: in Python, key_cache (KeyCache) is dict, with tuple(items) as key, tuple(ids) is value
// {
//     (<allpairspy.allpairs.Item object at 0x1015dd2e8>,): ('a0v0',), 
//     (<allpairspy.allpairs.Item object at 0x1015dd358>,): ('a0v1',), 
//     (<allpairspy.allpairs.Item object at 0x1015dd3c8>,): ('a0v2',)
// }
type KeyCache struct {
    PairItems [][]interface{} // [][]Item ->for Item
    PairIds [][]interface{} // [][]string -> for id
}

var keyCache KeyCache

// items -> []Item, return []string -> []Id
func GetPairIdsBk(items []interface{}) []interface{} {
    // Note: in Python, items is tuple(item, item, ...)
    // (<allpairspy.allpairs.Item object at 0x1015dd3c8>,) 
    for ind, cacheItems := range keyCache.PairItems {
        totalFound := 0
        fmt.Println("keyCache, items: ", keyCache, "::", items)

        for i, _ := range cacheItems {
            // to compare Item, need special method, use the reflect.DeepEqual
            // if cacheItems[i] == items[i] {
            if reflect.DeepEqual(cacheItems[i], items[i]) {
                totalFound = totalFound + 1
                if totalFound == len(items) {
                    return keyCache.PairIds[ind]
                }
            }
        }
    }

    var keyIds []interface{}
    for _, item := range items {
        keyIds = append(keyIds, item.(Item).Id)
    }
    // set keyCache
    keyCache.PairItems = append(keyCache.PairItems, items)
    keyCache.PairIds = append(keyCache.PairIds, keyIds)

    // Note: in Python, key_value (keyIds) is tuple(id, id, ...)
    // ('a0v2',) 
    return keyIds
}

func GetPairIds(items []interface{}) []interface{} {
    // Note: in Python, items is tuple(item, item, ...)
    // (<allpairspy.allpairs.Item object at 0x1015dd3c8>,) 
    var keyIds []interface{}
    for _, item := range items {
        keyIds = append(keyIds, item.(Item).Id)
    }
    // Note: in Python, key_value (keyIds) is tuple(id, id, ...)
    // ('a0v2',) 
    // in golang, keyIds = ['id', ...]
    return keyIds
}
///


// for the element already added / arranged
type Node struct {
    Id string
    Counter int
    InIds []string  // ensure no duplicate
    OutIds []string  // ensure no duplicate
}

// for the elements already added / arranged
type PairsStorage struct {
    PwLength int
    PwNodes map[string]Node // Note: in Python, this is dict -> {'id': Node, ...}
    PwCombsArr [][][]interface{} // ensure no duplicate
    //[ [[1-comb id], [1-combid], [...]], [[2-comb ids], [2-combids], [...]], ...]
    // Note: in Python PWCombsArr is array for set(id), like, [set(id), set(id), ...]:
    // [ {('a0v0',), ..., ('a3v1',), ('a1v0',)}, 
    // {('a1v1', 'a2v3'), ..., ('a0v1', 'a2v0')} ]
}


func (pairs PairsStorage) GetNodeInfo (item interface{}) Node {
    var node Node
    node.Id = item.(Item).Id
    node.Counter = 0
    for _, node_i := range pairs.PwNodes {
        if node_i.Id == item.(Item).Id {
            node = node_i
            break
        }
    }
    return node
}

func (pairs PairsStorage) GetCombs() [][][]interface{} {
    return pairs.PwCombsArr
}

func (pairs PairsStorage) Length() int {
    if len(pairs.PwCombsArr) > 0 {
        return len(pairs.PwCombsArr[len(pairs.PwCombsArr) - 1]) 
    } else {
        return 0
    }
}

// sequence -> [PwLength]item, which has been choseen for the next test data (case)
func (pairs PairsStorage) AddSequence(sequence []interface{}) {
    for i := 0; i < pairs.PwLength; i++ {
        for combination := range combins.CombinationsInterface(sequence, i + 1) {
            pairs.AddCombination(combination)
        }
    }     
}

// combination -> [1]item, ..., [PwLength]item
func (pairs PairsStorage) AddCombination(combination []interface{}) {
    // fmt.Println("combination: ", combination)
    n := len(combination)
    if n > 0 {
        ids := GetPairIds(combination)
        // Note: not duplicate: check if pairs.PwCombsArr[n - 1] has ids first, if not, then add
        if len(pairs.PwCombsArr[n - 1]) > 0 {
            var ifExists bool
            ifExists = false
            for _, set := range pairs.PwCombsArr[n - 1] {
                if compareSlice(set, ids) {
                    ifExists = true
                    break
                }
            }
            if ifExists == false {
                pairs.PwCombsArr[n - 1] = append(pairs.PwCombsArr[n - 1], ids)
            }
        } else {
            pairs.PwCombsArr[n - 1] = append(pairs.PwCombsArr[n - 1], ids)
        }
        
        // if n == 1 and combination[0].id not in self.__nodes:
        if n == 1 {
            var ifExists bool
            ifExists = false
            // to check if combination[0].(Item).Id already exists in keys of PwNodes (map[id]Node)
            for key, _ := range pairs.PwNodes {
                if combination[0].(Item).Id == key {
                    ifExists = true
                    break
                }
            }
            if ifExists == false {
                var node Node
                node.Id = combination[0].(Item).Id
                node.Counter = 0

                // fmt.Println("AddCombination: ", combination, combination[0].(Item).Id)
                pairs.PwNodes[combination[0].(Item).Id] = node

                return
            }
        }
        
        for i, id := range ids {
            // fmt.Println("ids: ", ids)
            var node Node

            node.Id = id.(string)

            curr := pairs.PwNodes[id.(string)] // curr is Node type
            // fmt.Println("curr -- 0: ", ids, curr)
            node.Counter = curr.Counter + 1
            

            tempInIds := curr.InIds
            for _, id_i := range ids[:i] {
                // ensure tempInIds has no duplicate elements
                var ifExists bool
                ifExists = false
                for _, id_ii := range tempInIds {
                    if id_ii == id_i.(string) {
                        ifExists = true
                        break
                    }
                }
                if ifExists == false {
                    tempInIds = append(tempInIds, id_i.(string))
                }
                
            }
            node.InIds = tempInIds

            tempOutIds := curr.OutIds
            for _, id_i := range ids[i + 1:] {
                // ensure tempOutIds has no duplicate elements
                var ifExists bool
                ifExists = false
                for _, id_ii := range tempInIds {
                    if id_ii == id_i.(string) {
                        ifExists = true
                        break
                    }
                }
                if ifExists == false {
                    tempOutIds = append(tempOutIds, id_i.(string))
                }
            }
            node.OutIds = tempOutIds

            // fmt.Println("curr -- 1: ", ids, curr, node)
            pairs.PwNodes[id.(string)] = node
        }
    }
}


// to get the total number of pairwise combinations
func GetMaxPairWiseCombinationNumber(combs [][]interface{}, PwLength int) int {
    // init -----------------
    var indexSlice []int
    for i := 0; i < len(combs); i++ {
        indexSlice = append(indexSlice, i)
    }
    // to get the combinations like [1 1][1 2][1 3] ...
    var pwLen int
    if PwLength > len(combs) {
        pwLen = len(combs)
    } else {
        pwLen = PwLength
    }
    indexCombs := combins.CombinationsInt(indexSlice, pwLen)
    //
    totalNumber := 0
    //
    for indexPair := range indexCombs {
        var combins_pw_index_slice [][]interface{}
        for _, ind_value := range indexPair {
            var indexSlice []interface{}
            for i, _ := range combs[ind_value] {
                indexSlice = append(indexSlice, i)
            }

            combins_pw_index_slice = append(combins_pw_index_slice, indexSlice)
        }
        //
        c := make(chan []interface{})
        go func(c chan []interface{}) {
            defer close(c)
            combins.CombinsSliceInterface(c, []interface{}{}, combins_pw_index_slice)
        }(c)

        // can not use len(c) to get the channel length, as len(c) is always 0 here, why?
        cLenght := 0
        for range c{
            cLenght = cLenght + 1
        }

        // fmt.Println("c: ", cLenght, len(c))
        totalNumber = totalNumber + cLenght
    }

    // fmt.Println("MaxPairWiseCombinationNumber: ", totalNumber)
    return totalNumber
}

// to get the item matrix with type Item
func GetWorkingItemMatrix(combs [][]interface{}) [][]Item {
    var workingItemMatrix [][]Item
    for i, combsSlice := range combs {
        var itemSlice []Item
        for j, value := range combsSlice {
            var strId string
            strId = "a" + fmt.Sprint(i) + "v" + fmt.Sprint(j)

            var item Item
            item.Id = strId
            item.Value = value

            // init the item.Weights, otherwise, will receive 'can not assign ...' error to change the weight later
            // var weights []int
            // item.Weights = weights

            itemSlice = append(itemSlice, item)
        }
        workingItemMatrix = append(workingItemMatrix, itemSlice)
    }
    // fmt.Println("workingItemMatrix: ", workingItemMatrix)
    return workingItemMatrix
}


// -------------------------------------------------------------------------
func (allPairs AllPairs) NextPairWiseTestCaseData() []interface{} {
    maxUniquePairsExpected := allPairs.MaxPairWiseCombinationNumber
    if allPairs.Pairs.Length() > maxUniquePairsExpected {
        fmt.Println("!! Error, added pairs more than maxUniquePairsExpected: ", allPairs.Pairs.Length(), maxUniquePairsExpected)
        return []interface {}{}
    }
    if allPairs.Pairs.Length() == maxUniquePairsExpected {
        fmt.Println("all pairs have been added: ", maxUniquePairsExpected)
        return []interface {}{}
    }
    workingItemMatrix := allPairs.WorkingItemMatrix

    previousUniquePairsCount := allPairs.Pairs.Length()
    chosenValuesArr := make([]interface{}, len(workingItemMatrix))
    indexes := make([]int, len(workingItemMatrix))

    direction := 1
    i := 0

    for {
        // to break the for if ...
        if i <= -1 || i >= len(workingItemMatrix) {
            break
        }
        //
        if direction == 1 {
            // move forward
            allPairs.resortWorkingArray(chosenValuesArr[:i], i)
            indexes[i] = 0
        } else if direction == 0 || direction == -1 {
            // scan current array or go back
            indexes[i] += 1
            if indexes[i] >= len(workingItemMatrix[i]) {
                direction = -1
                if i == 0 {
                    fmt.Println("stop and return")
                    return []interface {}{}
                }
                i += direction
                continue
            }
            direction = 0
        } else {
            fmt.Println("next(): unknown 'direction' code '{}'", direction)
        }

        chosenValuesArr[i] = workingItemMatrix[i][indexes[i]]

        if true {
            if direction <= -1 {
                return []interface {}{}
            }
            direction = 1
        } else {
            direction = 0
        }
        i += direction
        //
        // fmt.Println("indexes: ", i, indexes, chosenValuesArr)
    }
   
    if len(workingItemMatrix) != len(chosenValuesArr) {
        fmt.Println("stop and return")
        return []interface {}{}
    }

    allPairs.Pairs.AddSequence(chosenValuesArr)
    // fmt.Println("allPairs.Pairs.Length() vs. maxUniquePairsExpected: ", allPairs.Pairs.Length(), maxUniquePairsExpected)

    if allPairs.Pairs.Length() == previousUniquePairsCount {
        // could not find new unique pairs - stop
        fmt.Println("could not find new unique pairs - stop: ", previousUniquePairsCount)
        return []interface {}{}
    }

    // replace returned array elements with real values and return it
    var chosenValues []interface{}
    for _, item := range chosenValuesArr {
        chosenValues = append(chosenValues, item.(Item).Value)
    }
    // fmt.Println("-----> chosenValues: ", chosenValuesArr, chosenValues)
    return chosenValues
}


func (allPairs AllPairs) resortWorkingArray (chosenValuesArr []interface{}, num int) {
    for ii, item := range allPairs.WorkingItemMatrix[num] {
        dataNode := allPairs.Pairs.GetNodeInfo(item)

        // numbers of new combinations to be created if this item is
        // appended to array
        var newCombs [][][]interface{}
        // newCombs -> [][][]interface{} //[ [[1-comb id], [1-combid], [...]], [[2-comb ids], [2-combids], [...]], ...]
        for i := 0; i < allPairs.PwLength; i++ {
            var tempChosen []interface{}
            // copy(tempChosen, chosenValuesArr) => not sure why copy does not work here
            tempChosen = append(tempChosen, chosenValuesArr[0:]...)
            tempChosen = append(tempChosen, item)

            // fmt.Println("---> tempChosen : ", chosenValuesArr, item, tempChosen)

            var setPairIds [][]interface{}
            for z := range combins.CombinationsInterface(tempChosen, i + 1) {
                var idss []interface{}
                for _, item := range z {
                    idss = append(idss, item.(Item).Id)
                }

                if allPairs.Pairs.Length() > 0 {
                    var ifHas bool
                    ifHas = false
                    for _, comb := range allPairs.Pairs.GetCombs()[i] {
                        if compareSlice(idss, comb) {
                            ifHas = true
                            break
                        }
                    }
                    if ifHas == false {
                        setPairIds = append(setPairIds, idss)
                    }
                } else {
                    setPairIds = append(setPairIds, idss)
                }
            }
            newCombs = append(newCombs, setPairIds)
        }

        // reset the item.Weights
        weights := []int{}
        item.Weights = weights
        // (1). weighting the node
        // node that creates most of new pairs is the best
        item.Weights = append(item.Weights, -len(newCombs[len(newCombs) - 1]))

        // (2). less used outbound connections most likely to produce more new
        // pairs while search continues
        item.Weights = append(item.Weights, len(dataNode.OutIds))

        // (3). reverse the newCombs except the last [] as it is used in (1)
        var reversedNewCombs [][][]interface{}
        for i := len(newCombs) - 2; i >= 0; i-- {
            reversedNewCombs = append(reversedNewCombs, newCombs[i])
        }
        for _, x := range reversedNewCombs {
            item.Weights = append(item.Weights, len(x))
        }

        // (4). less used node is better
        item.Weights = append(item.Weights, -dataNode.Counter)

        // (5). otherwise we will prefer node with most of free inbound
        // connections; somehow it works out better ;)
        item.Weights = append(item.Weights, -len(dataNode.InIds))

        // re-assign the item.Weights to the allPairs
        allPairs.WorkingItemMatrix[num][ii].Weights = item.Weights
    }

    // workingItemMatrix[num].sort(key=cmp_to_key(cmp_item))
    // Sort: Ascending order
    // sort.Sort(allPairs.WorkingItemMatrix[num])
    var items Items
    items = allPairs.WorkingItemMatrix[num]
    // fmt.Println("items--before sort: ", items)
    sort.Stable(items)
    // fmt.Println("items--after sort: ", items, "\n")
    // fmt.Println("items--after sort : allPairs - PwNodes", allPairs.Pairs.PwNodes)
    // fmt.Println("items--after sort : allPairs - PwCombsArr", allPairs.Pairs.PwCombsArr, "\n\n")
}

func compareSlice(sliceA []interface{}, sliceB []interface{}) bool {
    var ifMatched bool
    ifMatched = true
    for i, valueA := range sliceA {
        // if valueA != sliceB[i] {
        if !reflect.DeepEqual(valueA, sliceB[i]) {
            ifMatched = false
            break
        }
    }
    return ifMatched
}


// Implements the Interface in sort package, used for Item sort
// type Interface interface {
//     Len() int
//     Less(i, j int) bool
//     Swap(i, j int)
// }
type Items []Item

func (items Items) Len() int { 
    return len(items) 
}

func (items Items) Less(i, j int) bool {
    // lenI := len(items[i].Weights)
    // lenJ := len(items[j].Weights)

    var ifLess bool
    // if lenI < lenJ {
    //     for ii, itemI := range items[i].Weights {
    //         if itemI > items[j].Weights[ii] {
    //             ifLess = false
    //             break
    //         }
    //     }
    // } else if lenI >= lenJ {
    //     for jj, itemJ := range items[j].Weights {
    //         if items[i].Weights[jj] > itemJ {
    //             ifLess = false
    //             break
    //         }
    //     }
    // }
    for ii, itemI := range items[i].Weights {
        if itemI < items[j].Weights[ii] {
            ifLess = true
            break
        } else if itemI > items[j].Weights[ii] {
            ifLess = false
            break
        }
    }

    return ifLess
}

func (items Items) Swap(i, j int) {
    items[i], items[j] = items[j], items[i] 
}




