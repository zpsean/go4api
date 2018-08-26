/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package fuzz

import (
    "fmt"
    "sort"
    // "sync"
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
func GetPairWiseValid11(fuzzData FuzzData, PwLength int) {
    var combos [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            combos = append(combos, validList)
        }
    }


    var allPairs AllPairs

    allPairs.PwLength = PwLength
    allPairs.MaxPairWiseCombinationNumber = GetMaxPairWiseCombinationNumber(combos, PwLength)
    allPairs.WorkingItemMatrix = GetWorkingItemMatrix(combos)


    allPairs.NextPairWiseTestCaseData()

    // NextPairWiseTestCaseData(combos, PwLength)
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
func (item Item) Append (weight int) {
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
func GetPairIds(items []interface{}) []interface{} {
    // Note: in Python, items is tuple(item, item, ...)
    // (<allpairspy.allpairs.Item object at 0x1015dd3c8>,) 
    for ind, cacheItems := range keyCache.PairItems {
        totalFound := 0
        for i, _ := range cacheItems {
            if cacheItems[i] == items[i] {
                totalFound = totalFound + 1
                if totalFound == len(items) {
                    return keyCache.PairIds[ind]
                }
            }
        }
    }

    var keyIds []interface{}
    for _, id := range items {
        keyIds = append(keyIds, id)
    }
    // set keyCache
    keyCache.PairItems = append(keyCache.PairItems, items)
    keyCache.PairIds = append(keyCache.PairIds, keyIds)

    // Note: in Python, key_value (keyIds) is tuple(id, id, ...)
    // ('a0v2',) 
    return keyIds
}
///


// for the element already added / arranged
type Node struct {
    Id string
    Counter int
    InIds []string
    OutIds []string
}

// for the elements already added / arranged
type PairsStorage struct {
    PwLength int
    PwNodes map[string]Node // Note: in Python, this is dict -> {'id': Node, ...}
    PwCombsArr [][][]interface{} //[ [[1-comb id], [1-combid], [...]], [[2-comb ids], [2-combids], [...]], ...]
    // Note: in Python PWCombsArr is array for set(id), like, [set(id), set(id), ...]:
    // [ {('a0v0',), ..., ('a3v1',), ('a1v0',)}, 
    // {('a1v1', 'a2v3'), ..., ('a0v1', 'a2v0')} ]
}


func (pairs PairsStorage) GetNodeInfo (item interface{}) Node {
    var node Node
    node.Id = item.(Item).Id
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
        for combination := range combinationsInterface(sequence, i + 1) {
            pairs.AddCombination(combination)
        }
    }     
}

// combination -> [1]item, ..., [PwLength]item
func (pairs PairsStorage) AddCombination(combination []interface{}) {
    n := len(combination)
    if n > 0 {
            pairs.PwCombsArr[n - 1] = append(pairs.PwCombsArr[n - 1], GetPairIds(combination))
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

                pairs.PwNodes[combination[0].(Item).Id] = node
            } 
        }
        
        var ids []string
        for _, item := range combination {
            ids = append(ids, item.(Item).Id)
        }
        for i, id := range ids {
            curr := pairs.PwNodes[id]
            curr.Counter = curr.Counter + 1

            tempInIds := curr.InIds
            for _, id_i := range ids[:i] {
                for _, id_ii := range curr.InIds {
                    if id_i == id_ii {
                        tempInIds = append(tempInIds, id_i)
                    }
                } 
            }

            tempOutIds := curr.OutIds
            for _, id_i := range ids[i + 1:] {
                for _, id_ii := range curr.OutIds {
                    if id_i == id_ii {
                        tempInIds = append(tempOutIds, id_i)
                    }
                } 
            }
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
    indexCombs := combinations(indexSlice, pwLen)
    //
    totalNumber := 0
    //
    for indexPair := range indexCombs {
        var combos_pw_index_slice [][]interface{}
        for _, ind_value := range indexPair {
            var indexSlice []interface{}
            for i, _ := range combs[ind_value] {
                indexSlice = append(indexSlice, i)
            }

            combos_pw_index_slice = append(combos_pw_index_slice, indexSlice)
        }
        //
        c := make(chan []interface{})
        go func(c chan []interface{}) {
            defer close(c)
            combosSliceString(c, []interface{}{}, combos_pw_index_slice)
        }(c)

        // can not use len(c) to get the channel length, as len(c) is always 0 here, why?
        cLenght := 0
        for range c{
            cLenght = cLenght + 1
        }

        // fmt.Println("c: ", cLenght, len(c))
        totalNumber = totalNumber + cLenght
    }

    fmt.Println("MaxPairWiseCombinationNumber: ", totalNumber)
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
    fmt.Println("workingItemMatrix: ", workingItemMatrix)
    return workingItemMatrix
}


// -------------------------------------------------------------------------
func (allPairs AllPairs) NextPairWiseTestCaseData() []interface{} {
    var pairs PairsStorage
    pairs.PwLength = allPairs.PwLength

    maxUniquePairsExpected := allPairs.MaxPairWiseCombinationNumber
    // if pairs.Length() > maxUniquePairsExpected {
    //     os.Exit(1)
    // }
    if pairs.Length() == maxUniquePairsExpected {
        return []interface {}{}
    }
    workingItemMatrix := allPairs.WorkingItemMatrix

    // previousUniquePairsCount = len(pairs)
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
            if direction < -1 {
                return []interface {}{}
            }
            direction = 1
        } else {
            direction = 0
        }
        i += direction
        //
        fmt.Println("indexes: ", i, indexes, chosenValuesArr)
    }
   
    if len(workingItemMatrix) != len(chosenValuesArr) {
        fmt.Println("stop and return")
        return []interface {}{}
    }

    pairs.AddSequence(chosenValuesArr)

    // if pairs.Length() == previousUniquePairsCount {
    //     // could not find new unique pairs - stop
    //     return []interface {}{}
    // }

    // replace returned array elements with real values and return it
    var chosenValues []interface{}
    for _, item := range chosenValuesArr {
        chosenValues = append(chosenValues, item.(Item).Value)
    }
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
            chosenValuesArr = append(chosenValuesArr, item)
        
            var setPairIds [][]interface{}
            for z := range combinationsInterface(chosenValuesArr, i + 1) {
                if allPairs.Pairs.Length() > 0 {
                    for _, comb := range allPairs.Pairs.GetCombs()[i] {
                        if compareSlice(GetPairIds(z), comb) {
                            setPairIds = append(setPairIds, GetPairIds(z))
                        }
                    }
                }
            }
            newCombs = append(newCombs, setPairIds)
        }

        // (1). weighting the node
        // node that creates most of new pairs is the best
        item.Weights = append(item.Weights, -len(newCombs[len(newCombs) - 1]))

        // (2). less used outbound connections most likely to produce more new
        // pairs while search continues
        item.Weights = append(item.Weights, len(dataNode.OutIds))

        var reversedNewCombs [][][]interface{}
        for i := len(newCombs) - 2; i > 0; i-- {
            reversedNewCombs = append(reversedNewCombs, newCombs[i])
        }
        // (3). 
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
        fmt.Println("weights--: ", allPairs.WorkingItemMatrix[num][ii], item.Weights)
    }

    // workingItemMatrix[num].sort(key=cmp_to_key(cmp_item))
    // Sort: Ascending order
    // sort.Sort(allPairs.WorkingItemMatrix[num])
    var items Items
    items = allPairs.WorkingItemMatrix[num]
    fmt.Println("items--before sort: ", items)
    sort.Sort(items)
    fmt.Println("items--after sort: ", items)
}

func compareSlice(sliceA []interface{}, sliceB []interface{}) bool {
    var ifMatched bool
    ifMatched = true
    for i, valueA := range sliceA {
        if valueA != sliceB[i] {
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
    lenI := len(items[i].Weights)
    lenJ := len(items[j].Weights)

    var ifLess bool
    ifLess = true
    if lenI < lenJ {
        for ii, itemI := range items[i].Weights {
            if itemI > items[j].Weights[ii] {
                ifLess = false
                break
            }
        }
    } else {
        for jj, itemJ := range items[j].Weights {
            if items[i].Weights[jj] >= itemJ {
                ifLess = false
                break
            }
        }
    }

    return ifLess
}

func (items Items) Swap(i, j int) {
    items[i], items[j] = items[j], items[i] 
}





