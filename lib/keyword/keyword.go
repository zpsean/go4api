/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package keyword

import (
    // "fmt"
    // "time"
    // "os"
    // "sort"
)

// for keyword to ts / tc execution, there are two options:
// option 1: convert to ts / tc (i.e. temp files) format then execution
// option 2: mapping the keywords with existing ts/tc (i.e. treat the ts/tc as library, but can accept params)
//
// here use option 2, steps are:
// (1). scan the keyword files, get all the testcases and their variables
// (2). lookup ts/tc, generate the cases to be executed
// (3). reporting
// i.e. go4api -run -K -kw xx/xx/*.keyword -tc xxxx/ -tsuite xxxxx -jsFuncs xxx -r xxx -tr xxx

func (kw *KWBlock) ddd () {
	
}