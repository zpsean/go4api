package protocal

import (
    // "fmt"
    "net/http" 
    "io/ioutil"
    "strconv"
    "reflect"
    "strings"
    "bytes"
)


func HttpGet(url string, apiMethod string, reqHeaders map[string]interface{}, reqBody *strings.Reader) (string, http.Header, string) { 
    //client 
    client := &http.Client{}

    // payload := reqBody
    //
    reqest, err := http.NewRequest(apiMethod, url, nil)

    //Header
    mv := reflect.ValueOf(reqHeaders)
    for _, k := range mv.MapKeys() {
        v := mv.MapIndex(k)
        reqest.Header.Add(k.Interface().(string), v.Interface().(string))
    }

    // for _, k := range mv.MapKeys() {
    //     v := mv.MapIndex(k)
    //     // fmt.Println("reqHeaders", k, v)
    //     if k.Interface().(string) == "authorization" {
    //         reqest.Header.Add(k.Interface().(string), "afafddsfas")
    //     } else {
    //         reqest.Header.Add(k.Interface().(string), v.Interface().(string))
    //     }
    // }

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

    // protocalChan <- 1

    return strconv.Itoa(response.StatusCode), response.Header, string(body)
}

func HttpPost(url string, apiMethod string, reqHeaders map[string]interface{}, reqBody *strings.Reader) (string, http.Header, string) { 
    //client 
    client := &http.Client{}
    //
    // payload := strings.NewReader("{\"mid\":\"550049154\"}")
    payload := reqBody
    //
    reqest, err := http.NewRequest(apiMethod, url, payload)

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

    return strconv.Itoa(response.StatusCode), response.Header, string(body)
}

func HttpPostMultipart(url string, apiMethod string, reqHeaders map[string]interface{}, reqBody *bytes.Buffer) (string, http.Header, string) { 
    //client 
    client := &http.Client{}
    //
    // payload := strings.NewReader("{\"mid\":\"550049154\"}")
    payload := reqBody
    //
    reqest, err := http.NewRequest(apiMethod, url, payload)

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

    return strconv.Itoa(response.StatusCode), response.Header, string(body)
}


func CallHttpMethod(m map[string]interface{}, name string, params ... interface{}) (string, http.Header, string) {
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

    return result[0].Interface().(string), result[1].Interface().(http.Header), result[2].Interface().(string)
}
