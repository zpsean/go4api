/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package mode

import (
    "fmt"
    "go4api/types"
)


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
        // remove this conditional to return all sets of length <=k
        if(length == 1){
            output := make([]string, len(newCombo))
            copy(output, newCombo)
            c <- output
        }
        combosString(c, newCombo, data, length - 1)
    }
}



func GetCombinationValid(fuzzData types.FuzzData) <-chan []string {
    var combos [][]string
    for _, validDataMap := range fuzzData.ValidData {
        for key, validList := range validDataMap {

            fmt.Println("validList: ", key, validList)

            combos = append(combos, validList)
        }
    }

    fmt.Println("valid comoss: ", combos)

    //
    c := make(chan []string)
    go func(c chan []string) {
        defer close(c)
        combosSliceString(c, []string{}, combos)
    }(c)

    return c
}


func combosSliceString(c chan []string, combo []string, data [][]string) {  
    // Check if we reached the length limit
    // If so, we just return without adding anything
    
    for _, i_v := range data[0] {
        var newCombo []string
        newCombo = append(newCombo, i_v)
        for _, j_v := range data[1] {
            output := make([]string, 1)
            copy(output, newCombo)
            output = append(output, j_v)
            c <- output
        }
    }
}




func GetCombinationInvalid(fuzzData types.FuzzData) {
    // comb 1, for invalid + valid
    for _, invalidDataMap := range fuzzData.InvalidData {
        for key, invalidList := range invalidDataMap {

            fmt.Println("invalidList: ", key, invalidList)
        }
    }

    // comb 1, for invalid + invalid
}






