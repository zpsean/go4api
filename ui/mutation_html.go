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

var Mutation = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
    <link href="style/go4api.css" rel="stylesheet" type="text/css"/>
    <script type="text/javascript" src="js/go4api.js"></script>
    <script type="text/javascript" src="js/results.js"></script>
    <script type="text/javascript" src="js/stats.js"></script>
    <script type="text/javascript" src="js/mutationstats.js"></script>
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
                      <div class="item "><a id="mindex_link" href="mindex.html">MutationOverview</a></div>
                      <div class="item selected"><a id="mutation_link" href="mutation.html">Mutation</a></div>
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
                                      <th id="col-2" class="header sortable"><span>Case ID</span></th>
                                      <th id="col-3" class="header sortable"><span>HttpUrl</span></th>
                                      <th id="col-4" class="header sortable"><span>HttpMethod</span></th>
                                      <th id="col-5" class="header sortable"><span>MutationPart</span></th>
                                      <th id="col-6" class="header sortable"><span>MutationCategory</span></th>
                                      <th id="col-7" class="header sortable"><span>MutationRule</span></th>
                                      <th id="col-8" class="header sortable"><span>HttpStatus</span></th>
                                      <th id="col-9" class="header sortable"><span>TestStatus</span></th>
                                      <th id="col-10" class="header sortable"><span>Count</span></th>
                                      <th id="col-11" class="header sortable"><span>MutationMessage</span></th>
                                  </tr>
                              </thead>
                              <thead>
                                  <tr>
                                      <th id="col-1-1"></th>

                                      <th id="col-2-1">
                                        <input type="text" id="caseid_input" list="caseid_datalist" size="10" onchange="dataListChange()" />
                                          <datalist id="caseid_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-3-1">
                                        <input type="text" id="httpurl_input" list="httpurl_datalist" size="10" onchange="dataListChange()" />
                                          <datalist id="httpurl_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-4-1">
                                        <input type="text" id="httmethod_input" list="httmethod_datalist" size="15" onchange="dataListChange()" />
                                          <datalist id="httmethod_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-5-1">
                                        <input type="text" id="MutationPart_input" list="MutationPart_datalist" size="20" onchange="dataListChange()" />
                                          <datalist id="MutationPart_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-6-1">
                                        <input type="text" id="MutationCategory_input" list="MutationCategory_datalist" size="20" onchange="dataListChange()" />
                                          <datalist id="MutationCategory_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-7-1">
                                        <input type="text" id="MutationRule_input" list="MutationRule_datalist" size="18" onchange="dataListChange()" />
                                          <datalist id="MutationRule_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-8-1">
                                        <input type="text" id="HttpStatus_input" list="HttpStatus_datalist" size="12" onchange="dataListChange()" />
                                          <datalist id="HttpStatus_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-9-1">
                                        <input type="text" id="TestStatus_input" list="TestStatus_datalist" size="12" onchange="dataListChange()" />
                                          <datalist id="TestStatus_datalist">
                                        </datalist>
                                      </th>

                                      <th id="col-10-1"><span></span></th>

                                      <th id="col-11-1">
                                        <input type="text" id="MutationMessage_input" list="MutationMessage_datalist" size="18" placeholder="Search" onchange="dataListChange()" />
                                          <datalist id="MutationMessage_datalist">
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
                                for (var i = 0;i < tcResults.length; i++)
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

                                    newTd0.width="30px";
                                    newTd1.width="120px";
                                    newTd2.width="165px";
                                    newTd3.width="50px";
                                    newTd4.width="70px";
                                    newTd5.width="95px";
                                    newTd6.width="185px";
                                    newTd7.width="50px";
                                    newTd8.width="60px";
                                    newTd9.width="40px";
                                    newTd10.width="254px";
                             
                                    newTd0.innerText = i + 1;
                                    newTd1.innerText = tcResults[i].TcName;
                                    newTd2.innerText = tcResults[i].Path;
                                    newTd3.innerText = tcResults[i].Method;
                                    newTd4.innerText = tcResults[i].MutationArea;
                                    newTd5.innerText = tcResults[i].MutationCategory;
                                    newTd6.innerText = tcResults[i].MutationRule;
                                    newTd7.innerText = tcResults[i].ActualStatusCode;
                                    newTd8.innerText = tcResults[i].TestResult;
                                    newTd9.innerText = 1;
                                    newTd10.innerText = JSON.stringify(tcResults[i].MutationInfo, null, 4);
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
    var list1 = new Array(9);
    var li2 = new Array(9);
    var list2 = new Array(9);

    list1[0] = "TcName";
    list1[1] = "Path";
    list1[2] = "Method";
    list1[3] = "MutationArea";
    list1[4] = "MutationCategory";
    list1[5] = "MutationRule";
    list1[6] = "ActualStatusCode";
    list1[7] = "TestResult";
    list1[8] = "Field";

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
    window.onload = function() {
      populateDatalist("caseid_datalist", 0)
      populateDatalist("httpurl_datalist", 1)
      populateDatalist("httmethod_datalist", 2)
      populateDatalist("MutationPart_datalist", 3)
      populateDatalist("MutationCategory_datalist", 4)
      populateDatalist("MutationRule_datalist", 5)
      populateDatalist("HttpStatus_datalist", 6)
      populateDatalist("TestStatus_datalist", 7)
      // populateDatalist("MutationMessage_datalist", 8)
    }

    function dataListChange(){
      var v0 = document.getElementById("caseid_input").value
      var v1 = document.getElementById("httpurl_input").value
      var v2 = document.getElementById("httmethod_input").value
      var v3 = document.getElementById("MutationPart_input").value
      var v4 = document.getElementById("MutationCategory_input").value
      var v5 = document.getElementById("MutationRule_input").value
      var v6 = document.getElementById("HttpStatus_input").value
      var v7 = document.getElementById("TestStatus_input").value
      var v8 = document.getElementById("MutationMessage_input").value

      clearRows()
      insertRows(v0, v1, v2, v3, v4, v5, v6, v7, v8)
    }


    function insertRows(v0, v1, v2, v3, v4, v5, v6, v7, v8){
      for (var i = 0; i < tcResults.length; i++)
      {   
          if (searchCriteria(i, v0, v1, v2, v3, v4, v5, v6, v7, v8)) {
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

            newTd0.width="30px";
            newTd1.width="120px";
            newTd2.width="165px";
            newTd3.width="50px";
            newTd4.width="70px";
            newTd5.width="95px";
            newTd6.width="185px";
            newTd7.width="50px";
            newTd8.width="60px";
            newTd9.width="40px";
            newTd10.width="254px";
     
            newTd0.innerText = i;
            newTd1.innerText = tcResults[i].TcName;
            newTd2.innerText = tcResults[i].Path;
            newTd3.innerText = tcResults[i].Method;
            newTd4.innerText = tcResults[i].MutationArea;
            newTd5.innerText = tcResults[i].MutationCategory;
            newTd6.innerText = tcResults[i].MutationRule;
            newTd7.innerText = tcResults[i].ActualStatusCode;
            newTd8.innerText = tcResults[i].TestResult;
            newTd9.innerText = 1;
            newTd10.innerText = JSON.stringify(tcResults[i].MutationInfo, null, 4);
          }   
      }
    }

    function searchCriteria(i, v0, v1, v2, v3, v4, v5, v6, v7, v8){
      if (searchStr(tcResults[i].TcName, v0) && searchStr(tcResults[i].Path, v1)
        && searchStr(tcResults[i].Method, v2) && searchStr(tcResults[i].MutationArea, v3)
        && searchStr(tcResults[i].MutationCategory, v4) && searchStr(tcResults[i].MutationRule, v5)
        && searchStr(tcResults[i].ActualStatusCode, v6) && searchStr(tcResults[i].TestResult, v7)
        && searchStr(JSON.stringify(tcResults[i].MutationInfo, null, 4), v8)) {
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
