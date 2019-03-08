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

func Run() {
    vm := goja.New()
    v, err := vm.RunString("2 + 2")
    if err != nil {
        panic(err)
    }
    if num := v.Export().(int64); num != 4 {
        panic(num)
    }
}

