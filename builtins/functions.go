/*
 * go4api - an api testing tool written in Go
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
var specialChar = []rune(" !@#$%^&*(){}|<>?~<>")


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

// { "Fn::Substitute" : [ String, { Var1Name: Var1Value, Var2Name: Var2Value } ] }
// { "Fn::Substitute" : [ "www.${var1}", { "var1": "value1"} ]}
func Substitute (param interface{}) string {
    switch param.(type) {
        case []interface{}:
            sourceStr := ""
            var err error
            paramSlice := reflect.ValueOf(param).Interface().([]interface{})
            
            if len(paramSlice) <= 1 {
                panic(err)
            } else if len(paramSlice) > 1 {
                sourceStr = paramSlice[0].(string)

                if err != nil {
                    panic(err)
                }

                var keyValueMap = make(map[string]interface{})
                keyValueMap = reflect.ValueOf(paramSlice[1]).Interface().(map[string]interface{})
                for key, value := range keyValueMap {
                    sourceStr = strings.Replace(sourceStr, "${" + key + "}", fmt.Sprint(value), -1)
                }
            }

            return sourceStr
        default:
            return ""
    }
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

                // cosider the string first, tbd, for enhance to interface{}
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

// Replace(s, old, new string, n int) string
// { "Fn::Replace" : ["2019-01-01" , "-", "/"] } => "2019/01/01"
func Replace (param interface{}) string {
    switch param.(type) {
        case []interface{}:
            var res string
            var err error

            paramSlice := reflect.ValueOf(param).Interface().([]interface{})

            if len(paramSlice) <= 2 {
                panic(err)
            } else if len(paramSlice) > 2 {
                res = strings.Replace(paramSlice[0].(string), paramSlice[1].(string), paramSlice[2].(string), -1)
            }

            return res
        default:
            return ""
    }
}

// SubString(s, start_position, end_position) string
// { "Fn::SubString" : ["2019-01-01" , 1, 3] } => "019"
func SubString (param interface{}) string {
    switch param.(type) {
        case []interface{}:
            var res string
            var err error

            paramSlice := reflect.ValueOf(param).Interface().([]interface{})

            if len(paramSlice) <= 2 {
                panic(err)
            } else if len(paramSlice) > 2 {
                start_pos := int(paramSlice[1].(float64))
                end_pos := int(paramSlice[2].(float64))

                res = paramSlice[0].(string)[start_pos:end_pos]
            }

            return res
        default:
            return ""
    }
}


//
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

//
func Length (param interface{}) int {
    var l int

    switch param.(type) {
    case []interface{}:
        l = len(param.([]interface{}))
    case string:
        // check if [] string
        l = len(param.(string))
    default:
        l = 0
    }

    return l
}

//
func ToInt (param interface{}) interface{} {
    switch param.(type) {
    case float64:
        // Note: tbd, to cosider format for int, float64 more
        return int(param.(float64))
    case string:
        // for null, _null_key_, _null_value_
        if fmt.Sprint(param) == "_null_key_" || fmt.Sprint(param) == "_null_value_" {
            return param
        }

        if fmt.Sprint(param) == "_ignore_assertion_" {
            return param
        }

        i, err := strconv.Atoi(param.(string))
        if err != nil {
            fmt.Println("ToInt param is: ", param)
            panic(err)
        }
        return i
    default:
        return 0
    }
}

func ToBool (param interface{}) bool {
    switch param.(type) {
        case float64:
            if param.(float64) == 0 {
                return false
            } else {
                return true
            }
        default:
            if fmt.Sprint(param) == "0" || strings.ToLower(fmt.Sprint(param)) == "false" {
                return false
            } else {
                return true
            }
    }
}


func CurrentTimeStampString (param interface{}) string {
    t := time.Now()
    format := "2006-01-02 15:04:05"

    switch strings.ToLower(param.(string)) {
    case "micro":
        format = "2006-01-02 15:04:05.999"
        return t.Format(format)
    case "milli":
        format = "2006-01-02 15:04:05.999999"
        return t.Format(format)
    case "nano":
        format = "2006-01-02 15:04:05.999999999"
        return t.Format(format)
    default:
        return t.Format(format)
        }
}

func CurrentTimeStampUnix (param interface{}) int64 {
	t := time.Now()

    switch strings.ToLower(param.(string)) {
    case "milli":
        return t.UnixNano() / 1000000
    case "micro":
        return t.UnixNano() / 1000
    case "nano":
        return t.UnixNano()
    default:
        return t.Unix()
    }
}

// 2006-01-02 00:00:00
func DayStart (param interface{}) interface{} {
    var ts int64

    switch param.(type) {
    // case string:
        //
    case int64:
        t := time.Unix(param.(int64) / 1000, 0)
        y, m, d := t.Date()
        tStart := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
        // fmt.Println("Go launched at start int64: ", t, tStart.Local())

        ts = tStart.Unix()
    case float64:
        t := time.Unix(int64(param.(float64)) / 1000, 0)
        y, m, d := t.Date()
        tStart := time.Date(y, m, d, 0, 0, 0, 0, time.Local)

        ts = tStart.Unix()
    default:
        fmt.Println("can not recognized type")
    }

    return ts * 1000
}

// 2006-01-02 23:59:59
func DayEnd (param interface{}) interface{} {
    var ts int64

    switch param.(type) {
    // case string:
        //
    case int64:
        t := time.Unix(param.(int64) / 1000, 0)
        y, m, d := t.Date()
        tStart := time.Date(y, m, d, 23, 59, 59, 0, time.Local)
        // fmt.Println("Go launched at end int64: ", t, tStart.Local())

        ts = tStart.Unix()
    case float64:
        t := time.Unix(int64(param.(float64)) / 1000, 0)
        y, m, d := t.Date()
        tStart := time.Date(y, m, d, 23, 59, 59, 0, time.Local)

        ts = tStart.Unix()
    default:
        fmt.Println("can not recognized type")
    }

    return ts * 1000
}

// param is time str
func ConvertTimeToUnix (param interface{}) int64 {
    format := "2006-01-02 15:04:05 +0800 CST"
    t, err := time.Parse(format, param.(string))
    if err != nil {
        panic(err)
    }

    fmt.Println(param.(string), t, t.UnixNano() / 1000000)

    return t.UnixNano() / 1000000
}

// param is time int, to str
func ConvertTimeToStr (param interface{}) string {
    format := "2006-01-02 15:04:05"
    timeStr := ""
   
    switch param.(type) {
    case string:
        i, err := strconv.ParseInt(param.(string), 10, 64)
        if err != nil {
            panic(err)
        }
        t := time.Unix(i / 1000, 0)

        timeStr = t.Format(format)
    case int64:
        t := time.Unix(param.(int64) / 1000, 0)

        timeStr = t.Format(format)
    case float64:
        t := time.Unix(int64(param.(float64)) / 1000, 0)
        
        timeStr = t.Format(format)
    default:
        fmt.Println("can not recognized time")
    }

    return timeStr
}

// { "Fn::TimeStampOffset" : [ "time" , "offset", "unit" ] }
func TimeStampUnixOffset (param interface{}) interface{} {
    // unit: years, months, days, hours, minutes, seconds, millis, micros, nanos
    var res int64

    switch param.(type) {
    case []interface{}:
        paramSlice := reflect.ValueOf(param).Interface().([]interface{})

        if len(paramSlice) <= 2 {
            fmt.Println("Not enough params for TimeStampUnixOffset")
        } else if len(paramSlice) > 2 {
            var timeSource int64
            var offSet int
            var err error

            switch paramSlice[0].(type) {
            case int:
                timeSource = int64(paramSlice[0].(int))
            case int64:
                timeSource = int64(paramSlice[0].(int64))
            case float64:
                timeSource = int64(paramSlice[0].(float64))
            }
            
            t := time.Unix(timeSource / 1000, 0)
            // fmt.Println("TimeStampUnixOffset day: ", t)
            
            offSet, err = strconv.Atoi(paramSlice[1].(string))
            if err != nil {
                panic(err)
            }

            switch paramSlice[2].(string) {
                // case years:
                // case months:
                case "day":
                    // oD := offSet * 24
                    // oDStr := strconv.ItoA(oD) + "h"
                    // d, _ := time.ParseDuration(oDStr)
                    resT := t.AddDate(0, 0, offSet)
                    // fmt.Println("TimeStampUnixOffset day: ", offSet, resT)

                    res = resT.Unix()
                // case hours:
                // case minutes:
                // case seconds:
                // case millis:
                // case micros:
                // case nanos:
                default:
                    fmt.Println("can not read the offset unit, it must be one of years, months, days, hours, minutes, seconds, millis, micros, nanos ")
            }
        }

        return res * 1000
    default:
        return res * 1000
    }

    return res * 1000
}

// conditions
func And (param interface{}) bool {
    return true
}

func Or (param interface{}) bool {

    return true
}

func If (param interface{}) bool {

    return true
}

func Not (param interface{}) bool {

    return true
}


