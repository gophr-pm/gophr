package main

import (
	"fmt"
	"net/http"

	"golang.org/x/tools/godoc"

	"github.com/gophr-pm/gophr/scheduler/indexer/github"
	"github.com/robfig/cron"
)

func main() {
	// Instantiate new cron instance.
	c := cron.New()

	// List cron jobs.
	c.AddFunc("0 0 * * * *", awesome.Index)
	c.AddFunc("0 0 * * * *", godoc.Index)
	c.AddFunc("0 0 0 * * *", github.Index)

	// Start the cron process.
	c.Start()

	// Create handlers for generating metrics.
	http.HandleFunc("/daily-downloads", dailyDownloadsHandler)

	// Start HTTP server.
	http.ListenAndServe(":8080", nil)
}

// dailyDownloadsHandler responsible running a metric job via a handler.
func dailyDownloadsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("placeholder")
}
