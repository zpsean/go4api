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
    "encoding/json"

    "go4api/lib/testcase"
)

// mutation is to mutate the valid data to working api, see if mutated invalid data still can be handled by the api
// two ways to mutate the testcase:
// Option 1: 
// copy the underlying fields and values to another TestCaseDataInfo, with mutation(s)
// the better way would be deep copy the TestCaseDataInfo, and change the value, but Golang standard
// Lib has no deepcopy, so that, here uses a plain way, that is, to re-sturct the TestCaseDataInfo
//
// Option 2:
// json.Marshal the tc in originTcArray, 
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

type MTestCaseDataInfo struct {
    OriginTcD      *testcase.TestCaseDataInfo
    HContentType   interface{}    // Content-Type: multipart/form-data
    Tc4MH          *testcase.TestCaseDataInfo  // for mt headers
    Tc4MQS         *testcase.TestCaseDataInfo  // for mt queryString
    Tc4MPL         *testcase.TestCaseDataInfo  // if multipartForm, payload without file
    TcPlType       string    // text, multipartForm, form, <others>
    PLMPForm       PLMPForm  // if multipartForm
    MFileIndex     []string
    MTcDs          []*testcase.TestCaseDataInfo  // mutated tcds
    NextTcPriority int
}

//
type MFieldDetails struct {
    FieldPath     []string
    CurrValue     interface{}
    FieldType     string  // the json supported types
    FieldSubType  string  // like ip/email/phone/etc.
    MutationType  string  // inner code like, MChar, MCharNumeric, etc.
    MutatedValues []*MutatedValue
}

type MutatedValue struct {
    MutationRule string
    MutatedValue interface{}
}

//
type MFunc struct {
    MPriority string
    MTcSuffix string
    MArea     string
    MCategory string
    // MFunc interface{}
}

//
type PLMPForm []*MPForm

type MPForm struct {
    Name        string                 `json:"name"`
    Value       string                 `json:"value"`
    Type        string                 `json:"type"`
    MIMEHeader  map[string]interface{} `json:"mIMEHeader"`
}

//
var mFuncs []*MFunc

//
func init() {
    mFuncs = append(mFuncs, &MFunc{"3", "-M-H-S-", "headers", "SetRequestHeader"})
    mFuncs = append(mFuncs, &MFunc{"4", "-M-H-D-", "headers", "DelRequestHeader"})
    mFuncs = append(mFuncs, &MFunc{"5", "-M-H-A-", "headers", "AddRequestHeader"})
    mFuncs = append(mFuncs, &MFunc{"6", "-M-H-D-", "headers", "DelAllRequestHeaders"})

    mFuncs = append(mFuncs, &MFunc{"7", "-M-QS-S-", "queryString", "SetRequestQueryString"})
    mFuncs = append(mFuncs, &MFunc{"8", "-M-QS-D-", "queryString", "DelRequestQueryString"})
    mFuncs = append(mFuncs, &MFunc{"9", "-M-QS-A-", "queryString", "AddRequestQueryString"})
    mFuncs = append(mFuncs, &MFunc{"10", "-M-QS-D-", "queryString", "DelAllRequestQueryStrings"})

    mFuncs = append(mFuncs, &MFunc{"11", "-M-PL-S-", "payload", "SetRequestPayload"})
    mFuncs = append(mFuncs, &MFunc{"12", "-M-PL-D-", "payload", "DelRequestPayload"})
    mFuncs = append(mFuncs, &MFunc{"13", "-M-PL-A-", "payload", "AddRequestPayloadNode"})
    mFuncs = append(mFuncs, &MFunc{"14", "-M-PL-D-", "payload", "DelWholeRequestPayloadNode"})
    mFuncs = append(mFuncs, &MFunc{"16", "-M-PL-D-F-", "payload", "MDelRequestPayloadMPFile"})
}

func MutateTcArray (originTcArray []*testcase.TestCaseDataInfo) []*testcase.TestCaseDataInfo {
    var mutatedTcArray []*testcase.TestCaseDataInfo

    for _, originTcData := range originTcArray {
        if originTcData.TestCase.IfGlobalSetUpTestCase() == true {
            continue
        }

        originTcData.TestCase.SetPriority(fmt.Sprint(1))
        // mutatedTcArray = append(mutatedTcArray, originTcData)
        // json, originTcData, multipartForm, form
        mTd := InitMTc(originTcData)

        mTd.NextTcPriority = 2
        // --- here to start the mutation
        // mTd.MRequestHeaders()
        // mTd.MRequestQueryString()
        mTd.MRequestPayload()

        mutatedTcArray = append(mutatedTcArray, mTd.MTcDs...)
    }
    // aa, _ := json.Marshal(mutatedTcArray[0:1])   
    // fmt.Println(string(aa))
    return mutatedTcArray
}

func InitMTc (originTcData *testcase.TestCaseDataInfo) (MTestCaseDataInfo) {
    var m MTestCaseDataInfo
    var lKey string
    //
    hContentType := originTcData.TestCase.ReqHeaders()["Content-Type"]
    for key, _ := range originTcData.TestCase.ReqPayload() {
        lKey = key
        break
    }
    //
    m = MTestCaseDataInfo {
        OriginTcD:     originTcData,
        HContentType:  hContentType,
        TcPlType:      lKey,
    }

    return m    
}

func getMutatedTcData (tcJson []byte, i int, mFunc *MFunc, mutationRule string, 
        mInfoStr string, tcMutationInfo testcase.MutationInfo) testcase.TestCaseDataInfo {
    //------
    tcSuffix := mFunc.MTcSuffix + fmt.Sprint(i)

    //-- set new info to mutated tc
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    // mTcData.TestCase.SetPriority(mFunc.MPriority)
    mTcData.MutationArea     = mFunc.MArea
    mTcData.MutationCategory = mFunc.MCategory
    mTcData.MutationRule     = mutationRule
    mTcData.MutationInfoStr  = mInfoStr
    mTcData.MutationInfo     = tcMutationInfo
    
    return mTcData
}

func getTcMutationInfo (mFieldDetails MFieldDetails, mutatedValue interface{}) testcase.MutationInfo {
    tcMutationInfo := testcase.MutationInfo {
        FieldPath:    mFieldDetails.FieldPath,
        CurrValue:    mFieldDetails.CurrValue,
        FieldType:    mFieldDetails.FieldType,
        FieldSubType: mFieldDetails.FieldSubType,
        MutatedValue: mutatedValue,
    }

    return tcMutationInfo
}


