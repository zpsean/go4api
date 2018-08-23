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
    // "os"
    // "time"
    // "fmt"
    // "path/filepath"
    // "strings"
    // "strconv"
    // "go4api/utils"  
)

// JSON Schema defines the following basic types:
// string
// Numeric -> float64
// boolean
// null
// object (raw)
// array (raw)

// to get the fuzz data table files with naming fuzzcase_fuzz_dt_valid.csv / fuzzcase_fuzz_dt_invalid.csv
func getChar(fieldName string, fieldType string, fieldMin int, fieldMax int) (map[string][]interface{}, map[string][]interface{}) {
    validValueMap := make(map[string][]interface{})
    invalidValueMap := make(map[string][]interface{})
    // get the Boundary (valid, invalid), Equivalence, etc.
    var validLenList []int
    var invalidLenList []int
    //
    validLenList = append(validLenList, fieldMin)
    validLenList = append(validLenList, fieldMin + 1)

    validLenList = append(validLenList, fieldMax)
    if fieldMax - 1 > fieldMin {
        validLenList = append(validLenList, fieldMax - 1)
    }
    //
    if fieldMin - 1 > 0 {
       invalidLenList = append(invalidLenList, fieldMin - 1) 
    }
    invalidLenList = append(invalidLenList, fieldMax + 1) 
    //

    fieldRands := []string{"RandStringRunes"} //, "RandStringCNRunes"}
    //
    for _, validLen := range validLenList {
        for _, randType := range fieldRands {
            // CallRands(randType, validLen)
            validValue := CallRands(randType, validLen)
            // fmt.Println("validLen, validValue: ", validLen, validValue)

            validValueMap[fieldName] = append(validValueMap[fieldName], validValue)
            // validValueMap[fieldName] = append(validValueMap[fieldName], fieldName + "a" + strconv.Itoa(i)) 
        }        
    }
    //
    for _, invalidLen := range invalidLenList {
        for _, randType := range fieldRands {
            invalidValue := CallRands(randType, invalidLen)
            // fmt.Println("invalidLen, invalidValue: ", invalidLen, invalidValue)

            invalidValueMap[fieldName] = append(invalidValueMap[fieldName], invalidValue)
        }
    }

    return validValueMap, invalidValueMap
}


func getInt(fieldName string, fieldType string, fieldMin int, fieldMax int) (map[string][]interface{}, map[string][]interface{}) {
    validValueMap := make(map[string][]interface{})
    invalidValueMap := make(map[string][]interface{})
    // get the Boundary (valid, invalid), Equivalence, etc.
    validValueMap[fieldName] = append(validValueMap[fieldName], fieldMin)
    validValueMap[fieldName] = append(validValueMap[fieldName], fieldMin + 1)
    validValueMap[fieldName] = append(validValueMap[fieldName], fieldMax)
    validValueMap[fieldName] = append(validValueMap[fieldName], fieldMax - 1)
    //
    invalidValueMap[fieldName] = append(invalidValueMap[fieldName], fieldMin - 1)
    invalidValueMap[fieldName] = append(invalidValueMap[fieldName], fieldMax + 1)
    //

    return validValueMap, invalidValueMap
}

