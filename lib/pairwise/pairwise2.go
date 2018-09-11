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
    "reflect"
    "sync"
)

// -------------------------------------------------------------------------
// Note: for pairewise, there are two kinds of algorithm:
// --> algorithm 1. use the pairWises to build the test case data (test case)
// --> algorithm 2. get the full combination first, the filter (remove) them based on the pairWises if repeated
// -------------------------------------------------------------------------
// Here uses the algorithm 2
// Warning: after testing, this algorithm has performance issue, whil also can not result minimal set of pairwise test data (case)

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

// -------------------------------------------------------------

func GetPairWise2(validVectors [][]interface{}, pwLength int) {
    //
    c := make(chan []interface{})
    go func(c chan []interface{}) {
        defer close(c)
        combinsSliceInterface(c, []interface{}{}, validVectors)
    }(c)

    var combinsFullValid[][]interface{}
    for subCombValid := range c {
        combinsFullValid = append(combinsFullValid, subCombValid) 
    }

    // init
    
    // to get the combinations like [1 1][1 2][1 3] ...
    var pwLen int
    if pwLength > len(validVectors) {
        pwLen = len(validVectors)
    } else {
        pwLen = pwLength
    }
    //
    var indexSlice []int
    for i, _ := range combinsFullValid[0] {
        indexSlice = append(indexSlice, i)
    }

    indexCombs := combinationsInt(indexSlice, pwLen)
    var indexPW [][]int
    combLen := 0
    for value := range indexCombs {
        indexPW = append(indexPW, value)
        combLen = combLen + 1
    }
    //

    GetPairWisedTestCases(combinsFullValid, pwLen, indexPW)
}


func GetPairWisedTestCases (combinsFullValidP [][]interface{}, pwLen int, indexPW [][]int) {
    // var pairWiseTestCaseDataSet [][]interface{}

    var combinsFullValid [][]interface{}
    combinsFullValid = combinsFullValidP
    loopDepth := 0

    miniLoop:
    for {
        var resultsCombosFullValid [][]interface{}

        tryIndex := 0
        for i, subCombValid := range combinsFullValid {
            // fmt.Println("!------subCombValid: ", len(combinsFullValid), i, subCombValid, combinsFullValid)
            // Step 1: get the PairWises for subCombValid
            var totalFoundPwCount []int

            co := make(chan int, len(indexPW))
            var wg sync.WaitGroup
            //
            for _, indvalue := range indexPW {
                wg.Add(1)
                go GetFoundPwCount(co, &wg, combinsFullValidP, pwLen, indvalue, subCombValid)
            }
            wg.Wait()
            close(co)

            for foundPwCount := range co {
                totalFoundPwCount = append(totalFoundPwCount, foundPwCount) 
            }
            // --> if yes, then remove the subCombValid from combinsFullValid
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
                // fmt.Println("len(combinsFullValid): ", len(combinsFullValid))
                resultsCombosFullValid = RemoveSliceItem(combinsFullValid, subCombValid)
                // fmt.Println("len(resultsCombosFullValid): ", len(resultsCombosFullValid))
                combinsFullValid = resultsCombosFullValid
                break
                // GetPairWisedTestCases(resultsCombosFullValid, pwLen, indexPW) 
            } else {
                // fmt.Println(" ---> to next ")
                combinsFullValid = combinsFullValid
            }
            tryIndex = i
        }
        // can not remove anymore
        if len(combinsFullValid) - 1 == tryIndex {
            break miniLoop
        }

        fmt.Println("len(combinsFullValid)", len(combinsFullValid))
        loopDepth = loopDepth + 1
        // if loopDepth == 3000 {
        //     break miniLoop
        // }
    }
    fmt.Println("touch the ending", len(combinsFullValidP), loopDepth)
}


func GetFoundPwCount(co chan int, wg *sync.WaitGroup, combinsFullValid [][]interface{}, pwLen int,
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
    // Step 2: check if the PairWises appears in combinsFullValid more than once, 
    foundPwCount := CheckSliceItemExistence(combinsFullValid, pairWise)

    co <- foundPwCount
}


func CheckSliceItemExistence(combinsFullValid [][]interface{}, pairWise PairWise) int {
    foundCount := 0

    for _, subCombValid := range combinsFullValid {
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


func GetPairWise22(validVectors [][]interface{}, pwLength int) {
    // init -----------------
    var indexSlice []int
    for i, _ := range validVectors {
        indexSlice = append(indexSlice, i)
    }
    // to get the combinations like [1 1][1 2][1 3] ...
    var pwLen int
    if pwLength > len(validVectors) {
        pwLen = len(validVectors)
    } else {
        pwLen = pwLength
    }
    indexCombs := combinationsInt(indexSlice, pwLen)
    GetPairWise12(indexCombs, validVectors, pwLength)
}

func GetPairWise12 (indexCombs chan []int, combins [][]interface{}, pwLength int) {
    var pairWises PairWises
    mapp := make(map[string][]PairWise)

    for indexPair := range indexCombs {
        var pairWise PairWise
        pairWise.PwLength = pwLength
        pairWise.PwVectorIndices = indexPair
        
        var combins_pw_index_slice [][]interface{}

        for _, ind_value := range indexPair {
            var indexSlice []interface{}
            for i, _ := range combins[ind_value] {
                indexSlice = append(indexSlice, i)
            }

            combins_pw_index_slice = append(combins_pw_index_slice, indexSlice)
        }

        // fmt.Println("combins_pairwise_index: ", combins_pw_index_slice, len(combins_pw_index_slice))
        
        c := make(chan []interface{})
        go func(c chan []interface{}) {
            defer close(c)
            combinsSliceInterface(c, []interface{}{}, combins_pw_index_slice)
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

                pairwiseValue = append(pairwiseValue, combins[indexPair[ii]][ind_value.(int)])
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


