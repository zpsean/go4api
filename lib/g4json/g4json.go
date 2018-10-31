/*
 * go4api - a api testing tool written in Go
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
    "strings"
    "reflect"
    // "encoding/csv"
    // "encoding/json"
)

type FieldDetails struct {
    FieldPath []string
    CurrValue interface{}
    FieldType string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
}

func GetFieldsDetails(value interface{}) []FieldDetails {
    c := make(chan FieldDetails)

    go func(c chan FieldDetails) {
        defer close(c)
        sturctFields(c, []string{}, value)
    }(c)

    var mFieldDetailsSlice []FieldDetails
    //
    for mFieldDetails := range c {
        mFieldDetailsSlice = append(mFieldDetailsSlice, mFieldDetails)
    }

    return mFieldDetailsSlice
}


func sturctFields(c chan FieldDetails, subPath []string, value interface{}) {
    switch reflect.TypeOf(value).Kind() {
        case reflect.Map: {
            // fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
            for key2, value2 := range reflect.ValueOf(value).Interface().(map[string]interface{}) {
                // fmt.Println("key2, value2: ", key2, reflect.TypeOf(value2))
                if value2 != nil {
                    switch reflect.TypeOf(value2).Kind() {
                        case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
                            subPathNew := append(subPath, key2)
                            output := make([]string, len(subPathNew))
                            copy(output, subPathNew)

                            mtD := FieldDetails{output, value2, reflect.TypeOf(value2).Kind().String(), ""}
                            c <- mtD
                        case reflect.Map:
                            subPathNew := append(subPath, key2)
                            sturctFields(c, subPathNew, value2)
                        case reflect.Array, reflect.Slice:
                            sturctFieldsSlice(c, subPath, value2, key2)
                    }
                } else {
                    subPathNew := append(subPath, key2)
                    output := make([]string, len(subPathNew))
                    copy(output, subPathNew)

                    mtD := FieldDetails{output, nil, "", ""}
                    c <- mtD
                }
            }     
        }
        case reflect.Array, reflect.Slice: {
            // fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
            for key2, value2 := range reflect.ValueOf(value).Interface().([]interface{}) {
                // fmt.Println("key2, value2: ", key2, reflect.TypeOf(value2))
                switch reflect.TypeOf(value2).Kind() {
                    case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
                        subPathNew := append(subPath, fmt.Sprint(key2))
                        output := make([]string, len(subPathNew))
                        copy(output, subPathNew)

                        mtD := FieldDetails{output, value2, reflect.TypeOf(value2).Kind().String(), ""}
                        c <- mtD
                    case reflect.Map:
                        subPathNew := append(subPath, fmt.Sprint(key2))
                        sturctFields(c, subPathNew, value2)
                    case reflect.Array, reflect.Slice:
                        sturctFieldsSlice(c, subPath, value2, key2)
                }
            } 
        }
    }
}

func sturctFieldsSlice (c chan FieldDetails, subPath []string, value2 interface{}, key2 interface{}) {
    // note: maybe the Array/Slice is the last node, if it contains concrete type, like [1, 2, 3, ...]
    if len(reflect.ValueOf(value2).Interface().([]interface{})) == 0 {
        subPathNew := append(subPath, fmt.Sprint(key2))
        output := make([]string, len(subPathNew))
        copy(output, subPathNew)

        mtD := FieldDetails{output, value2, reflect.TypeOf(value2).Kind().String(), ""}
        c <- mtD
    }

    for _, v := range reflect.ValueOf(value2).Interface().([]interface{}) {
        if v != nil { 
            switch reflect.TypeOf(v).Kind() {
                case reflect.Array, reflect.Slice, reflect.Map:
                    subPathNew := append(subPath, fmt.Sprint(key2))
                    sturctFields(c, subPathNew, value2)
                case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
                    subPathNew := append(subPath, fmt.Sprint(key2))
                    // subPathNew = append(subPathNew, fmt.Sprint(index))
                    output := make([]string, len(subPathNew))
                    copy(output, subPathNew)

                    mtD := FieldDetails{output, value2, reflect.TypeOf(value2).Kind().String(), ""}
                    c <- mtD
            }
            break
        } else {
            subPathNew := append(subPath, fmt.Sprint(key2))
            output := make([]string, len(subPathNew))
            copy(output, subPathNew)

            mtD := FieldDetails{output, nil, "", ""}
            c <- mtD
        }
    }
}


func getJsonNodePaths (FieldDetailsSlice []FieldDetails) ([]string, int) {
    // get the max level of the paths
    max := 0
    for _, fieldDetails := range FieldDetailsSlice {
        if len(fieldDetails.FieldPath) > max {
            max = len(fieldDetails.FieldPath)
        }
    }
    // 
    var nodePaths []string
    for i := max; i > 0; i-- {
        for _, fieldDetails := range FieldDetailsSlice {
            if len(fieldDetails.FieldPath) >= i {
                nodePathStr := strings.Join(fieldDetails.FieldPath[0:i], ".")

                ifExists := ""
                for _, str := range nodePaths {
                    if nodePathStr == str {
                        ifExists = "Y"
                        break
                    }
                }
                if ifExists == "" {
                    nodePaths = append(nodePaths, nodePathStr)
                }
            }
        }
    }

    return nodePaths, max
}



