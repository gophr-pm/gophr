package main

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron"
)

func main() {
	// Instantiate new cron instance.
	c := cron.New()

	// List cron jobs.
	c.AddFunc("* * * * * *", Github_indexer)

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
