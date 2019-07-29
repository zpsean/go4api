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
    // "os"
    "io/ioutil"
    "net/http"
    "fmt"
)


func SaveHttpRespFile (actualBody []byte, outputsFile string) {
    // filePath := cmd.Opt.Testresource

    ioutil.WriteFile(outputsFile, actualBody, 0644)

    // use the first 512 bytes for type
    buffer := make([]byte, 512)
    buffer = actualBody[0:512]

    contentType := http.DetectContentType(buffer)
    fmt.Println("contentType: ", contentType)
}


