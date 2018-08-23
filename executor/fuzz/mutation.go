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
    // "fmt"
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





