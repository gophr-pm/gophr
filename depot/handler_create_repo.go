package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
)

// CreateRepoHandler creates a new repository in the depot.
func CreateRepoHandler(
	conf *config.Config,
	dataDogClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingArgs := datadog.TrackTransactionArgs{
			Tags: []string{
				"create-repo",
				"internal",
			},
			Client:          dataDogClient,
			StartTime:       time.Now(),
			EventInfo:       []string{},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: "create.repo",
		}

		defer datadog.TrackTransaction(trackingArgs)

		// Get request metadata.
		vars, err := readURLVars(r)
		trackingArgs.EventInfo = append(
			trackingArgs.EventInfo,
			fmt.Sprintf("%v", vars),
		)
		if err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		alreadyExisted, err := createNewRepo(
			conf.DepotPath,
			vars.author,
			vars.repo,
			vars.sha)
		if err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// If the repo already existed, then return a 304.
		if alreadyExisted {
			trackingArgs.AlertType = datadog.Info
			w.WriteHeader(http.StatusNotModified)
			return
		}

		// Otherwise, the repo was created successfully.
		trackingArgs.AlertType = datadog.Success
		w.WriteHeader(http.StatusOK)
		return
	}
}
