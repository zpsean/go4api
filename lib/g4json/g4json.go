/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package g4json

import (
    // "os"
    "fmt"
    "sync"
    "strings"
    "reflect"
    // "encoding/csv"
    "encoding/json"
)

type FieldDetails struct {
    FieldPath []string
    CurrValue interface{}
    FieldType string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
}

func GetFieldsDetails(valueSource interface{}) []FieldDetails {
    var fieldDetailsSlice []FieldDetails
    c := make(chan FieldDetails)

    go func(c chan FieldDetails) {
        defer close(c)
        wg := &sync.WaitGroup{}

        wg.Add(1)
        TraverseFields(c, []string{}, valueSource, wg)

        wg.Wait()
    }(c)

    for fieldDetails := range c {
        fieldDetailsSlice = append(fieldDetailsSlice, fieldDetails)
    }
    
    return fieldDetailsSlice
}

// encoding/json: func Unmarshal(data []byte, v interface{}) error
// Bool                   => JSON bool
// float64                => JSON numbers
// string                 => JSON strings
// []interface{}          => JSON array
// map[string]interface{} => JSON object
// nil                    => JSON null

// Note: some issues exposed:
// 1. all Json numbers's type is marked as float64, can not distinguish int with float64
// 2. if Json number is 1234.00, then the value in Go is 1234 once json.Unmarshal, lost the .00
// 3. as Json mumbers are treated as float64, it may result in some issue, like use sci a.xxxe+1yy to represent timestamp
func TraverseFields (c chan FieldDetails, subPath []string, value interface{}, wg *sync.WaitGroup) {
    defer wg.Done()

    switch value.(type) {
        case nil:
            wg.Add(1)
            go fieldNull(c, subPath, value, "", wg)
        case float64:
            wg.Add(1)
            fieldPrimitive(c, subPath, value, "", wg)
        case bool:
            wg.Add(1)
            fieldPrimitive(c, subPath, value, "", wg)
        case string:
            ss := strings.TrimLeft(value.(string), "\n")
            ss = strings.TrimSpace(ss)

            if ss[0:1] == "{" || ss[1:2] == "{" {
                pMap := make(map[string]interface{})
                err := json.Unmarshal([]byte(value.(string)), &pMap)
                if err != nil {
                    panic(err)
                }

                wg.Add(1)
                go fieldMap(c, subPath, pMap, "", wg)
            } else if ss[0:1] == "[" || ss[1:2] == "[" {
                var pSlice []interface{}
                err := json.Unmarshal([]byte(value.(string)), &pSlice)
                if err != nil {
                    panic(err)
                }

                wg.Add(1)
                go fieldSlice(c, subPath, pSlice, "", wg)
            } else {
                wg.Add(1)
                fieldPrimitive(c, subPath, value, "", wg)
            }
        case map[string]interface{}:
            wg.Add(1)
            go fieldMap(c, subPath, value, "", wg)
        case []interface{}:
            wg.Add(1)
            go fieldSlice(c, subPath, value, "", wg)
        default:
            fmt.Println("!! Warning, unknown type to traverse.")
    }
}

func fieldNull (c chan FieldDetails, subPath []string, value interface{}, key interface{}, wg *sync.WaitGroup) {
    defer wg.Done()
    subPathNew := make([]string, len(subPath))
    if key == "" {
        copy(subPathNew, subPath)
    } else {
        copy(subPathNew, subPath)
        subPathNew = append(subPathNew, fmt.Sprint(key))
    }
 
    mtD := FieldDetails{subPathNew, value, "", ""}
    c <- mtD
}

func fieldPrimitive (c chan FieldDetails, subPath []string, value interface{}, key interface{}, wg *sync.WaitGroup) {
    defer wg.Done()
    subPathNew := make([]string, len(subPath))
    if key == "" {
        copy(subPathNew, subPath)
    } else {
        copy(subPathNew, subPath)
        subPathNew = append(subPathNew, fmt.Sprint(key))
    }

    mtD := FieldDetails{subPathNew, value, reflect.TypeOf(value).Kind().String(), ""}
    c <- mtD
}

func fieldMap (c chan FieldDetails, subPath []string, value interface{}, key interface{}, wg *sync.WaitGroup) {
    defer wg.Done()

    subPathNew := make([]string, len(subPath))
    if key == "" {
        copy(subPathNew, subPath)
    } else {
        copy(subPathNew, subPath)
        subPathNew = append(subPathNew, fmt.Sprint(key))
    }
    // once the value == {}
    if len(reflect.ValueOf(value).Interface().(map[string]interface{})) == 0 {
        mtD := FieldDetails{subPathNew, value, reflect.TypeOf(value).Kind().String(), ""}
        c <- mtD
    } else {
        mtD := FieldDetails{subPathNew, value, reflect.TypeOf(value).Kind().String(), ""}
        c <- mtD
        //
        for key2, value2 := range reflect.ValueOf(value).Interface().(map[string]interface{}) {
            switch value2.(type) {
                case nil:
                    wg.Add(1)
                    go fieldNull(c, subPathNew, value2, key2, wg)
                case string, float64, bool:
                    wg.Add(1)
                    go fieldPrimitive(c, subPathNew, value2, key2, wg)
                case map[string]interface{}:
                    wg.Add(1)
                    go fieldMap(c, subPathNew, value2, key2, wg)
                case []interface{}:
                    wg.Add(1)
                    go fieldSlice(c, subPathNew, value2, key2, wg)
                
            }
        }
    }
}

func fieldSlice (c chan FieldDetails, subPath []string, value interface{}, key interface{}, wg *sync.WaitGroup) {
    defer wg.Done()

    subPathNew := make([]string, len(subPath))
    if key == "" {
        copy(subPathNew, subPath)
    } else {
        copy(subPathNew, subPath)
        subPathNew = append(subPathNew, fmt.Sprint(key))
    }
    // once the value == []
    if len(reflect.ValueOf(value).Interface().([]interface{})) == 0 {
        mtD := FieldDetails{subPathNew, value, reflect.TypeOf(value).Kind().String(), ""}
        c <- mtD
    } else {
        mtD := FieldDetails{subPathNew, value, reflect.TypeOf(value).Kind().String(), ""}
        c <- mtD
        // loop all elments of the Slice, as it may contains different types
        for key2, value2 := range reflect.ValueOf(value).Interface().([]interface{}) {
            switch value2.(type) {
                case nil:
                    wg.Add(1)
                    go fieldNull(c, subPathNew, value2, fmt.Sprint(key2), wg)
                case string, float64, bool:
                    wg.Add(1)
                    go fieldPrimitive(c, subPathNew, value2, fmt.Sprint(key2), wg)
                case map[string]interface{}:
                    wg.Add(1)
                    go fieldMap(c, subPathNew, value2, fmt.Sprint(key2), wg)
                case []interface{}:
                    wg.Add(1)
                    go fieldSlice(c, subPathNew, value2, fmt.Sprint(key2), wg)
                
            }
        }
    }
}

func GetJsonNodesLevel (fieldDetailsSlice []FieldDetails) int {
    // get the max level of the paths
    max := 0
    for _, fieldDetails := range fieldDetailsSlice {
        if len(fieldDetails.FieldPath) > max {
            max = len(fieldDetails.FieldPath)
        }
    }

    return max
}

func GetJsonNodesPath (fieldDetailsSlice []FieldDetails) []string {
    var nodePaths []string
    for i, _ := range fieldDetailsSlice {
        nodePathStr := strings.Join(fieldDetailsSlice[i].FieldPath, ".")
        nodePaths = append(nodePaths, nodePathStr)
    }

    return nodePaths
}

// to get all the leaves, including the blank slice [], blank map {}
func GetJsonLeavesPath (fieldDetailsSlice []FieldDetails) []string {
    var leavesPath []string
    for i, _ := range fieldDetailsSlice {
        switch fieldDetailsSlice[i].FieldType {
            case "", "string", "float64", "bool":
                nodePathStr := strings.Join(fieldDetailsSlice[i].FieldPath, ".")
                leavesPath = append(leavesPath, nodePathStr)
            case "map":
                if len(reflect.ValueOf(fieldDetailsSlice[i].CurrValue).Interface().(map[string]interface{})) == 0 {
                    nodePathStr := strings.Join(fieldDetailsSlice[i].FieldPath, ".")
                    leavesPath = append(leavesPath, nodePathStr)
                }
            case "slice":
                if len(reflect.ValueOf(fieldDetailsSlice[i].CurrValue).Interface().([]interface{})) == 0 {
                    nodePathStr := strings.Join(fieldDetailsSlice[i].FieldPath, ".")
                    leavesPath = append(leavesPath, nodePathStr)
                }
        }
    }

    return leavesPath
}

// to get all the leaves, including the blank slice [], blank map {}
func GetJsonLeaves (fieldDetailsSlice []FieldDetails) []FieldDetails {
    var leaves []FieldDetails
    for i, _ := range fieldDetailsSlice {
        switch fieldDetailsSlice[i].FieldType {
            case "", "string", "float64", "bool":
                leaves = append(leaves, fieldDetailsSlice[i])
            case "map":
                if len(reflect.ValueOf(fieldDetailsSlice[i].CurrValue).Interface().(map[string]interface{})) == 0 {
                    leaves = append(leaves, fieldDetailsSlice[i])
                }
            case "slice":
                if len(reflect.ValueOf(fieldDetailsSlice[i].CurrValue).Interface().([]interface{})) == 0 {
                    leaves = append(leaves, fieldDetailsSlice[i])
                }
        }
    }

    return leaves
}

