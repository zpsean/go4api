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

func (mTd *MTestCaseDataInfo) MRequestQueryString () {
    var tc4MQS testcase.TestCaseDataInfo
    tj, _ := json.Marshal(mTd.OriginTcD)   
    json.Unmarshal(tj, &tc4MQS)

    mTd.Tc4MQS = &tc4MQS

    mTd.MSetRequestQueryString(mFuncs[4])
    mTd.MDelRequestQueryString(mFuncs[5])
    mTd.MAddRequestQueryString(mFuncs[6])
    mTd.MDelAllRequestQueryStrings(mFuncs[7])
}

func (mTd *MTestCaseDataInfo) MSetRequestQueryString (mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.Tc4MQS)

    i := 0
    for key, value := range mTd.Tc4MQS.TestCase.ReqQueryString() {
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
            mTcData.TestCase.SetRequestQueryString(key, fmt.Sprint(mtedValue.MutatedValue))

            mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
            mTd.NextTcPriority = mTd.NextTcPriority + 1

            mTd.MTcDs = append(mTd.MTcDs, &mTcData)
        }
    }
}


func (mTd *MTestCaseDataInfo) MDelRequestQueryString (mFunc *MFunc) {
    tcJson, _ := json.Marshal(mTd.Tc4MQS)
    i := 0
    for key, _ := range mTd.Tc4MQS.TestCase.ReqQueryString() {
        // del key
        i = i + 1
        mFd := MFieldDetails {
            FieldPath:     []string{key}, 
            CurrValue:     "", 
            FieldType:     reflect.TypeOf("").Kind().String(), 
            FieldSubType:  "", 
        }
        mInfoStr := "Remove querystring key: " + "`" + key + "`"

        tcMutationInfo := getTcMutationInfo(mFd, "")

        // del the key
        mTcData := getMutatedTcData(tcJson, i, mFunc, "Remove querystring key", mInfoStr, tcMutationInfo)
        mTcData.TestCase.DelRequestQueryString(key)

        mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
        mTd.NextTcPriority = mTd.NextTcPriority + 1

        mTd.MTcDs = append(mTd.MTcDs, &mTcData)
    }
}


func (mTd *MTestCaseDataInfo) MAddRequestQueryString (mFunc *MFunc) {
    // add new key: get rand key, get rand value, then Add()
    tcJson, _ := json.Marshal(mTd.Tc4MQS)
    i := 0

    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)

    mFd := MFieldDetails {
        FieldPath:     []string{}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
    }
    mInfoStr := "Add new rand QueryString key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    //
    mTcData := getMutatedTcData(tcJson, i, mFunc, "Add new rand querystring key", mInfoStr, tcMutationInfo)
    mTcData.TestCase.AddRequestQueryString(randKey, randValue)

    mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
    mTd.NextTcPriority = mTd.NextTcPriority + 1

    mTd.MTcDs = append(mTd.MTcDs, &mTcData)
}

func (mTd *MTestCaseDataInfo) MDelAllRequestQueryStrings (mFunc *MFunc) {
    // remove all querystring
    tcJson, _ := json.Marshal(mTd.Tc4MQS)
    i := 0
 
    qSFullPath := "TestCase." + mTd.Tc4MQS.TcName() + ".request." + "queryString"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), qSFullPath)
    mFd := MFieldDetails {
        FieldPath:     []string{}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
    }
    mInfoStr := "Remove all querystring"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, "Remove all querystring", mInfoStr, tcMutationInfo)

    mTcData.TestCase.SetPriority(fmt.Sprint(mTd.NextTcPriority))
    mTd.NextTcPriority = mTd.NextTcPriority + 1
    
    mTd.MTcDs = append(mTd.MTcDs, &mTcData)
}


