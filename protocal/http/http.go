package protocal

import (
    "fmt"
    "net/http"
    "io/ioutil"
    // "strconv"
    "reflect"
    "strings"
    "bytes"
)


func HttpGet(urlStr string, apiMethod string, reqHeaders map[string]interface{}, reqBody *strings.Reader) (int, http.Header, []byte) { 
    //client 
    client := &http.Client{}

    // payload := reqBody
    fmt.Println("urlStr", urlStr)
    //
    reqest, err := http.NewRequest(apiMethod, urlStr, nil)

    //Header
    mv := reflect.ValueOf(reqHeaders)
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        reqest.Header.Add(k.Interface().(string), v.Interface().(string))
    }
    // fmt.Println("url", url)
    // fmt.Println("payload", payload)
    // fmt.Println("reqest.Header", reqest.Header)

    //response
    response, err := client.Do(reqest)
    if err != nil {
        panic(err)
    } 
    defer response.Body.Close()

    //func ReadAll(r io.Reader) ([]byte, error)
    body, _ := ioutil.ReadAll(response.Body)

    return response.StatusCode, response.Header, body
}

func HttpPost(urlStr string, apiMethod string, reqHeaders map[string]interface{}, reqBody *strings.Reader) (int, http.Header, []byte) { 
    //client 
    client := &http.Client{}
    //
    payload := reqBody

    reqest, err := http.NewRequest(apiMethod, urlStr, payload)

    //Header
    mv := reflect.ValueOf(reqHeaders)
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        reqest.Header.Add(k.Interface().(string), v.Interface().(string))
    }
    // fmt.Println("url", url)
    // fmt.Println("payload", payload)
    // fmt.Println("reqest.Header", reqest.Header)

    //response
    response, err := client.Do(reqest)
    if err != nil {
        panic(err)
    }  
    defer response.Body.Close()
    body, _ := ioutil.ReadAll(response.Body)
    // fmt.Println(response.Header)
    // fmt.Println(response.StatusCode)
    // fmt.Println(string(body))

    return response.StatusCode, response.Header, body
}

func HttpPostForm(urlStr string, apiMethod string, reqHeaders map[string]interface{}, reqBody *strings.Reader) (int, http.Header, []byte) { 
    //client 
    client := &http.Client{}
    //
    payload := reqBody

    reqest, err := http.NewRequest(apiMethod, urlStr, payload)

    //Header
    mv := reflect.ValueOf(reqHeaders)
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        reqest.Header.Add(k.Interface().(string), v.Interface().(string))
    }
    // fmt.Println("url", url)
    // fmt.Println("payload", payload)
    // fmt.Println("reqest.Header", reqest.Header)

    //response
    response, err := client.Do(reqest)
    if err != nil {
        panic(err)
    }  
    defer response.Body.Close()
    body, _ := ioutil.ReadAll(response.Body)
    // fmt.Println(response.Header)
    // fmt.Println(response.StatusCode)
    // fmt.Println(string(body))

    return response.StatusCode, response.Header, body
}

func HttpPostMultipart(urlStr string, apiMethod string, reqHeaders map[string]interface{}, reqBody *bytes.Buffer) (int, http.Header, []byte) { 
    //client 
    client := &http.Client{}
    //
    payload := reqBody
    
    reqest, err := http.NewRequest(apiMethod, urlStr, payload)

    //Header
    mv := reflect.ValueOf(reqHeaders)
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        reqest.Header.Add(k.Interface().(string), v.Interface().(string))
    }
    // fmt.Println("url", url)
    // fmt.Println("payload", payload)
    // fmt.Println("reqest.Header", reqest.Header)

    //response
    response, err := client.Do(reqest)
    if err != nil {
        panic(err)
    }  
    defer response.Body.Close()
    body, _ := ioutil.ReadAll(response.Body)
    // fmt.Println(response.Header)
    // fmt.Println(response.StatusCode)
    // fmt.Println(string(body))

    // protocalChan <- "aaa"

    return response.StatusCode, response.Header, body
}


func CallHttpMethod(m map[string]interface{}, name string, params ... interface{}) (int, http.Header, []byte) {
    f := reflect.ValueOf(m[name])
    // if len(params) != f.Type().NumIn() {
    //     err = errors.New("The number of params is not adapted.")
    //     return
    // }
    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }

    result := f.Call(in)

    return result[0].Interface().(int), result[1].Interface().(http.Header), result[2].Interface().([]byte)
}
