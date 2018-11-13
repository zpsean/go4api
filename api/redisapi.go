/*
 * go4api - a api testing tool written in Go
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

 	gredis "go4api/redis"
)

func RunRedis (stmt string) (int, interface{}, string) {
   	keysCount, cmdResults, redExecStatus := gredis.Run(stmt)

    return keysCount, cmdResults, redExecStatus
}