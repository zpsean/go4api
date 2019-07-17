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

)

type State struct {
	Id             	string        `json:"id"`
	Type 			string   	  `json:"type"`     // optional: atomic, compound. mandatory: parallel, final, history
	Initial        	string 		  `json:"initial"`
	Entry			[]string      `json:"entry"`    // entry actions
	Exit 			[]string      `json:"exit"`     // exit actions
	On 				*Transition   `json:"on"`       // map["on"]...
	States		    States   	  `json:"states"`   // map["states"]...
	Activities 		[]string      `json:"activities"`
	// Invoke  		interface{}  // External Communications: send, cancel, invoke, finalize
}

//
type States	map[string]*State  // map[stateName]state

type Transition map[string]*Target  // map[event]target

type Target map[string]*TargetAttr  // map[targetName]attr

type TargetAttr struct {
	Cond 			string 		`json:"cond"`
	Actions 		[]string    `json:"actions"`  // Executable Content
}
