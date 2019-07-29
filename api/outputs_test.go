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
    "testing"
    "reflect"
    // "encoding/csv"
)

var actualBody []byte

func init() {
    actualBodyS := `
        {
            "count": 20,
            "start": 0,
            "total": 250,
            "subjects": [{
                    "rating": {
                        "max": 10,
                        "average": 9.6,
                        "stars": "50",
                        "min": 0
                    },
                    "title": "肖申克的救赎"
                },
                {
                    "rating": {
                        "max": 10,
                        "average": 9.5,
                        "stars": "50",
                        "min": 0
                    },
                    "title": "霸王别姬"
                },
                {
                    "rating": {
                        "max": 10,
                        "average": 9.4,
                        "stars": "50",
                        "min": 0
                    },
                    "title": "这个杀手不太冷"
                }
            ],
            "title": "豆瓣电影Top250",
            "dummykeyfornull": null
        }`
                    
    actualBody = []byte(actualBodyS)
}


func Test_convertSliceAsString(t *testing.T) {
    fmt.Println("\n--> test started")

    var aa []interface{}
    aa = append(aa, "1234")
    aa = append(aa, `1234`)
    aa = append(aa, `123,asdfasd"asdfasdf;4`)

    str := convertSliceAsString(aa)
    fmt.Println("==>", aa, str)
    if str != `["1234","1234","123,asdfasd\"asdfasdf;4"]` {
        t.Fatalf("Error, convert failed")
    }
    fmt.Println("\n--> test finished")
}

func Test_convertSliceAsString2(t *testing.T) {
    fmt.Println("\n--> test started")

    var aa []interface{}    
    aa = append(aa, `{"1324": 12423, "asfdsf": "kjhgfd"}`)

    str := convertSliceAsString(aa)
    fmt.Println("==>", aa, str)

    if str != `["{\"1324\": 12423, \"asfdsf\": \"kjhgfd\"}"]` {
        t.Fatalf("Error, convert failed")
    }
    fmt.Println("\n--> test finished")
}

func Test_convertSliceAsString3(t *testing.T) {
    fmt.Println("\n--> test started")

    var aa []interface{}    
    aa = append(aa, `["file1", "union", "file2", "join", "file3"]`)
    aa = append(aa, `[1234,134231,1324343]`)

    str := convertSliceAsString(aa)
    fmt.Println("==>", aa, str)

    if str != `["[\"file1\", \"union\", \"file2\", \"join\", \"file3\"]","[1234,134231,1324343]"]` {
        t.Fatalf("Error, convert failed")
    }
    fmt.Println("\n--> test finished")
}


func Test_GetActualValueByJsonPath(t *testing.T) {
    fmt.Println("\n--> test started")

    key := "$.subjects.#.title"

    res := GetActualValueByJsonPath(key, actualBody)
    resSlice := reflect.ValueOf(res).Interface().([]interface{})
    fmt.Println("==>", res, resSlice, len(resSlice))

    if len(resSlice) != 3 {
        t.Fatalf("Error, look up failed")
    }
    fmt.Println("\n--> test finished")
}


func Test_GetActualValueByJsonPath2(t *testing.T) {
    fmt.Println("\n--> test started")

    key := "$.subjects.#.titleBK"

    res := GetActualValueByJsonPath(key, actualBody)
    resSlice := reflect.ValueOf(res).Interface().([]interface{})
    fmt.Println("==>", res, resSlice, len(resSlice))

    if len(resSlice) != 0 {
        t.Fatalf("Error, look up failed")
    }
    fmt.Println("\n--> test finished")
}


func Test_GetActualValueByJsonPath3(t *testing.T) {
    fmt.Println("\n--> test started")

    key := "$.subjectsBK"

    res := GetActualValueByJsonPath(key, actualBody)
    fmt.Println("==>", res)

    if res != nil {
        t.Fatalf("Error, look up failed")
    }
    fmt.Println("\n--> test finished")
}

func Test_GetActualValueByJsonPath4(t *testing.T) {
    fmt.Println("\n--> test started")

    key := ".subjectsBK"

    res := GetActualValueByJsonPath(key, actualBody)
    fmt.Println("==>", "`" + fmt.Sprint(res) + "`")

    if res != key {
        t.Fatalf("Error, look up failed")
    }
    fmt.Println("\n--> test finished")
}

func Test_GetActualValueByJsonPath5(t *testing.T) {
    fmt.Println("\n--> test started")

    key := ""

    res := GetActualValueByJsonPath(key, actualBody)
    fmt.Println("==>", res)

    if res != key {
        t.Fatalf("Error, look up failed")
    }
    fmt.Println("\n--> test finished")
}

func Test_GetActualValueByJsonPath6(t *testing.T) {
    fmt.Println("\n--> test started")

    key := "$.dummykeyfornull"

    res := GetActualValueByJsonPath(key, actualBody)
    fmt.Println("==>", res)

    if res != nil {
        t.Fatalf("Error, dummykeyfornull look up failed")
    } else {
        t.Log("Error, dummykeyfornull look up passed")
    }
    fmt.Println("\n--> test finished")
}

