package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
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

	dataDogClient, err := datadog.NewClient(conf, "router.")
	if err != nil {
		log.Println(err)
	}

	// Instantiate the the github request service to pass into new
	// package requests.
	ghSvc, err := github.NewRequestService(github.RequestServiceArgs{
		Conf:             conf,
		Queryable:        client,
		ForScheduledJobs: false,
	})
	if err != nil {
		log.Fatalln("Failed to create Github API request service:", err)
	}

	// Instantiate the IO module for use in package downloading and versioning.
	io := io.NewIO()

	// Start serving.
	http.HandleFunc(wildcardHandlerPattern, RequestHandler(
		io,
		conf,
		creds,
		ghSvc,
		client,
		dataDogClient))
	log.Printf("Servicing HTTP requests on port %d.\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}
