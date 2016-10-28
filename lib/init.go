package lib

import (
	"log"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// Init initializes a gophr backend component. Returns all the pre-requisites to
// starting up successfully.
func Init() (*config.Config, db.Client) {
	// Get config then session to be returned later.
	conf := config.GetConfig()
	log.Println("Configuration:\n\n" + conf.String() + "\n")

	c, err := db.NewClient(conf)
	// Exit if anything goes wrong.
	if err != nil {
		log.Fatalln("Initialization failed:", err)
	}

	// Decide the replication factor according to the environment.
	var replicationFactor int
	if conf.IsDev {
		// Replication factor is one because there are only two node total.
		replicationFactor = 1
	} else {
		// Replication factor is two since there are at least three nodes.
		replicationFactor = 2
	}

	// Make sure that the keyspace exists.
	log.Println("Asserting the existence of the keyspace.")
	err = query.CreateKeyspaceIfNotExists().
		WithReplication("SimpleStrategy", replicationFactor).
		WithDurableWrites(true).
		Create(c).
		Exec()
	if err != nil {
		log.Fatalln("Initialization failed:", err)
	}

	return conf, c
}
