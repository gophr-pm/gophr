package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the API.
	config, client := lib.Init()

	// Ensure that the client is closed eventually.
	defer client.Close()

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/blob/{%s}/{%s}/{%s}/{%s}",
		urlVarAuthor,
		urlVarRepo,
		urlVarSHA,
		urlVarPath),
		BlobHandler()).Methods("GET")
	r.HandleFunc(
		"/packages/new",
		GetNewPackagesHandler(client)).Methods("GET")
	r.HandleFunc(
		"/packages/search",
		SearchPackagesHandler(client)).Methods("GET")
	r.HandleFunc(
		"/packages/trending",
		GetTrendingPackagesHandler(client)).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/packages/top/{%s}/{%s}",
		urlVarLimit,
		urlVarTimeSplit),
		GetTopPackagesHandler(client)).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/packages/{%s}/{%s}",
		urlVarAuthor,
		urlVarRepo),
		GetPackageHandler(client)).Methods("GET")

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
