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

    "go4api/lib/testcase"
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

var mCategoryPriorityMap = make(map[string][]string)

type MutationDetails struct {
    FieldPath []string
    CurrValue interface{}
    FieldType string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
    MutatedValues []interface{}
}

func init() {
    mCategoryPriorityMap["SetRequestHeader"] = []string{"1", "-M-H-S-"}
    mCategoryPriorityMap["DelRequestHeader"] = []string{"2", "-M-H-D-"}
    mCategoryPriorityMap["AddRequestHeader"] = []string{"3", "-M-H-A-"}
    mCategoryPriorityMap["DelAllRequestHeaders"] = []string{"4", "-M-H-D-"}

    mCategoryPriorityMap["SetRequestQueryString"] = []string{"5", "-M-QS-S-"}
    mCategoryPriorityMap["DelRequestQueryString"] = []string{"6", "-M-QS-D-"}
    mCategoryPriorityMap["AddRequestQueryString"] = []string{"7", "-M-QS-A-"}
    mCategoryPriorityMap["DelAllRequestQueryStrings"] = []string{"8", "-M-QS-D-"}

    mCategoryPriorityMap["SetRequestPayload"] = []string{"9", "-M-PL-S-"}
    mCategoryPriorityMap["DelRequestPayload"] = []string{"10", "-M-PL-D-"}
    mCategoryPriorityMap["AddRequestPayloadNode"] = []string{"11", "-M-PL-A-"}
    mCategoryPriorityMap["DelWholeRequestPayloadNode"] = []string{"12", "-M-PL-D-"}   
}


func MutateTcArray(originMutationTcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    for _, originTcData := range originMutationTcArray {
        tcJson, _ := json.Marshal(originTcData)
        // fmt.Println("tcJsonHTTP:", string(tcJson)) 
        mutatedTcArray = append(mutatedTcArray, originTcData)

        // --- here to start the mutation
        // headers
        mutatedTcArrayH := MutateRequestRequestHeader(originTcData, tcJson)
        mutatedTcArray = append(mutatedTcArray, mutatedTcArrayH[0:]...)

        // querystring
        mutatedTcArrayQS := MutateRequestQueryString(originTcData, tcJson)
        mutatedTcArray = append(mutatedTcArray, mutatedTcArrayQS[0:]...)
        
        // Payload
        mutatedTcArrayPL := MutateRequestPayload(originTcData, tcJson)
        mutatedTcArray = append(mutatedTcArray, mutatedTcArrayPL[0:]...)

    }
    // fmt.Println("\nmutatedTcArray: ", mutatedTcArray
    return mutatedTcArray
}

func getMutatedTcData (tcJson []byte, i int, mArea string, mCategory string, mutationRule string, mutationInfo string, tcMutationInfo testcase.MutationInfo) testcase.TestCaseDataInfo {
    tcSuffix := mCategoryPriorityMap[mCategory][1] + fmt.Sprint(i)

    //-- set new info to mutated tc
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    mTcData.TestCase.SetPriority(mCategoryPriorityMap[mCategory][0])
    mTcData.MutationArea = mArea
    mTcData.MutationCategory = mCategory
    mTcData.MutationRule = mutationRule
    mTcData.MutationInfoStr = mutationInfo
    mTcData.MutationInfo = tcMutationInfo
    
    return mTcData
}

func getTcMutationInfo (mutationDetails MutationDetails, mutatedValue interface{}) testcase.MutationInfo {
    tcMutationInfo := testcase.MutationInfo {
        FieldPath: mutationDetails.FieldPath,
        CurrValue: mutationDetails.CurrValue,
        FieldType: mutationDetails.FieldType,
        FieldSubType: mutationDetails.FieldSubType,
        MutatedValue: mutatedValue,
    }

    return tcMutationInfo
}

// RequestHeader
func MutateRequestRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    mSet := MutateSetRequestHeader(originTcData, tcJson, "headers", "SetRequestHeader")
    mutatedTcArray = append(mutatedTcArray, mSet[0:]...)

    mDel := MutateDelRequestHeader(originTcData, tcJson, "headers", "DelRequestHeader")
    mutatedTcArray = append(mutatedTcArray, mDel[0:]...)

    mAdd := MutateAddRequestHeader(originTcData, tcJson, "headers", "AddRequestHeader")
    mutatedTcArray = append(mutatedTcArray, mAdd[0:]...)

    mDelAll := MutateDelAllRequestHeaders(originTcData, tcJson, "headers", "DelAllRequestHeaders")
    mutatedTcArray = append(mutatedTcArray, mDelAll[0:]...)

    return mutatedTcArray
}

func MutateSetRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte, mArea string, mCategory string) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    for key, value := range originTcData.TestCase.ReqHeaders() {
        //
        mutationDetails := MutationDetails{[]string{}, value, reflect.TypeOf(value).Kind().String(), "", []interface{}{}}
        mType := mutationDetails.DetermineMutationType()
        mutatedValues := mutationDetails.CallMutationRules(mType)
        //
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            mutationInfo := "Update/Set header key: " + key + ", `" + fmt.Sprint(mutationDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mutatedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mutationDetails, mutatedValue.MutatedValue)
            
            //-- set new info to mutated tc
            mTcData := getMutatedTcData(tcJson, i, mArea, mCategory, mutatedValue.MutationRule, mutationInfo, tcMutationInfo)
            mTcData.TestCase.SetRequestHeader(key, fmt.Sprint(mutatedValue.MutatedValue))

            mutatedTcArray = append(mutatedTcArray, mTcData)
        }
    }

    return mutatedTcArray
}

func MutateDelRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte, mArea string, mCategory string) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    for key, _ := range originTcData.TestCase.ReqHeaders() {
        i = i + 1
        mutationDetails := MutationDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
        mutationInfo := "Remove header key: " + "`" + key + "`"

        tcMutationInfo := getTcMutationInfo(mutationDetails, "")

        // del the key
        mTcData := getMutatedTcData(tcJson, i, mArea, mCategory, "Remove header key", mutationInfo, tcMutationInfo)
        mTcData.TestCase.DelRequestHeader(key)

        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}

func MutateAddRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte, mArea string, mCategory string) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // add new key: get rand key, get rand value, then Add()
    i := 0

    randKey := RandStringRunes(5)
    randValue := RandStringRunes(5)
    mutationDetails := MutationDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Add new rand header key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mutationDetails, randValue)
    //
    mTcData := getMutatedTcData(tcJson, i, mArea, mCategory, "Add new rand header key", mutationInfo, tcMutationInfo)
    mTcData.TestCase.AddRequestHeader(randKey, randValue)

    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}


func MutateDelAllRequestHeaders (originTcData testcase.TestCaseDataInfo, tcJson []byte, mArea string, mCategory string) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // remove all headers
    i := 0
  
    hFullPath := "TestCase." + originTcData.TcName() + ".request." + "headers"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), hFullPath)
    mutationDetails := MutationDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Remove all headers"

    tcMutationInfo := getTcMutationInfo(mutationDetails, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mArea, mCategory, "Remove all headers", mutationInfo, tcMutationInfo)
    
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}


// RequestQueryString
func MutateRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    mSet := MutateSetRequestQueryString(originTcData, tcJson, "queryString", "SetRequestQueryString")
    mutatedTcArray = append(mutatedTcArray, mSet[0:]...)

    mDel := MutateDelRequestQueryString(originTcData, tcJson, "queryString", "DelRequestQueryString")
    mutatedTcArray = append(mutatedTcArray, mDel[0:]...)

    mAdd := MutateAddRequestQueryString(originTcData, tcJson, "queryString", "AddRequestQueryString")
    mutatedTcArray = append(mutatedTcArray, mAdd[0:]...)

    mDelAll := MutateDelAllRequestQueryStrings(originTcData, tcJson, "queryString", "DelAllRequestQueryStrings")
    mutatedTcArray = append(mutatedTcArray, mDelAll[0:]...)

    return mutatedTcArray
}


//
func MutateSetRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte, mArea string, mCategory string) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    for key, value := range originTcData.TestCase.ReqQueryString() {
        mutationDetails := MutationDetails{[]string{}, value, reflect.TypeOf(value).Kind().String(), "", []interface{}{}}
        mType := mutationDetails.DetermineMutationType()
        mutatedValues := mutationDetails.CallMutationRules(mType)
        // loop and mutate the value, set new value to key
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            mutationInfo := "Update/Set header key: " + key + ", `" + fmt.Sprint(mutationDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mutatedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mutationDetails, mutatedValue.MutatedValue)

            //-- set new info to mutated tc
            mTcData := getMutatedTcData(tcJson, i, mArea, mCategory, mutatedValue.MutationRule, mutationInfo, tcMutationInfo)
            mTcData.TestCase.SetRequestQueryString(key, fmt.Sprint(mutatedValue.MutatedValue))

            mutatedTcArray = append(mutatedTcArray, mTcData)
        }
    }

    return mutatedTcArray
}


func MutateDelRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte, mArea string, mCategory string) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    for key, _ := range originTcData.TestCase.ReqQueryString() {
        // del key
        i = i + 1
        mutationDetails := MutationDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
        mutationInfo := "Remove querystring key: " + "`" + key + "`"

        tcMutationInfo := getTcMutationInfo(mutationDetails, "")

        // del the key
        mTcData := getMutatedTcData(tcJson, i, mArea, mCategory, "Remove querystring key", mutationInfo, tcMutationInfo)
        mTcData.TestCase.DelRequestQueryString(key)

        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}


func MutateAddRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte, mArea string, mCategory string) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // add new key: get rand key, get rand value, then Add()
    i := 0

    randKey := RandStringRunes(5)
    randValue := RandStringRunes(5)

    mutationDetails := MutationDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Add new rand QueryString key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mutationDetails, "")
    //
    mTcData := getMutatedTcData(tcJson, i, mArea, mCategory, "Add new rand querystring key", mutationInfo, tcMutationInfo)
    mTcData.TestCase.AddRequestQueryString(randKey, randValue)

    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}

func MutateDelAllRequestQueryStrings (originTcData testcase.TestCaseDataInfo, tcJson []byte, mArea string, mCategory string) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // remove all querystring
    i := 0
 
    qSFullPath := "TestCase." + originTcData.TcName() + ".request." + "queryString"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), qSFullPath)
    mutationDetails := MutationDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Remove all querystring"

    tcMutationInfo := getTcMutationInfo(mutationDetails, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mArea, mCategory, "Remove all querystring", mutationInfo, tcMutationInfo)
    
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}

// MutateRequestPayload
func MutateRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    for key, value := range originTcData.TestCase.ReqPayload() {
        if key == "text" {
            // sturctFieldsDisplay(value)
            // to loop over the struct
            var mutationDetailsSlice []MutationDetails
            mutationDetailsSlice = getFieldsMutationDetails(value)
            //
            // (1) set
            mSet := MutateSetRequestPayload(originTcData, tcJson, key, mutationDetailsSlice, "payload", "SetRequestPayload")
            mutatedTcArray = append(mutatedTcArray, mSet[0:]...)
            // (2) del
            mDel := MutateDelRequestPayload(originTcData, tcJson, key, mutationDetailsSlice, "payload", "DelRequestPayload")
            mutatedTcArray = append(mutatedTcArray, mDel[0:]...)
            // (3) add
            mAdd := MutateAddRequestPayloadNode(originTcData, tcJson, key, mutationDetailsSlice, "payload", "AddRequestPayloadNode")
            mutatedTcArray = append(mutatedTcArray, mAdd[0:]...)
            // (4)
            mDelAll := MutateDelWholeRequestPayloadNode(originTcData, tcJson, key, mutationDetailsSlice, "payload", "DelWholeRequestPayloadNode")
            mutatedTcArray = append(mutatedTcArray, mDelAll[0:]...)
        }
    }
    return mutatedTcArray
}

func MutateSetRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mutationDetailsSlice []MutationDetails, 
        mArea string, mCategory string) []testcase.TestCaseDataInfo {
    //--------------------
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
     
            mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, mutatedValue.MutatedValue)
            mDJsonByte, _ := json.Marshal(mutationDetails)
            mutationInfo := fmt.Sprint(string(mDJsonByte)) + ", `" + fmt.Sprint(mutationDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mutatedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mutationDetails, mutatedValue.MutatedValue)
            //
            mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mArea, mCategory, mutatedValue.MutationRule, mutationInfo, tcMutationInfo)
    
            mutatedTcArray = append(mutatedTcArray, mTcData)
        }

        // add new node
    }
    return mutatedTcArray
}

func MutateDelRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mutationDetailsSlice []MutationDetails, 
        mArea string, mCategory string) []testcase.TestCaseDataInfo {
    //--------------------
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
    
        plPath := key + "." + pathStr
        plFullPath := "TestCase." + originTcData.TcName() + ".request.payload" + "." + plPath

        mutatedTcJson, _ := sjson.Delete(string(tcJson), plFullPath)
        mutationDetails := MutationDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
        mutationInfo := "Remove payload value on node: " + pathStr

        tcMutationInfo := getTcMutationInfo(mutationDetails, "")
        //
        mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mArea, mCategory, "Remove payload value on node", mutationInfo, tcMutationInfo)
    
        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}


func MutateAddRequestPayloadNode (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mutationDetailsSlice []MutationDetails, 
        mArea string, mCategory string) []testcase.TestCaseDataInfo {
    //--------------------
    var mutatedTcArray []testcase.TestCaseDataInfo
    // (3). add new node, for each node level

    return mutatedTcArray
}

func MutateDelWholeRequestPayloadNode (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mutationDetailsSlice []MutationDetails, 
        mArea string, mCategory string) []testcase.TestCaseDataInfo {
    //--------------------
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    // (4). remove whole payload, i.e. set to: "text" : {} or "text" : null
    plFullPath := "TestCase." + originTcData.TcName() + ".request.payload." + key
    mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, "")
    mutationDetails := MutationDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Remove whole post body"

    tcMutationInfo := getTcMutationInfo(mutationDetails, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mArea, mCategory, "Remove whole post body", mutationInfo, tcMutationInfo)
    
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
                fmt.Println("------> key2, value2: ", key2, value2, reflect.TypeOf(value2))
                // note, to deal with <nil>
                if value2 != nil {
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
                } else {
                    fmt.Println("------> key2, value2: is nil ", key2, value2, reflect.TypeOf(value2))
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
                if value2 != nil {
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
                } else {
                    subPathNew := append(subPath, key2)
                    output := make([]string, len(subPathNew))
                    copy(output, subPathNew)

                    mtD := MutationDetails{output, nil, "", "", []interface{}{}}
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

