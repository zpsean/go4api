
function mergeTcResults(setUpTcResults, normalTcResults, tearDownTcResults) {
	var c = setUpTcResults.concat(normalTcResults);
	var dest = c.concat(tearDownTcResults);

	return dest
}
