/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package statechart

import (
    "fmt"
    // "os"
    "strings"
    // "bufio"
    "io"
    // "path/filepath"
    "encoding/xml"

    "go4api/utils"
    "go4api/lib/testcase"
)

func InitFullScTcSlice (scfilePathSlice []string) []*testcase.TestCaseDataInfo {
    var fullScTcSlice []*testcase.TestCaseDataInfo
    // var fullKwJsPathSlice []string

    fmt.Println(scfilePathSlice)

    for i, _ := range scfilePathSlice {
        suiteFileListTemp, _ := utils.WalkPath(scfilePathSlice[i], ".scxml")

        for _, path := range suiteFileListTemp {
            content := utils.GetContentFromFile(path)

            XmlDecode(content) 
        }
    }

    return fullScTcSlice
}
    
func XmlDecode (data []byte) {
    decoder := xml.NewDecoder(strings.NewReader(string(data)))

    // result  := make(map[string]string)
    // key := ""

    for {
        token, err := decoder.Token()
    
        if err == io.EOF{
            fmt.Println("parse Finish")
            break
             // return result
        }

        if err != nil{
            fmt.Println("parse Fail:",err)
            break
            // return result
        }

        switch tp := token.(type) {
        case xml.StartElement:
            se := xml.StartElement(tp) 
 
            fmt.Println("se.Name:", se.Name) 
            fmt.Println("se.Attr:", se.Attr)
            fmt.Println()
        case xml.EndElement:
            ee := xml.EndElement(tp)

            fmt.Println("ee.Name:", ee.Name) 
            fmt.Println()
        case xml.CharData:
            // cd := xml.CharData(tp)

            // data := string(cd)
            // fmt.Println("cd.content:", data) 
            // fmt.Println()
        }
    }
}





