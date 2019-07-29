/*
 * go4api - an api testing tool written in Go
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
    d21,12,d23`
    gcsv := GetCsv(csvcontents)

    fmt.Println("gcsv: ", gcsv)
    fmt.Println("gcsv: ", gcsv.Header)
    fmt.Println("gcsv: ", gcsv.DataRows)
    fmt.Println("gcsv: ", gcsv.FieldsCount)
    fmt.Println("gcsv: ", gcsv.DataRowsCount)

    if gcsv.FieldsCount != 3 || gcsv.DataRowsCount != 2 {
        t.Fatalf("csv load unexpected")
    }

    fmt.Println("\n--> test finished")
}

func Test_GetCsv2(t *testing.T) {
    fmt.Println("\n--> test started")

    csvcontents := `h1,h2,h3`
    gcsv := GetCsv(csvcontents)

    fmt.Println("gcsv: ", gcsv)
    fmt.Println("gcsv: ", gcsv.Header)
    fmt.Println("gcsv: ", gcsv.DataRows)
    fmt.Println("gcsv: ", gcsv.FieldsCount)
    fmt.Println("gcsv: ", gcsv.DataRowsCount)

    if gcsv.FieldsCount != 3 || gcsv.DataRowsCount != 0 {
        t.Fatalf("csv load unexpected")
    }

    fmt.Println("\n--> test finished")
}

func Test_Join(t *testing.T) {
    fmt.Println("\n--> test started")

    gcsvcontents := `h1,h2,h3
    d11,d12,d13
    d21,12,d23`
    gcsv := GetCsv(gcsvcontents)

    lcsvcontents := `h1,h2,h3
    e11,e12,e13
    e21,23,e23`

    fmt.Println("gcsv: ", gcsv)
    gcsv.Join(lcsvcontents)
    fmt.Println("gcsv: ", gcsv)

    if gcsv.FieldsCount != 6 {
        t.Fatalf("csv joined FieldsCount not correct, header")
    }
    if len(gcsv.DataRows) != 4 {
        t.Fatalf("csv joined FieldsCount not correct, datarows")
    }

    fmt.Println("\n--> test finished")
}


func Test_Union(t *testing.T) {
    fmt.Println("\n--> test started")

    gcsvcontents := `h1,h2,h3
    d11,d12,d13
    d21,12,d23`
    gcsv := GetCsv(gcsvcontents)

    lcsvcontents := `h1,h2,h3
    e11,e12,e13
    e21,23,e23`

    fmt.Println("gcsv: ", gcsv)
    gcsv.Union(lcsvcontents)

    fmt.Println("gcsv: ", gcsv)
    fmt.Println("\n--> test finished")
}


func Test_Append(t *testing.T) {
    fmt.Println("\n--> test started")

    gcsvcontents := `h1,h2,h3
    d11,d12,d13
    d21,12,d23`
    gcsv := GetCsv(gcsvcontents)

    lcsvcontents := `h11,h22,h2
    e11,e12,e13`

    fmt.Println("gcsv - before: ", gcsv)
    gcsv.Append(lcsvcontents)
    fmt.Println("gcsv - after: ", gcsv)

    fmt.Println("\n--> test finished")
}

func Test_Append2(t *testing.T) {
    fmt.Println("\n--> test started")

    gcsvcontents := `h1,h2,h3
    d11,"[""aa"", ""bb"", ""cc""]",d13
    d21,12,d23`
    gcsv := GetCsv(gcsvcontents)

    lcsvcontents := `h11,h2,h23
    e11,56,e13`

    fmt.Println("gcsv - before: ", gcsv)
    fmt.Println("lcsv - before: ", lcsvcontents)
    gcsv.Append(lcsvcontents)
    fmt.Println("gcsv - after: ", gcsv)

    fmt.Println("\n--> test finished")
}