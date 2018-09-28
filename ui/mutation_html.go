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

var Mutation = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
    <link href="style/go4api.css" rel="stylesheet" type="text/css"/>
    <script type="text/javascript" src="js/go4api.js"></script>
    <script type="text/javascript" src="js/reslts.js"></script>
    <script type="text/javascript" src="js/stats.js"></script>
    <script type="text/javascript" src="js/executed.js"></script>
    <script type="text/javascript" src="js/notexecuted.js"></script>
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
                      <div class="item "><a id="graphic_link" href="graphic.html">Graphic</a></div>
                      <div class="item "><a id="details_link" href="details.html">Details</a></div>
                      <div class="item "><a id="mindex_link" href="mindex.html">MutationOverview</a></div>
                      <div class="item selected"><a id="mutation_link" href="mutation.html">Mutation</a></div>
                      <div class="item "><a id="fuzz_link" href="fuzz.html">Fuzz</a></div>

                      <script type="text/javascript">
                        // var timestamp = 1523957748602;
                        // var runStartHumanDate = moment(timestamp).format("YYYY-MM-DD HH:mm:ss Z");
                        document.writeln("<p class='sim_desc' title='" +"Started at 2018-xx-xx, duration : 10 seconds' data-content=''>");
                        document.writeln("<b>" + "Started at 2018-xx-xx, duration : 10 seconds </b>");
                        document.writeln("</p>");
                      </script>
                  </div>


                  <div class="content-in">
                    <div> 
                        <select id = "mySelect" onchange="selectKey(this);" style="width:150px;"></select>
                        <select id = "mySelect2" style="width:300px;"></select>

                        <input type="text" id="myInput" size="50" name="search_text" placeholder="Please enter search text here">
                        <button type="button" onClick="btnClick()">Search</button>
                    </div>

                    <h1><span>> </span>Overview Information</h1>
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
                              <tbody></tbody>
                          </table>

                          <div class="scrollable">
                              <table id="container_statistics_body" class="statistics-in extensible-geant">
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


  <script language='javascript'>
    function searchCriteria(i, x, y){
      // console.log("params: ", i, x, tcResults[i].Priority, y)
      switch(x)
      {
      case 0:
        if (tcResults[i].TcName == y) {
          return true
        } else {
          return false
        }
        break;
      case 1:
        if (tcResults[i].Path == y) {
          return true
        } else {
          return false
        }
        break;
      case 2:
        if (tcResults[i].Method == y) {
          return true
        } else {
          return false
        }
        break;
      case 3:
        if (tcResults[i].MutationArea == y) {
          return true
        } else {
          return false
        }
        break;
      case 4:
        if (tcResults[i].MutationCategory == y) {
          return true
        } else {
          return false
        }
        break;
      case 5:
        if (tcResults[i].MutationRule == y) {
          return true
        } else {
          return false
        }
        break;
      case 6:
        if (tcResults[i].ActualStatusCode == y) {
          return true
        } else {
          return false
        }
        break;
      case 7:
        if (tcResults[i].TestResult == y) {
          return true
        } else {
          return false
        }
        break;
      default:
        console.log("no valid select option selected")
        return false
      }
    }

    function btnClick(){
      var x = document.getElementById("mySelect")
      // var y = document.getElementById("myInput")
      var y = document.getElementById("mySelect2")

      clearRows()
      insertRows(x.selectedIndex, y.value)
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

    function insertRows(x, y){
      for (var i = 0; i < tcResults.length; i++)
      { 
          if (searchCriteria(i, x, y)) {
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
  </script>

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
      // console.log(li2[i])

      var distinctItems = Array.from(new Set(li2[i]))
      list2[i] = distinctItems
      // console.log(list2[i])
    }
    

    var firstSelect = document.getElementById("mySelect");
    var secondSelect = document.getElementById("mySelect2");

    for(var i = 0; i < list1.length; i++)
    {
      var option = document.createElement("option");
      option.appendChild(document.createTextNode(list1[i]));
      option.value = list1[i];
      firstSelect.appendChild(option);
    }

    var firstKeyValue = list2[0];
    // console.log(firstKeyValue)
    for (var j = 0; j < firstKeyValue.length; j++) {
        var option2 = document.createElement("option");
        option2.appendChild(document.createTextNode(firstKeyValue[j]));
        option2.value = firstKeyValue[j];
        secondSelect.appendChild(option2);
    }

    function indexof(obj, value)
    {
      var k = 0;
      for(; k < obj.length; k++)
      {
          if(obj[k] == value)
          return k;
      }
      return k;
    }

    function selectKey(obj) {
      secondSelect.options.length = 0;
      var index = indexof(list1,obj.value);
      var list2Element = list2[index];
      for(var i = 0; i < list2Element.length; i++)
      {
          var option = document.createElement("option");
          option.appendChild(document.createTextNode(list2Element[i]));
          option.value = list2Element[i];
          secondSelect.appendChild(option);
      }
    }
  </script>

</body>
</html>
`
