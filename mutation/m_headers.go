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
    // "strings"
    "reflect"
    "encoding/json"

    "go4api/lib/rands"
    "go4api/lib/testcase"

    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

func (mTd *MTestCaseDataInfo) MRequestHeaders () {
    // not to mutate the Content-Type
    if mTd.HContentType != nil {
        var tc4MH testcase.TestCaseDataInfo
        tj, _ := json.Marshal(mTd.OriginTcD)   
        json.Unmarshal(tj, &tc4MH)

        tc4MH.TestCase.DelRequestHeader("Content-Type")
        mTd.Tc4MH = &tc4MH
    }

    mTd.MSetRequestHeader(mFuncs[0])
    mTd.MDelRequestHeader(mFuncs[1])
    mTd.MAddRequestHeader(mFuncs[2])
    mTd.MDelAllRequestHeaders(mFuncs[3])
}

func (mTd *MTestCaseDataInfo) MSetRequestHeader (mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.Tc4MH)
    //
    i := 0
    for key, value := range mTd.Tc4MH.TestCase.ReqHeaders() {
        //
        mFd := MFieldDetails {
            FieldPath:     []string{key}, 
            CurrValue:     value, 
            FieldType:     reflect.TypeOf(value).Kind().String(),
            FieldSubType:  "", 
        }
        mFd.CallMutationRules()
        //
        for _, mtedValue := range mFd.MutatedValues {
            i = i + 1
            mInfoStr := "Update/Set header key: " + key + ", `" + 
                fmt.Sprint(mFd.CurrValue) + "`, `" + 
                fmt.Sprint(mtedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mtedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFd, mtedValue.MutatedValue)
            
            //-- set new info to mutated tc
            mTcData := getMutatedTcData(tcJson, i, mFunc, mtedValue.MutationRule, mInfoStr, tcMutationInfo)
            mTcData.TestCase.SetRequestHeader(key, fmt.Sprint(mtedValue.MutatedValue))

            mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
            mTd.NextTcPriority = mTd.NextTcPriority + 1
            //
            if mTd.HContentType != nil {
                mTcData.TestCase.SetRequestHeader("Content-Type", fmt.Sprint(mTd.HContentType))
            }

            mTd.MTcDs = append(mTd.MTcDs, &mTcData)
        }
    }
}

func (mTd *MTestCaseDataInfo) MDelRequestHeader (mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.OriginTcD)
    //
    i := 0
    for key, _ := range mTd.OriginTcD.TestCase.ReqHeaders() {
        i = i + 1
        mFd := MFieldDetails {
            FieldPath:     []string{key}, 
            CurrValue:     "", 
            FieldType:     reflect.TypeOf("").Kind().String(), 
            FieldSubType:  "", 
        }
        mInfoStr := "Remove header key: " + "`" + key + "`"

        tcMutationInfo := getTcMutationInfo(mFd, "")

        // del the key
        mTcData := getMutatedTcData(tcJson, i, mFunc, "Remove header key", mInfoStr, tcMutationInfo)
        mTcData.TestCase.DelRequestHeader(key)

        mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
        mTd.NextTcPriority = mTd.NextTcPriority + 1
        //
        if mTd.HContentType != nil {
            mTcData.TestCase.SetRequestHeader("Content-Type", fmt.Sprint(mTd.HContentType))
        }

        mTd.MTcDs = append(mTd.MTcDs, &mTcData)
    }
}

func (mTd *MTestCaseDataInfo) MAddRequestHeader (mFunc *MFunc) {
    // add new key: get rand key, get rand value, then Add()
    tcJson, _ := json.Marshal(mTd.OriginTcD)
    //
    i := 0
    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)
    mFd := MFieldDetails {
        FieldPath:     []string{randKey}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
    }
    mInfoStr := "Add new rand header key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFd, randValue)
    //
    mTcData := getMutatedTcData(tcJson, i, mFunc, "Add new rand header key", mInfoStr, tcMutationInfo)
    mTcData.TestCase.AddRequestHeader(randKey, randValue)

    mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
    mTd.NextTcPriority = mTd.NextTcPriority + 1
    //
    if mTd.HContentType != nil {
        mTcData.TestCase.SetRequestHeader("Content-Type", fmt.Sprint(mTd.HContentType))
    }

    mTd.MTcDs = append(mTd.MTcDs, &mTcData)
}


func (mTd *MTestCaseDataInfo) MDelAllRequestHeaders (mFunc *MFunc) {
    // remove all headers
    tcJson, _ := json.Marshal(mTd.OriginTcD)
    //
    i := 0
    hFullPath := "TestCase." + mTd.OriginTcD.TcName() + ".request." + "headers"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), hFullPath)
    mFd := MFieldDetails {
        FieldPath:     []string{}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
    }
    mInfoStr := "Remove all headers"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, "Remove all headers", mInfoStr, tcMutationInfo)

    mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
    mTd.NextTcPriority = mTd.NextTcPriority + 1
    //
    if mTd.HContentType != nil {
        mTcData.TestCase.AddRequestHeader("Content-Type", fmt.Sprint(mTd.HContentType))
    }
    
    mTd.MTcDs = append(mTd.MTcDs, &mTcData)
}



