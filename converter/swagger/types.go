/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package swagger

import ( 
)

// swagger type

type Swagger2 struct {
	Swagger string
	Info map[string]interface{}
	BasePath string
	Tags map[string]interface{}
	Schemes []string
	Paths map[string]Path
	SecurityDefinitions map[string]interface{}
	Definitions map[string]interface{}
	ExternalDocs map[string]interface{}
}

type Path map[string]PathDetails

type PathDetails struct {
	Method string
	Consumes []string // content-type?
	Produces []string // content-type?
	Parameters []map[string]interface{}
	Responses []string
	Security []map[string]interface{}
}

type Definitions map[string]Definition

type Definition struct {
	Type string
	Required []string 
	Properties map[string]interface{}
	Xml map[string]interface{}
}


