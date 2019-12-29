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

    gsql "go4api/db/sqldb"
    // gpg  "go4api/db/postgres"
)

// for mysql
func RunSql (tgtDb string, stmt string) (int, []string, []map[string]interface{}, string) {
    // gsql.Run will return: <impacted rows : int>, <rows for select : [][]interface{}{}>, <sql status : string>
    // status: SqlSuccess, SqlFailed
    rowsCount, rowsHeaders, rowsData, sqlExecStatus := gsql.Run("mysql", tgtDb, stmt)

    return rowsCount, rowsHeaders, rowsData, sqlExecStatus
}

// for postgresql
func RunPgSql (tgtDb string, stmt string) (int, []string, []map[string]interface{}, string) {
    // gsql.Run will return: <impacted rows : int>, <rows for select : [][]interface{}{}>, <sql status : string>
    // status: SqlSuccess, SqlFailed
    rowsCount, rowsHeaders, rowsData, sqlExecStatus := gsql.Run("postgres", tgtDb, stmt)

    return rowsCount, rowsHeaders, rowsData, sqlExecStatus
}

