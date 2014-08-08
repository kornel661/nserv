/*
Package nserv (nice server) provides a variation of the standard http.Server
implementation enhanced with graceful exit and throttling.

Usage:

import "gopkg.in/kornel661/nserv.v0"

or

go get gopkg.in/kornel661/nserv.v0

alternative path (to the developement version):
github.com/kornel661/nserv

Note:
If you use the standard net/http package you may want to switch to
gopkg.in/kornel661/http.v0
which is net/http with some new exported symbols. It integrates better with
nserv (as nserv uses it internally).

Needs Go 1.3.

*/
package nserv
