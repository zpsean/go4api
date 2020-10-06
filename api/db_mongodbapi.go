/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package api

import (
 	// "fmt"

 	gmongodb "go4api/db/mongodb"
)

func RunMongoDB (cmdStr string) (int, interface{}, string) {
   	keysCount, cmdResults, redExecStatus := gmongodb.Run(cmdStr)

    return keysCount, cmdResults, redExecStatus
}

