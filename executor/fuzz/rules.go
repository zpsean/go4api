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
    "fmt"
    "reflect"
    // "path/filepath"
    "strings"
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

//----------------------------------------------
//------- Mutation Rules - Types ---------------
//----------------------------------------------
//------ mutation: if normal char 
// (1) -> empty ("")
// (2) -> blank (" ")
// (3) -> prefix blank (" " + previousValue)
// (4) -> suffix blank (previousValue + " ")
// (5) -> mid blank (previousValue[0:2] + " " + previousValue[2:])
// (6) -> only one char (previousValue[0])
// (7) -> longlong string (strings.Repeat(previousValue, 50)
// (8) -> special char(s) (~!@#$%^&*()_+{}[]<>?)
// (9) -> null
// (10) -> change type (simple, i.e. to float64/bool...)
// (11) -> change type (complex, i.e. object, arrary)
// (12) -> remove this node

//------ mutation: if not normal char -> all number
// (1) -> empty ("")
// (2) -> blank (" ")
// (3) -> not all number (char + previousValue[1:] / previousValue[0:-1] + char)
// (4) -> null
// (5) -> change type (simple, i.e. to float64/bool...)
// (6) -> change type (complex, i.e. object, arrary)
// (7) -> remove this node


//------ mutation: if not normal char -> timestamp/date/time/...
// (1) -> empty ("")
// (2) -> blank (" ")
// (3) -> valid time format (1970/01/01)
// (4) -> valid time format (9999/01/01)
// (5) -> valid time format (1901/01/01)
// (6) -> invalid time format (1901/13/01)
// (7) -> invalid time format ("aaaaaaaa")
// (8) -> null
// (9) -> change type (simple, i.e. to float64/bool...)
// (10) -> change type (complex, i.e. object, arrary)
// (11) -> remove this node


//------ mutation: if not normal char -> email
// (1) -> empty ("")
// (2) -> blank (" ")
// (3) -> valid time format ("xxxxxxxxx" + previousValue)
// (4) -> invalid time format (" " + previousValue / previousValue + " ")
// (5) -> null
// (6) -> change type (simple, i.e. to float64/bool...)
// (7) -> change type (complex, i.e. object, arrary)
// (8) -> remove this node


//------ mutation: if Numeric -> float64
// (1) -> zero (0)
// (2) -> positive (1)
// (3) -> big positve (xxxxx)
// (4) -> negative (-1)
// (5) -> big negative (-xxxx)
// (9) -> null
// (6) -> change type (i.e. to string/float64/bool...)
// (7) -> change type (complex, i.e. object, arrary)

//------ mutation: if bool
// (1) -> true (0)
// (2) -> false (1)
// (3) -> null
// (4) -> change type (i.e. to string/float64/bool...)
// (5) -> change type (complex, i.e. object, arrary)

//------ mutation: if array
// (1) -> empty ([])
// (2) -> one element ([x])
// (3) -> more element (previousValue + [y])
// (4) -> null
// (5) -> change the element type ()
// (6) -> change type (i.e. to string/float64/bool...)
// (7) -> change type (complex, i.e. object, arrary)


//----------------------------------------------
//------- Mutation Rules - Header --------------
//----------------------------------------------
// (1) -> change the existing key/value - based on types
// (2) -> add new key/value
// (3) -> remove key/value (one each time)
// (4) -> remove key/value (all)

//----------------------------------------------
//------- Mutation Rules - QueryString ---------
//----------------------------------------------
// (1) -> change the existing key/value - based on types
// (2) -> add new key/value
// (3) -> remove key/value (one each time)
// (4) -> remove key/value (all)


func RulesMapping(key string) []interface{} {
    //
    RulesMapping := map[string][]interface{} {
        "MutateChar": []interface{}{MutateCharR1, MutateCharR2, MutateCharR3},
        "MutateCharNumeric": []interface{}{MutateCharR1, MutateCharR2, MutateCharR3},
        "MutateCharAlpha": []interface{}{MutateCharR1, MutateCharR2, MutateCharR3},
        "MutateCharAlphaNumeric": []interface{}{MutateCharR1, MutateCharR2, MutateCharR3},
        "MutateNumeric": []interface{}{MutateCharR1, MutateCharR2, MutateCharR3},
        "MutateBool": []interface{}{MutateCharR1, MutateCharR2, MutateCharR3},

        "FuzzCharValid": []interface{}{FuzzCharValidR1, FuzzCharValidR2, FuzzCharValidR3},
        "FuzzCharInvalid": []interface{}{FuzzCharInvalidR1},
        "FuzzCharNumericValid": []interface{}{FuzzCharValidR1, FuzzCharValidR2, FuzzCharValidR3},
        "FuzzCharAlphaValid": []interface{}{FuzzCharValidR1, FuzzCharValidR2, FuzzCharValidR3},
        
    }

    return RulesMapping[key]
}
//
func (mtD MutationDetails) DetermineMutationType() string {
    var mType string
    if mtD.FieldType == "string" {
        mType = "MutateChar"
    }
    
    return mType
}

// fuzz - mutation
func (mtD MutationDetails) CallMutationRules(key string) []interface{} {
    var mutatedValues []interface{}

    for _, ruleFunc := range RulesMapping(key) {
        f := reflect.ValueOf(ruleFunc)
        //
        in := make([]reflect.Value, 3)
        in[0] = reflect.ValueOf(mtD.CurrValue)
        in[1] = reflect.ValueOf(mtD.FieldType)
        in[2] = reflect.ValueOf(mtD.FieldSubType)
        //
        result := f.Call(in)

        // fmt.Println("result =>>>>: ", result[0])
        mutatedValues = append(mutatedValues, result[0].Interface())
    }

    return mutatedValues
}

// fuzz - random
func (fD FieldDefinition) DetermineFuzzType() string {
    var mType string

    switch strings.ToLower(fD.FieldType) {
        case "char":
            switch strings.ToLower(fD.FieldSubType) {
                case "numeric":
                    mType = "FuzzCharNumericValid"
                case "alpha":
                    mType = "FuzzCharAlphaValid"
                default: 
                    mType = "FuzzCharValid"
            }
        case "int":
            mType = "FuzzInt"
        default:
            fmt.Println("!! Error: No specific rules mapping matched")
    }
    
    return mType
}

func (fieldDefinition FieldDefinition) CallFuzzRules(key string) []interface{} {
    var values []interface{}

    for _, ruleFunc := range RulesMapping(key) {
        f := reflect.ValueOf(ruleFunc)
        //
        in := make([]reflect.Value, 4)
        in[0] = reflect.ValueOf(fieldDefinition.FieldName)
        in[1] = reflect.ValueOf(fieldDefinition.FieldType)
        in[2] = reflect.ValueOf(fieldDefinition.FieldMin)
        in[3] = reflect.ValueOf(fieldDefinition.FieldMax)
        //
        result := f.Call(in)

        // fmt.Println("result =>>>>: ", result[0])
        values = append(values, result[0].Interface())
    }

    return values
}

//----------------------------------------------
//------- Below are the rule functions ---------
//----------------------------------------------
// empty
func MutateCharR1(currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := ""

    return mutatedValue
}

// blank
func MutateCharR2(currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := " "

    return mutatedValue
}

// prefix blank
func MutateCharR3(currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := " " + fmt.Sprint(currValue)

    return mutatedValue
}

// suffix blank
func MutateCharR4(currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := fmt.Sprint(currValue) + " "

    return mutatedValue
}

// mid blank
func MutateCharR5(currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := fmt.Sprint(currValue)[0:2] + " " + fmt.Sprint(currValue)[2:]

    return mutatedValue
}




// -------- for the fuzz data based on the field definition -----
// to get the fuzz data table files with naming fuzzcase_fuzz_dt_valid.csv / fuzzcase_fuzz_dt_invalid.csv

// valid fieldMin, RandStringRunes
func FuzzCharValidR1(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    value := RandStringRunes(fieldMin)

    return value
}

func FuzzCharValidR2(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin + 1 < fieldMax
    value := RandStringRunes(fieldMin + 1)

    return value
}

func FuzzCharValidR3(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    value := RandStringRunes(fieldMax)

    return value
}

func FuzzCharValidR4(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := RandStringRunes(fieldMax - 1)

    return value
}

// invalid
func FuzzCharInvalidR1(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := RandStringRunes(fieldMax + 1)

    return value
}


// int ---
func FuzzNumValidR1(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := fieldMin

    return value
}

func FuzzNumValidR2(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := fieldMin + 1

    return value
}

func FuzzNumValidR3(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := fieldMax

    return value
}

func FuzzNumValidR4(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := fieldMax - 1

    return value
}

func FuzzNumValidR5(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := 0

    return value
}

func FuzzNumValidR6(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := -1

    return value
}



