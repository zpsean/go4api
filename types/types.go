/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package types

type TcRunResults struct {  
    TcName string
    ParentTestCase string
    TestResult string
    ActualStatusCode string
    JsonFile_Base string
    CsvFileBase string
    RowCsv string
    Start string
    End string
    TestMessages string
    Start_time_UnixNano int64
    End_time_UnixNano int64
    Duration_UnixNano int64
}


type TcReportResults struct {  
    Priority string
    TcRunRes TcRunResults
}

// valid, invalid data may have more than one field, but the map itself can not ensure the key sequence
// so that, here use slice
type FuzzData struct {  
    ValidData []map[string][]interface{}
    InvalidData []map[string][]interface{}
    ValidStatusCode int
    InvalidStatusCode int
}







