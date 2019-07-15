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
	"encoding/xml"

    // "go4api/lib/testcase"
)


type Recurlyservers struct {
    XMLName          xml.Name     `xml:"scxml"` 
    VersionAttr      string       `xml:"version,attr"` 
    DatamodelAttr    string       `xml:"datamodel,attr"`
    Datamodel        Datamodel    `xml:"datamodel"`
    State            []State      `xml:"state"`
    Description      string       `xml:",innerxml"` 
}

type Datamodel struct {
    Data []Data    
}

type Data struct {
    Id string       `xml:"id,attr"` 
}

type State struct {
    StateId string  `xml:"id,attr"`
}