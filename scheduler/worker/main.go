package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/scheduler/worker/deleter/downloads"
	"github.com/gophr-pm/gophr/scheduler/worker/indexer/awesome"
	"github.com/gophr-pm/gophr/scheduler/worker/indexer/gosearch"
	ghUpdater "github.com/gophr-pm/gophr/scheduler/worker/updater/github"
	"github.com/gophr-pm/gophr/scheduler/worker/updater/metrics"
	"github.com/gorilla/mux"
)

var (
	// updateMetricsWorkerThreads is the number of go routines elected to process
	// packages in the database.
	updateMetricsWorkerThreads = runtime.NumCPU() * 2
	// indexGoSearchWorkerThreads is the number of go routines elected to index
	// packages from go-search.org.
	indexGoSearchWorkerThreads = runtime.NumCPU() * 2
	// deleteOldDownloadsWorkerThreads is the number of go routines elected to
	// delete old downloads in the database.
	deleteOldDownloadsWorkerThreads = runtime.NumCPU()
)

func main() {
	// Initialize the db client and configuration.
	config, client := lib.Init()

	// Ensure that the client is closed eventually.
	defer client.Close()

	// Initialize datadog client.
	ddClient, err := datadog.NewClient(config, "scheduler-worker.")
	if err != nil {
		log.Fatalln("Failed to create the DataDog client:", err)
	}

	// Create an instance of the github request service.
	ghSvc, err := github.NewRequestService(github.RequestServiceArgs{
		Conf:             config,
		DDClient:         ddClient,
		Queryable:        client,
		ForScheduledJobs: true,
	})
	if err != nil {
		log.Fatalln("Failed to create the Github request service:", err)
	}

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	r.HandleFunc(
		"/update/metrics",
		metrics.UpdateHandler(
			client,
			ddClient,
			updateMetricsWorkerThreads)).Methods("GET")
	r.HandleFunc(
		"/update/github-metadata",
		ghUpdater.UpdateHandler(
			client,
			ghSvc,
			ddClient,
			updateMetricsWorkerThreads)).Methods("GET")
	r.HandleFunc(
		"/index/awesome",
		awesome.IndexHandler(client, ddClient)).Methods("GET")
	r.HandleFunc(
		"/index/go-search",
		gosearch.IndexHandler(
			client,
			config,
			ghSvc,
			ddClient,
			indexGoSearchWorkerThreads)).Methods("GET")
	r.HandleFunc(
		"/delete/old-downloads",
		downloads.DeleteHandler(
			client,
			ddClient,
			deleteOldDownloadsWorkerThreads)).Methods("GET")

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
