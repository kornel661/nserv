Example of zero downtime capable server using nserv
=====

See ```server.go``` for one of many possible ways of implementing a server that can be updated without any interruption to the connected clients. Try running the following in your terminal:
```
./ZeroDowntime-example -n=N
```
for some natural N.

It starts serving a simple web page. If N>0, after receiving a SIGINT signal, the process spawns a new server with n=N-1 that will continue to serve on the same port and it exits gracefully.
If N=0 the server terminates gracefully on receiving SIGINT.

One can easily implement other behavior, e.g., with a supervisor that manages the actual  server processes (it might be useful for some init systems).
