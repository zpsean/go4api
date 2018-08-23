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
type PairWises []PairWise


type PairWise struct {
    PwLength int
    PwVectorIndices []int
    PwElementIndices []int
    PwValues []interface{}
}

// type PairWise struct {
//     PwLength int
//     PwElementIndices []int
//     PwValues []interface{}
// }

func (pws PairWises) ContainsVectorIndex (pos int) (bool, int) {
    var ifContains bool
    ifContains = false
    var pairWiseIndex int

    for i, pairWise := range pws {
        for _, v_ind := range pairWise.PwVectorIndices {
            if pos == v_ind {
                ifContains = true
                pairWiseIndex = i
                break
            }
        }
        if ifContains == true {
            break
        }
    }
    return ifContains, pairWiseIndex
}

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

    for value := range indexCombs {
        var pairWise PairWise
        pairWise.PwLength = PwLength
        pairWise.PwVectorIndices = value
        
        var combos_pw_index_slice [][]interface{}

        for _, ind_value := range value {
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

                pairwiseValue = append(pairwiseValue, combos[value[ii]][ind_value.(int)])
            }
            pairWise.PwElementIndices = pwind
            pairWise.PwValues = pairwiseValue
            fmt.Println("pairwiseValue length: ", pairwiseValue, i)

            pairWises = append(pairWises, pairWise)
        }
    }
    fmt.Println("pairWises length: ", pairWises, len(pairWises))
    // return pairWises
}

func GeneratePairWiseTestData () {

}





// -------------------------------------------------------------
// Note: for pairewise, there are two kinds of algorithm:
// --> algorithm 1. use the pairWises to build the test case data (test case)
// --> algorithm 2. get the full combination first, the filter (remove) them based on the pairWises if repeated
// -------------------------------------------------------------
// Here uses the algorithm 2

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




