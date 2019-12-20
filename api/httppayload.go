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
    "os"
    "bytes"
    "net/textproto"
    "mime/multipart"
    "io"    
    "net/url"  
    "strings"
    "encoding/json"

    "go4api/cmd"
    "go4api/lib/testcase" 

    gjson "github.com/tidwall/gjson"
)

func GetPayloadInfo (tcData *testcase.TestCaseDataInfo) (string, string, *strings.Reader, *bytes.Buffer, string) {
    apiMethod := tcData.TestCase.ReqMethod()
    // request payload(body)
    reqPayload := tcData.TestCase.ReqPayload()
    //
    var bodyText *strings.Reader // init body
    bodyMultipart := &bytes.Buffer{}
    boundary := ""
    //
    apiMethodSelector := apiMethod
    // Note, has 3 conditions: text (json), form, or multipart file upload
    for key, _ := range reqPayload {
        switch key {
            case "multipart-form": {
                //multipart/form-data
                apiMethodSelector = "POSTMultipart"
                multipartFilePath := cmd.Opt.Testresource

                if string(multipartFilePath[len(multipartFilePath) - 1]) != "/" {
                    multipartFilePath = multipartFilePath + "/"
                }

                bodyMultipart, boundary, _ = PrepMultipart(reqPayload, multipartFilePath)
            }
            case "text": {
                //application/json
                bodyText = PrepPostPayload(reqPayload)
            }
            case "form": {
                //application/x-www-form-urlencoded
                bodyText = PrepPostFormPayload(reqPayload)
            }
            default: {
                bodyText = strings.NewReader("")
            }
        }
    }

    return apiMethodSelector, apiMethod, bodyText, bodyMultipart, boundary
}

func fileOpen (path string, fileName string) *os.File {
    fp, err := os.Open(path + fileName) 
    if err != nil {
        panic(err)
    }

    return fp
}

func PrepMultipart (reqPayload map[string]interface {}, path string) (*bytes.Buffer, string, error) {
    body := &bytes.Buffer{} // init body
    writer := multipart.NewWriter(body) // multipart
    //
    reqPayloadJsonBytes, _ := json.Marshal(reqPayload)
    reqPayloadJson := string(reqPayloadJsonBytes)

    var i int64
    var err error

    total := gjson.Get(reqPayloadJson, "multipart-form.#")
    for i = 0; i < total.Int(); i++ {
        item := gjson.Get(reqPayloadJson, "multipart-form." + fmt.Sprint(i))
        //
        ifFile := false
        iMap := item.Map()
        for k, v := range iMap {
            if k == "type" && v.String() == "file" {
                ifFile = true
                break
            }
        }
        //
        var fw io.Writer
        var r io.Reader
        if ifFile == true {
            n := iMap["name"].String()
            v := iMap["value"].String()

            h := make(textproto.MIMEHeader)
            h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, n, v))

            for k, v := range iMap {
                if k != "name" || k != "value" || k != "type" {
                    h.Set(k, v.String())
                }
            }

            if fw, err = writer.CreatePart(h); err != nil {
                return nil, "", err
            }

            fp := fileOpen(path, v)
            defer fp.Close()

            if _, err = io.Copy(fw, fp); err != nil {
                panic(err)
                return nil, "", err
            }
        } else {
            n := iMap["name"].String()
            v := iMap["value"].String()

            r = strings.NewReader(v)

            if fw, err = writer.CreateFormField(n); err != nil {
                return nil, "", err
            }

            if _, err = io.Copy(fw, r); err != nil {
                panic(err)
                return nil, "", err
            }
        }
    }

    err = writer.Close()
    if err != nil {
        return nil, "", err
    }
    // Note: do not forget this, to get the boundary value
    boundary := writer.FormDataContentType()
  
    return body, boundary, nil
}

func PrepPostPayload (reqPayload map[string]interface{}) *strings.Reader {
    var body *strings.Reader

    for key, value := range reqPayload {
        if key == "text" {
            repJson, _ := json.Marshal(value)
            body = strings.NewReader(string(repJson))
            break
        }
    }

    return body
}

func PrepPostFormPayload (reqPayload map[string]interface{}) *strings.Reader {
    var body *strings.Reader
    //
    reqPayloadJsonBytes, _ := json.Marshal(reqPayload)
    reqPayloadJson := string(reqPayloadJsonBytes)
    formMap := gjson.Get(reqPayloadJson, "form").Map()

    data := url.Values{}
    for k, v := range formMap {
        data.Set(k, v.String())
    }
    body = strings.NewReader(data.Encode())

    return body
}

