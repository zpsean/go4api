/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package fuzz

import (                                                                                                                                             
    // "os"
    // "time"
    "fmt"
    // "path/filepath"
    // "strings"
    // "strconv"
    "reflect"
    "encoding/json"
    "go4api/testcase"  
)

// mutation is to mutate the valid data to working api, see if mutated invalid data still can be handled by the api
// func (tcData testcase.TestCaseDataInfo) MutateRequestMethod () {
//     tcData.TestCase.SetRequestMethod("DELETE")

// }

// func (muTc testcase.TestCaseDataInfo) MutateRequestPath () {
//     tcData.TestCase.SetRequestPath("/aa/bb/cc")
// }

// two ways to mutate the testcase:
// Option 1: 
// copy the underlying fields and values to another TestCaseDataInfo, with mutation(s)
// the better way would be deep copy the TestCaseDataInfo, and change the value, but Golang standard
// Lib has no deepcopy, so that, here use a plain way, that is, to re-sturct the TestCaseDataInfo
//
// Option 2:
// json.Marshal the tc in originMutationTcArray, 
// then change the value(s) in the json
// then Unmarshal the to testcase, and add to mutatedTcArray
// then execute the mutatedTcArray

// focus on the Request to mutate
// type Request struct {  
//     Method string
//     Path string
//     Headers map[string]interface{}
//     QueryString map[string]interface{}
//     Payload map[string]interface{}
// }

func MutateTcArray(originMutationTcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    for _, originTcData := range originMutationTcArray {
        tcJson, _ := json.Marshal(originTcData)
        fmt.Println("tcJson:", string(tcJson)) 

        mutatedTcArray = append(mutatedTcArray, originTcData)

        // here to start the mutation
        // headers
        mutatedTcArray = append(mutatedTcArray, MutateSetRequestHeader(tcJson))
        mutatedTcArray = append(mutatedTcArray, MutateAddRequestHeader(tcJson))

        i := 0
        for k, _ := range originTcData.TestCase.ReqHeaders() {
            mutatedTcArray = append(mutatedTcArray, MutateDelRequestHeader(tcJson, k, i))
            i = i + 1
        }

        // querystring
        i = 0
        for key, value := range originTcData.TestCase.ReqQueryString() {
            fmt.Println(reflect.TypeOf(value))
            // if value match number mode
            mutatedTcArray = append(mutatedTcArray, MutateSetRequestQueryString(tcJson, key, fmt.Sprint(-1), key + fmt.Sprint(i)))
            i = i + 1
            mutatedTcArray = append(mutatedTcArray, MutateSetRequestQueryString(tcJson, key, fmt.Sprint(0), key + fmt.Sprint(i)))
            i = i + 1
            mutatedTcArray = append(mutatedTcArray, MutateSetRequestQueryString(tcJson, key, fmt.Sprint(10000), key + fmt.Sprint(i)))
            i = i + 1
        }


        mutatedTcArray = append(mutatedTcArray, MutateAddRequestQueryString(tcJson))

        i = 0
        for key, _ := range originTcData.TestCase.ReqQueryString() {
            mutatedTcArray = append(mutatedTcArray, MutateDelRequestQueryString(tcJson, key, i))
            i = i + 1
        }

        // Payload
        i = 0
        for key, _ := range originTcData.TestCase.ReqPayload() {
            mutatedTcArray = append(mutatedTcArray, MutateDelPayload(tcJson, key, i))
            i = i + 1
        }


    }
    // fmt.Println("\nmutatedTcArray: ", mutatedTcArray)

    return mutatedTcArray
}




func MutateSetRequestHeader (tcJson []byte) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-seth-" + fmt.Sprint(1))
    mTcData.TestCase.SetRequestHeader("aaaa", "dbddsdsfa")

    return mTcData
}


func MutateAddRequestHeader (tcJson []byte) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-addh-" + fmt.Sprint(2))
    mTcData.TestCase.AddRequestHeader("aaaakk", "dbddsdsfa")

    return mTcData
}


func MutateDelRequestHeader (tcJson []byte, k string, i int) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-Delh-" + fmt.Sprint(i))
    mTcData.TestCase.DelRequestHeader(k)

    return mTcData
}



//
func MutateSetRequestQueryString (tcJson []byte, key string, value string, suffix string) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-setq-" + suffix)
    mTcData.TestCase.SetRequestQueryString(key, value)
    fmt.Println(mTcData.TestCase.ComposeReqQueryString())

    return mTcData
}


func MutateAddRequestQueryString (tcJson []byte) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-addq-" + fmt.Sprint(2))
    fmt.Println("\n", mTcData.TestCase.ComposeReqQueryString())
    mTcData.TestCase.AddRequestQueryString("aaaakk", "dbddsdsfa")
    fmt.Println(mTcData.TestCase.ComposeReqQueryString())

    return mTcData
}


func MutateDelRequestQueryString (tcJson []byte, k string, i int) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-Delq-" + fmt.Sprint(i))
    fmt.Println("\n", mTcData.TestCase.ComposeReqQueryString(), k)
    mTcData.TestCase.DelRequestQueryString(k)
    fmt.Println(mTcData.TestCase.ComposeReqQueryString())

    return mTcData
}


//
func MutateDelPayload (tcJson []byte, k string, i int) testcase.TestCaseDataInfo {
    var mTcData testcase.TestCaseDataInfo
    json.Unmarshal(tcJson, &mTcData)

    originTcName := mTcData.TcName()
    mTcData.TestCase = mTcData.TestCase.UpdateTcName(originTcName + "-M-Delp-" + fmt.Sprint(i))

    fmt.Println("\n", mTcData.TestCase.ReqPayload(), k)
    mTcData.TestCase.DelReqPayload(k)
    fmt.Println(mTcData.TestCase.ReqPayload())

    return mTcData
}







