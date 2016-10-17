package main

import (
	"fmt"
	"net/http"

	"github.com/gophr-pm/gophr/scheduler/indexer/github"
	"github.com/robfig/cron"
)

/*

	1). Github Indexer
	2). Downloads Indexer
*/

func main() {
	// Instantiate new cron instance.
	c := cron.New()

	// List cron jobs.
	c.AddFunc("* * * * * *", github.Github_indexer)

	// Start the cron process.
	c.Start()

	// Create handlers for generating metrics.
	http.HandleFunc("/daily-downloads", dailyDownloadsHandler)

	// Start HTTP server.
	http.ListenAndServe(":8080", nil)
}

// dailyDownloadsHandler responsible running a metric job via a handlr.
func dailyDownloadsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("placeholder")
}
