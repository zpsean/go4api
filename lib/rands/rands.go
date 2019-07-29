/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package rands

import (
    // "fmt"
    // "reflect"
    "math/rand"
    "time"
)

// Note: this package is for the automatically generate the combinations from the data mode file
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var letterCNRunes = []rune("这是为了中文测试的一些字符集可以使用一二三四五六七八九十")
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numericSet = "0123456789"
// ip
// email


// func GetOriginModeData(dataModeFile string) {
//     // typeList = []
//     // targetValueList = [][]

//     // switch type {
//     //     case char()
//     //     case varchar()
//     //     case int
//     //     case numeric()
//     //     case email
//     //     case float
//     //     case list
//     //     ...

//     // }

//     // GenerateCombinations()
// }


// RandStringRunes - this can handle with non ASCII char
func RandStringRunes(n int) string {
    b := make([]rune, n)
    l := len(letterRunes)
    for i := range b {
        // [0,n)
        rand.Seed(time.Now().UnixNano())
        b[i] = letterRunes[rand.Intn(l)]
    }
    // fmt.Println("RandStringRunes: ", string(b))
    return string(b)
}

func RandStringCNRunes(n int) string {
    b := make([]rune, n)
    l := len(letterCNRunes)
    for i := range b {
        // [0,n)
        rand.Seed(time.Now().UnixNano())
        b[i] = letterCNRunes[rand.Intn(l)]
    }
    // fmt.Println("RandStringCNRunes: ", string(b))
    return string(b)
}

// RandNums
func RandNums(n int) string {
    b := make([]rune, n)
    l := len(numericSet)
    for i := range b {
        // [0,n)
        b[i] = letterRunes[rand.Intn(l)]
    }
    return string(b)
}
