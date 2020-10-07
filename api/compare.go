/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
    "fmt"

    "go4api/assertion" 
    "go4api/lib/testcase" 
)

//
func compareCommon (reponsePart string, key string, assertionKey string, actualValue interface{}, expValue interface{}) (bool, *testcase.TestMessage) {
    // Note: As get Go nil, for JSON null, need special care, two possibilities:
    // p1: expResult -> null, but can not find out actualValue, go set it to nil, i.e. null (assertion -> false)
    // p2: expResult -> null, actualValue can be founc, and its value --> null (assertion -> true)
    // but here can not distinguish them
    assertionResults := ""
    var testRes bool

    // reserved word: _ignore_assertion_
    if fmt.Sprint(expValue) == "_ignore_assertion_" {
        msg := testcase.TestMessage {
            AssertionResults: "Success",
            ReponsePart:      reponsePart,
            FieldName:        key,
            AssertionKey:     assertionKey,
            ActualValue:      actualValue,
            ExpValue:         expValue,   
        }
        
        testRes = true

        return testRes, &msg
    } 

    if actualValue == nil || expValue == nil {
        // if only one nil
        if actualValue != nil || expValue != nil {
            assertionResults = "Failed"
            testRes = false
        // both nil
        } else {
            assertionResults = "Success"
            testRes = true
        }
    // no nil
    } else {
        // call the assertion function
        testResult := assertion.CallAssertion(assertionKey, actualValue, expValue)

        if testResult == false {
            assertionResults = "Failed"
            testRes = false
        } else {
            assertionResults = "Success"
            testRes = true
        }
    }
    //
    msg := testcase.TestMessage {
        AssertionResults: assertionResults,
        ReponsePart:      reponsePart,
        FieldName:        key,
        AssertionKey:     assertionKey,
        ActualValue:      actualValue,
        ExpValue:         expValue,   
    }

    return testRes, &msg
}

