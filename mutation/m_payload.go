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

// MutateRequestPayload
func (mTd *MTestCaseDataInfo) MRequestPayload () {
    mTd.initTc4MPL()

    for _, value := range mTd.Tc4MPL.TestCase.ReqPayload() {
        switch mTd.TcPlType {
        case "text", "form", "multipartForm":
            mFds := getMFieldsDetails(value)

            if mTd.TcPlType == "multipartForm" {
                mTd.getTypeFileIndex(mFds)
            }

            mTd.MSetRequestPayload(mTd.TcPlType, mFds, mFuncs[8])
            mTd.MDelRequestPayload(mTd.TcPlType, mFds, mFuncs[9])
            mTd.MAddRequestPayloadNode(mTd.TcPlType, mFds, mFuncs[10])
            mTd.MDelWholeRequestPayloadNode(mTd.TcPlType, mFds, mFuncs[11])
        }

        break
    }
}

func (mTd *MTestCaseDataInfo) initTc4MPL () {
    switch mTd.TcPlType {
    case "text", "form", "", "multipartForm":
        var tc4MPL testcase.TestCaseDataInfo
        tj, _ := json.Marshal(mTd.OriginTcD)   
        json.Unmarshal(tj, &tc4MPL)
        mTd.Tc4MPL = &tc4MPL
    case "multipartFormM":
        var pLMPForm PLMPForm

        reqPayload := mTd.OriginTcD.TestCase.ReqPayload()["multipartForm"]
        reqPayloadJsonBytes, _ := json.Marshal(reqPayload)
        json.Unmarshal(reqPayloadJsonBytes, &pLMPForm)
        //
        var tc4MPL testcase.TestCaseDataInfo
        tj, _ := json.Marshal(mTd.OriginTcD)   
        json.Unmarshal(tj, &tc4MPL)
        //
        mTd.Tc4MPL         = &tc4MPL
        mTd.PLMPForm       = pLMPForm
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

func (mTd *MTestCaseDataInfo) getTypeFileIndex (mFds []MFieldDetails) {
    var ii []string
    for _, mFd := range mFds {
        if len(mFd.FieldPath) == 2 {
            if mFd.FieldPath[1] == "type" && mFd.CurrValue == "file" {
                ii = append(ii, mFd.FieldPath[0])
            }
        }
    }

    mTd.MFileIndex = ii
}

func ifItemExists (item string, items []string) bool {
    matched := false
    for _, it := range items {
        if it == item {
            matched = true
            break
        }
    }

    return matched
}

// MSetRequestPayload
func (mTd *MTestCaseDataInfo) MSetRequestPayload (key string, mFds []MFieldDetails, mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.Tc4MPL)
    //
    i := 0
    for _, mFd := range mFds {
        if key == "multipartForm" {
            switch len(mFd.FieldPath) {
            case 0:
                continue
            case 1:
                continue
            case 2:
                if mFd.FieldPath[1] == "name" {
                    continue
                }
                if ifItemExists(mFd.FieldPath[0], mTd.MFileIndex) {
                    continue
                }
            }
        }
        //
        plPath := key + "." + strings.Join(mFd.FieldPath, ".")
        plFullPath := "TestCase." + mTd.Tc4MPL.TcName() + ".request.payload" + "." + plPath
        // 
        mFd.CallMutationRules()
        //
        for _, mtedValue := range mFd.MutatedValues {
            i = i + 1
            
            mtedTcJson, _ := sjson.Set(string(tcJson), plFullPath, mtedValue.MutatedValue)
            // mDJsonByte, _ := json.Marshal(mFd)

            mInfoStr := "FieldPath: " + fmt.Sprint(mFd.FieldPath) + 
                "\nCurrValue: " + "`" + fmt.Sprint(mFd.CurrValue) + "`" + 
                "\nFieldType: " + "`" + fmt.Sprint(mFd.FieldType) + "`"+ 
                "\nFieldSubType: " + "`" + fmt.Sprint(mFd.FieldSubType) + "`" + 
                "\nMutatedValue: " + "`" + fmt.Sprint(mtedValue.MutatedValue) + "`" + 
                "\n=> Using Mutation Rule: " + mtedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFd, mtedValue.MutatedValue)
            //
            mTcData := getMutatedTcData([]byte(mtedTcJson), i, mFunc, 
                mtedValue.MutationRule, mInfoStr, tcMutationInfo)

            mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
            mTd.NextTcPriority = mTd.NextTcPriority + 1
            //
            // if key == "multipartForm" {
            //     mTd.reWrite4MPForm(&mTcData)
            // }

            mTd.MTcDs = append(mTd.MTcDs, &mTcData)
        }
    }
}

// MDelRequestPayload
func (mTd *MTestCaseDataInfo) MDelRequestPayload (key string, mFds []MFieldDetails, mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.Tc4MPL)
    nodePaths, _ := getPayloadNodePaths(mFds)

    aa, _ := json.Marshal(nodePaths)
    fmt.Println("nodePaths: ", string(aa))

    i := 0
    for _, pathStr := range nodePaths {
        if key == "multipartForm" {
            pp := strings.Split(pathStr, ".")
            switch len(pp) {
            case 0:
                continue
            case 2:
                if pp[1] == "name" || pp[1] == "value" {
                    continue
                }
            }
        }
        //
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

    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)

    if key == "multipartForm" {
        m := make(map[string]string)

        m["name"] = randKey
        m["value"] = randValue

        // value, _ := sjson.Set(`{"friends":["Andy","Carol"]}`, "friends.-1", "Sara")
        // sjson.Set(`{"key":true}`, "key", map[string]interface{}{"hello":"world"})

        // pp := strings.Split(pathStr, ".")
        // switch len(pp) {
        // case 0:
        //     continue
        // case 2:
        //     if pp[1] == "name" || pp[1] == "value" {
        //         continue
        //     }
        // }
    } else {
        i := 0
        
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

        mTd.MTcDs = append(mTd.MTcDs, &mTcData)
    }  
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
    if key == "multipartForm" {
        var pLMPForm PLMPForm
        mTcData.TestCase.DelReqPayload("multipartForm")
        mTcData.TestCase.SetRequestPayload("multipartForm", pLMPForm)
    }
    
    mTd.MTcDs = append(mTd.MTcDs, &mTcData)
}
