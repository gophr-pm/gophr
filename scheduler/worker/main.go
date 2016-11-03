package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/scheduler/worker/indexer/awesome"
	"github.com/gophr-pm/gophr/scheduler/worker/indexer/gosearch"
	"github.com/gophr-pm/gophr/scheduler/worker/updater/metrics"
	"github.com/gorilla/mux"
)

// updateMetricsWorkerThreads is the number of go routines elected to process
// packages in the database.
var updateMetricsWorkerThreads = runtime.NumCPU() * 2

func main() {
	// Initialize the db client and github service.
	var (
		config, client = lib.Init()
		ghSvc          = github.NewRequestService(github.RequestServiceArgs{
			Conf:       config,
			Queryable:  client,
			ForIndexer: true,
		})
	)

	// Ensure that the client is closed eventually.
	defer client.Close()

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	r.HandleFunc(
		"/update/metrics",
		metrics.UpdateHandler(
			client,
			ghSvc,
			updateMetricsWorkerThreads)).Methods("GET")
	r.HandleFunc(
		"/index/awesome",
		awesome.IndexHandler(client)).Methods("GET")
	r.HandleFunc(
		"/index/go-search",
		gosearch.IndexHandler(client)).Methods("GET")

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
