/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package js

var Results = `
var gStartUnixNano = {{.gStart_time}}
var gStart = {{.gStart}}
var pEndUnixNano = {{.PEnd_time}}
var pEnd = {{.PEnd}}
var tcResults = {{.TcReportStr}}
`