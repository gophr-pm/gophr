package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/db"
)

func main() {
	conf := config.GetConfig()
	log.Println("Configuration:\n\n" + conf.String() + "\n")
	session, err := db.OpenConnection(conf)

	// Exit if we can't connect to the database.
	if err != nil {
		log.Fatalln("Could not start the API:", err)
	}

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	r.HandleFunc("/search", SearchHandler(session)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/{%s}/{%s}/versions", urlVarAuthor, urlVarRepo), VersionsHandler(session)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/{%s}/{%s}/versions/latest", urlVarAuthor, urlVarRepo), LatestVersionHandler(session)).Methods("GET")

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), r)
}
