nserv
=====

nserv (nice server) Go package provides a variation of standard http.Server enhanced with *graceful exit* and *throttling*.
It's been inspired by the [manners package]().

The package is in its early stages of developement. Use at your own risk.


Usage
=====

```go
go get -u gopkg.in/kornel661/nserv.v0
```
or
```
import "gopkg.in/kornel661/nserv.v0"
```
See [package import site](http://gopkg.in/kornel661/nserv.v0) and [gopkg.in](http://labix.org/gopkg.in) for import path convention.


Goals
=====

* This package is intended as a light extension of Go's standard implementation.
* Performance is important but not by overcomplicating the code or reimplementing standard library (more performance-oriented version of the package may follow).
