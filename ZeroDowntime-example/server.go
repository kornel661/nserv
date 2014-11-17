// This program illustrates how to write a simple server with zero-downtime
// restarts using the nserv package. Try running the program with -n=X command
// line option for some natural number X>0.
// Sending SIGINT signal (usually ctrl+c) to the server will result in zero
// downtime restart if X>0.
package main // import "gopkg.in/kornel661/nserv.v0/ZeroDowntime-example"

import (
	"flag"
	"fmt"
	"gopkg.in/kornel661/nserv.v0"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// record the time of  launch
	startTime := time.Now()
	defer func() {
		log.Printf("Exitting instance started at %v.\n", startTime)
	}()

	// how many restarts left?
	numRestarts := flag.Int("n", 0, "-n=[number of zero-downtime restarts before exit]")
	if *numRestarts < 0 {
		*numRestarts = 0
	}

	// prepare for zero downtime restart
	nserv.InitializeZeroDowntime()
	flag.Parse()

	// set-up server:
	srv := &nserv.Server{}
	srv.Addr = "localhost:12345"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, I was launched at %v. Still %d zero-downtime restarts to go.", startTime, *numRestarts)
	})

	// catch signals:
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		<-signals
		if *numRestarts <= 0 { // stop
			log.Println("Caught signal. Shutting down gracefully.")
			srv.Stop()
		} else { // restart
			log.Println("Caught signal. Trying to restart with zero downtime.")
			// prepare the command-line argument (the number of restarts to go)
			arg := fmt.Sprintf("-n=%d", *numRestarts-1)
			// the line below stops the server (srv.Serve will terminate only
			// after all connections are finished) and runs this program with
			// the following arguments:
			//     * arg
			//     * some internal nserv flag that specifies which file descriptor
			//       to use for srv.ResumeServe()
			srv.ZeroDowntimeRestart(arg)
		}
	}()

	// start or resume serving:
	log.Printf("Hello, I was launched at %v. Still %d zero-downtime restarts to go.\n", startTime, *numRestarts)
	if !nserv.CanResume() { // start serving
		log.Printf("Serving at http://%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	} else { // try to resume serving
		log.Printf("Trying to resume serving at the previous address.\n")
		err := srv.ResumeAndServe()
		if err != nil {
			log.Println("Couldn't resume serving! Buu! (or other error)")
			log.Println(err)
		}
	}
}
