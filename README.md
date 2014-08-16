nserv [![GoDoc](https://godoc.org/gopkg.in/kornel661/nserv.v0?status.svg)](https://godoc.org/gopkg.in/kornel661/nserv.v0)
=====

[nserv](https://godoc.org/gopkg.in/kornel661/nserv.v0) (nice server) Go package provides a variation of standard http.Server enhanced with *graceful exit* and *throttling*.
*Throttling* makes it easier to build a DOS-attack-resistant server and *graceful exit* feature makes it easy to write a stoppable server with proper clean-up (e.g., closing database connections).
Nserv has been inspired by the [manners package](https://github.com/braintree/manners).

The package is in its early stages of development (in the sense that hasn't been tested extensively yet).
API of the v0 version might change without warning but once I get reports or comments that current API is reasonable and there are no obvious bugs, version v1 with stable public API will be published.
See [gopkg.in](https://gopkg.in/) on versioning scheme.


Features
========

* Full functionality of the standard http.Server.
* Limiting number of simultaneous connections.
  The limit can be dynamically changed while the server is running.
* Gracefull exit.


Usage
=====

```
go get -u gopkg.in/kornel661/nserv.v0
```
or
```go
import "gopkg.in/kornel661/nserv.v0"
```
See [package import site](https://gopkg.in/kornel661/nserv.v0) and [gopkg.in](https://labix.org/gopkg.in) for import path convention.


Goals
=====

* This package is intended as a light extension of Go's standard http server implementation.
* Performance is important but not by overcomplicating the code.

Changelog
=========

* 2014.08.16 (version v0): Testing & bug hunting season opened.
