/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package builtins

import (
	"math/rand"                                                                                                                                        
	"time"
)

var alphaNumeric = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var charSet = []rune("中文测试的些字符集可以使用一二三四五六七八九十abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")


func NextInt (min int, max int) int {
    l := max - min
    rand.Seed(time.Now().UnixNano())
    
    return rand.Intn(l) + min
}

func NextAlphaNumeric (n int) string {
    b := make([]rune, n)
    l := len(alphaNumeric)
    for i := range b {
        // [0,n)
        rand.Seed(time.Now().UnixNano())
        b[i] = alphaNumeric[rand.Intn(l)]
    }

    return string(b)
}

func CurrentTimeStampString () string {
	t := time.Now()

	return t.Format("2006-01-02 15:04:05")
}

func CurrentTimeStampMilliString () string {
	t := time.Now()

	return t.Format("2006-01-02 15:04:05.999")
}

func CurrentTimeStampMicroString () string {
	t := time.Now()

	return t.Format("2006-01-02 15:04:05.999999")
}

func CurrentTimeStampNanoString () string {
	t := time.Now()

	return t.Format("2006-01-02 15:04:05.999999999")
}

func CurrentTimeStampUnix () int64 {
	t := time.Now()

	return t.Unix()
}

func CurrentTimeStampUnixMilli () int64 {
	t := time.Now()

	return t.UnixNano() / 1000000
}

func CurrentTimeStampUnixMicro () int64 {
	t := time.Now()

	return t.UnixNano() / 1000
}

func CurrentTimeStampUnixNano () int64 {
	t := time.Now()

	return t.UnixNano()
}


