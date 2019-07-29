/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package js

var Results = `
var setUpStartUnixNano = {{.SetUpStartUnixNano}};
var setUpStart = {{.SetUpStart}};
var setUpEndUnixNano = {{.SetUpEndUnixNano}};
var setUpEnd = {{.SetUpEnd}};

var normalStartUnixNano = {{.NormalStartUnixNano}};
var normalStart = {{.NormalStart}};
var normalEndUnixNano = {{.NormalEndUnixNano}};
var normalEnd = {{.NormalEnd}};

var tearDownStartUnixNano = {{.TearDownStartUnixNano}};
var tearDownStart = {{.TearDownStart}};
var tearDownEndUnixNano = {{.TearDownEndUnixNano}};
var tearDownEnd = {{.TearDownEnd}};

var gStartUnixNano = {{.GStartUnixNano}};
var gStart = {{.GStart}};
var gEndUnixNano = {{.GEndUnixNano}};
var gEnd = {{.GEnd}};

var tcResults = {{.TcResults}}
`