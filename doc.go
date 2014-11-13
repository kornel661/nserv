/*
Package nserv (nice server) provides a variation of the standard http.Server
implementation enhanced with graceful exit, throttling and zero downtime
restarts.

Throttling enables you to limit number of simultaneous connections the server
accepts. This limit can be changed even after the server had been started.

Graceful exit means you can signal the server to stop and at that point the
server stops accepting new connections. Active connections run their natural
course and only after all connections are closed the server shuts down.

Zero downtime restarts allow the server to hand off responsibility of handling
incoming connections to another program (e.g., updated version of the server)
without interrupting the clients.

Usage:

	import "gopkg.in/kornel661/nserv.v0"

or

	go get gopkg.in/kornel661/nserv.v0

Replace v0 by the version you need (v0 is a development version, with no API
stability guaratees).

For up-to-date changelog and features list see [README]
(https://github.com/kornel661/nserv/blob/master/README.md).
*/
package nserv // import "gopkg.in/kornel661/nserv.v0"
