/*
Package nserv (nice server) provides a variation of the standard http.Server
implementation enhanced with graceful exit and throttling.

Throttling enables you to limit number of simultaneous connections the server
accepts. This limit can be changed even after the server had been started.
Graceful exit means you can signal the server to stop and at that point the
server stops accepting new connections. Active connections run their natural
course and only after all connections are closed the server shuts down.

Usage:

	import "gopkg.in/kornel661/nserv.v1"

or

	go get gopkg.in/kornel661/nserv.v1

Replace v1 by the version you need.

For up-to-date changelog and features list see [README]
(https://github.com/kornel661/nserv/blob/master/README.md).
*/
package nserv
