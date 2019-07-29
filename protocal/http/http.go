/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package g4http

import (
    "fmt"
    "net/http"
    "io/ioutil"
    // "strconv"
    "reflect"
    "strings"
    "bytes"
)

type HttpRequest interface {
    Request(urlStr string, apiMethod string, reqHeaders map[string]interface{}, reqBody interface{}) (int, http.Header, []byte)
}

type HttpRestful struct{}

func (httpRestful HttpRestful) Request(urlStr string, apiMethod string, reqHeaders map[string]interface{}, reqBody interface{}) (int, http.Header, []byte) { 
    //client 
    client := &http.Client{}
    //
    // type conversion to payload, based on reqBody, apiMethod
    var reqest *http.Request
    var err error
    switch reflect.TypeOf(reqBody).String() {
        case "*strings.Reader":
            if apiMethod == "GET" {
                reqest, err = http.NewRequest(apiMethod, urlStr, nil)
            } else {
                reqest, err = http.NewRequest(apiMethod, urlStr, reqBody.(*strings.Reader))
            }
        case "*bytes.Buffer":
            reqest, err = http.NewRequest(apiMethod, urlStr, reqBody.(*bytes.Buffer))
    }
    if err != nil {
        panic(err)
    }
    //Header
    for key, value := range reqHeaders {
        reqest.Header.Add(key, fmt.Sprint(value))
    }
    //response
    response, err := client.Do(reqest)
    if err != nil {
        panic(err)
    } 
    defer response.Body.Close()

    body, _ := ioutil.ReadAll(response.Body)
    // fmt.Print("-------> body: ", reqBody, string(body))

    var actualHeader = make(map[string][]string)
    for k, v := range response.Header {
        actualHeader[k] = v
    }

    return response.StatusCode, actualHeader, body
}

