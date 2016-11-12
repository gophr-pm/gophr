package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
)

func main() {
	// Get the config, sesh and creds.
	conf, client := lib.Init()

	// Ensure that the client is closed eventually.
	defer client.Close()

	creds, err := config.ReadCredentials(conf)
	if err != nil {
		log.Fatalln("Failed to read credentials secret:", err)
	}

	datadogClient, err := datadog.NewClient(conf, "router.")
	if err != nil {
		log.Println(err)
	}

	// Start serving.
	http.HandleFunc(wildcardHandlerPattern, RequestHandler(
		conf,
		client,
		creds,
		datadogClient))
	log.Printf("Servicing HTTP requests on port %d.\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}
