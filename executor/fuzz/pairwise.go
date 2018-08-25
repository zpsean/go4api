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
    "reflect"
    "sync"
)

// Below is for pairwise data
type PairWises struct {
    PairWises map[string][]PairWise
}

type PairWise struct {
    PwLength int
    PwVectorIndices []int
    PwElementIndices []int
    PwValues []interface{}
}

// to define the item for each element
type PairWiseTestCaseDataSet []PairWiseTestCaseData

type PairWiseTestCaseData []Item

type Item struct {
    Id string
    Value interface{}
    Weights []int
}

// for the element already added / arranged
type Node struct {
    Id string
    Counter int
    InIds []string
    OutIds []string
}

// for the elements already added / arranged
type Pairs struct {
    PwLength int
    PwNodes map[string]Node
    PwCombsArr [][]interface{}
}

//
type KeyCache struct {
    PairItems [][]interface{}
    PairIds [][]interface{}
}

var keyCache KeyCache

func GetPairIds(items []interface{}) []interface{} {
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

    return keyIds
}
///
    






// here is the entry for pairwise algorithm 1
func GetPairWiseValid11(fuzzData FuzzData, PwLength int) {
    var combos [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            combos = append(combos, validList)
        }
    }

    GetWorkingItemMatrix(combos)
    GetMaxPairWiseCombinationNumber(combos, PwLength)

    NextPairWiseTestCaseData(combos, PwLength)

    NextPairWiseTestCaseData(combos, PwLength)
}


func (pairs Pairs) AddSequence(sequence []interface{}) {
    indexSlice := make([]int, pairs.PwLength)
    for i, _ := range indexSlice {
        for combination := range combinationsInterface(sequence, i + 1) {
            pairs.AddCombination(combination)
        }
    }     
}

func (pairs Pairs) GetNodeInfo (item interface{}) string {
    var node Node
    node.Id = item.(Item).Id
    nodeInfo := node.Id
    for _, node := range pairs.PwNodes {
        if node.Id == item.(Item).Id {
            nodeInfo = node.Id
            break
        }
    }
    return nodeInfo
}

func (pairs Pairs) GetCombs() [][]interface{} {
    return pairs.PwCombsArr
}

func (pairs Pairs) Length() int {
    if len(pairs.PwCombsArr) > 0 {
        return len(pairs.PwCombsArr[len(pairs.PwCombsArr) - 1]) 
    } else {
        return 0
    }
}

func (pairs Pairs) AddCombination(combination []interface{}) {
    n := len(combination)
    if n > 0 {
            pairs.PwCombsArr[n - 1] = append(pairs.PwCombsArr[n - 1], GetPairIds(combination))
        if n == 1 {
            for key, _ := range pairs.PwNodes {
                if combination[0].(Item).Id == key {
                    var node Node
                    node.Id = combination[0].(Item).Id

                    pairs.PwNodes[combination[0].(Item).Id] = node
                    break
                }
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
    for i, _ := range combs {
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
func GetWorkingItemMatrix(combs [][]interface{}) [][]interface{} {
    var workingItemMatrix [][]interface{}
    for i, combsSlice := range combs {
        var itemSlice []interface{}
        for j, value := range combsSlice {
            var strId string
            strId = "a" + fmt.Sprint(i) + "v" + fmt.Sprint(j)

            var item Item
            item.Id = strId
            item.Value = value

            itemSlice = append(itemSlice, item)
        }
        workingItemMatrix = append(workingItemMatrix, itemSlice)
    }
    fmt.Println("workingItemMatrix: ", workingItemMatrix)
    return workingItemMatrix
}


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

func NextPairWiseTestCaseData(combs [][]interface{}, PwLength int) []interface{} {
    var pairs Pairs
    pairs.PwLength = PwLength

    maxUniquePairsExpected := GetMaxPairWiseCombinationNumber(combs, PwLength)
    // if pairs.Length() > maxUniquePairsExpected {
    //     os.Exit(1)
    // }
    if pairs.Length() == maxUniquePairsExpected {
        return []interface {}{}
    }
    workingItemMatrix := GetWorkingItemMatrix(combs)

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
            // resortWorkingArray(chosenValuesArr[:i], i)
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

    // pairs.add_sequence(chosenValuesArr)

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


// func resortWorkingArray (chosenValuesArr []interface{}, num int) {
//     for item in workingItemMatrix[num]:
//         data_node = self.__pairs.get_node_info(item)

//         new_combs = [
//             // numbers of new combinations to be created if this item is
//             // appended to array
//             set([key(z) for z in combinations(chosen_values_arr + [item], i + 1)])
//             - self.__pairs.get_combs()[i]
//             for i in range(0, self.__n)
//         ]

//         // weighting the node
//         // node that creates most of new pairs is the best
//         item.weights = [-len(new_combs[-1])]

//         // less used outbound connections most likely to produce more new
//         // pairs while search continues
//         item.weights += (
//             [len(data_node.out)]
//             + [len(x) for x in reversed(new_combs[:-1])]
//             + [-data_node.counter]  # less used node is better
//         )

//         // otherwise we will prefer node with most of free inbound
//         // connections; somehow it works out better ;)
//         item.weights.append(-len(data_node.in_))

//     workingItemMatrix[num].sort(key=cmp_to_key(cmp_item))
// }




// func (pws PairWises) ContainsVectorIndex (pos int) (bool, int) {
//     var ifContains bool
//     ifContains = false
//     var pairWiseIndex int

//     for i, pairWise := range pws {
//         for _, v_ind := range pairWise.PwVectorIndices {
//             if pos == v_ind {
//                 ifContains = true
//                 pairWiseIndex = i
//                 break
//             }
//         }
//         if ifContains == true {
//             break
//         }
//     }
//     return ifContains, pairWiseIndex
// }

func (pw PairWise) ContainsVectorIndex (pos int) bool {
    var ifContains bool
    ifContains = false
    for _, v_ind := range pw.PwVectorIndices {
        if pos == v_ind {
            ifContains = true
            break
        }
    }
    return ifContains
}


// -------------------------------------------------------------
// Note: for pairewise, there are two kinds of algorithm:
// --> algorithm 1. use the pairWises to build the test case data (test case)
// --> algorithm 2. get the full combination first, the filter (remove) them based on the pairWises if repeated
// -------------------------------------------------------------
// Here uses the algorithm 1

func GetPairWiseValid22(fuzzData FuzzData, PwLength int) {
    var combos [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            // fmt.Println("validList: ", key, validList)
            combos = append(combos, validList)
        }
    }
    fmt.Println("combos length", len(combos), "\n")

    // init -----------------
    var indexSlice []int
    for i, _ := range combos {
        indexSlice = append(indexSlice, i)
    }
    // to get the combinations like [1 1][1 2][1 3] ...
    var pwLen int
    if PwLength > len(combos) {
        pwLen = len(combos)
    } else {
        pwLen = PwLength
    }
    indexCombs := combinations(indexSlice, pwLen)
    GetPairWise12(indexCombs, combos, PwLength)
}

func GetPairWise12 (indexCombs chan []int, combos [][]interface{}, PwLength int) {
    var pairWises PairWises
    mapp := make(map[string][]PairWise)

    for indexPair := range indexCombs {
        var pairWise PairWise
        pairWise.PwLength = PwLength
        pairWise.PwVectorIndices = indexPair
        
        var combos_pw_index_slice [][]interface{}

        for _, ind_value := range indexPair {
            var indexSlice []interface{}
            for i, _ := range combos[ind_value] {
                indexSlice = append(indexSlice, i)
            }

            combos_pw_index_slice = append(combos_pw_index_slice, indexSlice)
        }

        // fmt.Println("combos_pairwise_index: ", combos_pw_index_slice, len(combos_pw_index_slice))
        
        c := make(chan []interface{})
        go func(c chan []interface{}) {
            defer close(c)
            combosSliceString(c, []interface{}{}, combos_pw_index_slice)
        }(c)

        fmt.Println("c: ", c, len(c))

        i := 0
        
        for pairwise := range c {
            i = i + 1
            // fmt.Println("results_pairwise: ", pairwise, len(pairwise), pairwise[0], pairwise[1])
            var pairwiseValue []interface{}
            var pwind []int
            for ii, ind_value := range pairwise {
                pwind = append(pwind, ind_value.(int))

                pairwiseValue = append(pairwiseValue, combos[indexPair[ii]][ind_value.(int)])
            }
            pairWise.PwElementIndices = pwind
            pairWise.PwValues = pairwiseValue
            fmt.Println("pairwiseValue length: ", pairwiseValue, i)


            ///
            str := ""
            for _, ind_value := range indexPair {
                str = str + fmt.Sprint(ind_value)
            }
            // fmt.Println("pairWises length: ", fmt.Sprint(indexPair))
            mapp[str] = append(mapp[str], pairWise)
            // pairWises.PairWises[str] = append(pairWises.PairWises[str], pairWise)
            pairWises.PairWises = mapp
        }
    }
    fmt.Println("pairWises length: ", pairWises)
    // return pairWises
}

// func GeneratePairWiseTestData () {
//     // now we have the informations about:
//     // 1. pairwise length (N), the test case length (M)
//     // 2. the pairWises, which is grouped by sub-combinations
//     // 3. the total number of the sub-combinations

// }





// -------------------------------------------------------------
// Note: for pairewise, there are two kinds of algorithm:
// --> algorithm 1. use the pairWises to build the test case data (test case)
// --> algorithm 2. get the full combination first, the filter (remove) them based on the pairWises if repeated
// -------------------------------------------------------------
// Here uses the algorithm 2
// Warning: after testing, this algorithm has performance issue, whil also can not result minimal set of pairwise test data (case)

func GetPairWiseValid(fuzzData FuzzData, PwLength int) {
    var combos [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            // fmt.Println("validList: ", key, validList)
            combos = append(combos, validList)
        }
    }
    fmt.Println("combos length", len(combos), "\n")

    //
    c := make(chan []interface{})
    go func(c chan []interface{}) {
        defer close(c)
        combosSliceString(c, []interface{}{}, combos)
    }(c)

    var combosFullValid[][]interface{}
    for subCombValid := range c {
        combosFullValid = append(combosFullValid, subCombValid) 
    }

    // init
    
    // to get the combinations like [1 1][1 2][1 3] ...
    var pwLen int
    if PwLength > len(combos) {
        pwLen = len(combos)
    } else {
        pwLen = PwLength
    }
    //
    var indexSlice []int
    for i, _ := range combosFullValid[0] {
        indexSlice = append(indexSlice, i)
    }

    indexCombs := combinations(indexSlice, pwLen)
    var indexPW [][]int
    combLen := 0
    for value := range indexCombs {
        indexPW = append(indexPW, value)
        combLen = combLen + 1
    }
    //

    GetPairWisedTestCases(combosFullValid, pwLen, indexPW)
}


func GetPairWisedTestCases (combosFullValidP [][]interface{}, pwLen int, indexPW [][]int) {
    // var pairWiseTestCaseDataSet [][]interface{}

    var combosFullValid [][]interface{}
    combosFullValid = combosFullValidP
    loopDepth := 0

    miniLoop:
    for {
        var resultsCombosFullValid [][]interface{}

        tryIndex := 0
        for i, subCombValid := range combosFullValid {
            // fmt.Println("!------subCombValid: ", len(combosFullValid), i, subCombValid, combosFullValid)
            // Step 1: get the PairWises for subCombValid
            var totalFoundPwCount []int

            co := make(chan int, len(indexPW))
            var wg sync.WaitGroup
            //
            for _, indvalue := range indexPW {
                wg.Add(1)
                go GetFoundPwCount(co, &wg, combosFullValidP, pwLen, indvalue, subCombValid)
            }
            wg.Wait()
            close(co)

            for foundPwCount := range co {
                totalFoundPwCount = append(totalFoundPwCount, foundPwCount) 
            }
            // --> if yes, then remove the subCombValid from combosFullValid
            // --> if no, then keep this subCombValid, and to next subCombValid
            // fmt.Println("totalFoundPwCount: ", totalFoundPwCount, indexPW)
            var ifSubRepeated bool
            ifSubRepeated = true
            for _, value := range totalFoundPwCount {
                if value < 2 {
                    ifSubRepeated = false
                    break
                }
            }
            if ifSubRepeated == true {
                // fmt.Println("len(combosFullValid): ", len(combosFullValid))
                resultsCombosFullValid = RemoveSliceItem(combosFullValid, subCombValid)
                // fmt.Println("len(resultsCombosFullValid): ", len(resultsCombosFullValid))
                combosFullValid = resultsCombosFullValid
                break
                // GetPairWisedTestCases(resultsCombosFullValid, pwLen, indexPW) 
            } else {
                // fmt.Println(" ---> to next ")
                combosFullValid = combosFullValid
            }
            tryIndex = i
        }
        // can not remove anymore
        if len(combosFullValid) - 1 == tryIndex {
            break miniLoop
        }

        fmt.Println("len(combosFullValid)", len(combosFullValid))
        loopDepth = loopDepth + 1
        // if loopDepth == 3000 {
        //     break miniLoop
        // }
    }
    fmt.Println("touch the ending", len(combosFullValidP), loopDepth)
}


func GetFoundPwCount(co chan int, wg *sync.WaitGroup, combosFullValid [][]interface{}, pwLen int,
        indvalue []int, subCombValid []interface{}) {
    //
    var pairWise PairWise
    var pairValues []interface{}

    defer wg.Done()

    pairWise.PwLength = pwLen
    pairWise.PwElementIndices = indvalue
    
    for _, v_i := range indvalue {
        pairValues = append(pairValues, subCombValid[v_i])
    }
    pairWise.PwValues = pairValues
    // fmt.Println("pairWise: ", pairWise)
    // Step 2: check if the PairWises appears in combosFullValid more than once, 
    foundPwCount := CheckSliceItemExistence(combosFullValid, pairWise)

    co <- foundPwCount
}


func CheckSliceItemExistence(combosFullValid [][]interface{}, pairWise PairWise) int {
    foundCount := 0

    for _, subCombValid := range combosFullValid {
        var sourcePairValues []interface{}
        for _, v_i := range pairWise.PwElementIndices {
            sourcePairValues = append(sourcePairValues, subCombValid[v_i])
        }

        if reflect.DeepEqual(sourcePairValues, pairWise.PwValues) {
            foundCount = foundCount + 1

            if foundCount == 2 {
                break
            }
        }
    }

    return foundCount
}


func RemoveSliceItem(sourceSlice [][]interface{}, item []interface{}) [][]interface{} {
    var resultSlice [][]interface{}

    for index, source_item := range sourceSlice {
        if reflect.DeepEqual(source_item, item) {
            // fmt.Println("reflect.DeepEqual(source_item, item): ", reflect.DeepEqual(source_item, item), source_item, item)
            if index == 0 {
                resultSlice = append(resultSlice, sourceSlice[index + 1:]...)
            } else if index == len(sourceSlice) - 1 {
                resultSlice = sourceSlice[:index]
            } else {
                resultSlice = append(sourceSlice[:index], sourceSlice[index + 1:]...)
            }
        }
    }

    return resultSlice
}

func ifHasNilElement (vector []interface{}) (bool, int) {
    var ifNil bool
    ifNil = false
    var pos int

    for i, v := range vector {
        if v == nil {
            ifNil = true
            pos = i
            break
        }
    }
    return ifNil, pos
}


func ifPairWiseHasPosElement (vector []interface{}) (bool, int) {
    var ifNil bool
    ifNil = false
    var pos int

    for i, v := range vector {
        if v == nil {
            ifNil = true
            pos = i
            break
        }
    }
    return ifNil, pos
}




