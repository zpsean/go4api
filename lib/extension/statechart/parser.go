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

            f1(content)

            f2(content) 
        }
    }

    return fullScTcSlice
}

func f1 (content []byte) {
    v := Recurlyservers{}

    err := xml.Unmarshal(content, &v)

    fmt.Println(err)
    fmt.Println(v.XMLName)
    fmt.Println(v.VersionAttr)
    fmt.Println(v.DatamodelAttr)
    fmt.Println(v.Datamodel)
    fmt.Println(v.State)
}
    
func f2 (data []byte) {
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
            if se.Name.Local != "xml" {
                // key = se.Name.Local
            }
            if len(se.Attr) != 0{ 
                fmt.Println("Attrs:", se.Attr)
            }
            fmt.Println("SE.NAME.SPACE:", se.Name.Space) 
            fmt.Println("SE.NAME.LOCAL:", se.Name.Local) 
            fmt.Println()
        case xml.EndElement:
            ee := xml.EndElement(tp)
            if ee.Name.Local == "xml" {
                // return result
            }
            fmt.Println("EE.NAME.SPACE:", ee.Name.Space)
            fmt.Println("EE.NAME.LOCAL:", ee.Name.Local)
        }
    }
}


