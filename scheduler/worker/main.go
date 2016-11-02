package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/scheduler/worker/updater/metrics"
	"github.com/gorilla/mux"
)

var (
	updateMetricsWorkerThreads = runtime.NumCPU() * 2
)

func main() {
	// Initialize the API.
	config, client := lib.Init()

	// Ensure that the client is closed eventually.
	defer client.Close()

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	r.HandleFunc(
		"/update/metrics",
		metrics.UpdateHandler(client, updateMetricsWorkerThreads)).Methods("GET")

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
