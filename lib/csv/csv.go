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
    "strings"
    "encoding/csv"
)

type Gcsv struct {
    Header []string
    DataRows [][]string
    FieldsCount int
    DataRowsCount int
}

func (gcsv *Gcsv) Join (lGcsv Gcsv) {

}

func (gcsv *Gcsv) Union (lGcsv Gcsv) {
    
}

func (gcsv *Gcsv) Append (lGcsv Gcsv) {
    
}

func GetCsv(r string) Gcsv {
    r2 := csv.NewReader(strings.NewReader(r))
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


