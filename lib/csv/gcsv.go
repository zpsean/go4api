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
    "os"
    "fmt"
    "strings"
    "reflect"
    "encoding/csv"
    "encoding/json"
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
    var tmpDataRows [][]string
    for _, gv := range gcsv.DataRows {
        for _, lv := range lcsv.DataRows {
            var tempgv []string
            tempgv = gv
            for _, v := range lv {
                tempgv = append(tempgv, v)
            }
            tmpDataRows = append(tmpDataRows, tempgv)
        }
    }
    gcsv.DataRows = tmpDataRows
    // count
    gcsv.FieldsCount = gcsv.FieldsCount + lcsv.FieldsCount
    gcsv.DataRowsCount = gcsv.DataRowsCount * lcsv.DataRowsCount
}

// like sql union
func (gcsv *Gcsv) Union (latterCsv string) {
    lcsv := GetCsv(latterCsv)
    // validate header count
    if gcsv.FieldsCount != lcsv.FieldsCount {
        fmt.Println("!! Error, two csv have different number of fields, can not Union")
        os.Exit(1)
    }
    // validata and match the hader value, consider the header may not in the same order
    var mFieldsSlice [][]int
    // header
    for i, gv := range gcsv.Header {
        for j, lv := range lcsv.Header {
            var mSlice []int
            if gv == lv {
                mSlice = append(mSlice, i)
                mSlice = append(mSlice, j)

                mFieldsSlice = append(mFieldsSlice, mSlice)
            }
        }
    }
    fmt.Println("mFieldsSlice: ", mFieldsSlice)
    if len(mFieldsSlice) != gcsv.FieldsCount {
        fmt.Println("!! Error, two csv have different header name, can not Union")
        os.Exit(1)
    }
    // union data
    for _, v := range lcsv.DataRows {
        tempSlice := make([]string, gcsv.FieldsCount)
        for _, mSlice := range mFieldsSlice {
            tempSlice[mSlice[0]] = v[mSlice[1]]
        }
        gcsv.DataRows = append(gcsv.DataRows, tempSlice)
    }
    // count
    gcsv.DataRowsCount = gcsv.DataRowsCount + lcsv.DataRowsCount
}

// append the latter csv field [] to gcsv field ([]), which has same field name
func (gcsv *Gcsv) Append (latterCsv string) {
    lcsv := GetCsv(latterCsv)
    // validate, the latter csv can not have more than one data row
    if gcsv.DataRowsCount == 0 {
        fmt.Println("!! Error, please make sure the csv has at least one data row to be appended, can not Append")
        os.Exit(1)
    }
    if lcsv.DataRowsCount != 1 {
        fmt.Println("!! Error, please make sure the latter csv has one and only one data row, can not Append")
        os.Exit(1)
    }
    var mFieldsSlice [][]int
    // header
    for i, gv := range gcsv.Header {
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
            // check the first
            var dataItem interface{}
            json.Unmarshal([]byte(gcsv.DataRows[i][mSlice[0]]), &dataItem)

            if reflect.TypeOf(dataItem) != nil {
                switch reflect.TypeOf(dataItem).Kind() {
                    case reflect.String, reflect.Float64:
                        gcsv.DataRows[i][mSlice[0]] = gcsv.DataRows[i][mSlice[0]] + lcsv.DataRows[0][mSlice[1]]
                    case reflect.Slice:
                        tempSlice := []string{}
                        for _, v := range dataItem.([]interface{}) {
                            tempSlice = append(tempSlice, v.(string))
                        }
                        if len(lcsv.DataRows[0][mSlice[1]]) > 0 {
                            tempSlice = append(tempSlice, lcsv.DataRows[0][mSlice[1]])
                        }
                        // conver the []string to string format
                        tempBytes, _ := json.Marshal(tempSlice)
                        gcsv.DataRows[i][mSlice[0]] = string(tempBytes)
                    // !!! Note: here missed the type check for latterCsv
                }
            // string
            } else {
                gcsv.DataRows[i][mSlice[0]] = gcsv.DataRows[i][mSlice[0]] + lcsv.DataRows[0][mSlice[1]]
            }
        }
    }
}

func (gcsv *Gcsv) WriteToCsvFile () {
    
}


