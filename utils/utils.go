/*
 * go4api - a api testing tool written in Go
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
    "io/ioutil"                                                                                                                                              
    "os"
    "io"
    "strings"
    // "bufio"
    "reflect"
    // "path"
    "bytes"
    "path/filepath"
    "encoding/csv"
    "strconv"
    // "regexp"
    simplejson "github.com/bitly/go-simplejson"
)

func GetCurrentDir() string{
    // get current dir, 
    // Note: here may be a bug if run the main.go on other path, need to use abs path
    currentDir, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    return currentDir
}


func GetCsvFromFile(filePath string) [][]string {
    fi,err := ioutil.ReadFile(filePath)
    if err != nil {
        panic(err)
    }
    r2 := csv.NewReader(strings.NewReader(string(fi)))
    csvRows, _ := r2.ReadAll()

    return csvRows
}

func GetJsonFromFile(filePath string) *simplejson.Json {
    fi, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer fi.Close()
    //
    fd, err := ioutil.ReadAll(fi)

    res, err := simplejson.NewJson([]byte(fd))
    if err != nil {
        panic(err)
    }

    return res
}

func GetContentFromFile(filePath string) []byte {
    fi,err := ioutil.ReadFile(filePath)
    if err != nil {
        panic(err)
    }
    // contents := strings.NewReader(string(fi))

    return fi
}

func GetBaseUrlFromConfig(filePath string, jsonKey string) string {
    res := GetJsonFromFile(filePath)

    baseUrl, err := res.Get(jsonKey).Get("baseUrl").String()
    if err != nil {
        panic(err)
    }

    return baseUrl
}

func GetTestCaseJsonFromTestDataFile(filePath string) []interface{} {
    res := GetJsonFromFile(filePath)

    testcases, err := res.Get("TestCases").Array()
    if err != nil {
        panic(err)
    }

    var tcJsons []interface{}

    for i, _ := range testcases {
      tc := res.Get("TestCases").GetIndex(i)
      tcJsons = append(tcJsons, tc)
    }

    // fmt.Println("testcases:", testcases, "\n")
    // fmt.Println("tcJsons:", tcJsons, "\n")

    return tcJsons
}

func GetTestCaseJsonFromTestData(fjson *bytes.Buffer) []interface{} {
    // fmt.Println("fjson", fjson)
    readBuf, _ := ioutil.ReadAll(fjson)
    res, err := simplejson.NewJson(readBuf)

    testcases, err := res.Get("TestCases").Array()
    if err != nil {
        panic(err)
    }

    var tcJsons []interface{}

    for i, _ := range testcases {
      tc := res.Get("TestCases").GetIndex(i)
      tcJsons = append(tcJsons, tc)
    }

    // fmt.Println("testcases:", testcases, "\n")
    // fmt.Println("tcJsons:", tcJsons, "\n")

    return tcJsons
}

func GetTestCaseBasicInfoFromTestData(testcase interface{}) []interface{} {
    // Big question, how to get the KEY for each json => using: func (j *Json) Map() (map[string]interface{}, error)
    tc := testcase.(*simplejson.Json)
    tc_map, _ := tc.Map()

    var tcInfo []interface{}
    for key, _ := range tc_map {
      tcInfo = append(tcInfo, key)
      tcInfo = append(tcInfo, tc.Get(key).Get("priority").MustString())
      tcInfo = append(tcInfo, tc.Get(key).Get("parentTestCase").MustString())
      // inputs if optionsl 
      // tcInfo = append(tcInfo, tc.Get(key).Get("inputs").MustString())
      // as one test case has only one value each above, breank
    }

    return tcInfo
}


func GetRequestForTC(tc *simplejson.Json, tcName string) (string, string) {
    request := tc.Get(tcName).Get("request")

    apiPath, _ := request.Get("path").String()
    apiMethod, _ := request.Get("method").String()

    return apiPath, apiMethod
}


func GetRequestHeadersForTC(tc *simplejson.Json, tcNname string) map[string]interface{} {
    requestHeaders, _ := tc.Get(tcNname).Get("request").Get("headers").Map()
   
    return requestHeaders
}

func GetRequestPayloadForTC(tc *simplejson.Json, tcName string) map[string]interface{} {
    // var requestPayload map[string]interface{}
    requestPayload, _ := tc.Get(tcName).Get("request").Get("payload").Map()

    return requestPayload
}

func GetExpectedResponseForTC(tc *simplejson.Json, tcNname string) (map[string]interface{}, map[string]interface{}, map[string]interface{}) {
    var expStatusCode, expHeader, expBody map[string]interface{}

    expResponse := tc.Get(tcNname).Get("response")
    expStatusCode, _ = expResponse.Get("status").Map()
    expHeader, _ = expResponse.Get("headers").Map()
    expBody, _ = expResponse.Get("body").Map()

    return expStatusCode, expHeader, expBody
}

func GetInputsFileNameForTC(tc *simplejson.Json, tcNname string) []interface{} {
    var expOutputs []interface{}

    expOutputs, _ = tc.Get(tcNname).Get("inputs").Array()

    return expOutputs
}

func GetExpectedOutputsFieldsForTC(tc *simplejson.Json, tcNname string) []interface{} {
    var expOutputs []interface{}

    expOutputs, _ = tc.Get(tcNname).Get("outputs").Array()

    return expOutputs
}

func GetJsonKeys(expBody interface{}) {
  
  mv := reflect.ValueOf(expBody)
  for _, k := range mv.MapKeys() {
      v := mv.MapIndex(k)
      fmt.Println("aaaaaaaaa")
      fmt.Println("aaaaaaaaa", k, v)
  }
}


// for the dir and sub-dir
func WalkPath(searchDir string, extension string) ([]string, error) {
    fileList := make([]string, 0)

    e := filepath.Walk(searchDir, func(subPath string, f os.FileInfo, err error) error {
        if filepath.Ext(subPath) == extension {
            fileList = append(fileList, subPath)
        }
        return err
    })
    
    if e != nil {
        panic(e)
    }

    // for _, file := range fileList {
    //     fmt.Println(file)
    // }
    return fileList, nil
}


func FileCopy(src string, dest string, info os.FileInfo) error {
    f, err := os.Create(dest)
    if err != nil {
      return err
    }
    defer f.Close()

    if err = os.Chmod(f.Name(), info.Mode())
    err != nil {
      return err
    }

    s, err := os.Open(src)
    if err != nil {
      return err
    }
    defer s.Close()

    _, err = io.Copy(f, s)
    return err
  }


func DirCopy(src string, dest string, info os.FileInfo) error {
    if err := os.MkdirAll(dest, info.Mode())
    err != nil {
      return err
    }

    infos, err := ioutil.ReadDir(src)
    if err != nil {
      return err
    }

    for _, info := range infos {
      if err := FileCopy(filepath.Join(src, info.Name()), filepath.Join(dest, info.Name()), info) 
      err != nil {
        return err
      }
    }

    return nil
}

func ConvertIntArrayToStringArray(intArray []int) []string {
    var stringArray []string
    for _, k := range intArray{
        ii := strconv.Itoa(k)
        stringArray = append(stringArray, ii)
    }

    return stringArray
}

func ConvertStringArrayToIntArray(stringArray []string) []int {
    var intArray []int
    for _, k := range stringArray{
        ii, _ := strconv.Atoi(k)
        intArray = append(intArray, ii)
    }

    return intArray
}

func RemoveArryaItem(sourceArray [][]interface{}, tc []interface{}) [][]interface{} {
    // fmt.Println("RemoveArryaItem", sourceArray, tc)
    var resultArray [][]interface{}
    // resultArray := append(sourceArray[:index], sourceArray[index + 1:]...)
    for index, tc_i := range sourceArray {
        if tc_i[0] == tc[0] {
            resultArray = append(sourceArray[:index], sourceArray[index + 1:]...)
            break
        }
    }

    return resultArray
}


func GenerateFileBasedOnVarAppend(strVar string, filePath string) {
    outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()

    outFile.WriteString(strVar)
}

func GenerateFileBasedOnVarOverride(strVar string, filePath string) {
    outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
       panic(err) 
    }
    defer outFile.Close()

    outFile.WriteString(strVar)
}
    
