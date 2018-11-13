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

    "go4api/lib/rands"
)

// fuzz - random - valid
func (fD FieldDefinition) DetermineFuzzValidType() string {
    var fuzzType string
    switch strings.ToLower(fD.FieldType) {
        case "char":
            switch strings.ToLower(fD.FieldSubType) {
                case "numeric":
                    fuzzType = "FCharNumericValid"
                case "alpha":
                    fuzzType = "FCharAlphaValid"
                case "alphanumeric":
                    fuzzType = "FCharAlphaNumericValid"
                case "time":
                    fuzzType = "FCharTimeValid"
                case "email":
                    fuzzType = "FCharEmailValid"
                case "ip":
                    fuzzType = "FCharIpValid"
                default: 
                    fuzzType = "FCharValid"
            }
        case "int":
            switch strings.ToLower(fD.FieldSubType) {
                case "time":
                    fuzzType = "FIntTimeValid"
                default:
                    fuzzType = "FIntValid"
                }
        case "float", "float64":
            fuzzType = "FFloatValid"
        case "bool":
            fuzzType = "FBoolValid"
        case "array":
            fuzzType = "FArrayValid"
        default:
            fmt.Println("!! Error: No specific rules mapping matched")
    }
    
    return fuzzType
}

// fuzz - random - invalid
func (fD FieldDefinition) DetermineFuzzInvalidType() string {
    var fuzzType string
    switch strings.ToLower(fD.FieldType) {
        case "char":
            switch strings.ToLower(fD.FieldSubType) {
                case "numeric":
                    fuzzType = "FCharNumericInvalid"
                case "alpha":
                    fuzzType = "FCharAlphaInvalid"
                case "alphanumeric":
                    fuzzType = "FCharAlphaNumericInvalid"
                case "time":
                    fuzzType = "FCharTimeInvalid"
                case "email":
                    fuzzType = "FCharEmailInvalid"
                case "ip":
                    fuzzType = "FCharIpInvalid"
                default: 
                    fuzzType = "FCharInvalid"
            }
        case "int":
            switch strings.ToLower(fD.FieldSubType) {
                case "time":
                    fuzzType = "FIntTimeInvalid"
                default:
                    fuzzType = "FIntInvalid"
                }
        case "float", "float64":
            fuzzType = "FFloatInvalid"
        case "bool":
            fuzzType = "FBoolInvalid"
        case "array":
            fuzzType = "FArrayInvalid"
        default:
            fmt.Println("!! Error: No specific rules mapping matched")
    }
    
    return fuzzType
}

func (fieldDefinition FieldDefinition) CallFuzzRules(key string) []interface{} {
    var values []interface{}

    for _, ruleFunc := range FuzzRulesMapping(key) {
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

// -------- for the fuzz data based on the field definition -----
// to get the fuzz data table files with naming FCase_fuzz_dt_valid.csv / FCase_fuzz_dt_invalid.csv

// valid fieldMin, rands.RandStringRunes
func FCharValidR1(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    value := rands.RandStringRunes(fieldMin)

    return value
}

func FCharValidR2(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin + 1 < fieldMax
    value := rands.RandStringRunes(fieldMin + 1)

    return value
}

func FCharValidR3(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    value := rands.RandStringRunes(fieldMax)

    return value
}

func FCharValidR4(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := rands.RandStringRunes(fieldMax - 1)

    return value
}

// invalid
func FCharInvalidR1(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := rands.RandStringRunes(fieldMax + 1)

    return value
}


// int ---
func FNumValidR1(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := fieldMin

    return value
}

func FNumValidR2(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := fieldMin + 1

    return value
}

func FNumValidR3(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := fieldMax

    return value
}

func FNumValidR4(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := fieldMax - 1

    return value
}

func FNumValidR5(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := 0

    return value
}

func FNumValidR6(fieldName string, fieldType string, fieldMin int, fieldMax int) interface{} {
    // if fieldMin < fieldMax - 1
    value := -1

    return value
}



