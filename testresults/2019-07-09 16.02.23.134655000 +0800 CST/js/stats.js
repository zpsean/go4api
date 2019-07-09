
var stats1 = [
	{
		"ReportKey": {
			"IfGlobalSetUpTearDown": "GlobalSetUp",
			"Priority": "2",
			"TestResult": "Fail"
		},
		"Count": 1,
		"PerformanceGauge": {
			"Min": 164,
			"P50": 164,
			"P75": 164,
			"P95": 164,
			"P99": 164,
			"Max": 164,
			"Mean": 164,
			"StdDev": 0
		}
	},
	{
		"ReportKey": {
			"IfGlobalSetUpTearDown": "GlobalSetUp",
			"Priority": "2",
			"TestResult": "ParentFailed"
		},
		"Count": 1,
		"PerformanceGauge": {
			"Min": 0,
			"P50": 0,
			"P75": 0,
			"P95": 0,
			"P99": 0,
			"Max": 0,
			"Mean": 0,
			"StdDev": 0
		}
	},
	{
		"ReportKey": {
			"IfGlobalSetUpTearDown": "GlobalSetUp",
			"Priority": "3",
			"TestResult": "ParentSkipped"
		},
		"Count": 2,
		"PerformanceGauge": {
			"Min": 0,
			"P50": 0,
			"P75": 0,
			"P95": 0,
			"P99": 0,
			"Max": 0,
			"Mean": 0,
			"StdDev": 0
		}
	},
	{
		"ReportKey": {
			"IfGlobalSetUpTearDown": "GlobalSetUp",
			"Priority": "2",
			"TestResult": "ALL"
		},
		"Count": 2,
		"PerformanceGauge": {
			"Min": 164,
			"P50": 0,
			"P75": 0,
			"P95": 0,
			"P99": 0,
			"Max": 0,
			"Mean": 82,
			"StdDev": 82
		}
	},
	{
		"ReportKey": {
			"IfGlobalSetUpTearDown": "GlobalSetUp",
			"Priority": "3",
			"TestResult": "ALL"
		},
		"Count": 2,
		"PerformanceGauge": {
			"Min": 0,
			"P50": 0,
			"P75": 0,
			"P95": 0,
			"P99": 0,
			"Max": 0,
			"Mean": 0,
			"StdDev": 0
		}
	},
	{
		"ReportKey": {
			"IfGlobalSetUpTearDown": "GlobalSetUp",
			"Priority": "ALL",
			"TestResult": "ALL"
		},
		"Count": 4,
		"PerformanceGauge": {
			"Min": 164,
			"P50": 0,
			"P75": 0,
			"P95": 0,
			"P99": 0,
			"Max": 0,
			"Mean": 41,
			"StdDev": 71
		}
	},
	{
		"ReportKey": {
			"IfGlobalSetUpTearDown": "ALL",
			"Priority": "ALL",
			"TestResult": "ALL"
		},
		"Count": 4,
		"PerformanceGauge": {
			"Min": 0,
			"P50": 0,
			"P75": 164,
			"P95": 164,
			"P99": 164,
			"Max": 164,
			"Mean": 41,
			"StdDev": 71
		}
	}
];

var stats2 = [
	{
		"ReportKey": 1562659343000,
		"Count": 4
	}
];

var stats2_success = [
	{
		"ReportKey": 1562659343000,
		"Count": 0
	}
];

var stats2_fail = [
	{
		"ReportKey": 1562659343000,
		"Count": 1
	}
];

var stats3_status = [
	{
		"ReportKey": {
			"TestResult": "Fail"
		},
		"Count": 1
	},
	{
		"ReportKey": {
			"TestResult": "ParentFailed"
		},
		"Count": 1
	},
	{
		"ReportKey": {
			"TestResult": "ParentSkipped"
		},
		"Count": 2
	}
];
