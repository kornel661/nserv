nserv [![GoDoc](https://godoc.org/gopkg.in/kornel661/nserv.v0?status.svg)](https://godoc.org/gopkg.in/kornel661/nserv.v0)
=====

[nserv](http://godoc.org/gopkg.in/kornel661/nserv.v0) (nice server) Go package provides a variation of standard http.Server enhanced with *graceful exit* and *throttling*.
*Throttling* makes it easier to build a DOS-attack-resistant server and *graceful exit* feature makes it easy to write a stopable server with proper clean-up (e.g., closing database connections).
Nserv has been inspired by the [manners package](https://github.com/braintree/manners).

The package is in its early stages of developement.
Use at your own risk (or better wait for a version that actually works, should be coming very soon).

Usage
=====

```
go get -u gopkg.in/kornel661/nserv.v0
```
or
```go
import "gopkg.in/kornel661/nserv.v0"
```
See [package import site](http://gopkg.in/kornel661/nserv.v0) and [gopkg.in](http://labix.org/gopkg.in) for import path convention.


Goals
=====

* This package is intended as a light extension of Go's standard http server implementation.
* Performance is important but not by overcomplicating the code.
