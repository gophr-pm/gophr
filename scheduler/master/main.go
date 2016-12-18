package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

func main() {
	// Read configuration to check what environment this binary is running in.
	conf := config.GetConfig()

	if conf.IsDev {
		log.Println(
			"Scheduler is NOT waiting for scheduled jobs since it is " +
				"running in development mode.")
	} else {
		// Instantiate new cron instance.
		c := cron.New()

		// Cron string format:
		// "seconds minutes hours dayOfMonth month dayOfWeek"

		// Index awesome packages once a week on Sunday at 5am in the morning.
		c.AddFunc("0 0 5 * * SUN", newJobRunner(indexAwesomePackages, http.Get))
		// Index go-search packages on the 1st of every month at 3am in the morning.
		c.AddFunc("0 0 3 1 * *", newJobRunner(indexGoSearchPackages, http.Get))
		// Delete old hourly downloads everyday at midnight.
		c.AddFunc("0 0 0 * * *", newJobRunner(deleteOldDownloadsPackages, http.Get))
		// Update Github metadata once a day at 5am.
		c.AddFunc("0 0 5 * * *", newJobRunner(updateGithubMetadata, http.Get))
		// Update package metrics three times a day.
		c.AddFunc(
			"0 0 7,14,21 * * *",
			newJobRunner(updatePackageMetrics, http.Get))

		// Start the cron process.
		log.Println("Scheduler now waiting for scheduled jobs.")
		c.Start()
	}

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	log.Printf("Servicing HTTP status requests on port %d.\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), r)
}
