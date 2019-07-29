/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package style

var Style = `
body {
	font-family: arial;
	margin: 0;
	padding: 0;
	background: #DCDBD4;
	color: #000;
	font-size: 12px;
}

.container {
	width: 1200px;
	margin: 0 auto;
	position: relative;
	overflow: visible;
}


a {
	text-decoration: none;
	color: #B79F39;
}

a:hover {
	color: #D1D1CE;
}

p {
	margin: 10px 0;
	padding: 0;
}

img {
	border: 0;
}

h1 {
	font-size: 18px;
	color: #B79F39;
	margin: 10px 0;
}

h1 span {
	color: #D1D1CE;
}

.main {
	width: 100%
}

.head {
	height: 100px;
	width: 600px;
	padding: 10px 0 0 160px;
	background: no-repeat;
	background-position: 70px 0;
	position: absolute;
	top: 0;
	left: 0;
}

.sim_desc {
	float: right;
	padding: 0;
	margin: 0
}

.foot {
	background: #92918C;
	width: 100%;
	color: #FFF;
	margin-top: 78px;
}

.foot a {
	width: 85px;
	height: 27px;
	margin: 3px auto;
	display: block
}

.skeleton {
	width: 1150px;
	padding: 38px 25px 0 25px;
	position: relative;
	top: 70px;
}


.content {
	margin-right: 5px;
	padding: 0 10px 40px 1px;
	background: #FFF;
	box-shadow: -4px -4px 3px #D1D1CE, 4px 4px 3px #D1D1CE, -4px 4px 3px #D1D1CE, 4px -4px 3px #D1D1CE;
	-moz-box-shadow: -4px -4px 3px #D1D1CE, 4px 4px 3px #D1D1CE, -4px 4px 3px #D1D1CE, 4px -4px 3px #D1D1CE;
	-webkit-box-shadow: -4px -4px 3px #D1D1CE, 4px 4px 3px #D1D1CE, -4px 4px 3px #D1D1CE, 4px -4px 3px #D1D1CE;

	border-top-left-radius: 8px;
	border-bottom-left-radius: 8px;
	border-bottom-right-radius: 8px;
	-moz-border-top-left-radius: 8px;
	-moz-border-bottom-left-radius: 8px;
	-moz-border-bottom-right-radius: 8px;
	-webkit-border-top-left-radius: 8px;
	-webkit-border-bottom-left-radius: 8px;
	-webkit-border-bottom-right-radius: 8px;

	border-color: #D1D1CE
}

.content-in {
	margin: 30px 30px 30px 30px;
}

.sous-menu {
	border-bottom: 2px #B79F39 solid;
	padding: 10px 0 5px 0;
	margin: 0 5px;
	z-index: -1;
}

.sous-menu a {
	color: #FFF;
	font-size: 14px;
	font-weight: bold;
}

.sous-menu .item {
	background: #D1D1CE url('arrow_right.png') 11px 3px no-repeat;
	display: inline;
	margin: 0 10px 0 0;
	padding: 5px 15px 5px 25px;
	border-top-right-radius: 8px;
	border-top-left-radius: 8px
}

.sous-menu .selected {
	background: #B79F39 url('arrow_down.png') 9px 6px no-repeat;
}

.article {
	position: relative;
}

.schema {
	background: #EAEAEA;
	border-radius: 8px;
	border: 1px solid #EAEAEA;
	margin-bottom: 20px
}


.p_left {
	width: 800px;
	height: 250px;
}

.p_right {
	width: 250px;
	height: 250px;
	position: absolute;
	right: -17px;
}

.geant {
	width: 1089px;
	height: 362px;
}

.extensible-geant {
	width: 1089px;
}



.chart_title {
	background: #b0b0a8;
	padding: 2px 10px;
	color: #FFF;
	font-weight: bold;
	border-radius: 8px;
	font-size: 16px;
}

.statistics {
	margin-bottom: 20px;
	color: white;
	display: block;
}

.statistics .title {
	height: 15px;
	text-align: left;
	font-weight: bold;
	background: #997A26;
	border-top-left-radius: 8px;
	border-top-right-radius: 8px;
	padding: 5px 15px 0px 15px;
	width: 1059px;
}

.title_collapsed {
	background: #997A26 0px 0px no-repeat;
	cursor: pointer;
	padding-left: 20px;
}

#container_statistics_head {
	background: #997A26;
	padding: 5px 5px 0px 5px;
}

#container_statistics_body {
	border-bottom-left-radius: 8px;
	border-bottom-right-radius: 8px;
	background: #B79F39 repeat-x;
	padding: 0px 5px 5px 5px;
}

#container_statistics_body2 {
	border-bottom-left-radius: 8px;
	border-bottom-right-radius: 8px;
	background: #B79F39 repeat-x;
	padding: 0px 5px 5px 5px;
}

#container_statistics_body3 {
	border-bottom-left-radius: 8px;
	border-bottom-right-radius: 8px;
	background: #B79F39 repeat-x;
	padding: 0px 5px 5px 5px;
}

.statistics-in {
	margin: 0;
	border-spacing: 4px;
}

.statistics .scrollable {
	max-height: 785px;
	overflow-y: auto;
	width: 1089px;
}

.statistics-in a {
	color: white;
	font-weight: bold;
}

.statistics-in .header {
	font-size: 13px;
	font-weight: bold;
	border-radius: 4px;
	padding: 4px;
	background-color: #4D6B00;
}

.sortable {
	cursor: pointer;
}

.sortable span {
	padding-right: 10px;
}

.executions {
	padding: 4px 4px 4px 25px;
}

.response-time {
	padding: 4px 4px 4px 25px;
}

.statistics-in td {
	padding: 4px;
	border-radius: 4px;
}

.statistics-in .col-1 {
	width: 175px;
}

.statistics-in .value {
	text-align: right;
	width: 50px;
}


.test-table {
	border-collapse:collapse;
}
`
