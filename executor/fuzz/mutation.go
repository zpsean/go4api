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
    // "path/filepath"
    "strings"
    // "strconv"
    "reflect"
    "encoding/json"
    "go4api/testcase"
    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

// mutation is to mutate the valid data to working api, see if mutated invalid data still can be handled by the api
// two ways to mutate the testcase:
// Option 1: 
// copy the underlying fields and values to another TestCaseDataInfo, with mutation(s)
// the better way would be deep copy the TestCaseDataInfo, and change the value, but Golang standard
// Lib has no deepcopy, so that, here use a plain way, that is, to re-sturct the TestCaseDataInfo
//
// Option 2:
// json.Marshal the tc in originMutationTcArray, 
// then change the value(s) in the json
// then Unmarshal the to testcase, and add to mutatedTcArray
// then execute the mutatedTcArray

// focus on the Request to mutate
// type Request struct {  
//     Method string
//     Path string
//     Headers map[string]interface{}
//     QueryString map[string]interface{}
//     Payload map[string]interface{}
// }

type PayloadInfo struct {
    FieldPath []string
    CurrValue interface{}
    FieldType string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
}


type MutationDetails struct {
    FieldPath []string
    CurrValue interface{}
    FieldType string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
    MutatedValues []interface{}
}

func MutateTcArray(originMutationTcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    for _, originTcData := range originMutationTcArray {
        tcJson, _ := json.Marshal(originTcData)
        // fmt.Println("tcJson:", string(tcJson)) 

        mutatedTcArray = append(mutatedTcArray, originTcData)

        // here to start the mutation
        // headers
        // mutatedTcArrayH := MutateRequestRequestHeader(originTcData, tcJson)
        

        // querystring
        // mutatedTcArrayQS := MutateRequestQueryString(originTcData, tcJson)
        

        // --------------------------------------------
        // Payload, loop all node, mutate it (include remove)
        //
        // payloadJson, _ := json.Marshal(originTcData.TestCase.ReqPayload())
        // fmt.Println("->payloadJson: ", string(payloadJson)) 

        // mutatedTcArrayPL := MutateRequestPayload(originTcData, tcJson)

    }
    // fmt.Println("\nmutatedTcArray: ", mutatedTcArray
    return mutatedTcArray
}

// RequestQueryString
func MutateRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    i := 0
    for key, value := range originTcData.TestCase.ReqQueryString() {
        // if value match number mode
        mutationDetails := MutationDetails{[]string{}, value, reflect.TypeOf(value).Kind().String(), "", []interface{}{}}
        mType := mutationDetails.DetermineMutationType()
        mutatedValues := mutationDetails.CallRules(mType)

        // loogp and mutate the value, set new value to key
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            mutationInfo := fmt.Sprint(mutationDetails) + "," + fmt.Sprint(mutationDetails.CurrValue) + ", `" + fmt.Sprint(mutatedValue) + "`"

            mTc := MutateSetRequestQueryString(tcJson, mutationInfo, key, mutatedValue.(string), key + fmt.Sprint(i))
            mutatedTcArray = append(mutatedTcArray, mTc)
        }

        
        // del key
        mTc := MutateDelRequestQueryString(tcJson, key, i)
        mutatedTcArray = append(mutatedTcArray, mTc)
        i = i + 1
    }

    // add new key
    // get rand key, get rand value, then Add()
    // mTc := MutateAddRequestQueryString(tcJson, mutationInfo, key, mutatedValue.(string), key + fmt.Sprint(i))
    // mutatedTcArray = append(mutatedTcArray, mTc)
    // mutatedTcArray = append(mutatedTcArray, MutateAddRequestQueryString(tcJson))

    // remove all querystring?
    // qSFullPath := "TestCase." + originTcData.TcName() + ".Request." + "QueryString"
    // mutatedTcJson, _ := sjson.Del(string(tcJson), qSFullPath)
    // mTc := MutateDelRequestQueryString([]byte(mutatedTcJson), mutationInfo, "1-" + fmt.Sprint(i))

    return mutatedTcArray
}

// RequestHeader
func MutateRequestRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    i := 0
    for key, value := range originTcData.TestCase.ReqHeaders() {
        // mutatedTcArray = append(mutatedTcArray, MutateDelRequestHeader(tcJson, k, i))
        // i = i + 1

       // if value match number mode
        mutationDetails := MutationDetails{[]string{}, value, reflect.TypeOf(value).Kind().String(), "", []interface{}{}}
        mType := mutationDetails.DetermineMutationType()
        mutatedValues := mutationDetails.CallRules(mType)

        // loogp and mutate the value, set new value to key
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            mutationInfo := fmt.Sprint(mutationDetails) + "," + fmt.Sprint(mutationDetails.CurrValue) + ", `" + fmt.Sprint(mutatedValue) + "`"

            mTc := MutateSetRequestHeader(tcJson, mutationInfo, key, mutatedValue.(string), key + fmt.Sprint(i))
            mutatedTcArray = append(mutatedTcArray, mTc)
        }

        // del key
        mTc := MutateDelRequestHeader(tcJson, key, i)
        mutatedTcArray = append(mutatedTcArray, mTc)
        i = i + 1
    }

    // add new key
    // get rand key, get rand value, then Add()

    // remove all headers?
    // qSFullPath := "TestCase." + originTcData.TcName() + ".Request." + "QueryString"
    // mutatedTcJson, _ := sjson.Del(string(tcJson), qSFullPath)
    // mTc := MutateDelRequestQueryString([]byte(mutatedTcJson), mutationInfo, "1-" + fmt.Sprint(i))

    return mutatedTcArray
}

// MutateRequestPayload
func MutateRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    i := 0
    for key, value := range originTcData.TestCase.ReqPayload() {
        if key == "text" {
            // to loop over the struct
            c := make(chan MutationDetails)

            go func(c chan MutationDetails) {
                defer close(c)
                sturctFieldsMutation(c, []string{}, value)
            }(c)

            for mutationDetails := range c {
                // set the value
                payloadPath := key + "." + strings.Join(mutationDetails.FieldPath, ".")
                payloadFullPath := "TestCase." + originTcData.TcName() + ".Request.Payload" + "." + payloadPath

                // mutate the value based on rules
                // (1). get values 
                mType := mutationDetails.DetermineMutationType()
                
                mutatedValues := mutationDetails.CallRules(mType)
                // fmt.Println("mutatedValues: ", mutatedValues, payloadFullPath, i)
                // (2). remove node

                for _, mutatedValue := range mutatedValues {
                    i = i + 1
                    mutatedTcJson, _ := sjson.Set(string(tcJson), payloadFullPath, mutatedValue)

                    mutationInfo := fmt.Sprint(mutationDetails) + "," + fmt.Sprint(mutationDetails.CurrValue) + ", `" + fmt.Sprint(mutatedValue) + "`"

                    mTc := MutateSetRequestPayload([]byte(mutatedTcJson), mutationInfo, "1-" + fmt.Sprint(i))
                    mutatedTcArray = append(mutatedTcArray, mTc)
                }

                // add new node
                // del node
                // del payload
            }
        }
    }

    // add new key
    // get rand key, get rand value, then Add()

    // remove all headers?
    // qSFullPath := "TestCase." + originTcData.TcName() + ".Request." + "QueryString"
    // mutatedTcJson, _ := sjson.Del(string(tcJson), qSFullPath)
    // mTc := MutateDelRequestQueryString([]byte(mutatedTcJson), mutationInfo, "1-" + fmt.Sprint(i))

    return mutatedTcArray
}

//
func MutateSetRequestHeader (tcJson []byte, mutationInfo interface{}, key string, value string, suffix string) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-seth-" + fmt.Sprint(1))
    mTcData.TestCase.SetRequestHeader("aaaa", "dbddsdsfa")

    mTcData.MutationInfo = mutationInfo

    return mTcData
}


func MutateAddRequestHeader (tcJson []byte) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-addh-" + fmt.Sprint(2))
    mTcData.TestCase.AddRequestHeader("aaaakk", "dbddsdsfa")

    return mTcData
}


func MutateDelRequestHeader (tcJson []byte, k string, i int) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-Delh-" + fmt.Sprint(i))
    mTcData.TestCase.DelRequestHeader(k)

    return mTcData
}



//
func MutateSetRequestQueryString (tcJson []byte, mutationInfo interface{}, key string, value string, suffix string) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-setq-" + suffix)
    mTcData.TestCase.SetRequestQueryString(key, value)

    mTcData.MutationInfo = mutationInfo

    return mTcData
}


func MutateAddRequestQueryString (tcJson []byte) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-addq-" + fmt.Sprint(2))

    mTcData.TestCase.AddRequestQueryString("aaaakk", "dbddsdsfa")

    return mTcData
}


func MutateDelRequestQueryString (tcJson []byte, k string, i int) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-Delq-" + fmt.Sprint(i))

    mTcData.TestCase.DelRequestQueryString(k)

    return mTcData
}


//
func MutateSetRequestPayload (tcJson []byte, mutationInfo interface{}, suffix string) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    // change the tc name
    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-P-" + suffix)
    // change the tc priority?

    mTcData.MutationInfo = mutationInfo

    return mTcData
}


func sturctFieldsDisplay(value interface{}) {
    switch reflect.TypeOf(value).Kind() {
        case reflect.String:
            fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
        case reflect.Int32:
            fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
        case reflect.Map: {
            // fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
            for key2, value2 := range reflect.ValueOf(value).Interface().(map[string]interface{}) {
                // fmt.Println("key2, value2: ", key2, reflect.TypeOf(value2))
                switch reflect.TypeOf(value2).Kind() {
                    case reflect.String:
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                    case reflect.Int32:
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                    case reflect.Map:
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                        sturctFieldsDisplay(value2)
                    case reflect.Array:
                        // note: maybe the Array/Slice is the last node, if it contains concrete type, like [1, 2, 3, ...]
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                        sturctFieldsDisplay(value2)
                    case reflect.Slice:
                        // note: maybe the Array/Slice is the last node, if it contains concrete type, like [1, 2, 3, ...]
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                        sturctFieldsDisplay(value2)
                }
            }     
        }
        case reflect.Array, reflect.Slice: {
            // fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
            for key2, value2 := range reflect.ValueOf(value).Interface().([]interface{}) {
                // fmt.Println("key2, value2: ", key2, reflect.TypeOf(value2))
                switch reflect.TypeOf(value2).Kind() {
                    case reflect.String:
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                    case reflect.Int32:
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                    case reflect.Map:
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                        sturctFieldsDisplay(value2)
                    case reflect.Array:
                        // note: maybe the Array/Slice is the last node, if it contains concrete type, like [1, 2, 3, ...]
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                        sturctFieldsDisplay(value2)
                    case reflect.Slice:
                        // note: maybe the Array/Slice is the last node, if it contains concrete type, like [1, 2, 3, ...]
                        fmt.Println("key2, value2: ", key2, value2, reflect.TypeOf(value2), reflect.TypeOf(value2).Kind())
                        sturctFieldsDisplay(value2)
                }
            }  
        }
    }
}


func sturctFieldsMutation(c chan MutationDetails, subPath []string, value interface{}) {
    switch reflect.TypeOf(value).Kind() {
        case reflect.Map: {
            // fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
            for key2, value2 := range reflect.ValueOf(value).Interface().(map[string]interface{}) {
                // fmt.Println("key2, value2: ", key2, reflect.TypeOf(value2))
                switch reflect.TypeOf(value2).Kind() {
                    case reflect.String, reflect.Int, reflect.Float64:
                        subPathNew := append(subPath, key2)
                        output := make([]string, len(subPathNew))
                        copy(output, subPathNew)

                        mtD := MutationDetails{output, value2, reflect.TypeOf(value2).Kind().String(), "", []interface{}{}}
                        c <- mtD
                    case reflect.Map:
                        subPathNew := append(subPath, key2)
                        sturctFieldsMutation(c, subPathNew, value2)
                    case reflect.Array, reflect.Slice:
                        // note: maybe the Array/Slice is the last node, if it contains concrete type, like [1, 2, 3, ...]
                        for _, v := range reflect.ValueOf(value2).Interface().([]interface{}) {
                            switch reflect.TypeOf(v).Kind() {
                                case reflect.Array, reflect.Slice, reflect.Map:
                                    subPathNew := append(subPath, fmt.Sprint(key2))
                                    sturctFieldsMutation(c, subPathNew, value2)
                            }
                            break
                        }
                }
            }     
        }
        case reflect.Array, reflect.Slice: {
            // fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
            for key2, value2 := range reflect.ValueOf(value).Interface().([]interface{}) {
                // fmt.Println("key2, value2: ", key2, reflect.TypeOf(value2))
                switch reflect.TypeOf(value2).Kind() {
                    case reflect.String, reflect.Int, reflect.Float64:
                        subPathNew := append(subPath, fmt.Sprint(key2))
                        output := make([]string, len(subPathNew))
                        copy(output, subPathNew)

                        mtD := MutationDetails{output, value2, reflect.TypeOf(value2).Kind().String(), "", []interface{}{}}
                        c <- mtD
                    case reflect.Map:
                        subPathNew := append(subPath, fmt.Sprint(key2))
                        sturctFieldsMutation(c, subPathNew, value2)
                    case reflect.Array, reflect.Slice:
                        // note: maybe the Array/Slice is the last node, if it contains concrete type, like [1, 2, 3, ...]
                        for _, v := range reflect.ValueOf(value2).Interface().([]interface{}) {
                            switch reflect.TypeOf(v).Kind() {
                                case reflect.Array, reflect.Slice, reflect.Map:
                                    subPathNew := append(subPath, fmt.Sprint(key2))
                                    sturctFieldsMutation(c, subPathNew, value2)
                            }
                        }
                        break
                }
            } 
        }
    }
}

