package main

import (
	"log"

	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/db"
)

func main() {
	// Initialize the migrator.
	config := config.GetConfig()

	// Execute the migrations.
	log.Println("Executing pending database migrations.")
	err := db.Migrate(config.DbAddress, config.MigrationsPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Migrations executed successfully.")
}
