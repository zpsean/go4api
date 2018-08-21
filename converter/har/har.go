/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package har

import (
    // "fmt"
    // "time"
    // "os"
    // "sort"
    // "sync"
    // "go4api/types" 
)

// test case type
type HarLog struct {
    Version string
    Creator map[string]string
    Pages []map[string]interface{}
    Entries []Entry
}

type Entries []Entry

type Entry struct {
    startedDateTime string
    time float
    Request Request
    Response Response
}

type Request struct {  
    Method string
    Url string
    HttpVersion string
    Headers []map[string]interface{}
    QueryString []map[string]interface{}
    Cookies []map[string]interface{}
    HeadersSize int
    BodySize int
    PostData map[string]interface{}
}


type Response struct {  
    Status int
    StatusText string
    HttpVersion string
    Headers map[string]interface{}
    Cookies []map[string]interface{}
    content map[string]interface{}
}
