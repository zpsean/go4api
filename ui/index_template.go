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

var Index = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <link href="style/go4api.css" rel="stylesheet" type="text/css"/>
    <script type="text/javascript" src="js/go4api.js"></script>
    <script type="text/javascript" src="js/reslts.js"></script>
    <title>Go4Api Stats</title>
</head>

<body>
    <h1>Go4Api Executions</h1>

    <div class="container details">

    <div class="head">
        <a href="https://github.com/zpsean/go4api" target="blank_" title="Go4Api Home Page"><img alt="Go4Api" src=""/></a>
    </div>

        <!-- tab header -->
        <table class="tabTitlesContainer">
            <tr id="tabTitles">
                <td class="tabTitleSelected" onclick="tabPageControl(0)">Executions List</td>
                <td class="tabTitleUnSelected" onclick="tabPageControl(1)">Executions Graphic</td>
                <td class="tabTitleUnSelected" onclick="tabPageControl(2)">Other</td>
            </tr>
        </table>

         <!-- tab content -->
        <table id="tabPagesContainer">
            <tbody class="tabPageSelected">
                <tr class="tabPage">
                    <td>
                        <table border="1" id="testTble">
                            <col width="20" />
                            <col width="20" />
                            <col width="300" />
                            <col width="20" />
                            <col width="500" />
                            <tr style='text-align: left'>
                                <th>#</th>
                                <th>Priority/th>
                                <th>Case ID</th>
                                <th>Status</th>
                                <th>Case File / Data Table / Data Row</th>
                            </tr>
                            
                            <script type="text/javascript">
                                for (var i = 0;i < tcResults.length; i++)
                                {
                                    var newTr = testTble.insertRow();
                                    
                                    var newTd0 = newTr.insertCell();
                                    var newTd1 = newTr.insertCell();
                                    var newTd2 = newTr.insertCell();
                                    var newTd3 = newTr.insertCell();
                                    var newTd4 = newTr.insertCell();
                             
                                    newTd0.innerText= i + 1;
                                    newTd1.innerText= tcResults[i][0];
                                    newTd2.innerText= tcResults[i][1];
                                    newTd3.innerText= tcResults[i][2];
                                    newTd4.innerText= tcResults[i][9];
                                }


                            </script>

                        </table>
                    </td>
                </tr>
            </tbody>

            <tbody class="tabPageUnSelected">
                <tr class="tabPage">
                    <td> 
                        
                        <svg width="1000" height="5000" id="mainSVG" version="1.1" xmlns="http://www.w3.org/2000/svg">


                            <line x1="0" y1="0" x2="0" y2="1000" style="stroke:rgb(99,99,99);stroke-width:2"/>
                            <line x1="0" y1="0" x2="1000" y2="0" style="stroke:rgb(99,99,99);stroke-width:2"/>
                            <line x1="0" y1="500" x2="600" y2="500" style="stroke:rgb(99,99,99);stroke-width:1"/>

                        </svg>

                        <script type="text/javascript">
                                
                                var svgRoot = document.getElementById('mainSVG');
                                var priority = "";

                                var tcPositions = {}

                                for (var i = 0; i < tcResults.length; i++)
                                {
                                    var c = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
                                    c.setAttribute('cx', (i % 10 + 1) * 50);
                                    c.setAttribute('cy', (tcResults[i][11] - pStartUnixNano + tcResults[i][13] / 2) / 10000000);
                                    c.r.baseVal.value = tcResults[i][13] / 100000000 + 5;

                                    var pos = [(i % 10 + 1) * 50, (tcResults[i][11] - pStartUnixNano + tcResults[i][13] / 2) / 10000000];
                                    tcPositions[tcResults[i][1]] = pos;

                                    if (tcResults[i][3] == "Success")
                                        {
                                        c.setAttribute('fill', 'green');
                                        }
                                    else if (tcResults[i][3] == "Fail")
                                        {
                                        c.setAttribute('fill', 'red');
                                        }
                                    else
                                        {
                                        c.setAttribute('fill', 'gray');
                                        }

                                    svgRoot.appendChild(c);


                                    if (tcResults[i][0] != priority)
                                    {
                                        var line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
                                        line.setAttribute('x1', 0);
                                        line.setAttribute('y1', (tcResults[i][11] - pStartUnixNano + tcResults[i][13] / 2) / 10000000);
                                        line.setAttribute('x2', 600);
                                        line.setAttribute('y2', (tcResults[i][11] - pStartUnixNano + tcResults[i][13] / 2) / 10000000);
                                        line.setAttribute('style', "stroke:rgb(99,99,99);stroke-width:1");
                                        
                                        svgRoot.appendChild(line);

                                        priority = tcResults[i][0];
                                    }
                                }

                                for (var parent = 0; parent < tcResults.length; parent++)
                                {
                                    for (var child = 0; child < tcResults.length; child++)
                                    {
                                        if (tcResults[parent][1] == tcResults[child][2])
                                        {
                                            var line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
                                            
                                            line.setAttribute('x1', tcPositions[tcResults[parent][1]][0]);
                                            line.setAttribute('y1', tcPositions[tcResults[parent][1]][1]);
                                            line.setAttribute('x2', tcPositions[tcResults[child][1]][0]);
                                            line.setAttribute('y2', tcPositions[tcResults[child][1]][1]);
                                            line.setAttribute('style', "stroke:rgb(250,99,99);stroke-width:1");
                                            
                                            svgRoot.appendChild(line);
                                        }
                                    }
                                }


                        </script>

                    </td>
                </tr>
            </tbody>

            <tbody class="tabPageUnSelected">
                <tr class="tabPage">
                    <td> plase holder 2</td>
                </tr>
            </tbody>
        </table>


    </div>

<div class="foot">
    <a href="https://github.com/zpsean/go4api" title="Go4Api Home Page"><img alt="Go4Api" src=""/></a>
</div>

</body>
</html>`