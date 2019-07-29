/*
 * go4api - an api testing tool written in Go
 * Created by: Ping Zhu 2018
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 *
 */

package ui

var Index = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
    <link href="style/go4api.css" rel="stylesheet" type="text/css"/>
    <script type="text/javascript" src="js/go4api.js"></script>
    <script type="text/javascript" src="js/results.js"></script>
    <script type="text/javascript" src="js/stats.js"></script>
    <script type="text/javascript" src="js/Chart.bundle.min.js"></script>
  <title>Go4Api Reports</title>
</head>
<body>  
  <div class="container">
      <div class="head">
          <a href="https://github.com/zpsean/go4api" target="blank_" title="Go4Api Home Page"><img alt="Go4Api" src="style/logo.png"/></a>
      </div>
      <div class="main">
          <div class="skeleton">
              <div class="content">
                  <div class="sous-menu">
                      <div class="item selected"><a href="index.html">Overview</a></div>
                      <div class="item "><a id="details_link" href="details.html">Details</a></div>
                      <div class="item "><a id="graphic_link" href="graphic.html">Graphic</a></div>
                      <div class="item "><a id="mindex_link" href="mindex.html">MutationOverview</a></div>
                      <div class="item "><a id="mutation_link" href="mutation.html">Mutation</a></div>
                      <div class="item "><a id="fuzz_link" href="fuzz.html">Fuzz</a></div>

                      <script type="text/javascript">
                        // var timestamp = 1523957748602;
                        // var runStartHumanDate = new Date(timestamp).format("YYYY-MM-DD HH:mm:ss Z");
                        var runStartHumanDate = gStart.substr(0, 19)
                        var runDuration = (gEndUnixNano - gStartUnixNano) / 1000000000
                        document.writeln("<p class='sim_desc'>");
                        document.writeln("<b>" + "Started at " + runStartHumanDate + ", duration: " + runDuration + " seconds </b>");
                        document.writeln("</p>");
                      </script>
                  </div>

                  <div class="content-in">
                    <h1><span>> </span>Overview Information</h1>
                    <div class="article">

                      <div class="statistics extensible-geant collapsed">
                          <div class="schema p_right">
                            <div id="priority_2_percentage">
                              <canvas id="myChart" width="100" height="100"></canvas>
                            </div>
                          </div>

                          <div class="schema p_left">
                            <div id="priority_2_line" class="p_left">
                              <canvas id="myChart2" width="100" height="30"></canvas>
                            </div>
                          </div>

                          <div class="title">
                              <div id="statistics_title" class="title_collapsed">Statistics</div>
                          </div>

                          <table id="container_statistics_head" class="statistics-in extensible-geant">
                              <thead>
                                  <tr>
                                      <th rowspan="2" id="col-1" width="139px" class="header sortable sorted-up"><span>Phase</span></th>
                                      <th colspan="3" width="95px" class="header"><span class="executions">Executions</span></th>
                                      <th colspan="8" width="95px" class="header"><span class="response-time">Response Time (ns)</span></th>
                                  </tr>
                                  <tr>
                                      <th id="col-2" width="95px" class="header sortable"><span>Priority</span></th>
                                      <th id="col-3" width="95px" class="header sortable"><span>Status</span></th>
                                      <th id="col-4" width="95px" class="header sortable"><span>Count</span></th>

                                      <th id="col-7" width="95px" class="header sortable"><span>Min</span></th>
                                      <th id="col-8" width="95px" class="header sortable"><span>50th pct</span></th>
                                      <th id="col-9" width="95px" class="header sortable"><span>75th pct</span></th>
                                      <th id="col-10" width="95px" class="header sortable"><span>95th pct</span></th>
                                      <th id="col-11" width="95px" class="header sortable"><span>99th pct</span></th>
                                      <th id="col-12" width="95px" class="header sortable"><span>Max</span></th>
                                      <th id="col-13" width="95px" class="header sortable"><span>Mean</span></th>
                                      <th id="col-14" width="95px" class="header sortable"><span>Std Dev</span></th>
                                  </tr>
                              </thead>
                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body" class="statistics-in extensible-geant">
                                  <tbody></tbody>
                              </table>

                              <script type="text/javascript">
                                for (var i = 0; i < stats1.length; i++)
                                {
                                  var newTr = container_statistics_body.insertRow();
                                  
                                  var newTd0 = newTr.insertCell();
                                  var newTd1 = newTr.insertCell();
                                  var newTd2 = newTr.insertCell();
                                  var newTd3 = newTr.insertCell();
                                  var newTd4 = newTr.insertCell();
                                  var newTd5 = newTr.insertCell();
                                  var newTd6 = newTr.insertCell();
                                  var newTd7 = newTr.insertCell();
                                  var newTd8 = newTr.insertCell();
                                  var newTd9 = newTr.insertCell();
                                  var newTd10 = newTr.insertCell();
                                  var newTd11 = newTr.insertCell();

                                  newTd0.width="139px";
                                  newTd1.width="95px";
                                  newTd2.width="95px";
                                  newTd3.width="95px";
                                  newTd4.width="95px";
                                  newTd5.width="95px";
                                  newTd6.width="95px";
                                  newTd7.width="95px";
                                  newTd8.width="95px";
                                  newTd9.width="95px";
                                  newTd10.width="95px";
                                  newTd11.width="95px";

                           
                                  newTd0.innerText = stats1[i].ReportKey.IfGlobalSetUpTearDown;
                                  newTd1.innerText = stats1[i].ReportKey.Priority;
                                  newTd2.innerText = stats1[i].ReportKey.TestResult;
                                  newTd3.innerText = stats1[i].Count;
                                  newTd4.innerText = stats1[i].PerformanceGauge.Min;
                                  newTd5.innerText = stats1[i].PerformanceGauge.P50;
                                  newTd6.innerText = stats1[i].PerformanceGauge.P75;
                                  newTd7.innerText = stats1[i].PerformanceGauge.P95;
                                  newTd8.innerText = stats1[i].PerformanceGauge.P99;
                                  newTd9.innerText = stats1[i].PerformanceGauge.Max;
                                  newTd10.innerText = stats1[i].PerformanceGauge.Mean;
                                  newTd11.innerText = stats1[i].PerformanceGauge.StdDev;
                                }
                              </script>
                          </div>
                      </div>


                      <div class="schema p_right">
                        <div id="priority_2_percentage">
                          <canvas id="myChart3" width="100" height="100"></canvas>
                        </div>
                      </div>

                      <div class="schema p_left">
                        <div id="priority_2_line" class="p_left">
                          <canvas id="myChart4" width="100" height="30"></canvas>
                        </div>
                      </div>


                      <div class="schema geant">
                        <canvas id="myChart5" width="1089" height="350"></canvas>
                        <div id="container" class="geant"></div>
                          
                      </div>

                    </div>
                  </div>
              </div>
          </div>
      </div>
  </div>
  <div class="foot">
      <a href="https://github.com/zpsean/go4api" title="Go4Api Home Page"><img alt="Go4Api" src="style/logosmall.png"/></a>
  </div>


  <script>
    var tcCountArray = new Array(3)

    for (var i = 0; i < stats3_status.length; i++)
    {
      if (stats3_status[i].ReportKey.TestResult == "Fail")
      {
        tcCountArray[0] = stats3_status[i].Count
      } else if (stats3_status[i].ReportKey.TestResult == "Success") {
        tcCountArray[1] = stats3_status[i].Count
      } else if (stats3_status[i].ReportKey.TestResult == "ParentFailed") {
        tcCountArray[2] = stats3_status[i].Count
      }
    }
    
    var data = {
            labels: [
                "Fail",
                "Success",
                "ParentFailed"
            ],
            datasets: [
                {
                    data: tcCountArray,
                    backgroundColor: [
                        "#FF6384",
                        "#36A2EB",
                        "#FFCE56"
                    ],
                    hoverBackgroundColor: [
                        "#FF6384",
                        "#36A2EB",
                        "#FFCE56"
                    ]
                }]
        };
    // Get the context of the canvas element we want to select
    var ctx = document.getElementById("myChart").getContext("2d");
    var myBarChart = new Chart(ctx, {
                                        type: 'pie',
                                        data: data,
                                        options: {
                                            title: {
                                              display: true,
                                              text: 'Overall Executions'
                                            }
                                        }
                                });
  </script>


  <script>
    var resultLabel = [];
    for (var i in stats2) {
      startDate = new Date(stats2[i].ReportKey)
      str = startDate.getHours() + ":" + startDate.getMinutes() + ":" + startDate.getSeconds()
      resultLabel.push(str)
    }

    var resultDataSuccess = [];
    var resultDataFail = [];
    for (var i in stats2) {
      resultDataSuccess.push(stats2_success[i].Count)
      resultDataFail.push(stats2_fail[i].Count)
    }

    console.log(resultDataSuccess)
    console.log(resultDataFail)

    var ctx = document.getElementById("myChart2").getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: resultLabel,
            datasets: [{
                label: '# of TestCase Started - Success',
                data: resultDataSuccess,
                backgroundColor: "#36A2EB",
                borderColor: [
                    'rgba(255,99,132,1)'
                ],
                borderWidth: 1
            },
            {
                label: '# of TestCase Started - Fail',
                data: resultDataFail,
                backgroundColor: "#FF6384",
                borderColor: [
                    'rgba(255,99,132,1)'
                ],
                borderWidth: 1
            }]
        },
        options: {
            title: {
              display: true,
              text: 'Overall Executions'
            },
            scales: {
              xAxes: [{
                stacked: true,
              }],
              yAxes: [{
                stacked: true
              }]
            }
        }
    });
  </script>


  <script>
    var tcCountArray = new Array(3)
    
    for (var i = 0; i < stats3_status.length; i++)
    {
      if (stats3_status[i].ReportKey.TestResult == "Fail")
      {
        tcCountArray[0] = stats3_status[i].Count
      } else if (stats3_status[i].ReportKey.TestResult == "Success") {
        tcCountArray[1] = stats3_status[i].Count
      } else if (stats3_status[i].ReportKey.TestResult == "ParentFailed") {
        tcCountArray[2] = stats3_status[i].Count
      }
    }

    var data = {
            labels: [
                "Fail",
                "Success",
                "ParentFailed"
            ],
            datasets: [
                {
                    data: tcCountArray,
                    backgroundColor: [
                        "#FF6384",
                        "#36A2EB",
                        "#FFCE56"
                    ],
                    hoverBackgroundColor: [
                        "#FF6384",
                        "#36A2EB",
                        "#FFCE56"
                    ]
                }]
        };
    // Get the context of the canvas element we want to select
    var ctx = document.getElementById("myChart3").getContext("2d");
    var myBarChart = new Chart(ctx, {
                                        type: 'pie',
                                        data: data,
                                        options: {
                                            title: {
                                              display: true,
                                              text: 'Overall Executions'
                                            }
                                        }
                                });
  </script>


  <script>
    var resultLabel = [];
    for (var i in stats2) {
      startDate = new Date(stats2[i].ReportKey)
      str = startDate.getHours() + ":" + startDate.getMinutes() + ":" + startDate.getSeconds()
      resultLabel.push(str)
    }

    console.log(resultLabel)

    var resultData = [];
    for (var i in stats2) {
      resultData.push(stats2[i].Count)
    }

    console.log(resultData)

    var ctx = document.getElementById("myChart4").getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: resultLabel,
            datasets: [{
                label: '# of TestCase Started',
                data: resultData,
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero:true
                    }
                }]
            }
        }
    });
  </script>



  <script>
    var DATA_COUNT = 12;
    var MIN_XY = -150;
    var MAX_XY = 100;

    function generateData(stat) {
      var data = [
{"x":1537868438468841870, "y":1537868439569804655, "v":1100962785},
{"x":1537868440071394552, "y":1537868440494010301, "v":422615749},
{"x":1537868440997248915, "y":1537868441411196761, "v":413947846},
{"x":1537868441501379732, "y":1537868442058686728, "v":557306996},
{"x":1537868441501399696, "y":1537868441882226595, "v":380826899},
{"x":1537868441501403272, "y":1537868442042055322, "v":540652050},
{"x":1537868441501452035, "y":1537868442035993969, "v":534541934},
{"x":1537868442564148953, "y":1537868443053607182, "v":489458229},
{"x":1537868442564160947, "y":1537868443051335659, "v":487174712},
{"x":1537868443053782066, "y":1537868443053782066, "v":0},
{"x":1537868443554011422, "y":1537868444070319763, "v":516308341},
{"x":1537868443554017272, "y":1537868444070373388, "v":516356116}];
      // var i;

      // for (i = 0; i < DATA_COUNT; ++i) {
      //   data.push({
      //     x: i + 1,
      //     y: i,
      //     v: stat[i].Count
      //   });
      // }
      return data;
    }


    var resultLabel = [];
    for (var i in stats2) {
      startDate = new Date(stats2[i].ReportKey)
      str = startDate.getHours() + ":" + startDate.getMinutes() + ":" + startDate.getSeconds()
      resultLabel.push(str)
    }

    var resultDataSuccess = [];
    var resultDataFail = [];
    for (var i in stats2) {
      resultDataSuccess.push(stats2_success[i].Count)
      resultDataFail.push(stats2_fail[i].Count)
    }

    var data = {
      datasets: [{
        data: generateData(stats2_success)
      }]
    };

    console.log(data)

    var options = {
      aspectRatio: 1,
      legend: false,
      tooltips: false,

      elements: {
        point: {
          borderWidth: function(context) {
            return Math.min(Math.max(1, context.datasetIndex + 1), 8);
          },

          hoverBackgroundColor: 'transparent',

          hoverBorderWidth: function(context) {
            var value = context.dataset.data[context.dataIndex];
            return Math.round(8 * value.v / 1000);
          },

          radius: function(context) {
            var value = context.dataset.data[context.dataIndex];
            var size = context.chart.width;
            var base = Math.abs(value.v) / 1000;
            return (size / 24) * base;
          }
        }
      }
    };

    var ctx = document.getElementById("myChart5").getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'bubble',
        data: data,
        options: {
            title: {
              display: true,
              text: 'Overall Executions - Bubble'
            }
        }
    });
  </script>

</body>
</html>
`
