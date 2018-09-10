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
    "os"
    "fmt"
    "strings"
    "reflect"
    "encoding/csv"
)

type Gcsv struct {
    Header []string
    DataRows [][]string
    FieldsCount int
    DataRowsCount int
}

func GetCsv(csvcontents string) Gcsv {
    r2 := csv.NewReader(strings.NewReader(csvcontents))
    csvRows, _ := r2.ReadAll()

    var gcsv Gcsv
    var dataRows [][]string
    for i, _ := range csvRows {
        if i == 0 {
            gcsv.Header = csvRows[i]
        } else {
            dataRows = append(dataRows, csvRows[i])
        }
    }
    gcsv.DataRows = dataRows
    gcsv.FieldsCount = len(gcsv.Header)
    gcsv.DataRowsCount = len(gcsv.DataRows)

    return gcsv
}

// like sql join
func (gcsv *Gcsv) Join (latterCsv string) {
    lcsv := GetCsv(latterCsv)
    // join header
    for _, v := range lcsv.Header {
        gcsv.Header = append(gcsv.Header, v)
    }
    // join data, N * M
    for i, _ := range gcsv.DataRows {
        for _, lv := range lcsv.DataRows {
            for _, v := range lv {
                gcsv.DataRows[i] = append(gcsv.DataRows[i], v)
            }
        }
    }
    // count
    gcsv.FieldsCount = gcsv.FieldsCount + lcsv.FieldsCount
    gcsv.DataRowsCount = gcsv.DataRowsCount * lcsv.DataRowsCount
}

// like sql union
func (gcsv *Gcsv) Union (latterCsv string) {
    lcsv := GetCsv(latterCsv)
    // validate
    if gcsv.FieldsCount != lcsv.FieldsCount {
        fmt.Println("!! Error, two csv have different number of fields, can not join")
        os.Exit(1)
    }
    for i := 0; i < gcsv.FieldsCount; i++ {
        if gcsv.Header[i] != lcsv.Header[i] {
            fmt.Println("!! Error, two csv have different header name, can not join")
            os.Exit(1)
        }
    }
    // union data
    for _, v := range lcsv.DataRows {
        gcsv.DataRows = append(gcsv.DataRows, v)
    }
    // count
    gcsv.DataRowsCount = gcsv.DataRowsCount + lcsv.DataRowsCount
}

// append the latter csv field [] to gcsv field ([]), which has same field name
func (gcsv *Gcsv) Append (latterCsv string) {
    lcsv := GetCsv(latterCsv)
    // validate, the latter csv can not have more than one data row
    if gcsv.DataRowsCount == 0 {
        fmt.Println("!! Error, please make sure the csv has at least one data row to be appended, can not append")
        os.Exit(1)
    }
    if lcsv.DataRowsCount != 1 {
        fmt.Println("!! Error, please make sure the latter csv has one and only one data row, can not append")
        os.Exit(1)
    }
    var mFieldsSlice [][]int
    // header
    for i, gv := range lcsv.Header {
        for j, lv := range lcsv.Header {
            var mSlice []int
            if gv == lv {
                mSlice = append(mSlice, i)
                mSlice = append(mSlice, j)

                mFieldsSlice = append(mFieldsSlice, mSlice)
            }
        }
    } 
    // append
    for _, mSlice := range mFieldsSlice {
        for i, _ := range gcsv.DataRows {
            switch reflect.TypeOf(gcsv.DataRows[i][mSlice[0]]).Kind() {
                case reflect.String:
                    gcsv.DataRows[i][mSlice[0]] = gcsv.DataRows[i][mSlice[0]] + lcsv.DataRows[i][mSlice[1]]
                    // gcsv.DataRows[mSlice[0]] = append(gcsv.DataRows[mSlice[0]], lcsv.DataRows[mSlice[1]])
                // case reflect.Slice:
            } 
        }
    }
}

func (gcsv *Gcsv) WriteToCsvFile () {
    
}


