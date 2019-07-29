/*
 * go4api - an api testing tool written in Go
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
    "sync"
    // "time"
    // "os"
    // "sort"
    // "strings"
)

// Note: Moore -> Mealy 100%, but Mealy -> Moore possible 100%
// Here take care Mealy only

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

func (st *State) GetStateTransitions (ch chan *Transition, wg *sync.WaitGroup) {	
	defer wg.Done()

	for event, v1 := range st.On {
		for target, v2 := range v1 {
			transition := &Transition {
				FromState:  st.Id,
				Event:      event,
				ToState:    target,
				Cond:       v2.Cond,
				Actions:    v2.Actions,
			}
			
			ch <- transition
		}
	}

	// next
	for _, v := range st.States {
		wg.Add(1)
		go v.GetStateTransitions(ch, wg)
	}
}

// get state - event (json node: On)
func (st *State) GetStateEvents () {

}

// get Transition event - target names
func (st *State) GetTransitionInfos () {

}

