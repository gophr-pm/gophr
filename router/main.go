package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/config"
)

func main() {
	// Get the config, sesh and creds.
	conf, db := common.Init()
	creds, err := config.ReadCredentials(conf)
	if err != nil {
		log.Fatalln("Failed to read credentials secret:", err)
	}

	// Ensure that the session is closed eventually.
	defer db.Close()

	// Start serving.
	http.HandleFunc(wildcardHandlerPattern, RequestHandler(conf, db, creds))
	log.Printf("Servicing HTTP requests on port %d.\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}
