package main

import (
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/robfig/cron"
)

func main() {
	// Read configuration to check what environment this binary is running in.
	conf := config.GetConfig()

	if conf.IsDev {
		log.Println(
			"Scheduler is NOT waiting for scheduled jobs since it is " +
				"running in development mode.")

		// Block forever.
		select {}
	}

	// Instantiate new cron instance.
	c := cron.New()

	// Index awesome packages once a week on Sunday at 5am in the morning.
	c.AddFunc("0 5 * * 7 *", newJobRunner(indexAwesomePackages, http.Get))
	// Index go-search packages on the 1st of every month at 3am in the morning.
	c.AddFunc("0 3 1 * * *", newJobRunner(indexGoSearchPackages, http.Get))
	// Update package metrics everyday at 7am in the morning.
	c.AddFunc("0 7 * * * *", newJobRunner(updatePackageMetrics, http.Get))

	// Start the cron process.
	log.Println("Scheduler now waiting for scheduled jobs.")
	c.Start()
}
