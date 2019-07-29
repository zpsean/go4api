/*
 * go4api - an api testing tool written in Go
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
    // "time"
    // "os"
    // "sort"
    "io"
    "strings"
    "encoding/xml"
)

// convert the scxml format to xstate

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
