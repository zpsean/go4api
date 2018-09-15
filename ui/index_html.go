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

var Index = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
    <link href="style/style.css" rel="stylesheet" type="text/css" />
    <link href="style/go4api.css" rel="stylesheet" type="text/css"/>
    <script type="text/javascript" src="js/go4api.js"></script>
    <script type="text/javascript" src="js/reslts.js"></script>
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
              <h1>Go4Api Executions</h1>

              <div class="content">
                  <div class="sous-menu">
                      <div class="item selected"><a href="index.html">Overview</a></div>
                      <div class="item "><a id="graphic_link" href="graphic.html">Graphic</a></div>
                      <div class="item "><a id="details_link" href="details.html">Details</a></div>

                      <script type="text/javascript">
                        var timestamp = 1523957748602;
                        var runStartHumanDate = moment(timestamp).format("YYYY-MM-DD HH:mm:ss Z");
                        document.writeln("<p class='sim_desc' title='"+ runStartHumanDate +", duration : 1 seconds' data-content=''>");
                        document.writeln("<b>" + runStartHumanDate + ", duration : 10 seconds </b>");
                        document.writeln("</p>");
                      </script>
                  </div>

                  <div class="content-in">
                    <h1><span>> </span>Overview Information</h1>
                    <div class="article">

                      <div class="statistics extensible-geant collapsed">
                          <div class="title">
                              <div id="statistics_title" class="title_collapsed">Statistics</div>
                          </div>

                          <table id="container_statistics_head" class="statistics-in extensible-geant">
                              <thead>
                                  <tr>
                                      <th rowspan="2" id="col-1" class="header sortable sorted-up"><span>Priority</span></th>
                                      <th colspan="1" class="header"><span class="executions">Executions</span></th>
                                      <th colspan="8" class="header"><span class="response-time">Response Time (ns)</span></th>
                                  </tr>
                                  <tr>
                                      <th id="col-2" class="header sortable"><span>Status</span></th>

                                      <th id="col-7" class="header sortable"><span>Min</span></th>
                                      <th id="col-8" class="header sortable"><span>50th pct</span></th>
                                      <th id="col-9" class="header sortable"><span>75th pct</span></th>
                                      <th id="col-10" class="header sortable"><span>95th pct</span></th>
                                      <th id="col-11" class="header sortable"><span>99th pct</span></th>
                                      <th id="col-12" class="header sortable"><span>Max</span></th>
                                      <th id="col-13" class="header sortable"><span>Mean</span></th>
                                      <th id="col-14" class="header sortable"><span>Std Dev</span></th>
                                  </tr>
                              </thead>
                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body" class="statistics-in extensible-geant">
                                  <tbody></tbody>
                              </table>
                          </div>
                      </div>


                      <div class="schema p_right">
                        <div id="priority_1_percentage"></div>
                        <svg width="100px" height="100px" viewBox="0 0 100 100">
                          <circle r="25" cx="50" cy="50" fill="none" stroke="#399C2B" stroke-width="50" stroke-dasharray="16 158" />
                          <circle r="25" cx="50" cy="50" fill="none" stroke="#9A4324" stroke-width="50" stroke-dasharray="48 158" stroke-dashoffset="-16"/>
                          <circle r="25" cx="50" cy="50" fill="none" stroke="#9C9CB2" stroke-width="50" stroke-dasharray="79 158" stroke-dashoffset="-64"/>
                        </svg>
                      </div>

                      <div class="schema p_left">
                        <div id="priority_1_line" class="p_left"></div>
                      </div>


                      <div class="schema p_right">
                        <div id="priority_2_percentage"></div>
                      </div>

                      <div class="schema p_left">
                        <div id="priority_2_line" class="p_left"></div>
                      </div>


                      <div class="schema geant">
                        <a name="active_users"></a>
                        <div id="container_active_users" class="geant"></div>
                      </div>

                      <div class="schema geant">
                        <div id="container" class="geant"></div>
                      </div>

                    </div>
                  </div>
              </div>
          </div>
      </div>
  </div>
  <div class="foot">
      <a href="https://github.com/zpsean/go4api" title="Go4Api Home Page"><img alt="Go4Api" src="style/logo.png"/></a>
  </div>
</body>
</html>
`
