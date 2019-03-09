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
    "fmt"
    // "os"
    "strings"
    "path/filepath"

    "go4api/cmd"
    "go4api/utils"

    goja "github.com/dop251/goja"
)

var JsFunctions []GJsBasics

// trial code for js
func InitJsFunctions () {
    jsFileList, _ := utils.WalkPath(cmd.Opt.Testcase, ".js")
    // fmt.Println(jsFileList)

    for i, _ := range jsFileList {
        srcBytes := utils.GetContentFromFile(jsFileList[i])
        src := string(srcBytes)

        p, err := goja.Compile("", src, false)
        if err != nil {
            panic(err)
        }

        jsFileName := strings.TrimSuffix(filepath.Base(jsFileList[i]), ".js")

        jsFunc := GJsBasics {
            JsSourceFilePath: jsFileList[i],
            JsSourceFileName: filepath.Base(jsFileList[i]),
            JsFunctionName: jsFileName,
            JsProgram: p,
        }

        JsFunctions = append(JsFunctions, jsFunc)
    }   
}

// for testing
func CallJsFuncs (funcName string, funcParams interface{}) interface{} {
    // fmt.Println(JsFunctions)
    for ii, _ := range JsFunctions {
        for i := 1; i < 10; i++ {
            go RunProgram(JsFunctions[ii].JsProgram, funcParams)
        }
    } 

    return 1
}

// for testing
func CallJsFunc (funcName string, funcParams interface{}) interface{} {
    // fmt.Println(JsFunctions)
    idx := -1
    var returnValue interface{}

    for i, _ := range JsFunctions {
        if JsFunctions[i].JsFunctionName == funcName {
            idx = i
            break
        }
    } 

    if idx != -1 {
        returnValue = RunProgram(JsFunctions[idx].JsProgram, funcParams)
    } else {
        fmt.Println("! Error, no js function found")
    }

    return returnValue
}

func RunProgram(p *goja.Program, funcParams interface{}) interface{} {
    vm := goja.New()

    vm.Set("funcParams", funcParams)
    v, err := vm.RunProgram(p)

    if err != nil {
        panic(err)
    }

    fmt.Println("the sum results is:", v.Export())

    return v.Export()
}
