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
    TcType         string    // text, multipart-form, form, "" (others)
    HContentType   string    // Content-Type: multipart/form-data
    PLMPForm       PLMPForm  // if multipart-form
    PLMPFormNoFile PLMPForm
    PLMPFormFile   PLMPForm
    Tc4M           *testcase.TestCaseDataInfo  // if multipart-form, payload without file
    IMMTcs         []*testcase.TestCaseDataInfo  // intermediate
    FinalMTcs      []*testcase.TestCaseDataInfo  // final
}

type MFieldDetails struct {
    FieldPath     []string
    CurrValue     interface{}
    FieldType     string  // the json supported types
    FieldSubType  string  // like ip/email/phone/etc.
    MutatedValues []interface{}
}

type MFunc struct {
    MPriority string
    MTcSuffix string
    MArea     string
    MCategory string
    // MFunc interface{}
}

//
type PLForm []*Form
type Form struct {
    Name  string
    Value string
}

//
type PLMPForm []*MPForm
type MPForm struct {
    Name  string
    Value string
    Type  string
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
}

func MutateTcArray (originTcArray []*testcase.TestCaseDataInfo) []*testcase.TestCaseDataInfo {
    var mutatedTcArray []*testcase.TestCaseDataInfo

    for _, originTcData := range originTcArray {
        if originTcData.TestCase.IfGlobalSetUpTestCase() == true {
            continue
        }

        mutatedTcArray = append(mutatedTcArray, originTcData)
        // json, originTcData, multipart-form, form
        mTd := InitMTc(originTcData)
        tcJson, _ := json.Marshal(mTd)
        // --- here to start the mutation
        mTd.MRequestHeaders(tcJson)
        mTd.MRequestQueryString(tcJson)
        mTd.MRequestPayload(tcJson)

        // json, originTcData, multipart-form, form
        mTd.ReBuildTC()

        mutatedTcArray = append(mutatedTcArray, mTd.FinalMTcs...)
    }
    // aa, _ := json.Marshal(mutatedTcArray)
    // fmt.Println("\nmutatedTcArray: ", string(aa))
    return mutatedTcArray
}

func InitMTc (originTcData *testcase.TestCaseDataInfo) (MTestCaseDataInfo) {
    var m MTestCaseDataInfo
    //
    for key, _ := range originTcData.TestCase.ReqPayload() {
        lKey := strings.ToLower(key)
        switch lKey {
        case "text", "form", "" :
            tc4M := *originTcData
            m = MTestCaseDataInfo {
                OriginTcD: originTcData,
                TcType:    lKey,
                Tc4M:      &tc4M,
            }
        case "multipart-form":
            var pLMPForm PLMPForm
            var pLnf PLMPForm
            var pLf PLMPForm

            reqPayload := originTcData.TestCase.ReqPayload()
            reqPayloadJsonBytes, _ := json.Marshal(reqPayload)
            // reqPayloadJson := string(reqPayloadJsonBytes)

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
                pl[pLnf[i].Name] = pl[pLnf[i].Value]
            }
            tc4M := *originTcData
            hContentType := tc4M.TestCase.ReqHeaders()["Content-Type"].(string)
            tc4M.TestCase.DelRequestHeader("Content-Type") //multipart/form-data

            tc4M.TestCase.DelReqPayload("multipart-form")
            tc4M.TestCase.SetRequestPayload("text", pl)
            //
            m = MTestCaseDataInfo {
                OriginTcD:      originTcData,
                TcType:         "multipart-form",
                HContentType:   hContentType,
                PLMPForm:       pLMPForm,  // if multipart-form
                PLMPFormNoFile: pLnf,
                PLMPFormFile:   pLf,
                Tc4M:           &tc4M,
            }
        }
    }

    return m
}

func (mTd *MTestCaseDataInfo) ReBuildTC () {
    switch mTd.TcType {
    case "text", "form", "" :
        mTd.FinalMTcs = append(mTd.FinalMTcs, mTd.IMMTcs...)
    case "multipart-form":

    }
}

func getMutatedTcData (tcJson []byte, i int, mFunc *MFunc, mutationRule string, 
        mutationInfo string, tcMutationInfo testcase.MutationInfo) testcase.TestCaseDataInfo {
    //------
    tcSuffix := mFunc.MTcSuffix + fmt.Sprint(i)

    //-- set new info to mutated tc
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    mTcData.TestCase.UpdateTcName(mTcData.TcName() + tcSuffix)
    mTcData.TestCase.SetPriority(mFunc.MPriority)
    mTcData.MutationArea     = mFunc.MArea
    mTcData.MutationCategory = mFunc.MCategory
    mTcData.MutationRule     = mutationRule
    mTcData.MutationInfoStr  = mutationInfo
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


