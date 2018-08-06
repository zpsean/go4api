/*
 * go4api - a api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package ui

var Index_template = `<h1>Go4Api Executions</h1>
<body>
    <table border="1">
        <col width="20" />
        <col width="300" />
        <col width="20" />
        <col width="500" />
        <tr style='text-align: left'>
            <th>#</th>
            <th>Case ID</th>
            <th>Status</th>
            <th>Case File / Data Table / Data Row</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.Seq}}</td>
            <td>{{.CaseID}}</td>
            <td>{{.Status}}</td>
            <td>{{.CasePath}}</td>
        </tr>
        {{end}}
    </table>

</body>
</html>`