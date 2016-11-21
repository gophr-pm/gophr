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
	maxSearchQueryLength = 50
	// ddEventName is the name of the custom datadog event for this handler.
	ddEventSearchPackages  = "api.search.packages"
	maxSearchPackagesLimit = 20
)

// searchPackagesRequestArgs is the args struct for package search requests.
type searchPackagesRequestArgs struct {
	limit       int
	searchQuery string
}

// String serializes the arguments of the search packages handler into a
// representative string.
func (args searchPackagesRequestArgs) String() string {
	return fmt.Sprintf(
		`{ limit: %d, searchQuery: "%s" }`,
		args.limit,
		args.searchQuery)
}

// SearchPackagesHandler creates an HTTP request handler that responds to top
// packages get requests.
func SearchPackagesHandler(
	q db.Client,
	dataDogClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err          error
			args         searchPackagesRequestArgs
			json         []byte
			results      pkg.Summaries
			trackingArgs = datadog.TrackTransactionArgs{
				Tags:            []string{apiDDTag, datadog.TagExternal},
				Client:          dataDogClient,
				AlertType:       datadog.Success,
				StartTime:       time.Now(),
				MetricName:      datadog.MetricRequestDuration,
				CreateEvent:     statsd.NewEvent,
				CustomEventName: ddEventSearchPackages,
			}
		)

		// Track the request with DataDog.
		defer datadog.TrackTransaction(trackingArgs)

		// Parse out the args.
		if args, err = extractSearchPackagesRequestArgs(r); err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(
				trackingArgs.EventInfo,
				args.String(),
				err.Error())
			errors.RespondWithError(w, err)
			return
		}

		// Track request metadata.
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, args.String())

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
