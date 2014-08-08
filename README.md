nserv
=====

[nserv](http://godoc.org/gopkg.in/kornel661/nserv.v0) (nice server) Go package provides a variation of standard http.Server enhanced with *graceful exit* and *throttling*.
It's been inspired by the [manners package](https://github.com/braintree/manners).

The package is in its early stages of developement.
Use at your own risk (or better wait for a version that actually works, should be coming very soon).

Needs Go 1.3.

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

* This package is intended as a light extension of Go's standard implementation.
* Performance is important but not by overcomplicating the code.
* Right now nserv uses internally slightly modified version of net/http package. Unfortunately, this approach demands some type conversions. Once nserv is more mature and tested this situation will hopefully improve.
