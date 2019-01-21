/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package gexcel

import (
    "fmt"
    "os"
    "strings"

    "go4api/cmd"

    xlsx "github.com/tealeg/xlsx"
)

func Run() {
    var file *xlsx.File
    var sheet *xlsx.Sheet
    var row *xlsx.Row
    var cell *xlsx.Cell
    var err error

    file = xlsx.NewFile()
    sheet, err = file.AddSheet("Sheet1")
    if err != nil {
        fmt.Printf(err.Error())
    }
    row = sheet.AddRow()
    cell = row.AddCell()
    cell.Value = "I am a cell!"
    err = file.Save("MyXLSXFile.xlsx")
    if err != nil {
        fmt.Printf(err.Error())
    }
}

