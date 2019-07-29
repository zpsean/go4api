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

var MIndex = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
    <link href="style/go4api.css" rel="stylesheet" type="text/css"/>
    <script type="text/javascript" src="js/go4api.js"></script>
    <script type="text/javascript" src="js/results.js"></script>
    <script type="text/javascript" src="js/stats.js"></script>
    <script type="text/javascript" src="js/mutationstats.js"></script>
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
                      <div class="item "><a href="index.html">Overview</a></div>
                      <div class="item "><a id="details_link" href="details.html">Details</a></div>
                      <div class="item "><a id="graphic_link" href="graphic.html">Graphic</a></div>
                      <div class="item selected"><a id="mindex_link" href="mindex.html">MutationOverview</a></div>
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

                          <table id="container_statistics_head" class="statistics-in extensible-geant">
                              <thead>
                                  <tr>
                                      <th id="col-1" width="50px" class="header sortable"><span>#</span></th>
                                      <th id="col-2" width="300px" class="header sortable"><span>HttpUrl</span></th>
                                      <th id="col-2" width="300px" class="header sortable"><span>HttpMethod</span></th>
                                      <th id="col-4" width="300px" class="header sortable"><span>HttpStatus</span></th>
                                      <th id="col-5" width="140px" class="header sortable"><span>Count</span></th>
                                  </tr>
                              </thead>
                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body" class="statistics-in extensible-geant">
                                  <tbody></tbody>
                              </table>


                              <script type="text/javascript">
                                for (var i = 0;i < mutationStats1.length; i++)
                                  {
                                    var newTr = container_statistics_body.insertRow();
                                    
                                    var newTd0 = newTr.insertCell();
                                    var newTd1 = newTr.insertCell();
                                    var newTd2 = newTr.insertCell();
                                    var newTd3 = newTr.insertCell();
                                    var newTd4 = newTr.insertCell();

                                    newTd0.width="50px";
                                    newTd1.width="300px";
                                    newTd2.width="300px";
                                    newTd3.width="300px";
                                    newTd4.width="140px";

                                    newTd0.innerText = i;
                                    newTd1.innerText = mutationStats1[i].ReportKey.Path;
                                    newTd2.innerText = mutationStats1[i].ReportKey.Method;
                                    newTd3.innerText = mutationStats1[i].ReportKey.ActualStatusCode;
                                    newTd4.innerText = mutationStats1[i].Count;
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

                      <div class="statistics extensible-geant collapsed">
                          <table id="container_statistics_head" class="statistics-in extensible-geant">
                              <thead>
                                  <tr>
                                      <th id="col-1" width="50px"  class="header sortable"><span>#</span></th>
                                      <th id="col-2" width="300px" class="header sortable"><span>HttpUrl</span></th>
                                      <th id="col-2" width="300px" class="header sortable"><span>HttpMethod</span></th>
                                      <th id="col-3" width="300px" class="header sortable"><span>MutationPart</span></th>
                                      <th id="col-4" width="70px" class="header sortable"><span>HttpStatus</span></th>
                                      <th id="col-5" width="70px" class="header sortable"><span>Count</span></th>
                                  </tr>
                              </thead>
                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body2" class="statistics-in extensible-geant">
                                  <tbody></tbody>
                              </table>


                              <script type="text/javascript">
                                for (var i = 0;i < mutationStats2.length; i++)
                                  {
                                    var newTr = container_statistics_body2.insertRow();
                                    
                                    var newTd0 = newTr.insertCell();
                                    var newTd1 = newTr.insertCell();
                                    var newTd2 = newTr.insertCell();
                                    var newTd3 = newTr.insertCell();
                                    var newTd4 = newTr.insertCell();
                                    var newTd5 = newTr.insertCell();

                                    newTd0.width="50px";
                                    newTd1.width="300px";
                                    newTd2.width="300px";
                                    newTd3.width="300px";
                                    newTd4.width="70px";
                                    newTd5.width="70px";

                                    newTd0.innerText = i;
                                    newTd1.innerText = mutationStats2[i].ReportKey.Path;
                                    newTd2.innerText = mutationStats2[i].ReportKey.Method;
                                    newTd3.innerText = mutationStats2[i].ReportKey.MutationArea;
                                    newTd4.innerText = mutationStats2[i].ReportKey.ActualStatusCode;
                                    newTd5.innerText = mutationStats2[i].Count;
                                  }

                              </script>
                          </div>
                      </div>
                      <div class="schema geant">
                        <canvas id="myChart6" width="1089" height="350"></canvas>
                        <div id="container" class="geant"></div>
                          
                      </div>

                      <div class="statistics extensible-geant collapsed">
                          <table id="container_statistics_head" class="statistics-in extensible-geant">
                              <thead>
                                  <tr>
                                      <th id="col-1" width="50px"  class="header sortable"><span>#</span></th>
                                      <th id="col-2" width="300px" class="header sortable"><span>HttpUrl</span></th>
                                      <th id="col-2" width="260px" class="header sortable"><span>HttpMethod</span></th>
                                      <th id="col-3" width="260px" class="header sortable"><span>MutationPart</span></th>
                                      <th id="col-4" width="70px"  class="header sortable"><span>MutationCategory</span></th>
                                      <th id="col-5" width="70px"  class="header sortable"><span>HttpStatus</span></th>
                                      <th id="col-6" width="70px" class="header sortable"><span>Count</span></th>
                                  </tr>
                              </thead>
                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body3" class="statistics-in extensible-geant">
                                  <tbody></tbody>
                              </table>


                              <script type="text/javascript">
                                for (var i = 0;i < mutationStats3.length; i++)
                                  {
                                    var newTr = container_statistics_body3.insertRow();
                                    
                                    var newTd0 = newTr.insertCell();
                                    var newTd1 = newTr.insertCell();
                                    var newTd2 = newTr.insertCell();
                                    var newTd3 = newTr.insertCell();
                                    var newTd4 = newTr.insertCell();
                                    var newTd5 = newTr.insertCell();
                                    var newTd6 = newTr.insertCell();

                                    newTd0.width="50px";
                                    newTd1.width="300px";
                                    newTd2.width="300px";
                                    newTd3.width="300px";
                                    newTd4.width="70px";
                                    newTd5.width="70px";
                                    newTd6.width="70px";

                                    newTd0.innerText = i;
                                    newTd1.innerText = mutationStats3[i].ReportKey.Path;
                                    newTd2.innerText = mutationStats3[i].ReportKey.Method;
                                    newTd3.innerText = mutationStats3[i].ReportKey.MutationArea;
                                    newTd4.innerText = mutationStats3[i].ReportKey.MutationCategory
                                    newTd5.innerText = mutationStats3[i].ReportKey.ActualStatusCode;
                                    newTd6.innerText = mutationStats3[i].Count;
                                  }
                              </script>

                          </div>
                      </div>

                    </div>

                  </div>
              </div>
          </div>
      </div>
  </div>



  <script>
    var tcCountArray = new Array(3)
    tcCountArray[0] = stats1.Overall.Fail
    tcCountArray[1] = stats1.Overall.Success
    tcCountArray[2] = stats1.Overall.ParentFailed

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
    var resultLabel = [];
    for (var i in mutationStats1) {
      resultLabel.push(mutationStats1[i].ReportKey.Method + ":" + mutationStats1[i].ReportKey.ActualStatusCode)
    }

    var resultData = [];
    for (var i in mutationStats1) {
      resultData.push(mutationStats1[i].Count)
    }

    var data = {
            labels: resultLabel,
            datasets: [
                {
                    data: resultData
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
                                              text: 'Overall Executions - Method, StatusCode'
                                            }
                                        }
                                });
  </script>

  <script>
    var resultLabel = [];
    for (var i in mutationStats2) {
      resultLabel.push(mutationStats2[i].ReportKey.Method + ":" + mutationStats2[i].ReportKey.MutationArea + ":" + mutationStats2[i].ReportKey.ActualStatusCode)
    }

    var resultData = [];
    for (var i in mutationStats2) {
      resultData.push(mutationStats2[i].Count)
    }

    var data = {
            labels: resultLabel,
            datasets: [
                {
                    data: resultData
                }]
        };
    // Get the context of the canvas element we want to select
    var ctx = document.getElementById("myChart4").getContext("2d");
    var myBarChart = new Chart(ctx, {
                                        type: 'pie',
                                        data: data,
                                        options: {
                                            title: {
                                              display: true,
                                              text: 'Overall Executions - Method, StatusCode, MutationPart'
                                            }
                                        }
                                });
  </script>


  <script>
    var resultLabel = [];
    for (var i in mutationStats3) {
      resultLabel.push(mutationStats3[i].ReportKey.Method + ":" + mutationStats3[i].ReportKey.MutationArea + ":" + mutationStats3[i].ReportKey.MutationCategory + ":" + mutationStats3[i].ReportKey.ActualStatusCode)
    }

    var resultData = [];
    for (var i in mutationStats3) {
      resultData.push(mutationStats3[i].Count)
    }

    var data = {
            labels: resultLabel,
            datasets: [
                {
                    data: resultData
                }]
        };
    // Get the context of the canvas element we want to select
    var ctx = document.getElementById("myChart6").getContext("2d");
    var myBarChart = new Chart(ctx, {
                                        type: 'pie',
                                        data: data,
                                        options: {
                                            title: {
                                              display: true,
                                              text: 'Overall Executions - Method, StatusCode, MutationPart, MutationCategory'
                                            }
                                        }
                                });
  </script>
  

  <div class="foot">
      <a href="https://github.com/zpsean/go4api" title="Go4Api Home Page"><img alt="Go4Api" src="style/logosmall.png"/></a>
  </div>
</body>
</html>
`
