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
    "fmt"
    "strings"
    "reflect"
    "encoding/json"

    // "go4api/lib/testcase"
    "go4api/lib/rands"
    "go4api/lib/g4json"

    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

// -- payload has two type of format
// -- 1: for header: "Content-Type": "multipart/form-data"
// "payload": {
//   "multipart-form": [
//           {
//             "name": "name",
//             "value": "zp-c-01"
//           },
//           {
//             "name": "cover",
//             "value": "banner3-符合.jpeg",
//             "type": "file"
//           }
//         ]
// }
//
// -- 2: for header: "Content-Type": "application/json;charset=UTF-8"
// "payload": {
//   "text": {
//             "username":"${user}",
//             "password":"${password}"
//           }
// }
// -- 3: "Content-Type": application/x-www-form-urlencoded


// MutateRequestPayload
func (mTd *MTestCaseDataInfo) MRequestPayload (tcJson []byte) {
    for key, value := range mTd.Tc4M.TestCase.ReqPayload() {
        lKey := strings.ToLower(key)
        switch lKey {
        case "text", "form":
            mFds := getMFieldsDetails(value)

            mTd.MSetRequestPayload(tcJson, lKey, mFds, mFuncs[8])
            mTd.MDelRequestPayload(tcJson, lKey, mFds, mFuncs[9])
            mTd.MAddRequestPayloadNode(tcJson, lKey, mFds, mFuncs[10])
            mTd.MDelWholeRequestPayloadNode(tcJson, lKey, mFds, mFuncs[11])
        }
    }
}

func getMFieldsDetails(value interface{}) []MFieldDetails {
    var mFds []MFieldDetails

    fdSlice := g4json.GetFieldsDetails(value)

    for i, _ := range fdSlice {
        mFd := MFieldDetails {
            FieldPath:     fdSlice[i].FieldPath,
            CurrValue:     fdSlice[i].CurrValue,
            FieldType:     fdSlice[i].FieldType,
            FieldSubType:  fdSlice[i].FieldSubType,
            MutatedValues: []interface{}{},
        }
        mFds = append(mFds, mFd)
    }

    return mFds
}

// MSetRequestPayload
func (mTd *MTestCaseDataInfo) MSetRequestPayload (tcJson []byte, key string, mFds []MFieldDetails, mFunc *MFunc) {
    //--------------------
    // (1), set node value
    i := 0
    for _, mFd := range mFds {
        // set the value
        plPath := key + "." + strings.Join(mFd.FieldPath, ".")
        plFullPath := "TestCase." + mTd.Tc4M.TcName() + ".request.payload" + "." + plPath
        // mutate the value based on rules 
        mType := mFd.DetermineMutationType()
        mtedValues := mFd.CallMutationRules(mType)
        // (1). set node
        for _, mtedValue := range mtedValues {
            i = i + 1
     
            mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, mtedValue.MutatedValue)
            mDJsonByte, _ := json.Marshal(mFd)
            mutationInfo := fmt.Sprint(string(mDJsonByte)) + ", `" + fmt.Sprint(mFd.CurrValue) + "`, `" + 
                fmt.Sprint(mtedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mtedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFd, mtedValue.MutatedValue)
            //
            mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, mtedValue.MutationRule, mutationInfo, tcMutationInfo)
    
           mTd.IMMTcs = append(mTd.IMMTcs, &mTcData)
        }

        // add new node
    }
}

// MDelRequestPayload
func (mTd *MTestCaseDataInfo) MDelRequestPayload (tcJson []byte, key string, mFds []MFieldDetails, mFunc *MFunc) {
    //--------------------
    nodePaths, _ := getPayloadNodePaths(mFds)

    i := 0
    for _, pathStr := range nodePaths {
        i = i + 1
    
        plPath := key + "." + pathStr
        plFullPath := "TestCase." + mTd.Tc4M.TcName() + ".request.payload" + "." + plPath

        // (2). del node
        mutatedTcJson, _ := sjson.Delete(string(tcJson), plFullPath)
        mFd := MFieldDetails {
            FieldPath:     []string{pathStr}, 
            CurrValue:     "", 
            FieldType:     reflect.TypeOf("").Kind().String(), 
            FieldSubType:  "", 
            MutatedValues: []interface{}{},
        }
        mutationInfo := "Remove payload value on node: " + pathStr

        tcMutationInfo := getTcMutationInfo(mFd, "")
        //
        mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, "Remove payload value on node", mutationInfo, tcMutationInfo)
    
        mTd.IMMTcs = append(mTd.IMMTcs, &mTcData)
    }
}

func getPayloadNodePaths (mFds []MFieldDetails) ([]string, int) {
    // get the max level of the paths
    max := 0
    for _, mFd := range mFds {
        if len(mFd.FieldPath) > max {
            max = len(mFd.FieldPath)
        }
    }
    // 
    var nodePaths []string
    for i := max; i > 0; i-- {
        for _, mFd := range mFds {
            if len(mFd.FieldPath) >= i {
                nodePathStr := strings.Join(mFd.FieldPath[0:i], ".")

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

//
func (mTd *MTestCaseDataInfo) MAddRequestPayloadNode (tcJson []byte, key string, mFds []MFieldDetails, mFunc *MFunc) {
    //--------------------
    // (3). add new node, for each node level
    i := 0

    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)

    // set the value
    plPath := key + "." + randKey
    plFullPath := "TestCase." + mTd.Tc4M.TcName() + ".request.payload" + "." + plPath

    mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, randValue)

    mFd := MFieldDetails {
        FieldPath:     []string{randKey}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
        MutatedValues: []interface{}{},
    }
    mutationInfo := "Add new rand payload key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, "Add new rand payload key", mutationInfo, tcMutationInfo)

    mTd.IMMTcs = append(mTd.IMMTcs, &mTcData)
}

func (mTd *MTestCaseDataInfo) MDelWholeRequestPayloadNode (tcJson []byte, key string, mFds []MFieldDetails, mFunc *MFunc) {
    //--------------------
    i := 0
    // (4). remove whole payload, i.e. set to: "text" : {} or "text" : null
    plFullPath := "TestCase." + mTd.Tc4M.TcName() + ".request.payload." + key
    mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, "")
    mFd := MFieldDetails {
        FieldPath:     []string{}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
        MutatedValues: []interface{}{},
    }
    mutationInfo := "Remove whole post body"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, "Remove whole post body", mutationInfo, tcMutationInfo)
    
    mTd.IMMTcs = append(mTd.IMMTcs, &mTcData)
}

