<h3 align="center">Go4Api - an API testing tool written in Go</h3>
<p align="center">Implementing Data-Driven Test, Mutation Test, Fuzz Test.</p>

---

**Go4Api** is a tool focusing on API testing, which is targeting the huge test cases and test data, with execution concurrently based on Priority and Dependency.

<p align="center">
  <img width="600" src="https://cdn.rawgit.com/zpsean/go4api/master/demo5.svg">
</p>


Menu
----
- [Why another API Testing Tool?](#why-another-api-testing-tool)
- [Features](#features)
- [Install](#install)
- [Basic Concepts](#basic-concepts)
- [Quick Start](#quick-start)
- [Need help or want to contribute?](#need-help-or-want-to-contribute)
  

Why another API Testing Tool?
--------
Go4Api aims to the testing difficulty faced to QA, which is different from Developers. That is, plan and execute a single api test is easy, a bunch of tools can help on it. But how about if we have hundreds of API(s) and thousands of test data to manage and execute, and more, regression them during the API(s) lifetime?

Features
--------

- **Using pure Json format to represent test case(s)**: contains the all info about API request, response and assertion 
- **Each Case has its own setUp and tearDown** 
- **Test Cases Json file itself can be template**: its variable(s) can be rendered by data feeder 
- **Test Cases executed concurrently**: based on Priority and Dependency
- **Fuzz / Mutation Testing**: includes Mutation and Random testing (embedded pairwise algorithm implementation) 
- **Scenario Testing**: when APIs have data dependency
- **Convert HAR file / Swagger API file** 
- **DB manipulation support**: MySql, PostgreSql, MongoDB, Redis, and easy extension to other database(s)  
- **Built-in functions** 
- **User defined functions**: use ECMAScript 5.1(+)  

More information, refer to [wiki](https://github.com/zpsean/go4api/wiki) and Blog:  
[（上篇）API测试在测些什么，这些方向你值得拥有](https://www.jianshu.com/p/d5909638c82a).  
[（下篇）吃自己的狗粮：Go4Api之API测试项目实践](https://www.jianshu.com/p/cc536d6a1898).  
[测试与自动化测试，记测试工具Go4Api的诞生](https://www.jianshu.com/p/1d9241ec0c3e).  


Install and Run
------

### Option 1: Mac: Using the binary package, Run

Grab a prebuilt binary from [the Releases page](https://github.com/zpsean/go4api/releases).

Copy the binary in your _PATH_ to run go4api from any location.


### Option 2: Build from source, Run
To build from source you need **[Go](https://golang.org/doc/install)** (1.10 or newer). Follow these instructions:

- Run `go get github.com/zpsean/go4api` which will:
  - git clone the repo and put the source in `$GOPATH/src/github.com/zpsean/go4api`
  - build a `go4api` binary and put it in `$GOPATH/bin`, if for linux, use `env GOOS=linux GOARCH=amd64 go build`
- Make sure you have `$GOPATH/bin` in your PATH
- You can now run go4api using `$./go4api ...`


### Option 3: Run from source
To run from source you need **[Go](https://golang.org/doc/install)** (1.10 or newer). Follow these instructions:

- Run `go get github.com/zpsean/go4api` which will:
  - git clone the repo and put the source in `$GOPATH/src/github.com/zpsean/go4api`
  - Move to the path: `cd $GOPATH/src/github.com/zpsean/go4api`
- Make sure you have `$GOPATH` in your PATH
- You can now run go4api using `$go run main.go ...`


Basic Concepts
-----------
<p align="center">
  <img width="900" src="https://cdn.rawgit.com/zpsean/go4api/master/doc/1-CaseStructure.jpeg">
</p>
<p align="center">
  <img width="900" src="https://cdn.rawgit.com/zpsean/go4api/master/doc/2-CasesRelationship.jpeg">
</p>
<p align="center">
  <img width="900" src="https://cdn.rawgit.com/zpsean/go4api/master/doc/3-BigPicture.jpeg">
</p>

Quick start
-----------

Note: You can prepare many many test cases based on below examples to let Go4Api run for you.

### Your testing workspace may like below:
```js
samples/
├── mutation
│   └── MutationTeseCase.json
├── scenarios
│   └── scenario1
│       ├── s1ChildChildChildTeseCase.json
│       ├── s1ChildChildTeseCase.json
│       ├── s1ChildTeseCase.json
│       ├── s1ParentTeseCase.json
│       └── temp
│           ├── _join.csv
│           ├── s1ParentTestCase_out.csv
│           └── s1ParentTestCase_out2.csv
├── testconfig
│   └── config.json
├── testdata
│   └── Demo
│       ├── FirstTeseCase.json
│       ├── SecondTeseCase.json
│       ├── SecondTeseCase_dt1.csv
│       └── SecondTeseCase_dt2.csv
└── testresource
    └── swagger.json

testresults/
└── 2018-09-10\ 07:42:20.804070777\ +0800\ CST\ m=+0.001524050
    ├── 2018-09-10\ 07:42:20.804070777\ +0800\ CST\ m=+0.001524050.log
    ├── index.html
    ├── js
    └── style
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
      "response": [
        {
          "$(status).statusCode": {
            "Equals": 200
          }
        },
        {
          "$(headers).Content-Type": {
            "Contains": "application/json;charset=UTF-8"
          }
        },
        {
          "$(body).start": {
            "GreaterOrEquals": 0
          }
        },
        {
          "$(body).subjects.#": {
            "Equals": 20
          }
        },
        {
          "$(body).total": {
            "Equals": 250
          }
        },
        {
          "$(body).subjects.0.title": {
            "Contains": "肖申克的救赎"
          }
        }
      ]
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
    "SecondTestCase-${tc}": {
      "priority": "${priority}",
      "parentTestCase": "root",
      "request": {
        "method": "GET",
        "path": "/v2/movie/top250",
        "headers": {
          "authorization": "${authorization}"
        },
        "queryString": {
          "pageIndex": "1",
          "pageSize": "12"
        }
      },
      "response": [
        {
          "$(status).statusCode": {
            "Equals": {"Fn::ToInt": "${statuscode}"}
          }
        },
        {
          "$(headers).Content-Type": {
            "Contains": "application/json;charset=UTF-8"
          }
        }
      ]
    }
  }
]
```

SecondTeseCase_dt1.csv
```js
tc,priority,statuscode
dt1-1,1,500
dt1-2,2,500
```

#### Running go4api

```js
$./go4api -run -baseUrl https://api.douban.com -c /<you Path>/testconfig  -tc  /<you Path>/testdata -tr /<you Path>/testresource -r /<you Path>/testresults
```


More code snippets
-----------

#### Config file sample
```js
{
  "Local": {
      "baseUrl": "...",
      "mysql": {
          "master": {
              "ip": "",
              ...
          }
      },
      "postgresql": {
          "master": {
              "ip": "",
              ...
          }
      },
      "redis": {
          "master": {
              ...
          }
      },
      "mongoDB": {
          "master": {
              ...
          }
      }
  },
  "Dev": ...,
  "QA": ...,
  "UAT": ...
}
```

#### POST, application/json
```js
[
  {
    "tcname ...": {
      ...,
      "request": {
        "method": "POST",
        "path": "...",
        "headers": {
          "Content-Type": "application/json;charset=UTF-8"
        },
        "payload": {
          "text": {
                    "f_1": "v_1",
                    "f_2": "v_2"
                  }
        }
      },
      "response": ...
    }
  }
]
```

#### POST, multipart/form-data
```js
[
  {
    "tcname ...": {
      ...,
      "request": {
        "method": "POST",
        "path": "...",
        "headers": {
          "Content-Type": "multipart/form-data"
        },
        "payload": {
          "multipartForm": [
                  { 
                    "name": "fname_1",
                    "value": "fvalue_1"
                  },
                  {
                    "name": "fname_2",
                    "value": "fvalue_2.csv",
                    "type": "file",
                    "mIMEHeader": {
                      "content-type": "text/csv"
                    }
                  }
                ]
        }
      },
      "response": ...
    }
  }
]

```




Html Reporting
-----------
<p align="center">
  <img width="900" src="https://cdn.rawgit.com/zpsean/go4api/master/doc/4-html-report.png">
</p>

---


Reference
--------------------------------
- [gjson](https://github.com/tidwall/gjson)
- [sjson](https://github.com/tidwall/sjson)
- [go-linq](https://github.com/ahmetb/go-linq)
- [mysql](https://github.com/go-sql-driver/mysql)
- [redigo](https://github.com/gomodule/redigo)
- [xlsx](https://github.com/tealeg/xlsx)
- [goja](https://github.com/dop251/goja)
- [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)
- [pq](https://github.com/lib/pq)

