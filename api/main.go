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
	config, session := common.Init()

	// Ensure that the session is closed eventually.
	defer session.Close()

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	// r.HandleFunc("/packages/search", SearchHandler(session)).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/blob/{%s}/{%s}/{%s}/{%s}",
		blobHandlerURLVarAuthor,
		blobHandlerURLVarRepo,
		blobHandlerURLVarSHA,
		blobHandlerURLVarPath),
		BlobHandler()).Methods("GET")

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
