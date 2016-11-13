package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gorilla/mux"
)

// DepotRequestHandlerArgs lol
type DepotRequestHandlerArgs struct {
	config        *config.Config
	dataDogClient *statsd.Client
}

func main() {
	var (
		r        = mux.NewRouter()
		conf     = config.GetConfig()
		endpoint = fmt.Sprintf(
			"/repos/{%s}/{%s}/{%s}",
			urlVarAuthor,
			urlVarRepo,
			urlVarSHA)
	)

	// Initialize datadog client.
	dataDogClient, err := datadog.NewClient(conf, "depot.")
	if err != nil {
		log.Println(err)
	}

	// Register the status route.
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	// Register all the remaining routes for the main endpoint.
	r.HandleFunc(endpoint, RepoExistsHandler(conf, dataDogClient)).Methods("GET")
	r.HandleFunc(endpoint, CreateRepoHandler(conf, dataDogClient)).Methods("POST")
	r.HandleFunc(endpoint, DeleteRepoHandler(conf, dataDogClient)).Methods("DELETE")

	// Start tailing the nginx logs.
	if err := tailNginxLogs(); err != nil {
		log.Fatalln("Failed to start a tail on the nginx logs:", err)
	}

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), r)
}
