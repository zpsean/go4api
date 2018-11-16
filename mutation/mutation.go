/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package mutation

import ( 
    "fmt"
    "strings"
    "reflect"
    "encoding/json"

    "go4api/lib/testcase"
    "go4api/lib/rands"
    "go4api/lib/g4json"

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

var mCategoryFuncSlice []*MCategoryFuncMap

type MFieldDetails struct {
    FieldPath []string
    CurrValue interface{}
    FieldType string // the json supported types
    FieldSubType string  // like ip/email/phone/etc.
    MutatedValues []interface{}
}

type MCategoryFuncMap struct {
    MPriority string
    MTcSuffix string
    MArea string
    MCategory string
    MFunc interface{}
}


func init() {
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"3", "-M-H-S-", "headers", "SetRequestHeader", MutateSetRequestHeader})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"4", "-M-H-D-", "headers", "DelRequestHeader", MutateDelRequestHeader})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"5", "-M-H-A-", "headers", "AddRequestHeader", MutateAddRequestHeader})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"6", "-M-H-D-", "headers", "DelAllRequestHeaders", MutateDelAllRequestHeaders})

    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"7", "-M-QS-S-", "queryString", "SetRequestQueryString", MutateSetRequestQueryString})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"8", "-M-QS-D-", "queryString", "DelRequestQueryString", MutateDelRequestQueryString})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"9", "-M-QS-A-", "queryString", "AddRequestQueryString", MutateAddRequestQueryString})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"10", "-M-QS-D-", "queryString", "DelAllRequestQueryStrings", MutateDelAllRequestQueryStrings})

    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"11", "-M-PL-S-", "payload", "SetRequestPayload", MutateSetRequestPayload})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"12", "-M-PL-D-", "payload", "DelRequestPayload", MutateDelRequestPayload})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"13", "-M-PL-A-", "payload", "AddRequestPayloadNode", MutateAddRequestPayloadNode})
    mCategoryFuncSlice = append(mCategoryFuncSlice, &MCategoryFuncMap{"14", "-M-PL-D-", "payload", "DelWholeRequestPayloadNode", MutateDelWholeRequestPayloadNode})
}


func MutateTcArray(originMutationTcArray []*testcase.TestCaseDataInfo) []*testcase.TestCaseDataInfo {
    var mutatedTcArray []*testcase.TestCaseDataInfo

    for _, originTcData := range originMutationTcArray {
        if originTcData.TestCase.IfGlobalSetUpTestCase() == true {
            // mutatedTcArray = append(mutatedTcArray, originTcData)
            continue
        }

        tcJson, _ := json.Marshal(*originTcData)
        mutatedTcArray = append(mutatedTcArray, originTcData)

        // --- here to start the mutation
        for i, _ := range mCategoryFuncSlice {
            if mCategoryFuncSlice[i].MArea != "payload" {
                f := reflect.ValueOf(mCategoryFuncSlice[i].MFunc)

                in := make([]reflect.Value, 3)
                in[0] = reflect.ValueOf(*originTcData)
                in[1] = reflect.ValueOf(tcJson)
                in[2] = reflect.ValueOf(mCategoryFuncSlice[i])

                result := f.Call(in)

                mTcArray := result[0].Interface().([]testcase.TestCaseDataInfo)
                for ii, _ := range mTcArray {
                    mutatedTcArray = append(mutatedTcArray, &mTcArray[ii])
                }
            }
        }
        // for payload
        mTcArray := MutateRequestPayload(*originTcData, tcJson)
        for ii, _ := range mTcArray {
            mutatedTcArray = append(mutatedTcArray, &mTcArray[ii])
        }
    }
    // aa, _ := json.Marshal(mutatedTcArray)
    // fmt.Println("\nmutatedTcArray: ", string(aa))
    return mutatedTcArray
}

func getMutatedTcData (tcJson []byte, i int, mCategoryFuncMap *MCategoryFuncMap, mutationRule string, mutationInfo string, tcMutationInfo testcase.MutationInfo) testcase.TestCaseDataInfo {
    tcSuffix := mCategoryFuncMap.MTcSuffix + fmt.Sprint(i)

    //-- set new info to mutated tc
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    mTcData.TestCase.SetPriority(mCategoryFuncMap.MPriority)
    mTcData.MutationArea = mCategoryFuncMap.MArea
    mTcData.MutationCategory = mCategoryFuncMap.MCategory
    mTcData.MutationRule = mutationRule
    mTcData.MutationInfoStr = mutationInfo
    mTcData.MutationInfo = tcMutationInfo
    
    return mTcData
}

func getTcMutationInfo (mFieldDetails MFieldDetails, mutatedValue interface{}) testcase.MutationInfo {
    tcMutationInfo := testcase.MutationInfo {
        FieldPath: mFieldDetails.FieldPath,
        CurrValue: mFieldDetails.CurrValue,
        FieldType: mFieldDetails.FieldType,
        FieldSubType: mFieldDetails.FieldSubType,
        MutatedValue: mutatedValue,
    }

    return tcMutationInfo
}

// headers
func MutateSetRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte, mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    for key, value := range originTcData.TestCase.ReqHeaders() {
        //
        mFieldDetails := MFieldDetails{[]string{key}, value, reflect.TypeOf(value).Kind().String(), "", []interface{}{}}
        mType := mFieldDetails.DetermineMutationType()
        mutatedValues := mFieldDetails.CallMutationRules(mType)
        //
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            mutationInfo := "Update/Set header key: " + key + ", `" + fmt.Sprint(mFieldDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mutatedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFieldDetails, mutatedValue.MutatedValue)
            
            //-- set new info to mutated tc
            mTcData := getMutatedTcData(tcJson, i, mCategoryFuncMap, mutatedValue.MutationRule, mutationInfo, tcMutationInfo)
            mTcData.TestCase.SetRequestHeader(key, fmt.Sprint(mutatedValue.MutatedValue))

            mutatedTcArray = append(mutatedTcArray, mTcData)
        }
    }

    return mutatedTcArray
}

func MutateDelRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte, mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    for key, _ := range originTcData.TestCase.ReqHeaders() {
        i = i + 1
        mFieldDetails := MFieldDetails{[]string{key}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
        mutationInfo := "Remove header key: " + "`" + key + "`"

        tcMutationInfo := getTcMutationInfo(mFieldDetails, "")

        // del the key
        mTcData := getMutatedTcData(tcJson, i, mCategoryFuncMap, "Remove header key", mutationInfo, tcMutationInfo)
        mTcData.TestCase.DelRequestHeader(key)

        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}

func MutateAddRequestHeader (originTcData testcase.TestCaseDataInfo, tcJson []byte, mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // add new key: get rand key, get rand value, then Add()
    i := 0

    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)
    mFieldDetails := MFieldDetails{[]string{randKey}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Add new rand header key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFieldDetails, randValue)
    //
    mTcData := getMutatedTcData(tcJson, i, mCategoryFuncMap, "Add new rand header key", mutationInfo, tcMutationInfo)
    mTcData.TestCase.AddRequestHeader(randKey, randValue)

    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}


func MutateDelAllRequestHeaders (originTcData testcase.TestCaseDataInfo, tcJson []byte, mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // remove all headers
    i := 0
  
    hFullPath := "TestCase." + originTcData.TcName() + ".request." + "headers"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), hFullPath)
    mFieldDetails := MFieldDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Remove all headers"

    tcMutationInfo := getTcMutationInfo(mFieldDetails, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mCategoryFuncMap, "Remove all headers", mutationInfo, tcMutationInfo)
    
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}


// QueryString
func MutateSetRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte, mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    for key, value := range originTcData.TestCase.ReqQueryString() {
        mFieldDetails := MFieldDetails{[]string{key}, value, reflect.TypeOf(value).Kind().String(), "", []interface{}{}}
        mType := mFieldDetails.DetermineMutationType()
        mutatedValues := mFieldDetails.CallMutationRules(mType)
        // loop and mutate the value, set new value to key
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            mutationInfo := "Update/Set header key: " + key + ", `" + fmt.Sprint(mFieldDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mutatedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFieldDetails, mutatedValue.MutatedValue)

            //-- set new info to mutated tc
            mTcData := getMutatedTcData(tcJson, i, mCategoryFuncMap, mutatedValue.MutationRule, mutationInfo, tcMutationInfo)
            mTcData.TestCase.SetRequestQueryString(key, fmt.Sprint(mutatedValue.MutatedValue))

            mutatedTcArray = append(mutatedTcArray, mTcData)
        }
    }

    return mutatedTcArray
}


func MutateDelRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte, mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    for key, _ := range originTcData.TestCase.ReqQueryString() {
        // del key
        i = i + 1
        mFieldDetails := MFieldDetails{[]string{key}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
        mutationInfo := "Remove querystring key: " + "`" + key + "`"

        tcMutationInfo := getTcMutationInfo(mFieldDetails, "")

        // del the key
        mTcData := getMutatedTcData(tcJson, i, mCategoryFuncMap, "Remove querystring key", mutationInfo, tcMutationInfo)
        mTcData.TestCase.DelRequestQueryString(key)

        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}


func MutateAddRequestQueryString (originTcData testcase.TestCaseDataInfo, tcJson []byte, mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // add new key: get rand key, get rand value, then Add()
    i := 0

    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)

    mFieldDetails := MFieldDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Add new rand QueryString key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFieldDetails, "")
    //
    mTcData := getMutatedTcData(tcJson, i, mCategoryFuncMap, "Add new rand querystring key", mutationInfo, tcMutationInfo)
    mTcData.TestCase.AddRequestQueryString(randKey, randValue)

    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}

func MutateDelAllRequestQueryStrings (originTcData testcase.TestCaseDataInfo, tcJson []byte, mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    // remove all querystring
    i := 0
 
    qSFullPath := "TestCase." + originTcData.TcName() + ".request." + "queryString"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), qSFullPath)
    mFieldDetails := MFieldDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Remove all querystring"

    tcMutationInfo := getTcMutationInfo(mFieldDetails, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mCategoryFuncMap, "Remove all querystring", mutationInfo, tcMutationInfo)
    
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}

// MutateRequestPayload
func MutateRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo
    for key, value := range originTcData.TestCase.ReqPayload() {
        if key == "text" {
            // to loop over the struct
            var mFieldDetailsSlice []MFieldDetails
            mFieldDetailsSlice = getFieldsMutationDetails(value)
            //
            for i, _ := range mCategoryFuncSlice {
                if mCategoryFuncSlice[i].MArea == "payload" {
                    f := reflect.ValueOf(mCategoryFuncSlice[i].MFunc)

                    in := make([]reflect.Value, 5)
                    in[0] = reflect.ValueOf(originTcData)
                    in[1] = reflect.ValueOf(tcJson)
                    in[2] = reflect.ValueOf(key)
                    in[3] = reflect.ValueOf(mFieldDetailsSlice)
                    in[4] = reflect.ValueOf(mCategoryFuncSlice[i])

                    result := f.Call(in)
                    mTcArray := result[0].Interface().([]testcase.TestCaseDataInfo)

                    mutatedTcArray = append(mutatedTcArray, mTcArray[0:]...)
                }
            }
        }
    }
    return mutatedTcArray
}

func MutateSetRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mFieldDetailsSlice []MFieldDetails, 
        mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    //--------------------
    var mutatedTcArray []testcase.TestCaseDataInfo
    // (1), set node value
    i := 0
    for _, mFieldDetails := range mFieldDetailsSlice {
        // set the value
        plPath := key + "." + strings.Join(mFieldDetails.FieldPath, ".")
        plFullPath := "TestCase." + originTcData.TcName() + ".request.payload" + "." + plPath
        // mutate the value based on rules 
        mType := mFieldDetails.DetermineMutationType()
        mutatedValues := mFieldDetails.CallMutationRules(mType)
        // (1). set node
        for _, mutatedValue := range mutatedValues {
            i = i + 1
     
            mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, mutatedValue.MutatedValue)
            mDJsonByte, _ := json.Marshal(mFieldDetails)
            mutationInfo := fmt.Sprint(string(mDJsonByte)) + ", `" + fmt.Sprint(mFieldDetails.CurrValue) + "`, `" + fmt.Sprint(mutatedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mutatedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFieldDetails, mutatedValue.MutatedValue)
            //
            mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mCategoryFuncMap, mutatedValue.MutationRule, mutationInfo, tcMutationInfo)
    
            mutatedTcArray = append(mutatedTcArray, mTcData)
        }

        // add new node
    }
    return mutatedTcArray
}

func getPayloadNodePaths (mFieldDetailsSlice []MFieldDetails) ([]string, int) {
    // get the max level of the paths
    max := 0
    for _, mFieldDetails := range mFieldDetailsSlice {
        if len(mFieldDetails.FieldPath) > max {
            max = len(mFieldDetails.FieldPath)
        }
    }
    // 
    var nodePaths []string
    for i := max; i > 0; i-- {
        for _, mFieldDetails := range mFieldDetailsSlice {
            if len(mFieldDetails.FieldPath) >= i {
                nodePathStr := strings.Join(mFieldDetails.FieldPath[0:i], ".")

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


func MutateDelRequestPayload (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mFieldDetailsSlice []MFieldDetails, 
        mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    //--------------------
    var mutatedTcArray []testcase.TestCaseDataInfo
    nodePaths, _ := getPayloadNodePaths(mFieldDetailsSlice)

    i := 0
    for _, pathStr := range nodePaths {
        i = i + 1
    
        plPath := key + "." + pathStr
        plFullPath := "TestCase." + originTcData.TcName() + ".request.payload" + "." + plPath

        // (2). del node
        mutatedTcJson, _ := sjson.Delete(string(tcJson), plFullPath)
        mFieldDetails := MFieldDetails{[]string{pathStr}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
        mutationInfo := "Remove payload value on node: " + pathStr

        tcMutationInfo := getTcMutationInfo(mFieldDetails, "")
        //
        mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mCategoryFuncMap, "Remove payload value on node", mutationInfo, tcMutationInfo)
    
        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}


func MutateAddRequestPayloadNode (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mFieldDetailsSlice []MFieldDetails, 
        mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    //--------------------
    var mutatedTcArray []testcase.TestCaseDataInfo
    // (3). add new node, for each node level
    i := 0

    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)

    // set the value
    plPath := key + "." + randKey
    plFullPath := "TestCase." + originTcData.TcName() + ".request.payload" + "." + plPath

    mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, randValue)

    mFieldDetails := MFieldDetails{[]string{randKey}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Add new rand payload key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFieldDetails, "")
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mCategoryFuncMap, "Add new rand payload key", mutationInfo, tcMutationInfo)

    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}

func MutateDelWholeRequestPayloadNode (originTcData testcase.TestCaseDataInfo, tcJson []byte, key string, mFieldDetailsSlice []MFieldDetails, 
        mCategoryFuncMap *MCategoryFuncMap) []testcase.TestCaseDataInfo {
    //--------------------
    var mutatedTcArray []testcase.TestCaseDataInfo
    i := 0
    // (4). remove whole payload, i.e. set to: "text" : {} or "text" : null
    plFullPath := "TestCase." + originTcData.TcName() + ".request.payload." + key
    mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, "")
    mFieldDetails := MFieldDetails{[]string{}, "", reflect.TypeOf("").Kind().String(), "", []interface{}{}}
    mutationInfo := "Remove whole post body"

    tcMutationInfo := getTcMutationInfo(mFieldDetails, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mCategoryFuncMap, "Remove whole post body", mutationInfo, tcMutationInfo)
    
    mutatedTcArray = append(mutatedTcArray, mTcData)

    return mutatedTcArray
}


func getFieldsMutationDetails(value interface{}) []MFieldDetails {
    var mFieldDetailsSlice []MFieldDetails

    fieldDetailsSlice := g4json.GetFieldsDetails(value)

    for i, _ := range fieldDetailsSlice {
        mFieldDetails := MFieldDetails {
            FieldPath: fieldDetailsSlice[i].FieldPath,
            CurrValue: fieldDetailsSlice[i].CurrValue,
            FieldType: fieldDetailsSlice[i].FieldType,
            FieldSubType: fieldDetailsSlice[i].FieldSubType,
            MutatedValues: []interface{}{},
        }

        mFieldDetailsSlice = append(mFieldDetailsSlice, mFieldDetails)
    }

    return mFieldDetailsSlice
}


