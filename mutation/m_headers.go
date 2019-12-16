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
    // "encoding/json"
    "go4api/lib/rands"

    // gjson "github.com/tidwall/gjson"
    sjson "github.com/tidwall/sjson"
)

func (mTd *MTestCaseDataInfo) MRequestHeaders (tcJson []byte) {
    mTd.MSetRequestHeader(tcJson, mFuncs[0])
    mTd.MDelRequestHeader(tcJson, mFuncs[1])
    mTd.MAddRequestHeader(tcJson, mFuncs[2])
    mTd.MDelAllRequestHeaders(tcJson, mFuncs[3])
}

func (mTd *MTestCaseDataInfo) MSetRequestHeader (tcJson []byte, mFunc *MFunc) {
    i := 0
    //
    for key, value := range mTd.Tc4M.TestCase.ReqHeaders() {
        //
        mFd := MFieldDetails {
            FieldPath:     []string{key}, 
            CurrValue:     value, 
            FieldType:     reflect.TypeOf(value).Kind().String(),
            FieldSubType:  "", 
            MutatedValues: []interface{}{},
        }
        mType := mFd.DetermineMutationType()
        mutatedValues := mFd.CallMutationRules(mType)
        //
        for _, mtedValue := range mutatedValues {
            i = i + 1
            mutationInfo := "Update/Set header key: " + key + ", `" + fmt.Sprint(mFd.CurrValue) + "`, `" + 
                fmt.Sprint(mtedValue.MutatedValue) + "`" +
                "\nUsing Mutation Rule: " + mtedValue.MutationRule

            tcMutationInfo := getTcMutationInfo(mFd, mtedValue.MutatedValue)
            
            //-- set new info to mutated tc
            mTcData := getMutatedTcData(tcJson, i, mFunc, mtedValue.MutationRule, mutationInfo, tcMutationInfo)
            mTcData.TestCase.SetRequestHeader(key, fmt.Sprint(mtedValue.MutatedValue))

            mTd.IMMTcs = append(mTd.IMMTcs, &mTcData)
        }
    }
}

func (mTd *MTestCaseDataInfo) MDelRequestHeader (tcJson []byte, mFunc *MFunc) {
    i := 0
    for key, _ := range mTd.Tc4M.TestCase.ReqHeaders() {
        i = i + 1
        mFd := MFieldDetails {
            FieldPath:     []string{key}, 
            CurrValue:     "", 
            FieldType:     reflect.TypeOf("").Kind().String(), 
            FieldSubType:  "", 
            MutatedValues: []interface{}{},
        }
        mutationInfo := "Remove header key: " + "`" + key + "`"

        tcMutationInfo := getTcMutationInfo(mFd, "")

        // del the key
        mTcData := getMutatedTcData(tcJson, i, mFunc, "Remove header key", mutationInfo, tcMutationInfo)
        mTcData.TestCase.DelRequestHeader(key)

        mTd.IMMTcs = append(mTd.IMMTcs, &mTcData)
    }
}

func (mTd *MTestCaseDataInfo) MAddRequestHeader (tcJson []byte, mFunc *MFunc) {
    // add new key: get rand key, get rand value, then Add()
    i := 0

    randKey := rands.RandStringRunes(5)
    randValue := rands.RandStringRunes(5)
    mFd := MFieldDetails {
        FieldPath:     []string{randKey}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
        MutatedValues: []interface{}{},
    }
    mutationInfo := "Add new rand header key: " + " `" + fmt.Sprint(randKey) + "`, `" + fmt.Sprint(randValue) + "`"

    tcMutationInfo := getTcMutationInfo(mFd, randValue)
    //
    mTcData := getMutatedTcData(tcJson, i, mFunc, "Add new rand header key", mutationInfo, tcMutationInfo)
    mTcData.TestCase.AddRequestHeader(randKey, randValue)

    mTd.IMMTcs = append(mTd.IMMTcs, &mTcData)
}


func (mTd *MTestCaseDataInfo) MDelAllRequestHeaders (tcJson []byte, mFunc *MFunc) {
    // remove all headers
    i := 0
  
    hFullPath := "TestCase." + mTd.Tc4M.TcName() + ".request." + "headers"
    mutatedTcJson, _ := sjson.Delete(string(tcJson), hFullPath)
    mFd := MFieldDetails {
        FieldPath:     []string{}, 
        CurrValue:     "", 
        FieldType:     reflect.TypeOf("").Kind().String(), 
        FieldSubType:  "", 
        MutatedValues: []interface{}{},
    }
    mutationInfo := "Remove all headers"

    tcMutationInfo := getTcMutationInfo(mFd, "")
    //
    mTcData := getMutatedTcData([]byte(mutatedTcJson), i, mFunc, "Remove all headers", mutationInfo, tcMutationInfo)
    
    mTd.IMMTcs = append(mTd.IMMTcs, &mTcData)
}



