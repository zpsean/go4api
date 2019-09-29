/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2019.07
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package testsuite

import (
    // "fmt"

    gsession "go4api/lib/session"
)

// set the session info for test suite
func (ts *TestSuite) WriteSession () {
    gsession.WriteTcSession(ts.TsName(), (*ts)[ts.TsName()].Parameters)
}
