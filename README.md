# Tideland Go Libray

[![GitHub release](https://img.shields.io/github/release/tideland/go.svg)](https://github.com/tideland/go)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/tideland/go?status.svg)](https://godoc.org/github.com/tideland/go)
[![Sourcegraph](https://sourcegraph.com/github.com/tideland/go/-/badge.svg)](https://sourcegraph.com/github.com/tideland/go?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go)](https://goreportcard.com/report/tideland.one/go)

## Description

**Tideland Go Library** provides a number of helpful packages for different areas. They are grouped by higher level packages.

- `audit` provides packages to support testing
    + `asserts` provides routines for assertions helpful in tests and validation
    + `capture` allows capturing of STDOUT and STDERR
    + `environments` provides setting of environment variables, creation of temporary directories, and running web servers for tests
    + `generators` simplifies generation of test data; with a fixed random on demand even repeatable
- `db` provides database clients
    + `couchdb` realizes a client for the CouchDB and
    + `redis` for Redis
- `dsa` contains data structures and algorithms
    + `collections` contains collection types like a ring buffer, stacks, sets and trees
    + `identifier` allows the generation of UUIDs in different versions as well as other identifier
    + `mapreduce` provides a generic map/reduce algorithm
    + `sort` contains a parallel quicksort
    + `timex` helps working with times
    + `version` helps managing semantic versioning
- `net` groups the work with the network
    + `jwt` implements a complete JSON Web Token plus caching
    + `webbox` enhances the standard HTTP multiplexing of request to handlers and functions
- `text` simplifies life with text data
    + `etc` manages configurations including internel references to environment variables, cross references, and extraction of subtrees for using types; syntax is `sml` (see below)
    + `gjp` is the generic JSON processing without static type marshalling
    + `scroller` helps analyzing a continuously written line by line text content like log files
    + `sml` is the simple markup language, a LISP like notation using curly braces
    + `stringex` enhances the functionality of the standard library package `strings`
- `together` focusses on goroutines and how to manage them more convenient and reliable
    + `actor` runs a backend goroutine processing anonymous functions for the serialization of changes, e.g. in a structure
    + `crontab` allows running functions at configured times and in chronological order
    + `limiter` limits the number of parallel executing goroutines in its scope
    + `loop` helps running a controlled endless `select` loop for goroutine backends
    + `notifier` helps at the coordination of multiple goroutines
    + `wait` provides a flexible and controlled waiting for conditions by polling
- `trace` helps running applications and servers
    + `errors` is a more powerful error management than the standard package
    + `location` allows to retrieve current file and line, helpful for errors and logging
    + `logging` is a more controllable logging with an exchangeable backend, e.g. syslog
    + `monitor` allows to measure runtimes and monitor variables

I hope you like it. ;)

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)

## License

**Tideland Go Library** is distributed under the terms of the BSD 3-Clause license.
