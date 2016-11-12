package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
)

// DeleteRepoHandler creates a new repository in the depot.
func DeleteRepoHandler(
	conf *config.Config,
	datadogClient *statsd.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingArgs := datadog.TrackTranscationArgs{
			Tags: []string{
				"repo-delete",
				"internal",
			},
			Client:          datadogClient,
			StartTime:       time.Now(),
			EventInfo:       []string{},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: "repo.delete",
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

		err = destroyRepo(
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

		// Otherwise, the repo was deleted successfully.
		trackingArgs.AlertType = datadog.Success
		defer datadog.TrackTranscation(trackingArgs)
		w.WriteHeader(http.StatusOK)
		return
	}
}
