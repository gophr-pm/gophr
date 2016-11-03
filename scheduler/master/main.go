package main

import (
	"log"
	"net/http"

	"github.com/robfig/cron"
)

func main() {
	// Instantiate new cron instance.
	c := cron.New()

	// TODO(shikkic): take a look - these cron times are all wrong.
	c.AddFunc("0 0 * * * *", newJobRunner(indexAwesomePackages, http.Get))
	c.AddFunc("0 0 * * * *", newJobRunner(indexGoSearchPackages, http.Get))
	c.AddFunc("0 0 * * * *", newJobRunner(updatePackageMetrics, http.Get))

	// Start the cron process.
	log.Println("Scheduler now waiting for scheduled jobs.")
	c.Start()
}
