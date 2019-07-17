/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2019
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package statechart

import (
    "fmt"
    // "time"
    // "os"
    // "sort"
    // "strings"
)

// set
func (st *State) SetStateIds () {
	for k, _ := range st.States {
		st.States[k].Id = k

		st.States[k].SetStateIds()
	}
}

// get
func (st *State) GetStateIds () {
	fmt.Println("--> id: ", st.Id)
	for k, v := range st.States {
		fmt.Println("-->: ", k, v)

		v.GetStateIds()
	}
}

func (st *State) GetStateTransitions () {
	fmt.Println("--> tx: ", st.Id, st.On)
	for _, v := range st.States {
		// fmt.Println("-->: ", k, v)

		v.GetStateTransitions()
	}
}

