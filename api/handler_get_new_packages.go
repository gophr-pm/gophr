package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/errors"
)

const (
	ddEventGetNewPackage = "api.get.new.packages"
	maxNewPackagesLimit  = 200
)

// getNewPackagesRequestArgs is the args struct for new packages requests.
type getNewPackagesRequestArgs struct {
	limit int
}

// GetNewPackagesHandler creates an HTTP request handler that responds to top
// packages get requests.
func GetNewPackagesHandler(
	q db.Client,
	dataDogClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			args    getNewPackagesRequestArgs
			json    []byte
			results pkg.Summaries
		)

		trackingArgs := datadog.TrackTransactionArgs{
			Tags: []string{
				"api-get-new-packages",
				"external",
			},
			Client:          dataDogClient,
			StartTime:       time.Now(),
			EventInfo:       []string{},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: ddEventGetNewPackage,
		}

		defer datadog.TrackTransaction(trackingArgs)

		// Parse out the args.
		if args, err = extractGetNewPackagesRequestArgs(r); err != nil {
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
		if results, err = pkg.GetNew(q, args.limit); err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			errors.RespondWithError(w, err)
			return
		}
		// Turn the result into JSON.
		if json, err = results.ToJSON(); err != nil {
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

// extractGetNewPackagesRequestArgs validates and extracts the necessary
// parameters for a get new packages request.
func extractGetNewPackagesRequestArgs(
	r *http.Request,
) (getNewPackagesRequestArgs, error) {
	var (
		err      error
		args     getNewPackagesRequestArgs
		limitStr = r.URL.Query().Get(urlVarLimit)
	)

	if len(limitStr) == 0 {
		args.limit = maxNewPackagesLimit
		return args, nil
	}

	if args.limit, err = strconv.Atoi(limitStr); err != nil {
		return args, NewInvalidQueryStringParameterError(urlVarLimit, limitStr)
	}
	if args.limit > maxNewPackagesLimit {
		args.limit = maxNewPackagesLimit
	}

	return args, nil
}
