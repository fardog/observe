package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/fardog/observe"
)

var (
	logLevel = flag.String(
		"level",
		"info",
		"Log level, one of: debug, info, warn, error, fatal, panic",
	)

	listenAddress = flag.String(
		"listen", ":80", "listen address, as `[host]:port`",
	)
	shutdownTimeout = flag.Int(
		"timeout", 10, "time in seconds to hold shutdown for connected clients",
	)
	serverHeader = flag.String(
		"server-header",
		"observe/1.0 (+https://github.com/fardog/observe)",
		`Value to send in the Server header; if set to an empty string, no
        header will be sent.`,
	)

	gcloudProject = flag.String(
		"gcloud-project-id",
		"default-project",
		"The Google Cloud Project ID in which the bigquery dataset is stored.",
	)
	bigQueryTable = flag.String(
		"bigquery-table",
		"observe.observations",
		"The Google BigQuery table to use within the project.",
	)
)

func serve(server *http.Server) {
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	// serve until exit
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Infoln("shutting down on interrupt")
	timeout := time.Duration(*shutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("got unexpected error %s", err.Error())
	}

	<-ctx.Done()
}

func main() {
	flag.Usage = func() {
		_, exe := filepath.Split(os.Args[0])
		fmt.Fprint(os.Stderr, "Simple traffic analytics collection for static websites.")
		fmt.Fprintf(os.Stderr, "Usage:\n\n  %s [options]\n\nOptions:\n\n", exe)
		flag.PrintDefaults()
	}
	flag.Parse()

	// seed the global random number generator, used in some utilities and the
	// google provider
	rand.Seed(time.Now().UTC().UnixNano())

	// set the loglevel
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("invalid log level: %s", err.Error())
	}
	log.SetLevel(level)

	opts := &observe.BigQueryOptions{}
	store, err := observe.NewBigQuery(*gcloudProject, *bigQueryTable, opts)
	if err != nil {
		log.Fatalf("unable to instantiate bigquery client: %s", err.Error())
	}

	handler := observe.NewHandler(store, &observe.HandlerOptions{
		ServerHeader: *serverHeader,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/observe.gif", handler.Handle)
	server := &http.Server{
		Addr:    *listenAddress,
		Handler: mux,
	}

	// start the servers
	servers := make(chan bool)
	go func() {
		serve(server)
		servers <- true
	}()

	log.Infof("server started on %v", *listenAddress)
	<-servers
	log.Infoln("servers exited, stopping")

}
