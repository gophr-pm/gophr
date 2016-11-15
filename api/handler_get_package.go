package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gorilla/mux"
)

// getPackageRequestArgs is the args struct for get package requests.
type getPackageRequestArgs struct {
	repo   string
	author string
}

// GetPackageHandler creates an HTTP request handler that responds to individual
// package get requests.
func GetPackageHandler(
	q db.Client,
	dataDogClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			args   getPackageRequestArgs
			json   []byte
			result pkg.Details
		)

		trackingArgs := datadog.TrackTransactionArgs{
			Tags: []string{
				"api-get-package",
				"external",
			},
			Client:          dataDogClient,
			StartTime:       time.Now(),
			EventInfo:       []string{},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: "api.get.package",
		}

		defer datadog.TrackTransaction(trackingArgs)

		// Parse out the args.
		if args, err = extractGetPackageRequestArgs(r); err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			errors.RespondWithError(w, err)
			return
		}

		// Track request metadata.
		trackingArgs.EventInfo = append(
			trackingArgs.EventInfo,
			fmt.Sprintf("%v", args),
		)

		// Get from the database.
		if result, err = pkg.Get(q, args.author, args.repo); err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			errors.RespondWithError(w, err)
			return
		}
		// Turn the result into JSON.
		if json, err = result.ToJSON(); err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			errors.RespondWithError(w, err)
			return
		}

		trackingArgs.AlertType = datadog.Success
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, string(json))
		respondWithJSON(w, json)
	}
}

// extractGetPackageArgs validates and extracts the necessary parameters for a
// get package request.
func extractGetPackageRequestArgs(
	r *http.Request,
) (getPackageRequestArgs, error) {
	var (
		vars = mux.Vars(r)
		args getPackageRequestArgs
	)

	if args.author = vars[urlVarAuthor]; len(args.author) < 1 {
		return args, NewInvalidURLParameterError(urlVarAuthor, args.author)
	}
	if args.repo = vars[urlVarRepo]; len(args.repo) < 1 {
		return args, NewInvalidURLParameterError(urlVarRepo, args.repo)
	}

	return args, nil
}
