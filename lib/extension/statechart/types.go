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

)

type State struct {
	Id             	string        `json:"id"`
	Type 			string   	  `json:"type"`     // optional: atomic, compound. mandatory: parallel, final, history
	Initial        	string 		  `json:"initial"`
	Entry			[]string      `json:"entry"`    // entry actions
	Exit 			[]string      `json:"exit"`     // exit actions
	On 				map[string]map[string]*TargetAttr   `json:"on"` // map[event]map[targetName]attr
	States		    States   	  `json:"states"`   // map["states"]...
	Activities 		[]string      `json:"activities"`
	// Invoke  		interface{}  // External Communications: send, cancel, invoke, finalize
}

//
type States	map[string]*State  // map[stateName]state

type TargetAttr struct {
	Cond 			string 		`json:"cond"`
	Actions 		[]string    `json:"actions"`  // Executable Content
}

type Transition struct {
	FromState       string
	Event           string
	ToState         string
	Cond            string
	Actions         []string
}


// type Target map[string]*TargetAttr  // map[targetName]attr
