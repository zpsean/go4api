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

var Stats = `
var stats1 = {{.StatsStr_1}};

var stats2 = {{.StatsStr_2}};

var stats2_success = {{.StatsStr_Success}};

var stats2_fail = {{.StatsStr_Fail}};

var stats3_status = {{.StatsStr_Status}};
`