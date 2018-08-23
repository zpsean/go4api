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
    // "fmt"
)

// Refer to Python:
// product('ABCD', repeat=2)   AA AB AC AD BA BB BC BD CA CB CC CD DA DB DC DD => cartesian product
// permutations('ABCD', 2)   AB AC AD BA BC BD CA CB CD DA DB DC
// combinations('ABCD', 2)   AB AC AD BC BD CD
// combinations_with_replacement('ABCD', 2)   AA AB AC AD BB BC BD CC CD DD


// func GenerateCombinations(alphabet string, length int) <-chan string {
//     c := make(chan string)

//     // Starting a separate goroutine that will create all the combinations,
//     // feeding them to the channel c
//     go func(c chan string) {
//         defer close(c) // Once the iteration function is finished, we close the channel

//         // This is where the iteration will take place
//         // Your teacher's pseudo code uses recursion
//         // which mean you might want to create a separate function
//         // that can call itself.
//     }(c)

//     return c // Return the channel to the calling function
// }



// combinations([]int{1, 2, 3, 4}, 2) =>
// [1 2]
// [1 3]
// [1 4]
// [2 3]
// [2 4]
// [3 4]
func combinations(list []int, length int) (c chan []int) {
    c = make(chan []int)
    go func() {
        defer close(c)
        switch {
            case length == 0:
                c <- []int{}
            case length == len(list):
                c <- list
            case len(list) < length:
                return
            default:
                for i := 0; i < len(list); i++ {
                    for sub_comb := range combinations(list[i + 1:], length - 1) {
                        c <- append([]int{list[i]}, sub_comb...)
                    }
                }
            }
    }()
    return
}


func GenerateCombinationsString(data []string, length int) <-chan []string {  
    c := make(chan []string)
    go func(c chan []string) {
        defer close(c)
        combosString(c, []string{}, data, length)
    }(c)
    return c
}


func combosString(c chan []string, combo []string, data []string, length int) {  
    // Check if we reached the length limit
    // If so, we just return without adding anything
    if length <= 0 {
        return
    }
    var newCombo []string
    for _, ch := range data {
        newCombo = append(combo, ch)
        // remove this conditional to return all sets of length <= k
        if(length == 1){
            output := make([]string, len(newCombo))
            copy(output, newCombo)
            c <- output
        }
        combosString(c, newCombo, data, length - 1)
    }
}


// GenerateCombinationsInt([]int{1,2,3,4}, 2) =>
// [1 1][1 2][1 3][1 4][2 1][2 2][2 3][2 4][3 1][3 2][3 3][3 4][4 1][4 2][4 3][4 4]
func GenerateCombinationsInt(data []int, length int) <-chan []int {  
    c := make(chan []int)
    go func(c chan []int) {
        defer close(c)
        combosInt(c, []int{}, data, length)
    }(c)
    return c
}


func combosInt(c chan []int, combo []int, data []int, length int) {  
    // Check if we reached the length limit
    // If so, we just return without adding anything
    if length <= 0 {
        return
    }
    var newCombo []int
    for _, ch := range data {
        newCombo = append(combo, ch)
        // remove this conditional to return all sets of length <=k
        if(length == 1){
            output := make([]int, len(newCombo))
            copy(output, newCombo)
            c <- output
        }
        combosInt(c, newCombo, data, length - 1)
    }
}

///
func GetCombinationValid(fuzzData FuzzData) <-chan []interface{} {
    var combos [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            // fmt.Println("validList: ", key, validList)
            combos = append(combos, validList)
        }
    }
    //
    GetPairWiseValid22(fuzzData, 2)

    //
    c := make(chan []interface{})

    // combosSliceString(c1, []interface{}{}, combos)

    go func(c chan []interface{}) {
        defer close(c)
        combosSliceString(c, []interface{}{}, combos)
    }(c)

    return c
}


// -- for the fuzz data
func GetCombinationInvalid(fuzzData FuzzData) <-chan []interface{} {
    var validCombos [][]interface{}
    for _, validDataMap := range fuzzData.ValidData {
        for _, validList := range validDataMap {
            // fmt.Println("validList: ", key, validList)
            validCombos = append(validCombos, validList)
        }
    }

    var invalidCombos [][]interface{}
    for _, invalidDataMap := range fuzzData.InvalidData {
        for _, invalidList := range invalidDataMap {
            // fmt.Println("invalidList: ", key, invalidList)
            invalidCombos = append(invalidCombos, invalidList)
        }
    }

    // fmt.Println("invalid combos: ", invalidCombos)

    //
    c := make(chan []interface{})

    // comb type 1, for invalid + valid mix
    for i, invalid := range invalidCombos {
        var combos [][]interface{}

        if i == 0 {
            combos = append(combos, invalid)
            combos = append(combos, validCombos[i + 1: ]...)
        } else if i < len(invalidCombos) - 1 {
            combos = append(combos, validCombos[:i]...)
            combos = append(combos, invalid)
            combos = append(combos, validCombos[i + 1: ]...)
        } else {
            combos = append(combos, validCombos[:i]...)
            combos = append(combos, invalid)
        }
        // fmt.Println("invalid combos - 2: ", combos)

        // go combosSliceString(c2, []interface{}{}, combos)
    }

    // comb type 2, for invalid + invalid    
    go func(c chan []interface{}) {
        defer close(c)
        combosSliceString(c, []interface{}{}, invalidCombos)
    }(c)

    // defer close(c)

    return c
}


func combosSliceString(c chan []interface{}, combo []interface{}, data [][]interface{}) {  
    if len(data) > 1 {
        var newCombo []interface{}
        for _, i_v := range data[0] {
            newCombo = append(combo, i_v)

            combosSliceString(c, newCombo, data[1:])
        }

    } else if len(data) == 1 {
        for _, j_v := range data[0] {
            output := make([]interface{}, len(combo))
            copy(output, combo)

            output = append(output, j_v)
            // fmt.Println("output: ", output)
            c <- output
        }
    }
}

