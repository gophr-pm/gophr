package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/common"
	"github.com/gophr-pm/gophr/common/config"
	"github.com/gophr-pm/gophr/common/newrelic"
	newrelic "github.com/newrelic/go-agent"
)

func main() {
	// Get the config, sesh and creds.
	conf, db := common.Init()
	creds, err := config.ReadCredentials(conf)
	if err != nil {
		log.Fatalln("Failed to read credentials secret:", err)
	}

	var app newrelic.Application
	if !conf.IsDev {
		newRelicKey, err := nr.GenerateKey(conf)
		if err != nil {
			log.Fatalln("Failed to read newrelic credentials secret:", err)
		}
		config := newrelic.NewConfig("Gophr", newRelicKey)
		app, err = newrelic.NewApplication(config)
		if err != nil {
			log.Fatalln("Failed to create new relic monitoring application:", err)
		}
	}

	// Ensure that the session is closed eventually.
	defer db.Close()

	// Start serving.
	http.HandleFunc(wildcardHandlerPattern, RequestHandler(conf, db, creds, app))
	log.Printf("Servicing HTTP requests on port %d.\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}
