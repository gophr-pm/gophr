package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
)

// RepoExistsHandler returns a 200 if the repo exists, or 404 if it doesn't.
func RepoExistsHandler(
	conf *config.Config,
	datadogClient *statsd.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingArgs := datadog.TrackTranscationArgs{
			Tags: []string{
				"repo-exists",
				"internal",
			},
			Client:          datadogClient,
			StartTime:       time.Now(),
			EventInfo:       []string{},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: "repo.exists",
		}

		// Get request metadata.
		vars, err := readURLVars(r)
		// Track request metadata.
		trackingArgs.EventInfo = append(
			trackingArgs.EventInfo,
			fmt.Sprintf("%v", vars),
		)
		if err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			defer datadog.TrackTranscation(trackingArgs)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		exists, err := repoExists(
			conf.DepotPath,
			vars.author,
			vars.repo,
			vars.sha)
		if err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			defer datadog.TrackTranscation(trackingArgs)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// If the repo doesn't exist, then return a 404.
		if !exists {
			trackingArgs.AlertType = datadog.Info
			trackingArgs.Tags = append(trackingArgs.Tags, "404")
			defer datadog.TrackTranscation(trackingArgs)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Otherwise, the repo exists.
		trackingArgs.AlertType = datadog.Success
		defer datadog.TrackTranscation(trackingArgs)
		w.WriteHeader(http.StatusOK)
		return
	}
}
