package main

import (
	"log"

	"github.com/gophr-pm/gophr/lib"
)

func main() {
	conf, client := lib.Init()

	// Immediately discard the client since it won't be used.
	client.Close()

	// Execute the migrations.
	log.Println("Executing pending database migrations.")
	err := upSync(conf.IsDev, conf.DbAddress, conf.MigrationsPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Migrations executed successfully.")
}
