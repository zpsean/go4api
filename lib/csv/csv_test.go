/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gcsv

import (
    "fmt"
    "testing"
    // "encoding/csv"
)

func Test_GetCsv(t *testing.T) {
    fmt.Println("\n--> test started")

    csvcontents := `h1,h2,h3
    d11,d12,d13
    d21,d22,d23`
    gcsv := GetCsv(csvcontents)

    fmt.Println("gcsv: ", gcsv)
    fmt.Println("gcsv: ", gcsv.Header)
    fmt.Println("gcsv: ", gcsv.DataRows)
    fmt.Println("gcsv: ", gcsv.FieldsCount)
    fmt.Println("gcsv: ", gcsv.DataRowsCount)

    fmt.Println("\n--> test finished")
}


