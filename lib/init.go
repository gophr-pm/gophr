package common

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
)

// Init initializes a gophr backend component. Returns all the pre-requisites to
// starting up successfully.
func Init() (*config.Config, *gocql.Session) {
	// Get config then session to be returned later.
	conf := config.GetConfig()
	log.Println("Configuration:\n\n" + conf.String() + "\n")
	session, err := db.OpenConnection(conf)

	// Exit if anything goes wrong.
	if err != nil {
		log.Fatalln("Initialization failed:", err)
	}

	return conf, session
}
