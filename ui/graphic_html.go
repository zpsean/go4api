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

var Graphic = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
    <link href="style/go4api.css" rel="stylesheet" type="text/css"/>
    <script type="text/javascript" src="js/go4api.js"></script>
    <script type="text/javascript" src="js/results.js"></script>
    <script type="text/javascript" src="js/stats.js"></script>
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
                      <div class="item selected"><a id="graphic_link" href="graphic.html">Graphic</a></div>
                      <div class="item "><a id="details_link" href="details.html">Details</a></div>
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


                      <svg width="1150" height="2000" id="mainSVG" version="1.1" xmlns="http://www.w3.org/2000/svg">
                          <line x1="0" y1="0" x2="0" y2="1150" style="stroke:rgb(99,99,99);stroke-width:2"/>
                          <line x1="0" y1="0" x2="1080" y2="0" style="stroke:rgb(99,99,99);stroke-width:2"/>
                          <line x1="1080" y1="0" x2="1080" y2="1150" style="stroke:rgb(99,99,99);stroke-width:1"/>
                      </svg>

                      <script type="text/javascript">
                              
                              var svgRoot = document.getElementById('mainSVG');
                              var priority = "";

                              var tcPositions = {}


                              for (var i = 0; i < tcResults.length; i++)
                              {
                                  var c = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
                                  c.setAttribute('cx', (i % 15 + 1) * 50);
                                  c.setAttribute('cy', (tcResults[i].StartTimeUnixNano - gStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000);
                                  c.r.baseVal.value = tcResults[i].DurationUnixNano / 100000000 + 5;


                                  var text = document.createElementNS('http://www.w3.org/2000/svg', 'text');

                                  strJons = JSON.stringify(tcResults[i], null, 4)
                                  text.textContent = strJons

                                  text.setAttribute('x', (i % 15 + 1) * 50);
                                  text.setAttribute('y', (tcResults[i].StartTimeUnixNano - gStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000 - 20)
                                  text.style.width = '400px'
                                  text.setAttribute('fill', 'transparent')
                                  svgRoot.appendChild(text)


                                  ;(function(text) {
                                      c.addEventListener('mouseenter', () => {
                                          text.setAttribute('fill', 'red')
                                      })
                                      c.addEventListener('mouseleave', () => {
                                          text.setAttribute('fill', 'transparent')
                                      })
                                  })(text)




                                  var pos = [(i % 15 + 1) * 50, (tcResults[i].StartTimeUnixNano - gStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000];
                                  tcPositions[tcResults[i].TcName] = pos;

                                  if (tcResults[i].TestResult == "Success")
                                      {
                                      c.setAttribute('fill', 'green');
                                      }
                                  else if (tcResults[i].TestResult == "Fail")
                                      {
                                      c.setAttribute('fill', 'red');
                                      }
                                  else
                                      {
                                      c.setAttribute('fill', 'gray');
                                      }

                                  svgRoot.appendChild(c);


                                  if (tcResults[i].Priority != priority)
                                  {
                                      var line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
                                      line.setAttribute('x1', 0);
                                      line.setAttribute('y1', (tcResults[i].StartTimeUnixNano - gStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000);
                                      line.setAttribute('x2', 1080);
                                      line.setAttribute('y2', (tcResults[i].StartTimeUnixNano - gStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000);
                                      line.setAttribute('style', "stroke:rgb(99,99,99);stroke-width:1");
                                      
                                      svgRoot.appendChild(line);

                                      priority = tcResults[i].Priority;
                                  }
                              }

                              for (var parent = 0; parent < tcResults.length; parent++)
                              {
                                  for (var child = 0; child < tcResults.length; child++)
                                  {
                                      if (tcResults[parent].TcName == tcResults[child].ParentTestCase)
                                      {
                                          var line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
                                          
                                          line.setAttribute('x1', tcPositions[tcResults[parent].TcName][0]);
                                          line.setAttribute('y1', tcPositions[tcResults[parent].TcName][1]);
                                          line.setAttribute('x2', tcPositions[tcResults[child].TcName][0]);
                                          line.setAttribute('y2', tcPositions[tcResults[child].TcName][1]);
                                          line.setAttribute('style', "stroke:rgb(250,99,99);stroke-width:1");
                                          
                                          svgRoot.appendChild(line);
                                      }
                                  }
                              }


                      </script>


                    </div>
                  </div>
              </div>
          </div>
      </div>
  </div>
  <div class="foot">
      <a href="https://github.com/zpsean/go4api" title="Go4Api Home Page"><img alt="Go4Api" src="style/logosmall.png"/></a>
  </div>
</body>
</html>
`
