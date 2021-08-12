# Go Ethereum proxy server

[![GoDoc](https://godoc.org/github.com/qiangxue/go-rest-api?status.png)](http://godoc.org/github.com/qiangxue/go-rest-api)
[![Build Status](https://github.com/qiangxue/go-rest-api/workflows/build/badge.svg)](https://github.com/qiangxue/go-rest-api/actions?query=workflow%3Abuild)
[![Code Coverage](https://codecov.io/gh/qiangxue/go-rest-api/branch/master/graph/badge.svg)](https://codecov.io/gh/qiangxue/go-rest-api)
[![Go Report](https://goreportcard.com/badge/github.com/qiangxue/go-rest-api)](https://goreportcard.com/report/github.com/qiangxue/go-rest-api)

This software promotes hexagonal architecture, which is imho best for Go microservices or similar solutions. It implements best practices that follow the [SOLID principles](https://en.wikipedia.org/wiki/SOLID)
and [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). 
It encourages writing clean and idiomatic Go code. 

Functionalities:

* Endpoint for Ethereum latest block proxy: /block/latest
* Endpoint for Ethereum block by number proxy: /block/123456
* Endpoint for Ethereum transaction by block number and transaction index proxy: /block/123456/transaction/3
* EthereumClient & EthereumBlockCache are abstracted at import with IEthereumClient & IEthereumCache interfaces
* Go routine is implemented to fetch Latest Block Number every 3 seconds
* Go routine is implemented to clear expired items from cache every 1 second
* Implements data validation and bottom up error handling and logging
* Only test that is implemented is integration test: /internal/application/controller_test.go
* Test executes configured large number of requests showing latency, cache capacity, number of requests per block etc
* However, application is maximally decoupled so unit testing is easy to do
* Apache ab tests for heavy load testing are in /cmd/ab directory
* Application runs about 200% faster than example application
* Graceful shutdown is implemented by 10 seconds grace period executed by kill SIGNAL 1
* Healthcheck is implemented by /healthcheck
* There are Dockefile and docker-compose.yml attached
* See Makefile for available commands
 
## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://golang.org/doc/install) to
install Go on your computer. The kit requires **Go 1.16 or above**.

[Docker](https://www.docker.com/get-started) is also needed if you want to try the kit without setting up your
own server. The kit requires **Docker 17.05 or higher** for the multi-stage build support.

After installing Go and Docker, run the following commands to start experiencing this starter kit:

```shell
# download the starter kit
git clone https://github.com/divilla/ethproxy.git

cd ethproxy

# run the RESTful API server
make run

# or run the API server with live reloading, which is useful during development
# requires fswatch (https://github.com/emcrisostomo/fswatch)
make run-live
```

At this time, you have a RESTful API server running at `http://127.0.0.1:8080`. It provides the following endpoints:

* `GET /healthcheck`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
* `GET /block/latest`: latest Ethereum block
* `GET /block/:bnr`: Ethereum block, by integer block number
* `GET /block/:bnr/transaction/:tid`: Ethereum transaction, by integer block number and integer transaction index

Try the URL `http://localhost:8080/healthcheck` in a browser, and you should see something like `"OK v1.0.0"` displayed.

Heavy load testing is prepared via `Apache ab` that can be downloaded at [Postman](https://www.getpostman.com/)), you may try the following 
more complex scenarios:

```shell
# Ubuntu / Debian
sudo apt-get install apache2-utils
# CentOS
sudo yum install httpd-tools
# Fedora
sudo dnf install httpd-tools
```

To use ab load tester tool go to cmd/ab and check the scripts. -n sets total number of requests and -c number of concurrent requests.


## Project Layout

The starter kit uses the following project layout:
 
```
.
├── cmd                  main applications of the project
│   └── server           main file
├── config               configuration file
├── internal             private application and library code
│   ├── application      controller and service of main application
│   ├── healthcheck      healthcheck feature
└── pkg                  reusable packages made from scratch
   ├── ethcache         decoupled caching package
   ├── ethclient        client for fetching Ethereum / disabled multirequest for same resource
   └── jsonclient       decoupled json client with isolated Poster interface
```

The top level directories `cmd`, `internal`, `pkg` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Within `internal` and `pkg`, packages are structured by features in order to achieve the so-called
[screaming architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html). For example, 
the `album` directory contains the application logic related with the album feature. 

Within each feature package, code are organized in layers (API, service, repository), following the dependency guidelines
as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

### Managing Configurations

Due to lack of time config is implemented as package constants