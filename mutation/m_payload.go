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

    "go4api/lib/testcase"
    "go4api/lib/rands"
    "go4api/lib/g4json"

    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

// -- payload types:
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
// "payload": {
//   "form": {
//             "username":"${user}",
//             "password":"${password}"
//           }
// }

// MutateRequestPayload
func (mTd *MTestCaseDataInfo) MRequestPayload () {
    mTd.initTc4MPL()

    for _, value := range mTd.Tc4MPL.TestCase.ReqPayload() {
        switch mTd.TcPlType {
        case "text", "form", "multipart-form":
            mFds := getMFieldsDetails(value)

            mTd.MSetRequestPayload(mTd.TcPlType, mFds, mFuncs[8])
            // mTd.MDelRequestPayload(mTd.TcPlType, mFds, mFuncs[9])
            // mTd.MAddRequestPayloadNode(mTd.TcPlType, mFds, mFuncs[10])
            // mTd.MDelWholeRequestPayloadNode(mTd.TcPlType, mFds, mFuncs[11])
            //
            if mTd.TcPlType == "multipart-form" {
                mTd.MDelRequestPayloadMPFile(mTd.TcPlType, mFds, mFuncs[12])
            }
        }

        break
    }
}

func (mTd *MTestCaseDataInfo) initTc4MPL () {
    switch mTd.TcPlType {
    case "text", "form", "" :
        var tc4MPL testcase.TestCaseDataInfo
        tj, _ := json.Marshal(mTd.OriginTcD)   
        json.Unmarshal(tj, &tc4MPL)
        mTd.Tc4MPL = &tc4MPL
    case "multipart-form":
        var pLMPForm PLMPForm
        var pLnf PLMPForm
        var pLf PLMPForm

        reqPayload := mTd.OriginTcD.TestCase.ReqPayload()["multipart-form"]
        reqPayloadJsonBytes, _ := json.Marshal(reqPayload)
        json.Unmarshal(reqPayloadJsonBytes, &pLMPForm)

        for i, _ := range pLMPForm {
            if len(pLMPForm[i].Type) == 0 {
                pLnf = append(pLnf, pLMPForm[i])
            } else {
                pLf = append(pLf, pLMPForm[i])
            }
        }
        //
        var pl = make(map[string]string)
        for i, _ := range pLnf {
            pl[pLnf[i].Name] = pLnf[i].Value
        }
        //
        var tc4MPL testcase.TestCaseDataInfo
        tj, _ := json.Marshal(mTd.OriginTcD)   
        json.Unmarshal(tj, &tc4MPL)

        tc4MPL.TestCase.DelReqPayload("multipart-form")
        tc4MPL.TestCase.SetRequestPayload("multipart-form", pl)
        //
        mTd.Tc4MPL         = &tc4MPL
        mTd.PLMPForm       = pLMPForm
        mTd.PLMPFormNoFile = pLnf
        mTd.PLMPFormFile   = pLf
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
        }
        mFds = append(mFds, mFd)
    }

    return mFds
}

// MSetRequestPayload
func (mTd *MTestCaseDataInfo) MSetRequestPayload (key string, mFds []MFieldDetails, mFunc *MFunc) {
    
    //
    i := 0
    for _, mFd := range mFds {
        plPath := key + "." + strings.Join(mFd.FieldPath, ".")
        plFullPath := "TestCase." + mTd.Tc4MPL.TcName() + ".request.payload" + "." + plPath
        // 
        mFd.CallMutationRules()
        //
        for _, mtedValue := range mFd.MutatedValues {
            i = i + 1
            
            tcJson, _ := json.Marshal(mTd.Tc4MPL)
            mtedTcJson, _ := sjson.Set(string(tcJson), plFullPath, mtedValue.MutatedValue)
            // mDJsonByte, _ := json.Marshal(mFd)

            mInfoStr := fmt.Sprint(mFd.FieldPath) + ", `" + 
                fmt.Sprint(mFd.CurrValue) + "`, `" + 
                fmt.Sprint(mFd.FieldType) + "`, `" + 
                fmt.Sprint(mFd.FieldSubType) + "`, `" + 
                fmt.Sprint(mtedValue.MutatedValue) + "`" + 
                "\n=> Using Mutation Rule: " + mtedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFd, mtedValue.MutatedValue)
            //
            mTcData := getMutatedTcData([]byte(mtedTcJson), i, mFunc, 
                mtedValue.MutationRule, mInfoStr, tcMutationInfo)

            mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
            mTd.NextTcPriority = mTd.NextTcPriority + 1
            //
            if key == "multipart-form" {
                fmt.Println("--->>>>: ", "reached here")
                mTd.reWrite4MPForm(&mTcData)
            }

            tt, _ := json.Marshal(mTcData)
            fmt.Println("mInfoStr: ", key, mInfoStr)
            fmt.Println("--->: ", string(tt))
            fmt.Println("")
            mTd.MTcDs = append(mTd.MTcDs, &mTcData)
        }
    }
}

// MDelRequestPayload
func (mTd *MTestCaseDataInfo) MDelRequestPayload (key string, mFds []MFieldDetails, mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.Tc4MPL)
    nodePaths, _ := getPayloadNodePaths(mFds)

    i := 0
    for _, pathStr := range nodePaths {
        i = i + 1
    
        plPath := key + "." + pathStr
        plFullPath := "TestCase." + mTd.Tc4MPL.TcName() + ".request.payload" + "." + plPath

        // (2). del node
        mutatedTcJson, _ := sjson.Delete(string(tcJson), plFullPath)
        mFd := MFieldDetails {
            FieldPath:     []string{pathStr}, 
            CurrValue:     "", 
            FieldType:     reflect.TypeOf("").Kind().String(), 
            FieldSubType:  "", 
        }
        mInfoStr := "Remove payload value on node: " + pathStr

        tcMutationInfo := getTcMutationInfo(mFd, "")
        //
        mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, 
            "Remove payload value on node", mInfoStr, tcMutationInfo)

        mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
        mTd.NextTcPriority = mTd.NextTcPriority + 1
        //
        if key == "multipart-form" {
            mTd.reWrite4MPForm(&mTcData)
        }
    
        mTd.MTcDs = append(mTd.MTcDs, &mTcData)
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
func (mTd *MTestCaseDataInfo) MAddRequestPayloadNode (key string, mFds []MFieldDetails, mFunc *MFunc) {
    // (3). add new node, for each node level
    tcJson, _ := json.Marshal(mTd.Tc4MPL)

    i := 0

    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)

    // set the value
    plPath := key + "." + randKey
    plFullPath := "TestCase." + mTd.Tc4MPL.TcName() + ".request.payload" + "." + plPath

    mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, randValue)

    mFd := MFieldDetails {
        FieldPath:     []string{randKey}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
    }
    mInfoStr := "Add new rand payload key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, 
        "Add new rand payload key", mInfoStr, tcMutationInfo)

    mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
    mTd.NextTcPriority = mTd.NextTcPriority + 1
    //
    if key == "multipart-form" {
        mTd.reWrite4MPForm(&mTcData)
    }

    mTd.MTcDs = append(mTd.MTcDs, &mTcData)
}

func (mTd *MTestCaseDataInfo) MDelWholeRequestPayloadNode (key string, mFds []MFieldDetails, mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.Tc4MPL)

    i := 0
    // (4). remove whole payload, i.e. set to: "text" : {} or "text" : null
    plFullPath := "TestCase." + mTd.Tc4MPL.TcName() + ".request.payload." + key
    mutatedTcJson, _ := sjson.Set(string(tcJson), plFullPath, "")
    mFd := MFieldDetails {
        FieldPath:     []string{}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
    }
    mInfoStr := "Remove whole post body"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, 
        "Remove whole post body", mInfoStr, tcMutationInfo)

    mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
    mTd.NextTcPriority = mTd.NextTcPriority + 1
    //
    if key == "multipart-form" {
        var pLMPForm PLMPForm
        mTcData.TestCase.DelReqPayload("multipart-form")
        mTcData.TestCase.SetRequestPayload("multipart-form", pLMPForm)
    }
    
    mTd.MTcDs = append(mTd.MTcDs, &mTcData)
}

// MDelRequestPayloadMPFile, especially for multipart-form
// this is to remove the file node one by one for multipart-form
func (mTd *MTestCaseDataInfo) MDelRequestPayloadMPFile (key string, mFds []MFieldDetails, mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.Tc4MPL)

    i := 0

    for ii, v := range mTd.PLMPFormFile {
        i = i + 1
        
        mFd := MFieldDetails {
            FieldPath:     []string{v.Name}, 
            CurrValue:     v.Value, 
            FieldType:     reflect.TypeOf("").Kind().String(), 
            FieldSubType:  "", 
        } 
        //
        var pLMPForm PLMPForm
        for jj, vj := range mTd.PLMPFormFile {
            if ii != jj {
                pLMPForm = append(pLMPForm, vj)
            }
        }

        mInfoStr := "Remove one file field: " + " `" + fmt.Sprint(v.Name) + 
            "`, `" + fmt.Sprint(v.Value) + "`"
        tcMutationInfo := getTcMutationInfo(mFd, "")

        mTcData := getMutatedTcData([]byte(tcJson), i, mFunc, 
            "Remove at lease one of the file field", mInfoStr, tcMutationInfo)

        mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
        mTd.NextTcPriority = mTd.NextTcPriority + 1
        //
        mTcData.TestCase.DelReqPayload("multipart-form")
        mTcData.TestCase.SetRequestPayload("multipart-form", pLMPForm)
        //
        mTd.MTcDs = append(mTd.MTcDs, &mTcData)
    }    
}

//
func (mTd *MTestCaseDataInfo) reWrite4MPForm (mTcData *testcase.TestCaseDataInfo) {
    var pLMPForm PLMPForm
    mPl := mTcData.TestCase.ReqPayload()["multipart-form"]
    for k, v := range mPl.(map[string]interface{}) {
        mPForm := MPForm {
            Name:  k,
            Value: fmt.Sprint(v),
        }
        pLMPForm = append(pLMPForm, &mPForm)
    }
    pLMPForm = append(pLMPForm, mTd.PLMPFormFile...)
    
    mTcData.TestCase.DelReqPayload("multipart-form")
    mTcData.TestCase.SetRequestPayload("multipart-form", pLMPForm)

    mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
    mTd.NextTcPriority = mTd.NextTcPriority + 1
}


