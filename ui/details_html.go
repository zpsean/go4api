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

var Details = `<!DOCTYPE html>
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
                      <div class="item selected"><a id="details_link" href="details.html">Details</a></div>
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
                    <h1><span>> </span>Details Information</h1>
                    <div class="article">

                      <div class="statistics extensible-geant collapsed">
                          
                          <table id="container_statistics_head" class="statistics-in extensible-geant">
                              <thead>
                                  <tr style="word-wrap:break-word;word-break:normal";>
                                      <th id="col-1" width="30px"; class="header sortable"><span>#</span></th>
                                      <th id="col-2" width="95px"; class="header sortable"><span>Phase</span></th>
                                      <th id="col-3" width="40px"; class="header sortable"><span>Pri-ority</span></th>
                                      <th id="col-4" width="120px"; class="header sortable"><span>Case ID</span></th>
                                      <th id="col-5" width="120px"; class="header sortable"><span>Parent-TestCase</span></th>
                                      <th id="col-6" width="95px"; class="header sortable"><span>Http-Status</span></th>
                                      <th id="col-7" width="95px"; class="header sortable"><span>Case-Status</span></th>
                                      <th id="col-8" width="50px"; class="header sortable"><span>Duration-(ms)</span></th>
                                      <th id="col-9" width="160px"; class="header sortable"><span>Case Data</span></th>
                                      <th id="col-10" width="314px";  class="header sortable"><span>Message</span></th>
                                  </tr>
                              </thead>
                              <thead>
                                  <tr>
                                      <th id="col-1-1" ></th>

                                      <th id="col-2-1">
                                        <input type="text" id="Phase_input" list="Phase_datalist" style="width:80px" onchange="dataListChange()" />
                                          <datalist id="Phase_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-3-1">
                                        <input type="text" id="Priority_input" list="Priority_datalist" style="width:38px" onchange="dataListChange()" />
                                          <datalist id="Priority_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-4-1">
                                        <input type="text" id="caseid_input" list="caseid_datalist" style="width:110px" onchange="dataListChange()" />
                                          <datalist id="caseid_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-5-1">
                                        <input type="text" id="ParentTestCase_input" list="ParentTestCase_datalist" style="width:110px" onchange="dataListChange()" />
                                          <datalist id="ParentTestCase_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-6-1">
                                        <input type="text" id="HttpResult_input" list="HttpResult_datalist" style="width:85px" onchange="dataListChange()" />
                                          <datalist id="HttpResult_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-7-1">
                                        <input type="text" id="TestResult_input" list="TestResult_datalist" style="width:85px" onchange="dataListChange()" />
                                          <datalist id="TestResult_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-8-1">
                                        <input type="text" id="Duration_input" list="Duration_datalist" style="width:50px" onchange="dataListChange()" />
                                          <datalist id="Duration_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-9-1">
                                        <input type="text" id="CaseData_input" list="CaseData_datalist" style="width:150px" placeholder="Search" onchange="dataListChange()" />
                                          <datalist id="CaseData_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-10-1">
                                        <input type="text" id="Message_input" list="Message_datalist" style="width:210px" placeholder="Search" onchange="dataListChange()" />
                                          <datalist id="Message_datalist">
                                        </datalist>
                                      </th>

                                  </tr>
                              </thead>

                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body" style="word-wrap:break-word;word-break:break-all"; class="statistics-in extensible-geant">
                                  <tbody></tbody>
                              </table>

                              <script type="text/javascript">
                                for (var i = 0; i < tcResults.length; i++)
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

                                    newTd0.width="30px";
                                    newTd1.width="95px";
                                    newTd2.width="40px";
                                    newTd3.width="120px";
                                    newTd4.width="120px";
                                    newTd5.width="95px";
                                    newTd6.width="95px";
                                    newTd7.width="50px";
                                    newTd8.width="160px";
                                    newTd9.width="314px";
                             
                                    newTd0.innerText= i + 1;
                                    newTd1.innerText = tcResults[i].IfGlobalSetUpTearDown;
                                    newTd2.innerText = tcResults[i].Priority;
                                    newTd3.innerText = tcResults[i].TcName;
                                    newTd4.innerText = tcResults[i].ParentTestCase;
                                    newTd5.innerText = tcResults[i].ActualStatusCode;
                                    newTd6.innerText = tcResults[i].TestResult;
                                    newTd7.innerText = tcResults[i].DurationUnixNano / 1000000;
                                    newTd8.innerText = tcResults[i].JsonFilePath + "\n\n" + tcResults[i].CsvFile  + "\n\n" + tcResults[i].CsvRow;
                                    newTd9.innerText = JSON.stringify(tcResults[i].HttpTestMessages, null, 4);
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
  <div class="foot">
      <a href="https://github.com/zpsean/go4api" title="Go4Api Home Page"><img alt="Go4Api" src="style/logosmall.png"/></a>
  </div>


  <script type="text/javascript">
    var list1 = new Array(4);
    var li2 = new Array(4);
    var list2 = new Array(4);

    list1[0] = "IfGlobalSetUpTearDown";
    list1[1] = "Priority";
    list1[2] = "TcName";
    list1[3] = "ParentTestCase";
    list1[4] = "ActualStatusCode";
    list1[5] = "TestResult";
    list1[6] = "DurationUnixNano";
    list1[7] = "CaseData";
    list1[8] = "TestMessages";

    li2[0] = new Array;
    li2[1] = new Array;
    li2[2] = new Array;
    li2[3] = new Array;
    li2[4] = new Array;
    li2[5] = new Array;
    li2[6] = new Array;
    li2[7] = new Array;
    li2[8] = new Array;

    list2[0] = new Array;
    list2[1] = new Array;
    list2[2] = new Array;
    list2[3] = new Array;
    list2[4] = new Array;
    list2[5] = new Array;
    list2[6] = new Array;
    list2[7] = new Array;
    list2[8] = new Array;

    for(var i = 0; i < list1.length; i++)
    { 
      if (i == 6) {
        for (var j = 0; j < tcResults.length; j++)
        { 
          li2[i].push(tcResults[j][list1[i]] / 1000000)
        }
      } else {
        for (var j = 0; j < tcResults.length; j++)
        { 
          li2[i].push(tcResults[j][list1[i]])
        }
      }

      var distinctItems = Array.from(new Set(li2[i]))
      list2[i] = distinctItems
    }
   
    function populateDatalist(id, index) {
      var dataList = document.getElementById(id);
      var list2Element = list2[index];

      for(var i = 0; i < list2Element.length; i++)
      {
          var option = document.createElement("option");
          option.appendChild(document.createTextNode(list2Element[i]));
          option.value = list2Element[i];
          dataList.appendChild(option);
      }

      console.log(dataList)
    }
  </script>


  <script type="text/javascript">
    var tcResults;

    window.onload = function() {
      populateDatalist("Phase_datalist", 0)
      populateDatalist("Priority_datalist", 1)
      populateDatalist("caseid_datalist", 2)
      populateDatalist("ParentTestCase_datalist", 3)
      populateDatalist("HttpResult_datalist", 4)
      populateDatalist("TestResult_datalist", 5)
      populateDatalist("Duration_datalist", 6)
      // populateDatalist("CaseData_datalist", 4)
      // populateDatalist("Message_datalist", 5)
    }

    function dataListChange(){
      var v1 = document.getElementById("Phase_input").value
      var v2 = document.getElementById("Priority_input").value
      var v3 = document.getElementById("caseid_input").value
      var v4 = document.getElementById("ParentTestCase_input").value
      var v5 = document.getElementById("HttpResult_input").value
      var v6 = document.getElementById("TestResult_input").value
      var v7 = document.getElementById("Duration_input").value
      var v8 = document.getElementById("CaseData_input").value
      var v9 = document.getElementById("Message_input").value

      clearRows()
      insertRows(v1, v2, v3, v4, v5, v6, v7, v8, v9)
    }

    function insertRows(v1, v2, v3, v4, v5, v6, v7, v8, v9){
      for (var i = 0; i < tcResults.length; i++)
      {   
          if (searchCriteria(i, v1, v2, v3, v4, v5, v6, v7, v8, v9)) {
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

            newTd0.width="30px";
            newTd1.width="95px";
            newTd2.width="40px";
            newTd3.width="120px";
            newTd4.width="120px";
            newTd5.width="95px";
            newTd6.width="95px";
            newTd7.width="50px";
            newTd8.width="160px";
            newTd9.width="314px";
     
            newTd0.innerText= i + 1;
            newTd1.innerText = tcResults[i].IfGlobalSetUpTearDown;
            newTd2.innerText = tcResults[i].Priority;
            newTd3.innerText = tcResults[i].TcName;
            newTd4.innerText = tcResults[i].ParentTestCase;
            newTd5.innerText = tcResults[i].ActualStatusCode;
            newTd6.innerText = tcResults[i].TestResult;
            newTd7.innerText = tcResults[i].DurationUnixNano / 1000000;
            newTd8.innerText = tcResults[i].JsonFilePath + "\n\n" + tcResults[i].CsvFile  + "\n\n" + tcResults[i].CsvRow;
            newTd9.innerText = JSON.stringify(tcResults[i].HttpTestMessages, null, 4);
          }   
      }
    }

    function searchCriteria(i, v1, v2, v3, v4, v5, v6, v7, v8, v9){
      if (searchStr(tcResults[i].IfGlobalSetUpTearDown, v1) && searchStr(tcResults[i].Priority, v2) 
        && searchStr(tcResults[i].TcName, v3)
        && searchStr(tcResults[i].ParentTestCase, v4) && searchStr(tcResults[i].ActualStatusCode, v5)
        && searchStr(tcResults[i].TestResult, v6)
        && searchStr(tcResults[i].DurationUnixNano / 1000000, v7)
        && searchStr(tcResults[i].JsonFilePath + " / " + tcResults[i].CsvFile  + " / " + tcResults[i].CsvRow, v8)
        && searchStr(JSON.stringify(tcResults[i].HttpTestMessages), v9)) {
        return true
      } else {
        return false
      }
    }

    function searchStr(originStr, subStr) {
      if (subStr.trim() == "") {
        return true
      }

      if (originStr.toString().search(subStr) != -1 ) {
        return true
      } else {
        return false
      }
    }

    function clearRows(){
      var tb = document.getElementById("container_statistics_body");
      var rowNum = tb.rows.length;
      for (i = 0; i < rowNum; i++)
      {
          tb.deleteRow(i);
          rowNum = rowNum - 1;
          i = i - 1;
      }
    }

  </script>

</body>
</html>
`
