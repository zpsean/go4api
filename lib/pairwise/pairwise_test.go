/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */
 
package pairwise

import (
    "fmt"
    "testing"
)

func Test_GetPairWise(t *testing.T) {
	fmt.Println("\n--> test started")

	combins := [][]interface{} {
	    {"Brand X", "Brand Y", "Brand Z", "Brand ZZ"},
	    {"98", "NT", "2000", "XP"},
	    {"Internal", "Modem", "Modem2", "Modem3"},
	    {"Salaried", "Hourly", "Part-Time", "Contr."},
	}

	c := make(chan []interface{})

    go func(c chan []interface{}) {
        defer close(c)
        GetPairWise(c, combins, 3)
    }(c)


    for tcData := range c {
        fmt.Println(tcData)
    }

    fmt.Println("\n--> test finished")
}
