package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/common"
	"github.com/gophr-pm/gophr/common/config"
	"github.com/gophr-pm/gophr/common/newrelic"
)

func main() {
	// Get the config, sesh and creds.
	conf, db := common.Init()
	creds, err := config.ReadCredentials(conf)
	if err != nil {
		log.Fatalln("Failed to read credentials secret:", err)
	}

	// Create new relic app for monitoring.
	newRelicApp, err := nr.CreateNewRelicApp(conf)
	if err != nil {
		log.Fatalln(err)
	}

	// Ensure that the session is closed eventually.
	defer db.Close()

	// Start serving.
	http.HandleFunc(wildcardHandlerPattern, RequestHandler(conf, db, creds, newRelicApp))
	log.Printf("Servicing HTTP requests on port %d.\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}
