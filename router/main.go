package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/skeswa/gophr/common"
)

func main() {
	// Initialize the router.
	config, session := common.Init()

	// Ensure that the session is closed eventually.
	defer session.Close()

	// Start serving.
	http.HandleFunc(wildcardHandlerPattern, RequestHandler(config, session))
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
