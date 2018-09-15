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

var Js = `
function tabPageControl(n) {
	for (var i = 0; i < tabTitles.cells.length; i++) {
		tabTitles.cells[i].className = "tabTitleUnSelected";
	}
	tabTitles.cells[n].className = "tabTitleSelected";

	for (var i = 0; i < tabPagesContainer.tBodies.length; i++) {
		tabPagesContainer.tBodies[i].className = "tabPageUnSelected";
	}
	tabPagesContainer.tBodies[n].className = "tabPageSelected";
}
`