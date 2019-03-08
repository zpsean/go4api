/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gjs

import (
    // "fmt"
    // "os"
    // "strings"

    // "go4api/cmd"

    // goja "github.com/dop251/goja"
)

// mvp, one js file has only one function
type GJsBasics struct {
    JsSourceFilePath string
    JsSourceFileName string
    JsFunctionName string
    JsFunctionInParams interface{}
    JsFunctionOut interface{}
}

type GJsSet []GJsBasics