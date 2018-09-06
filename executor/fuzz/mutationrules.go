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
    "math"
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

//
func (mtD MutationDetails) DetermineMutationType() string {
    var mType string
    switch strings.ToLower(mtD.FieldType) {
        case "char", "string":
            switch strings.ToLower(mtD.FieldSubType) {
                case "numeric":
                    mType = "MCharNumeric"
                case "alpha":
                    mType = "MCharAlpha"
                case "alphanumeric":
                    mType = "MCharAlphaNumeric"
                case "time":
                    mType = "MCharTime"
                case "email":
                    mType = "MCharEmail"
                case "ip":
                    mType = "MCharIp"
                default: 
                    mType = "MChar"
            }
        case "int":
            switch strings.ToLower(mtD.FieldSubType) {
                case "time":
                    mType = "MIntTime"
                default:
                    mType = "MInt"
                }
        case "float":
            mType = "MFloat"
        case "bool":
            mType = "MBool"
        case "array":
            mType = "MArray"
        default:
            fmt.Println("!! Error: No specific rules mapping matched: ", strings.ToLower(mtD.FieldType))
    }
    
    return mType
}


// fuzz - mutation
func (mtD MutationDetails) CallMutationRules(key string) []interface{} {
    var mutatedValues []interface{}

    for _, ruleFunc := range MutationRulesMapping(key) {
        f := reflect.ValueOf(ruleFunc)
        //
        in := make([]reflect.Value, 3)
        in[0] = reflect.ValueOf(mtD.CurrValue)
        in[1] = reflect.ValueOf(mtD.FieldType)
        in[2] = reflect.ValueOf(mtD.FieldSubType)
        //
        result := f.Call(in)

        mutatedValues = append(mutatedValues, result[0].Interface())
    }

    return mutatedValues
}


//----------------------------------------------
//------- Below are the rule functions ---------
//----------------------------------------------
// empty
func MCharR1 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := ""
    return mutatedValue
}

// blank
func MCharR2 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := " "
    return mutatedValue
}

// prefix blank
func MCharR3 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := " " + fmt.Sprint(currValue)
    return mutatedValue
}

func MCharR4 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = " " + string(currValueRune[1:])
    }
    return mutatedValue
}

// suffix blank
func MCharR5 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := fmt.Sprint(currValue) + " "
    return mutatedValue
}

func MCharR6 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = string(currValueRune[0:len(currValueRune) - 1]) + " "
    }

    return mutatedValue
}

// mid blank
func MCharR7 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 2 {
        mutatedValue = string(currValueRune[0:1]) + " " + string(currValueRune[1:len(currValueRune)])
    }

    return mutatedValue
}

func MCharR8 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 2 {
        mutatedValue = string(currValueRune[0:1]) + " " + string(currValueRune[2:len(currValueRune)])
    }

    return mutatedValue
}

// only one char
func MCharR9 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 0 {
        mutatedValue = string(currValueRune[0])
    }

    return mutatedValue
}

// longlong string
func MCharR10 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := strings.Repeat(currValue.(string), 50)

    return mutatedValue
}

// special char(s) (~!@#$%^&*()_+{}[]<>?)
func MCharR11 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = "%" + string(currValueRune[1:])
    }
    return mutatedValue
}

func MCharR12 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = "*" + string(currValueRune[1:])
    }
    return mutatedValue
}

// null
func MCharR13 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}

    return mutatedValue
}

// change type (simple, i.e. to float64/bool...)
func MCharR14 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 123

    return mutatedValue
}

func MCharR15 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := true

    return mutatedValue
}

//  change type (complex, i.e. object, arrary)
func MCharR16 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    mutatedValue := []string{string(currValueRune[0])}

    return mutatedValue
}


// < --------------- Int ------------->
func MIntR1 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 0
    return mutatedValue
}

func MIntR2 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 1
    return mutatedValue
}

func MIntR3 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := -1
    return mutatedValue
}

func MIntR4 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MaxInt32
    return mutatedValue
}

func MIntR5 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MinInt32
    return mutatedValue
}

func MIntR6 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MaxInt64
    return mutatedValue
}

func MIntR7 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MinInt64
    return mutatedValue
}



