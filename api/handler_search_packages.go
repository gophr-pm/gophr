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
	maxSearchQueryLength   = 50
	ddEventSearchPackages  = "api.search.packages"
	maxSearchPackagesLimit = 20
)

// searchPackagesRequestArgs is the args struct for package search requests.
type searchPackagesRequestArgs struct {
	limit       int
	searchQuery string
}

// SearchPackagesHandler creates an HTTP request handler that responds to top
// packages get requests.
func SearchPackagesHandler(
	q db.Client,
	dataDogClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			args    searchPackagesRequestArgs
			json    []byte
			results pkg.Summaries
		)

		trackingArgs := datadog.TrackTransactionArgs{
			Tags: []string{
				"api-get-trending-packages",
				"external",
			},
			Client:          dataDogClient,
			StartTime:       time.Now(),
			EventInfo:       []string{},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: ddEventSearchPackages,
		}

		defer datadog.TrackTransaction(trackingArgs)

		// Parse out the args.
		if args, err = extractSearchPackagesRequestArgs(r); err != nil {
			errors.RespondWithError(w, err)
			return
		}

		// Track request metadata.
		trackingArgs.EventInfo = append(
			trackingArgs.EventInfo,
			fmt.Sprintf("%v", args),
		)

		// Get from the database.
		if results, err = pkg.Search(q, args.searchQuery, args.limit); err != nil {
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

// extractSearchPackagesRequestArgs validates and extracts the necessary
// parameters for a search packages request.
func extractSearchPackagesRequestArgs(
	r *http.Request,
) (searchPackagesRequestArgs, error) {
	var (
		err         error
		args        searchPackagesRequestArgs
		limitStr    = r.URL.Query().Get(urlVarLimit)
		searchQuery = r.URL.Query().Get(urlVarSearchQuery)
	)

	if len(limitStr) == 0 {
		args.limit = maxTrendingPackagesLimit
	} else if args.limit, err = strconv.Atoi(limitStr); err != nil {
		return args, NewInvalidQueryStringParameterError(urlVarLimit, limitStr)
	}
	if args.limit > maxSearchPackagesLimit {
		args.limit = maxSearchPackagesLimit
	}

	if len(searchQuery) < 1 || len(searchQuery) > maxSearchQueryLength {
		return args, NewInvalidQueryStringParameterError(
			urlVarSearchQuery,
			searchQuery)
	}

	args.searchQuery = searchQuery

	return args, nil
}
