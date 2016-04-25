package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/skeswa/gophr/common"
)

func main() {
	var (
		isDev     = common.ReadEnvIsDev()
		httpPort  = common.ReadEnvHTTPPort()
		dbAddress = common.ReadEnvDatabaseAddress()
	)

	dbSession, err := common.OpenDBConnection(dbAddress, isDev)
	if err != nil {
		log.Fatalf(
			"Failed to open a connection with the database (%s): %v\n",
			dbAddress,
			err,
		)
	} else {
		// Ensure that the session ends when the connection exits.
		defer common.CloseDBConnection(dbSession)
	}

	router := mux.NewRouter()
	router.HandleFunc("/status", StatusHandler()).Methods("GET")
	router.HandleFunc("/packages/installs/record", RecordInstallHandler(dbSession)).Methods("POST")

	log.Printf("Server is listening on port %d for HTTP requests...\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%d", httpPort), router)
}
