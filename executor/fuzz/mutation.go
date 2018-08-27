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
    "go4api/testcase"  
)

// mutation is to mutate the valid data to working api, see if mutated invalid data still can be handled by the api
type MutationField struct {
    MuChar string
    MuInt int64
}


type MutationTestCase struct {
    TestCase testcase.TestCase
}

func (muTc MutationTestCase) MutateRequestMethod () {
    muTc.TestCase.SetRequestMethod("DELETE")

    // "request": {
    //     "method": "GET",
    //     "path": "/api/operation/delivery-terms",
    //     "headers": {
    //       "authorization": "{{.authorization}}"
    //     },
    //     "queryString": {
    //       "pageIndex": "1",
    //       "pageSize": "12"
    //     }
}

func (muTc MutationTestCase) MutateRequestPath () {
    muTc.TestCase.SetRequestPath("/aa/bb/cc")
}


func MutateTcArray(originMutationTcArray []testcase.TestCaseDataInfo) []testcase.TestCaseDataInfo {
    var mutatedTcArray []testcase.TestCaseDataInfo

    for _, originTcData := range originMutationTcArray {
        // here copy the underlying fields and values to another TestCaseDataInfo, with mutation(s)
        // the better way would be deep copy the TestCaseDataInfo, and change the value, but Golang standard
        // Lib has no deepcopy, so that, here use a plain way, that is, to re-sturct the TestCaseDataInfo


        // type Request struct {  
        //     Method string
        //     Path string
        //     Headers map[string]interface{}
        //     QueryString map[string]interface{}
        //     Payload map[string]interface{}
        // }
        // var mtcRequest testcase.Request
        // mtcRequest.Method = originTcData.TestCase.ReqMethod()
        // mtcRequest.Path = originTcData.TestCase.ReqPath()

        // mtcRequest.Headers = originTcData.TestCase.DelRequestHeader("authorization")

        // mtcRequest.QueryString = originTcData.TestCase[originTcData.TcName()].Request.QueryString
        // mtcRequest.Payload = originTcData.TestCase.ReqPayload()


        var mTcBasics testcase.TestCaseBasics
        
        mTcBasics.Priority = originTcData.TestCase[originTcData.TcName()].Priority
        mTcBasics.ParentTestCase = originTcData.TestCase[originTcData.TcName()].ParentTestCase
        mTcBasics.Inputs = originTcData.TestCase[originTcData.TcName()].Inputs
        // Request
        mTcBasics.Request = originTcData.TestCase[originTcData.TcName()].Request
        // Response
        mTcBasics.Response = originTcData.TestCase[originTcData.TcName()].Response
        mTcBasics.Outputs = originTcData.TestCase[originTcData.TcName()].Outputs



        originTcData.TestCase.DelRequestHeader("authorization")

        mTc := testcase.TestCase{}
        mTc[originTcData.TcName()] = mTcBasics

        var mTcData testcase.TestCaseDataInfo
        mTcData.TestCase = mTc
        mTcData.JsonFilePath = originTcData.JsonFilePath
        mTcData.CsvFile = originTcData.CsvFile
        mTcData.CsvRow = originTcData.CsvRow

        fmt.Println("TcData: ", mTcData, originTcData)
        mutatedTcArray = append(mutatedTcArray, mTcData)
    }

    return mutatedTcArray
}









