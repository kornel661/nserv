nserv
=====

[nserv](https://gopkg.in/kornel661/nserv.v0) (nice server) Go package provides a variation of standard http.Server enhanced with *graceful exit* and *throttling*.
*Throttling* makes it easier to build a DOS-attack-resistant server and *graceful exit* feature makes it easy to write a stoppable server with proper clean-up (e.g., closing database connections).
Nserv has been inspired by the [manners](https://github.com/braintree/manners) package.

The package is in its early stages of development (in the sense that hasn't been tested extensively yet).
API of the v0 version might change without warning, v0 should be considered as a unstable/development version.

See [gopkg.in](https://gopkg.in/) on versioning scheme.

For up-to-date changelog and features list see [README](https://github.com/kornel661/nserv/blob/master/README.md).


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
Replace v0 by the version you need, see [package import site](https://gopkg.in/kornel661/nserv.v0) and [gopkg.in](https://labix.org/gopkg.in) for import path convention.


Versions
========

* Bleeding-edge development version (github.com/kornel661/nserv)
  [![GoDoc](https://godoc.org/github.com/kornel661/nserv?status.svg)](https://godoc.org/github.com/kornel661/nserv)  [![GoWalker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/kornel661/nserv)
* Development version (v0)
  [![GoDoc](https://godoc.org/gopkg.in/kornel661/nserv.v0?status.svg)](https://godoc.org/gopkg.in/kornel661/nserv.v0)  [![GoWalker](https://gowalker.org/api/v1/badge)](https://gowalker.org/gopkg.in/kornel661/nserv.v0) [![GoCover](http://gocover.io/_badge/gopkg.in/kornel661/nserv.v0)](http://gocover.io/gopkg.in/kornel661/nserv.v0)
* Initial version with stable API (v1)
  [![GoDoc](https://godoc.org/gopkg.in/kornel661/nserv.v1?status.svg)](https://godoc.org/gopkg.in/kornel661/nserv.v1)  [![GoWalker](https://gowalker.org/api/v1/badge)](https://gowalker.org/gopkg.in/kornel661/nserv.v1) [![GoCover](http://gocover.io/_badge/gopkg.in/kornel661/nserv.v1)](http://gocover.io/gopkg.in/kornel661/nserv.v1)


Goals
=====

* This package is intended as a light extension of Go's standard http server implementation.
* Performance is important but not by overcomplicating the code.


Changelog
=========

* 2014.08.18 (version v1): Created version v1 - its API should be stable, though
  it isn't well-tested yet. Methods & fields can be added to the Server struct.
* 2014.08.16 (version v0): Testing & bug hunting season opened.
