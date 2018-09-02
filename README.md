<h3 align="center">Go4Api - a API testing tool written in Go</h3>
<p align="center">Implementing Data-Driven Test, written in Go.</p>

---

**Go4Api** is a tool focusing on API testing, which is targeting the huge test cases and test data, with execution concurrently based on Priority and Dependency.

<p align="center">
  <img width="600" src="https://cdn.rawgit.com/zpsean/go4api/master/demo2.svg">
</p>


Menu
----

- [Features](#features)
- [Install](#install)
- [Quick Start](#quick-start)
- [v1.0 Target](#v1.0-target)
- [Need help or want to contribute?](#need-help-or-want-to-contribute)

Features
--------

- **Using Json format to represent API test case(s)**: contains the all info about API request, response and assertion
- **Test Cases Json file itself can be template**: which will be rendered by csv data table(s)
- **Test Cases executed concurrently**: based on Priority and Dependency
- **Support API Fuzz Testing**: includes Mutation and Random testing (embedded pairwise algorithm implementation)
- **Support Scenario Testing**: when APIs have data dependency and exchange
- **Convert HAR file / Swagger API file**


Install
------

### Mac: Using the binary package

Grab a prebuilt binary from [the Releases page](https://github.com/zpsean/go4api/releases).

Copy the binary in your _PATH_ to run go4api from any location.


### Build from source
To build from source you need **[Go](https://golang.org/doc/install)** (1.10 or newer). Follow these instructions:

- Run `go get github.com/zpsean/go4api` which will:
  - git clone the repo and put the source in `$GOPATH/src/github.com/zpsean/go4api`
  - build a `go4api` binary and put it in `$GOPATH/bin`
- Make sure you have `$GOPATH/bin` in your PATH
- You can now run go4api using `go4api`

Quick start
-----------

Note: You can prepare many many test cases based on below examples to let go4api run for you.

### Your testing workspace will like below:
```js
├── testhome
│   ├── testconfig
│   │   └── testconfig.json
│   ├── testdata
│   │   ├── FirstTeseCase.json
│   │   ├── SecondTeseCase.json
│   │   ├── SecondTeseCase_dt1.csv
│   │   └── SecondTeseCase_dt2.csv
│   └── testresource
│       ├── excelforupload.xlsx
│       └── image.png
└── testresults
    └── 2018-08-06\ 10:55:55.853034228\ +0800\ CST\ m=+0.001557642
```

### A simple case, with hard-coded Json:
#### Prepare the Json:

```js
[
  {
    "FirstTestCase-001": {
      "priority": "3",
      "parentTestCase": "root",
      "request": {
        "method": "GET",
        "path": "https://api.douban.com/v2/movie/top250",
        "headers": {
          "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"
        },
        "queryString": {
          "pageIndex": "1",
          "pageSize": "12"
        }
      },
      "response": {
        "status": {
          "Equals": 200
        },
        "headers": {
          "Content-Type": {
            "Contains": "application/json;charset=UTF-8"
          }
        },
        "body": {
          "total": {
            "Equals": 250
          }
        }
      }
    }
  }
]
```

#### Running go4api

```js
$./go4api -run -c /<you Path>/testconfig  -tc  /<you Path>/testdata -tr /<you Path>/testresource -r /<you Path>/testresults
```


### A much more real case, with variables in Json:
#### Prepare the Json:

SecondTeseCase.json
```js
[
  {
    "SecondTestCase-{{.tc}}": {
      "priority": "{{.priority}}",
      "parentTestCase": "root",
      "request": {
        "method": "GET",
        "path": "/v2/movie/top250",
        "headers": {
          "authorization": "{{.authorization}}"
        },
        "queryString": {
          "pageIndex": "1",
          "pageSize": "12"
        }
      },
      "response": {
        "status": {
          "Equals": {{.statuscode}}
        },
        "headers": {
          "Content-Type": {
            "Contains": "application/json;charset=UTF-8"
          }
        }
      }
    }
  }
]
```

SecondTeseCase_dt1.csv
```js
tc,priority,statuscode
dt1-1,1,500
dt1-2,1,500
```

#### Running go4api

```js
$./go4api -run -baseUrl https://api.douban.com -c /<you Path>/testconfig  -tc  /<you Path>/testdata -tr /<you Path>/testresource -r /<you Path>/testresults
```

---

CheatSheet
--------------------------------
<p align="center">
  <img width="900" src="https://cdn.rawgit.com/zpsean/go4api/master/CheatSheet.png">
</p>

---

v1.0 Target
--------------------------------

v1.0 is planning to have:

- Fully support the HTTP method on Get, Post, Put, Delete
- More options to control the test cases execution
- Fully coverage on Assertion on Equals, Contains, etc.


---

Need help or want to contribute?
--------------------------------

Types of questions and where to ask:

- How do I?