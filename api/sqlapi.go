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
    
    gsql "go4api/sql"
)

func RunSql (stmt string) (int, []string, []map[string]interface{}, string) {
    // gsql.Run will return: <impacted rows : int>, <rows for select : [][]interface{}{}>, <sql status : string>
    // status: SqlSuccess, SqlFailed
    rowsCount, rowsHeaders, rowsData, sqlExecStatus := gsql.Run(stmt)

    return rowsCount, rowsHeaders, rowsData, sqlExecStatus
}


