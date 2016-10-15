package main

import (
	"log"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
)

func main() {
	// Initialize the migrator.
	config := config.GetConfig()

	// Open up a connection to the database in order to assert the existence of
	// the keyspace.
	if conn, err := db.OpenConnection(config); err != nil {
		log.Fatalln(err)
	} else {
		// Close the conneciton since the migrator creates its own.
		conn.Close()
	}

	// Execute the migrations.
	log.Println("Executing pending database migrations.")
	err := db.Migrate(config.IsDev, config.DbAddress, config.MigrationsPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Migrations executed successfully.")
}
