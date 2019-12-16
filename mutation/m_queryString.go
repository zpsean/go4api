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

    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

func (mTd *MTestCaseDataInfo) MRequestQueryString () {
    tc4MQS := *mTd.OriginTcD
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
            MutatedValues: []interface{}{},
        }
        mType := mFd.DetermineMutationType()
        mutatedValues := mFd.CallMutationRules(mType)
        // loop and mutate the value, set new value to key
        for _, mutatedValue := range mutatedValues {
            i = i + 1
            mutationInfo := "Update/Set header key: " + key + ", `" + fmt.Sprint(mFd.CurrValue) + "`, `" + 
                fmt.Sprint(mutatedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mutatedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFd, mutatedValue.MutatedValue)

            //-- set new info to mutated tc
            mTcData := getMutatedTcData(tcJson, i, mFunc, mutatedValue.MutationRule, mutationInfo, tcMutationInfo)
            mTcData.TestCase.SetRequestQueryString(key, fmt.Sprint(mutatedValue.MutatedValue))

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
            MutatedValues: []interface{}{},
        }
        mutationInfo := "Remove querystring key: " + "`" + key + "`"

        tcMutationInfo := getTcMutationInfo(mFd, "")

        // del the key
        mTcData := getMutatedTcData(tcJson, i, mFunc, "Remove querystring key", mutationInfo, tcMutationInfo)
        mTcData.TestCase.DelRequestQueryString(key)

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
        MutatedValues: []interface{}{},
    }
    mutationInfo := "Add new rand QueryString key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    //
    mTcData := getMutatedTcData(tcJson, i, mFunc, "Add new rand querystring key", mutationInfo, tcMutationInfo)
    mTcData.TestCase.AddRequestQueryString(randKey, randValue)

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
        MutatedValues: []interface{}{},
    }
    mutationInfo := "Remove all querystring"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, "Remove all querystring", mutationInfo, tcMutationInfo)
    
    mTd.MTcDs = append(mTd.MTcDs, &mTcData)
}


