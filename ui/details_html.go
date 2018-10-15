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
                      <div class="item "><a id="graphic_link" href="graphic.html">Graphic</a></div>
                      <div class="item selected"><a id="details_link" href="details.html">Details</a></div>
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
                                  <tr>
                                      <th id="col-1" class="header sortable"><span>#</span></th>
                                      <th id="col-2" class="header sortable"><span>Priority</span></th>
                                      <th id="col-3" class="header sortable"><span>Case ID</span></th>
                                      <th id="col-4" class="header sortable"><span>ParentTestCase</span></th>
                                      <th id="col-5" class="header sortable"><span>Status</span></th>
                                      <th id="col-6" class="header sortable"><span>Case Data</span></th>
                                      <th id="col-7" class="header sortable"><span>Message</span></th>
                                  </tr>
                              </thead>
                              <thead>
                                  <tr>
                                      <th id="col-1-1"></th>

                                      <th id="col-2-1">
                                        <input type="text" id="Priority_input" list="Priority_datalist" size="10" onchange="dataListChange()" />
                                          <datalist id="Priority_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-3-1">
                                        <input type="text" id="caseid_input" list="caseid_datalist" size="30" onchange="dataListChange()" />
                                          <datalist id="caseid_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-4-1">
                                        <input type="text" id="ParentTestCase_input" list="ParentTestCase_datalist" size="30" onchange="dataListChange()" />
                                          <datalist id="ParentTestCase_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-5-1">
                                        <input type="text" id="TestResult_input" list="TestResult_datalist" size="20" onchange="dataListChange()" />
                                          <datalist id="TestResult_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-6-1">
                                        <input type="text" id="CaseData_input" list="CaseData_datalist" size="30" placeholder="Search" onchange="dataListChange()" />
                                          <datalist id="CaseData_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-7-1">
                                        <input type="text" id="Message_input" list="Message_datalist" size="30" placeholder="Search" onchange="dataListChange()" />
                                          <datalist id="Message_datalist">
                                        </datalist>
                                      </th>

                                  </tr>
                              </thead>

                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body" class="statistics-in extensible-geant">
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
                             
                                    newTd0.innerText= i + 1;
                                    newTd1.innerText = tcResults[i].Priority;
                                    newTd2.innerText = tcResults[i].TcName;
                                    newTd3.innerText = tcResults[i].ParentTestCase;
                                    newTd4.innerText = tcResults[i].TestResult;
                                    newTd5.innerText = tcResults[i].JsonFilePath + " / " + tcResults[i].CsvFile  + " / " + tcResults[i].CsvRow;
                                    newTd6.innerText = JSON.stringify(tcResults[i].TestMessages, null, 4);
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

    list1[0] = "Priority";
    list1[1] = "TcName";
    list1[2] = "ParentTestCase";
    list1[3] = "TestResult";
    list1[4] = "CaseData";
    list1[5] = "Message";

    li2[0] = new Array;
    li2[1] = new Array;
    li2[2] = new Array;
    li2[3] = new Array;
    li2[4] = new Array;
    li2[5] = new Array;

    list2[0] = new Array;
    list2[1] = new Array;
    list2[2] = new Array;
    list2[3] = new Array;
    list2[4] = new Array;
    list2[5] = new Array;

    for(var i = 0; i < list1.length; i++)
    {
      for (var j = 0; j < tcResults.length; j++)
      { 
        li2[i].push(tcResults[j][list1[i]])
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
    }
  </script>


  <script type="text/javascript">
    var tcResults;

    window.onload = function() {
      populateDatalist("Priority_datalist", 0)
      populateDatalist("caseid_datalist", 1)
      populateDatalist("ParentTestCase_datalist", 2)
      populateDatalist("TestResult_datalist", 3)
      // populateDatalist("CaseData_datalist", 4)
      // populateDatalist("Message_datalist", 5)
    }

    function dataListChange(){
      var v0 = document.getElementById("Priority_input").value
      var v1 = document.getElementById("caseid_input").value
      var v2 = document.getElementById("ParentTestCase_input").value
      var v3 = document.getElementById("TestResult_input").value
      var v4 = document.getElementById("CaseData_input").value
      var v5 = document.getElementById("Message_input").value

      clearRows()
      insertRows(v0, v1, v2, v3, v4, v5)
    }


    function insertRows(v0, v1, v2, v3, v4, v5){
      for (var i = 0; i < tcResults.length; i++)
      {   
          if (searchCriteria(i, v0, v1, v2, v3, v4, v5)) {
            var newTr = container_statistics_body.insertRow();
                                    
            var newTd0 = newTr.insertCell();
            var newTd1 = newTr.insertCell();
            var newTd2 = newTr.insertCell();
            var newTd3 = newTr.insertCell();
            var newTd4 = newTr.insertCell();
            var newTd5 = newTr.insertCell();
            var newTd6 = newTr.insertCell();

            newTd0.innerText = i + 1;
            newTd1.innerText = tcResults[i].Priority;
            newTd2.innerText = tcResults[i].TcName;
            newTd3.innerText = tcResults[i].ParentTestCase;
            newTd4.innerText = tcResults[i].TestResult;
            newTd5.innerText = tcResults[i].JsonFilePath + " / " + tcResults[i].CsvFile  + " / " + tcResults[i].CsvRow;
            newTd6.innerText = JSON.stringify(tcResults[i].TestMessages, null, 4);
          }   
      }
    }

    function searchCriteria(i, v0, v1, v2, v3, v4, v5){
      if (searchStr(tcResults[i].Priority, v0) && searchStr(tcResults[i].TcName, v1)
        && searchStr(tcResults[i].ParentTestCase, v2) && searchStr(tcResults[i].TestResult, v3)
        && searchStr(tcResults[i].JsonFilePath + " / " + tcResults[i].CsvFile  + " / " + tcResults[i].CsvRow, v4)
        && searchStr(JSON.stringify(tcResults[i].TestMessages), v5)) {
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
