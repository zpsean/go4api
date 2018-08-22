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
                            <col width="200" />
                            <col width="200" />
                            <col width="50" />
                            <col width="300" />
                            <col width="300" />
                            <tr style='text-align: left'>
                                <th>#</th>
                                <th>Priority</th>
                                <th>Case ID</th>
                                <th>ParentTestCase</th>
                                <th>Status</th>
                                <th>Case File</th>
                                <th>Message</th>
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
                                    var newTd5 = newTr.insertCell();
                                    var newTd6 = newTr.insertCell();
                             
                                    newTd0.innerText= i + 1;
                                    newTd1.innerText = tcResults[i].Priority;
                                    newTd2.innerText = tcResults[i].TcName;
                                    newTd3.innerText = tcResults[i].ParentTestCase;
                                    newTd4.innerText = tcResults[i].TestResult;
                                    newTd5.innerText = tcResults[i].JsonFilePath + " / " + tcResults[i].CsvFile  + " / " + tcResults[i].CsvRow;
                                    newTd6.innerText = tcResults[i].TestMessages;
                                }


                            </script>

                        </table>
                    </td>
                </tr>
            </tbody>

            <tbody class="tabPageUnSelected">
                <tr class="tabPage">
                    <td> 
                        
                        <svg width="1200" height="4000" id="mainSVG" version="1.1" xmlns="http://www.w3.org/2000/svg">


                            <line x1="0" y1="0" x2="0" y2="1500" style="stroke:rgb(99,99,99);stroke-width:2"/>
                            <line x1="0" y1="0" x2="1200" y2="0" style="stroke:rgb(99,99,99);stroke-width:2"/>

                        </svg>

                        <script type="text/javascript">
                                
                                var svgRoot = document.getElementById('mainSVG');
                                var priority = "";

                                var tcPositions = {}


                                for (var i = 0; i < tcResults.length; i++)
                                {
                                    var c = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
                                    c.setAttribute('cx', (i % 15 + 1) * 50);
                                    c.setAttribute('cy', (tcResults[i].StartTimeUnixNano - pStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000);
                                    c.r.baseVal.value = tcResults[i].DurationUnixNano / 100000000 + 5;


                                    var text = document.createElementNS('http://www.w3.org/2000/svg', 'text');

                                    strJons = JSON.stringify(tcResults[i], null, 4)
                                    text.textContent = strJons

                                    text.setAttribute('x', (i % 15 + 1) * 50);
                                    text.setAttribute('y', (tcResults[i].StartTimeUnixNano - pStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000 - 20)
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




                                    var pos = [(i % 15 + 1) * 50, (tcResults[i].StartTimeUnixNano - pStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000];
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
                                        line.setAttribute('y1', (tcResults[i].StartTimeUnixNano - pStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000);
                                        line.setAttribute('x2', 1200);
                                        line.setAttribute('y2', (tcResults[i].StartTimeUnixNano - pStartUnixNano + tcResults[i].DurationUnixNano / 2) / 10000000);
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