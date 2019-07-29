/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package utils

import (
    "fmt"
    "math"
    "strconv"

)

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', -1, 64)
}

func CheckFloat64SubType(f float64) {
    // get current dir, 
    fmt.Println(int64(f))
    a, b := math.Modf(f)

    fmt.Println(a, b)
    fmt.Println(f)
    fmt.Println(fmt.Sprint(f))
    fmt.Println(FloatToString(f))
}


