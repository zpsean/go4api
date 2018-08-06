<h3 align="center">Go4Api - a API testing tool written in Go</h3>
<p align="center">Implementing Data-Driven Test, written in Go.</p>

---

**Go4Api** is a tool focusing on API testing, which is targeting the huge test cases and test data, with execution concurrently based on Priority and Dependency.

<p align="center">
  <img width="600" src="https://">
</p>

Menu
----

- [Features](#features)
- [Install](#install)
- [Quick Start](#quick-start)
- [Need help or want to contribute?](#need-help-or-want-to-contribute)

Features
--------

- **Using the Json to represents the all info for API test case**
- **Json can be template wich render from csv data table(s)**
- **Json structured in tree with Priority and Dependency**
- **Test Cases executed concurrently based on Priority and Dependency**


Install
------

### Using the binary package

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

### Prepare the Json:

```js
{
  "TestCases": [
    {
      "Sku-TC-{{.tc}}": {
        "priority": "{{.priority}}",
        "parentTestCase": "root",
        "request": {
          "method": "GET",
          "path": "/api/operation/skus?pageIndex=1&pageSize=12",
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
}
```

### Running go4api

First run:

```js
$go4api --testEnv QA --testhome /<you Path>/go/run/testhome --testresults /<you Path>/go/run/testresults
```

---

Need help or want to contribute?
--------------------------------

Types of questions and where to ask:

- How do I?