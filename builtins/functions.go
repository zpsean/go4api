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
    "fmt"
	"math/rand"                                                                                                                                        
	"time"
    "strings"
    "strconv"
    "reflect"
)

var alphaNumeric = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var numeric = []rune("0123456789")
var charSet = []rune("中文测试的些字符集可以使用一二三四五六七八九十abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")


func NextInt (param interface{}) int {
    switch param.(type) {
        case []interface{}:
            var min, max int
            var err error

            paramSlice := reflect.ValueOf(param).Interface().([]interface{})

            if len(paramSlice) == 0 {
                panic(err)
            } else if len(paramSlice) == 1 {
                min = int(paramSlice[0].(float64))
                max = min
            } else if len(paramSlice) >= 1 {
                min = int(paramSlice[0].(float64))
                max = int(paramSlice[1].(float64))
            }

            l := max - min
            if l <= 0 {
                return min
            }
            rand.Seed(time.Now().UnixNano())
            
            return rand.Intn(l) + min
        default:
            return 0
    }
}

func NextAlphaNumeric (param interface{}) string {
    switch param.(type) {
        case float64:
        	n := int(param.(float64))

            b := make([]rune, n)
            l := len(alphaNumeric)
            if n <= 0 {
            	return ""
            } 
            
            for i := range b {
                // [0,n)
                rand.Seed(time.Now().UnixNano())
                b[i] = alphaNumeric[rand.Intn(l)]
            }

            return string(b)
        default:
            return ""
    }
}

func NextStringNumeric (param interface{}) string {
    switch param.(type) {
        case float64:
            n := int(param.(float64))

            b := make([]rune, n)
            l := len(numeric)
            if n <= 0 {
                return ""
            } 
            
            for i := range b {
                // [0,n)
                rand.Seed(time.Now().UnixNano())
                b[i] = numeric[rand.Intn(l)]
            }

            return string(b)
        default:
            return ""
    }
}

func Substitute (target string, dddd string) {

}

// { "Fn::Select" : [ "1", [ "apples", "grapes", "oranges", "mangoes" ] ] }
func Select (param interface{}) string {
    switch param.(type) {
        case []interface{}:
            res := ""
            var err error
            paramSlice := reflect.ValueOf(param).Interface().([]interface{})
            
            if len(paramSlice) <= 1 {
                panic(err)
            } else if len(paramSlice) > 1 {
                index, err := strconv.Atoi(paramSlice[0].(string))

                if err != nil {
                    panic(err)
                }

                var listValues []string
                valueSlice := reflect.ValueOf(paramSlice[1]).Interface().([]interface{})
                for i, _ := range valueSlice {
                    listValues = append(listValues, valueSlice[i].(string))
                }

                res = listValues[index]
            }

            return res
        default:
            return ""
    }
}

// { "Fn::Join" : [ ":", [ "a", "b", "c" ] ] }
func Join (param interface{}) string {
    switch param.(type) {
        case []interface{}:
            res := ""
            var err error
            paramSlice := reflect.ValueOf(param).Interface().([]interface{})
            
            if len(paramSlice) <= 1 {
                panic(err)
            } else if len(paramSlice) > 1 {
                delimiter := paramSlice[0].(string)

                var listValues []string
                valueSlice := reflect.ValueOf(paramSlice[1]).Interface().([]interface{})
                for i, _ := range valueSlice {
                    listValues = append(listValues, valueSlice[i].(string))
                }

                res = strings.Join(listValues, delimiter)
            }

            return res
        default:
            return ""
    }
}

// { "Fn::Split" : [ "|" , "a|b|c" ] }
func Split (param interface{}) []string {
    switch param.(type) {
        case []interface{}:
            var res []string
            var err error

            paramSlice := reflect.ValueOf(param).Interface().([]interface{})

            if len(paramSlice) <= 1 {
                panic(err)
            } else if len(paramSlice) > 1 {
                res = strings.Split(paramSlice[1].(string), paramSlice[0].(string))
            }

            return res
        default:
            return []string{}
    }
}

func ToString (param interface{}) string {
    // fmt.Println(">>>>>>>>>>>: ", param)
    // fmt.Println(">>>>>>>>>>>: ", int(param.(float64)))
    switch param.(type) {
        case float64:
            // Note: tbd, to cosider format for int, float64 more
            return fmt.Sprint(int(param.(float64)))
        default:
            return fmt.Sprint(param)
    }
}

func ToInt (param interface{}) int {
    switch param.(type) {
        case float64:
            // Note: tbd, to cosider format for int, float64 more
            return int(param.(float64))
        default:
            return 0
    }
}

func CurrentTimeStampString (param interface{}) string {
    t := time.Now()

    format := "2006-01-02 15:04:05"
    if len(fmt.Sprint(param)) > 0{
        format = fmt.Sprint(param)
    }
    // note: tbd, for more format according to param using yyyy mm dd hh mm ss 
    // "2006-01-02 15:04:05"
    // "2006-01-02 15:04:05.999" MilliString
    // "2006-01-02 15:04:05.999999" MicroString
    // "2006-01-02 15:04:05.999999999" NanoString

    return t.Format(format)
}

func CurrentTimeStampUnix (param interface{}) int64 {
	t := time.Now()

	return t.Unix()
}

func CurrentTimeStampUnixMilli (param interface{}) int64 {
	t := time.Now()

	return t.UnixNano() / 1000000
}

func CurrentTimeStampUnixMicro (param interface{}) int64 {
	t := time.Now()

	return t.UnixNano() / 1000
}

func CurrentTimeStampUnixNano (param interface{}) int64 {
	t := time.Now()

	return t.UnixNano()
}


