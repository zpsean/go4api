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
    "fmt"
    "strings"
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
// Lib has no deepcopy, so that, here uses a plain way, that is, to re-sturct the TestCaseDataInfo
//
// Option 2:
// json.Marshal the tc in originMutationTcArray, 
// then change the value(s) in the json
// then Unmarshal the to testcase, and add to mutatedTcArray
// then execute the mutatedTcArray

// focus on the Request to mutate
// Note: for convinence to the results analysis, will distinguish the mutation by priority:
// type Request struct {  
//     Method string
//     Path string
//     Headers map[string]interface{} => change the priority: set 1, del 2, add 3, del all 4
//     QueryString map[string]interface{} => change the priority: set 5, del 6, add 7, del all 8
//     Payload map[string]interface{} => change the priority: set 9, del 10, add 11, del all 12
// }


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

        // --- here to start the mutation
        // querystring
        mutatedTcArrayQS := MutateRequestQueryString(originTcData, tcJson)
        mutatedTcArray = append(mutatedTcArray, mutatedTcArrayQS[0:]...)

        // headers
        mutatedTcArrayH := MutateRequestRequestHeader(originTcData, tcJson)
        mutatedTcArray = append(mutatedTcArray, mutatedTcArrayH[0:]...)
        
        // Payload
        // mutatedTcArrayPL := MutateRequestPayload(originTcData, tcJson)
        // mutatedTcArray = append(mutatedTcArray, mutatedTcArrayPL[0:]...)

    }
    // fmt.Println("\nmutatedTcArray: ", mutatedTcArray
    return mutatedTcArray
}

// RequestHeader
func MutateRequestRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    mSet := MutateSetRequestHeader(originTcData, tcJson)
    mutatedTcArray = append(mutatedTcArray, mSet[0:]...)

    mDel := MutateDelRequestHeader(originTcData, tcJson)
    mutatedTcArray = append(mutatedTcArray, mDel[0:]...)

    mAdd := MutateAddRequestHeader(originTcData, tcJson)
    mutatedTcArray = append(mutatedTcArray, mAdd[0:]...)

    mDelAll := MutateDelAllRequestHeaders(originTcData, tcJson)
    mutatedTcArray = append(mutatedTcArray, mDelAll[0:]...)

    return mutatedTcArray
}

func MutateSetRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    tcSuffix := ""
    for key, value := range originTcData.TestCase.ReqHeaders() {
        //
        mutationDetails := MutationDetails{[]string{}, value, reflect.TypeOf(value).Kind().String(), "", []interface{}{}}
        mType := mutationDetails.DetermineMutationType()
        mutatedValues := mutationDetails.CallMutationRules(mType)
        //
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            tcSuffix = "-M-H-S-" + fmt.Sprint(i)
            mutationInfo := "Update/Set header key: " + key + ", `" + fmt.Sprint(mutationDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue) + "`"
            
            //-- set new info to mutated tc
            var mTcData testcase.TestCaseDataInfo
            json.Unmarshal(tcJson, &mTcData)

            mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
            mTcData.TestCase.SetPriority(fmt.Sprint(1))
            mTcData.MutationInfo = mutationInfo

            mTcData.TestCase.SetRequestHeader(key, fmt.Sprint(mutatedValue))
            //
            mutatedTcArray = append(mutatedTcArray, mTcData)
        }
    }

    return mutatedTcArray
}

func MutateDelRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    tcSuffix := ""
    for key, _ := range originTcData.TestCase.ReqHeaders() {
        i = i + 1
        tcSuffix = "-M-H-D-" + fmt.Sprint(i)
        mutationInfo := "Remove header key: " + "`" + key + "`"

        // del the key
        var mTcData testcase.TestCaseDataInfo
        json.Unmarshal(tcJson, &mTcData)

        mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
        mTcData.TestCase.SetPriority(fmt.Sprint(2))
        mTcData.MutationInfo = mutationInfo

        mTcData.TestCase.DelRequestHeader(key)
        //
        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}

func MutateAddRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // add new key: get rand key, get rand value, then Add()
    i := 0
    tcSuffix := "-M-H-A-" + fmt.Sprint(i)

    randKey := RandStringRunes(5)
    randValue := RandStringRunes(5)
    mutationInfo := "Add new rand header key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"
    //
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    mTcData.TestCase.SetPriority(fmt.Sprint(3))
    mTcData.MutationInfo = mutationInfo

    mTcData.TestCase.AddRequestHeader(randKey, randValue)
    //
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}


func MutateDelAllRequestHeaders (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // remove all headers
    i := 0
    tcSuffix := "-M-H-D-" + fmt.Sprint(i)

    hFullPath := "TestCase." + originTcData.TcName() + ".request." + "headers"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), hFullPath)
    mutationInfo := "Remove all headers"
    //
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal([]byte(mutatedTcJson), &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    mTcData.TestCase.SetPriority(fmt.Sprint(4))
    mTcData.MutationInfo = mutationInfo
    //
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}


// RequestQueryString
func MutateRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    mSet := MutateSetRequestQueryString(originTcData, tcJson)
    mutatedTcArray = append(mutatedTcArray, mSet[0:]...)

    mDel := MutateDelRequestQueryString(originTcData, tcJson)
    mutatedTcArray = append(mutatedTcArray, mDel[0:]...)

    mAdd := MutateAddRequestQueryString(originTcData, tcJson)
    mutatedTcArray = append(mutatedTcArray, mAdd[0:]...)

    mDelAll := MutateDelAllRequestQueryStrings(originTcData, tcJson)
    mutatedTcArray = append(mutatedTcArray, mDelAll[0:]...)

    return mutatedTcArray
}


//
func MutateSetRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    tcSuffix := ""
    for key, value := range originTcData.TestCase.ReqQueryString() {
        mutationDetails := MutationDetails{[]string{}, value, reflect.TypeOf(value).Kind().String(), "", []interface{}{}}
        mType := mutationDetails.DetermineMutationType()
        mutatedValues := mutationDetails.CallMutationRules(mType)
        // loop and mutate the value, set new value to key
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            tcSuffix = "-M-QS-S-" + fmt.Sprint(i)
            mutationInfo := "Update/Set header key: " + key + ", `" + fmt.Sprint(mutationDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue) + "`"

            //-- set new info to mutated tc
            var mTcData testcase.TestCaseDataInfo
            json.Unmarshal(tcJson, &mTcData)

            mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
            mTcData.TestCase.SetPriority(fmt.Sprint(5))
            mTcData.MutationInfo = mutationInfo

            mTcData.TestCase.SetRequestQueryString(key, fmt.Sprint(mutatedValue))
            //
            mutatedTcArray = append(mutatedTcArray, mTcData)
        }
    }

    return mutatedTcArray
}


func MutateDelRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    tcSuffix := ""
    for key, _ := range originTcData.TestCase.ReqQueryString() {
        // del key
        i = i + 1
        tcSuffix = "-M-QS-D-" + fmt.Sprint(i)
        mutationInfo := "Remove querystring key: " + "`" + key + "`"

        // del the key
        var mTcData testcase.TestCaseDataInfo
        json.Unmarshal(tcJson, &mTcData)

        mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
        mTcData.TestCase.SetPriority(fmt.Sprint(6))
        mTcData.MutationInfo = mutationInfo

        mTcData.TestCase.DelRequestQueryString(key)
        //
        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}


func MutateAddRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // add new key: get rand key, get rand value, then Add()
    i := 0
    tcSuffix := "-M-QS-A-" + fmt.Sprint(i)

    randKey := RandStringRunes(5)
    randValue := RandStringRunes(5)

    mutationInfo := "Add new rand header key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    mTcData.TestCase.SetPriority(fmt.Sprint(7))
    mTcData.MutationInfo = mutationInfo

    mTcData.TestCase.AddRequestQueryString(randKey, randValue)
    //
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}

func MutateDelAllRequestQueryStrings (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // remove all querystring
    i := 0
    tcSuffix := "-M-QS-D-" + fmt.Sprint(i)

    qSFullPath := "TestCase." + originTcData.TcName() + ".request." + "queryString"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), qSFullPath)
    mutationInfo := "Remove all querystring"
    //
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal([]byte(mutatedTcJson), &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    mTcData.TestCase.SetPriority(fmt.Sprint(8))
    mTcData.MutationInfo = mutationInfo
    //
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}

// MutateRequestPayload
func MutateRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    for key, value := range originTcData.TestCase.ReqPayload() {
        if key == "text" {
            sturctFieldsDisplay(value)
            // to loop over the struct
            mutationDetailsSlice := getFieldsMutationDetails(value)
            //
            // (1) set
            mSet := MutateSetRequestPayload(originTcData, tcJson, key, mutationDetailsSlice)
            mutatedTcArray = append(mutatedTcArray, mSet[0:]...)
            // (2) del
            mDel := MutateDelRequestPayload(originTcData, tcJson, key, mutationDetailsSlice)
            mutatedTcArray = append(mutatedTcArray, mDel[0:]...)
            // (3) add
            mAdd := MutateAddRequestPayloadNode(originTcData, tcJson, key, mutationDetailsSlice)
            mutatedTcArray = append(mutatedTcArray, mAdd[0:]...)
            // (4)
            mDelAll := MutateDelWholeRequestPayloadNode(originTcData, tcJson, key, mutationDetailsSlice)
            mutatedTcArray = append(mutatedTcArray, mDelAll[0:]...)
        }
    }
    return mutatedTcArray
}

func MutateSetRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mutationDetailsSlice []MutationDetails) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // (1), set node value
    i := 0
    for _, mutationDetails := range mutationDetailsSlice {
        // set the value
        plPath := key + "." + strings.Join(mutationDetails.FieldPath, ".")
        plFullPath := "TestCase." + originTcData.TcName() + ".request.payload" + "." + plPath
        // mutate the value based on rules 
        mType := mutationDetails.DetermineMutationType()
        mutatedValues := mutationDetails.CallMutationRules(mType)
        // (1). set node
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            tcSuffix := "-M-PL-S-" + fmt.Sprint(i)

            mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, mutatedValue)
            mutationInfo := fmt.Sprint(mutationDetails) + ", `" + fmt.Sprint(mutationDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue) + "`"
            //
            var mTcData testcase.TestCaseDataInfo
            json.Unmarshal([]byte(mutatedTcJson), &mTcData)

            mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
            mTcData.TestCase.SetPriority(fmt.Sprint(9))
            mTcData.MutationInfo = mutationInfo
            //
            mutatedTcArray = append(mutatedTcArray, mTcData)
        }

        // add new node
    }
    return mutatedTcArray
}

func MutateDelRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mutationDetailsSlice []MutationDetails) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    // get the max level of the paths
    max := 0
    for _, mutationDetails := range mutationDetailsSlice {
        if len(mutationDetails.FieldPath) > max {
            max = len(mutationDetails.FieldPath)
        }
    }
    // (2). del node
    var nodePaths []string
    for i := max; i > 0; i-- {
        for _, mutationDetails := range mutationDetailsSlice {
            if len(mutationDetails.FieldPath) >= i {
                nodePathStr := strings.Join(mutationDetails.FieldPath[0:i], ".")

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
    for _, pathStr := range nodePaths {
        i = i + 1
        tcSuffix := "-M-PL-D-" + fmt.Sprint(i)

        plPath := key + "." + pathStr
        plFullPath := "TestCase." + originTcData.TcName() + ".request.payload" + "." + plPath

        mutatedTcJson, _ := sjson.Delete(string(tcJson), plFullPath)
        mutationInfo := "Remove payload value on node: " + pathStr
        //
        var mTcData testcase.TestCaseDataInfo
        json.Unmarshal([]byte(mutatedTcJson), &mTcData)

        mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
        mTcData.TestCase.SetPriority(fmt.Sprint(10))
        mTcData.MutationInfo = mutationInfo
        //
        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}


func MutateAddRequestPayloadNode (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mutationDetailsSlice []MutationDetails) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // (3). add new node, for each node level

    return mutatedTcArray
}

func MutateDelWholeRequestPayloadNode (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mutationDetailsSlice []MutationDetails) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    // (4). remove whole payload, i.e. set to: "text" : {} or "text" : null
    tcSuffix := "-M-PL-D-" + fmt.Sprint(i)

    plFullPath := "TestCase." + originTcData.TcName() + ".request.payload." + key
    // mutatedTcJson, _ := sjson.Delete(string(tcJson), plFullPath)
    mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, "")
    mutationInfo := "Remove whole post body"
    //
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal([]byte(mutatedTcJson), &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    mTcData.TestCase.SetPriority(fmt.Sprint(12))
    mTcData.MutationInfo = mutationInfo
    //
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
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


func getFieldsMutationDetails(value interface{}) []MutationDetails {
    c := make(chan MutationDetails)

    go func(c chan MutationDetails) {
        defer close(c)
        sturctFieldsMutation(c, []string{}, value)
    }(c)

    var mutationDetailsSlice []MutationDetails
    //
    for mutationDetails := range c {
        mutationDetailsSlice = append(mutationDetailsSlice, mutationDetails)
    }

    return mutationDetailsSlice
}


func sturctFieldsMutation(c chan MutationDetails, subPath []string, value interface{}) {
    switch reflect.TypeOf(value).Kind() {
        case reflect.Map: {
            // fmt.Println("value: ", value, reflect.TypeOf(value), reflect.TypeOf(value).Kind())
            for key2, value2 := range reflect.ValueOf(value).Interface().(map[string]interface{}) {
                // fmt.Println("key2, value2: ", key2, reflect.TypeOf(value2))
                switch reflect.TypeOf(value2).Kind() {
                    case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
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
                    case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
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

