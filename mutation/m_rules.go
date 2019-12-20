/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package mutation

import (                                                                                                                                             
    // "os"
    "time"
    "fmt"
    "math"
    "reflect"
    "runtime"
    // "path/filepath"
    "strings"
    "math/rand"
    // "strconv"
    // "go4api/utils"  

    // "go4api/lib/testcase"
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


var (
    nilNull *int
)

//
func (mFd *MFieldDetails) DetermineMutationType() {
    var mType string
    switch strings.ToLower(mFd.FieldType) {
        case "char", "string":
            switch strings.ToLower(mFd.FieldSubType) {
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
            switch strings.ToLower(mFd.FieldSubType) {
                case "time":
                    mType = "MIntTime"
                default:
                    mType = "MInt"
                }
        case "float", "float64":
            mType = "MFloat"
        case "bool":
            mType = "MBool"
        case "array", "slice":
            mType = "MArray"
        case "map":
            mType = "MMap"
        default:
            fmt.Println("!! Error: No specific rules mapping matched: ", strings.ToLower(mFd.FieldType))
    }

    mFd.MutationType = mType
}


// fuzz - mutation
func (mFd *MFieldDetails) CallMutationRules () {
    // set mFd.MutationType
    mFd.DetermineMutationType()
    mType := mFd.MutationType

    //
    var mutatedValues []*MutatedValue
    for _, ruleFunc := range MutationRulesMapping(mType) {

        f := reflect.ValueOf(ruleFunc)
        //
        in := make([]reflect.Value, 3)
        in[0] = reflect.ValueOf(mFd.CurrValue)
        in[1] = reflect.ValueOf(mFd.FieldType)
        in[2] = reflect.ValueOf(mFd.FieldSubType)
        //
        result := f.Call(in)

        mutatedValue := MutatedValue {
            MutationRule: runtime.FuncForPC(f.Pointer()).Name(),
            MutatedValue: result[0].Interface(),
        }

        mutatedValues = append(mutatedValues, &mutatedValue)
    }

    mFd.MutatedValues = mutatedValues
}


//----------------------------------------------
//------- Common Rules ---------
//----------------------------------------------
// empty
func M_Common_Set_To_Empty (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := ""
    return mutatedValue
}

// blank
func M_Common_Set_To_One_Blank (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := " "
    return mutatedValue
}

func M_Common_Set_To_Null (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = nilNull

    return mutatedValue
}

// set the [] to bool - true
func M_Common_Set_To_Bool_True (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = true

    return mutatedValue
}

// set the [] to bool - false
func M_Common_Set_To_Bool_False (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = false

    return mutatedValue
}

// set to !
func M_Common_Set_To_Single_Exclamation_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `!`

    return mutatedValue
}

// set to @
func M_Common_Set_To_Single_At_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `@`

    return mutatedValue
}

// set to #
func M_Common_Set_To_Single_Number_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `#`

    return mutatedValue
}

// set to $
func M_Common_Set_To_Single_Dollar_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `$`

    return mutatedValue
}

// set to %
func M_Common_Set_To_Single_Percent_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `%`

    return mutatedValue
}

// set to ^
func M_Common_Set_To_Single_Caret_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `^`

    return mutatedValue
}

// set to &
func M_Common_Set_To_Single_Ampersand_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `&`

    return mutatedValue
}

// set to *
func M_Common_Set_To_Single_Asterisk_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `*`

    return mutatedValue
}

// set to .
func M_Common_Set_To_Single_Dot_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `.`

    return mutatedValue
}

// set to +
func M_Common_Set_To_Single_Plus_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `+`

    return mutatedValue
}

// set to -
func M_Common_Set_To_Single_Minus_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `-`

    return mutatedValue
}

// set to e/E, this is for MCharNumeric, which can accept only Numeric, but `e, E, -, +` sign sometimes are negeative test
func M_Common_Set_To_Single_e (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `e`

    return mutatedValue
}

// set to e/E, this is for MCharNumeric, which can accept only Numeric, but `e, E, -, +` sign sometimes are negeative test
func M_Common_Set_To_Single_E (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = `E`

    return mutatedValue
}


//----------------------------------------------
//------- other Rules
//----------------------------------------------
// prefix blank
func M_Char_Add_Prefix_One_Blank (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := " " + fmt.Sprint(currValue)
    return mutatedValue
}

func M_Char_Replace_All_Blank (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    for range currValueRune {
        mutatedValue = mutatedValue + " "
    }
    return mutatedValue
}

func M_Char_Replace_Prefix_One_Blank (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = " " + string(currValueRune[1:])
    }
    return mutatedValue
}

func M_Char_Replace_Prefix_None_ASCII (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = `中` + string(currValueRune[1:])
    }
    return mutatedValue
}

// suffix blank
func M_Char_Add_Suffix_One_Blank (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := fmt.Sprint(currValue) + " "
    return mutatedValue
}

func M_Char_Replace_Suffix_One_Blank (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = string(currValueRune[0:len(currValueRune) - 1]) + " "
    }

    return mutatedValue
}

// mid blank
func M_Char_Add_Mid_One_Blank (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) >= 2 {
        mutatedValue = string(currValueRune[0:1]) + " " + string(currValueRune[1:len(currValueRune)])
    }

    return mutatedValue
}

// mid none-ascii
func M_Char_Add_Mid_None_ASCII (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) >= 2 {
        mutatedValue = string(currValueRune[0:1]) + "中" + string(currValueRune[1:len(currValueRune)])
    }

    return mutatedValue
}

func M_Char_Replace_Mid_One_Blank (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 2 {
        mutatedValue = string(currValueRune[0:1]) + " " + string(currValueRune[2:len(currValueRune)])
    }

    return mutatedValue
}

func M_Char_Replace_Mid_One_E (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 2 {
        mutatedValue = string(currValueRune[0:1]) + "E" + string(currValueRune[2:len(currValueRune)])
    }

    return mutatedValue
}

func M_Char_Replace_Mid_One_Negative_Sign (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 2 {
        mutatedValue = string(currValueRune[0:1]) + "-" + string(currValueRune[2:len(currValueRune)])
    }

    return mutatedValue
}

func M_Char_Replace_Mid_None_ASCII (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 2 {
        mutatedValue = string(currValueRune[0:1]) + "中" + string(currValueRune[2:len(currValueRune)])
    }

    return mutatedValue
}

// only one char
func M_Char_Set_To_One_Char (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 0 {
        mutatedValue = string(currValueRune[0])
    }

    return mutatedValue
}

// longlong string
func M_Char_Repeat_50_Times (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := strings.Repeat(currValue.(string), 50)

    return mutatedValue
}

// special char(s) (~!@#$%^&*()_+{}[]<>?)
func M_Char_Replace_Prefix_Percentage (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = `%` + string(currValueRune[1:])
    }
    return mutatedValue
}

func M_Char_Replace_Prefix_Point (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = `.` + string(currValueRune[1:])
    }
    return mutatedValue
}

func M_Char_Replace_Prefix_Caret (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = `^` + string(currValueRune[1:])
    }
    return mutatedValue
}

func M_Char_Replace_Prefix_Dollar (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = `$` + string(currValueRune[1:])
    }
    return mutatedValue
}

func M_Char_Replace_Prefix_Star (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))
    var mutatedValue string

    if len(currValueRune) > 1 {
        mutatedValue = "*" + string(currValueRune[1:])
    }
    return mutatedValue
}


// change type (simple, i.e. to float64/bool...)
func M_Char_Set_To_Int (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 123

    return mutatedValue
}

func M_Char_Set_To_Float (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 123.12

    return mutatedValue
}


//  change type (complex, i.e. object, arrary)
func M_Char_Set_To_Array (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    currValueRune := []rune(currValue.(string))

    var mutatedValue []string

    if len(currValueRune) > 0 {
        mutatedValue = []string{string(currValueRune[0])}
    }

    return mutatedValue
}


// < --------------- Int ------------->
func M_Int_Set_To_Zero (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 0
    return mutatedValue
}

func M_Int_Set_To_One (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 1
    return mutatedValue
}

func M_Int_Set_To_Negative_One (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := -1
    return mutatedValue
}

func M_Int_Set_To_MaxInt32 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MaxInt32
    return mutatedValue
}

func M_Int_Set_To_MinInt32 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MinInt32
    return mutatedValue
}

func M_Int_Set_To_MaxInt64 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MaxInt64
    return mutatedValue
}

func M_Int_Set_To_MinInt64 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MinInt64
    return mutatedValue
}


// < --------------- Float ------------->
func M_Float_Set_To_E (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.E
    return mutatedValue
}

func M_Float_Set_To_Postive_Float (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 1.1
    return mutatedValue
}

func M_Float_Set_To_Negative_Float (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := -1.1
    return mutatedValue
}

func M_Float_Set_To_MaxFloat32 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MaxFloat32
    return mutatedValue
}

func M_Float_Set_To_SmallestNonzeroFloat32 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.SmallestNonzeroFloat32
    return mutatedValue
}

func M_Float_Set_To_MaxFloat64 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.MaxFloat64
    return mutatedValue
}

func M_Float_Set_To_SmallestNonzeroFloat64 (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := math.SmallestNonzeroFloat64
    return mutatedValue
}


// < --------------- Bool ------------->
func M_Bool_Set_To_Zero (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    mutatedValue := 0
    return mutatedValue
}

// < --------------- Array ------------->
// blank []
func M_Array_Set_To_Empty_Array (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}

    return mutatedValue
}

// remove one (random)
func M_Array_Remove_One_Item_Random (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    length := len(reflect.ValueOf(currValue).Interface().([]interface{}))
    if length > 0 {
        rand.Seed(time.Now().UnixNano())
        randInt := rand.Intn(length)

        for i, v := range reflect.ValueOf(currValue).Interface().([]interface{}) {
            if i != randInt {
                mutatedValue = append(mutatedValue, v)
            }
        }
    } else {
        mutatedValue = append(mutatedValue, currValue)
    }
    return mutatedValue
}

// just keep one item (random)
func M_Array_Set_Only_One_Item (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    length := len(reflect.ValueOf(currValue).Interface().([]interface{}))
    if length > 0 {
        rand.Seed(time.Now().UnixNano())
        randInt := rand.Intn(length)

        for i, v := range reflect.ValueOf(currValue).Interface().([]interface{}) {
            if i == randInt {
                mutatedValue = append(mutatedValue, v)
                break
            }
        }
    } else {
        mutatedValue = append(mutatedValue, currValue)
    }
    return mutatedValue
}

// repeat one item (random)
func M_Array_Duplicate_One_Item_Random (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    length := len(reflect.ValueOf(currValue).Interface().([]interface{}))
    if length > 0 {
        rand.Seed(time.Now().UnixNano())
        randInt := rand.Intn(length)

        for i, v := range reflect.ValueOf(currValue).Interface().([]interface{}) {
            if i == randInt {
                mutatedValue = append(mutatedValue, v)
            }
            mutatedValue = append(mutatedValue, v)
        }
    }
    return mutatedValue
}

// append one but another type, if int, then append string item, vice verse
func M_Array_Append_Another_Type_Item (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    length := len(reflect.ValueOf(currValue).Interface().([]interface{}))
    if length > 0 {
        // currValueSlice := reflect.ValueOf(currValue).Interface().([]interface{})
        for _, v := range reflect.ValueOf(currValue).Interface().([]interface{}) {
            mutatedValue = append(mutatedValue, v)
        }
 
        switch reflect.TypeOf(mutatedValue[0]).Kind() {
            case reflect.Int:
                mutatedValue = append(mutatedValue, "123")
            case reflect.String:
                mutatedValue = append(mutatedValue, 124)
            default:
                mutatedValue = append(mutatedValue, "125")
        }
    }
    return mutatedValue
}

// replace one (random) but another type, if int, then append string item, vice verse
func M_Array_Replace_Another_Type_Item (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    length := len(reflect.ValueOf(currValue).Interface().([]interface{}))
    if length > 0 {
        for _, v := range reflect.ValueOf(currValue).Interface().([]interface{}) {
            mutatedValue = append(mutatedValue, v)
        }
 
        rand.Seed(time.Now().UnixNano())
        randInt := rand.Intn(length)

        switch reflect.TypeOf(mutatedValue[0]).Kind() {
            case reflect.Int:
                mutatedValue[randInt] = "123"
            case reflect.String:
                mutatedValue[randInt] = 124
            default:
                mutatedValue[randInt] = "125"
        }
    }
    return mutatedValue
}

// replace one (random) to nil (i.e. json null)
func M_Array_Replace_One_Item_Null (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    length := len(reflect.ValueOf(currValue).Interface().([]interface{}))
    if length > 0 {
        for _, v := range reflect.ValueOf(currValue).Interface().([]interface{}) {
            mutatedValue = append(mutatedValue, v)
        }
 
        rand.Seed(time.Now().UnixNano())
        randInt := rand.Intn(length)

        mutatedValue[randInt] = nilNull
    }
    return mutatedValue
}

// replace one (random) to bool - true
func M_Array_Replace_One_Item_Bool_True (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    length := len(reflect.ValueOf(currValue).Interface().([]interface{}))
    if length > 0 {
        for _, v := range reflect.ValueOf(currValue).Interface().([]interface{}) {
            mutatedValue = append(mutatedValue, v)
        }
 
        rand.Seed(time.Now().UnixNano())
        randInt := rand.Intn(length)

        mutatedValue[randInt] = true
    }
    return mutatedValue
}

// replace one (random) to bool - false
func M_Array_Replace_One_Item_Bool_False (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    length := len(reflect.ValueOf(currValue).Interface().([]interface{}))
    if length > 0 {
        for _, v := range reflect.ValueOf(currValue).Interface().([]interface{}) {
            mutatedValue = append(mutatedValue, v)
        }
 
        rand.Seed(time.Now().UnixNano())
        randInt := rand.Intn(length)

        mutatedValue[randInt] = false
    }
    return mutatedValue
}

// set the [] has only item nil (i.e. json null)
func M_Array_Set_To_Only_One_Null (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    mutatedValue = append(mutatedValue, nilNull)

    return mutatedValue
}

// set the [] has only item int
func M_Array_Set_To_Only_One_Int (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    mutatedValue = append(mutatedValue, 123)

    return mutatedValue
}

// set the [] has only item string
func M_Array_Set_To_Only_One_String (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    mutatedValue = append(mutatedValue, "123")

    return mutatedValue
}

// set the [] has only item bool - true
func M_Array_Set_To_Only_One_Bool_True (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    mutatedValue = append(mutatedValue, true)

    return mutatedValue
}

// set the [] has only item bool - false
func M_Array_Set_To_Only_One_Bool_False (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue []interface{}
    mutatedValue = append(mutatedValue, false)

    return mutatedValue
}


// set the [] to int
func M_Array_Set_To_Int (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = 123

    return mutatedValue
}

// set the [] to string
func M_Array_Set_To_String (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue interface{}
    mutatedValue = "123"

    return mutatedValue
}

// < --------------- Map ------------->
// blank {}
func M_Map_Set_To_Empty_Map (currValue interface{}, fieldType string, fieldSubType string) interface{} {
    var mutatedValue = make(map[string]interface{})

    return mutatedValue
}

