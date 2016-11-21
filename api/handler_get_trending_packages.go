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
	maxTrendingPackagesLimit = 200
	// ddEventName is the name of the custom datadog event for this handler.
	ddEventGetTrendingPackages = "api.get.trending.packages"
)

// getTrendingPackagesRequestArgs is the args struct for get trending packages
// requests.
type getTrendingPackagesRequestArgs struct {
	limit int
}

// String serializes the arguments of the get trending packages handler into a
// representative string.
func (args getTrendingPackagesRequestArgs) String() string {
	return fmt.Sprintf(`{ limit: %d }`, args.limit)
}

// GetTrendingPackagesHandler creates an HTTP request handler that responds to
// top packages get requests.
func GetTrendingPackagesHandler(
	q db.Client,
	dataDogClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err          error
			args         getTrendingPackagesRequestArgs
			json         []byte
			results      pkg.Summaries
			trackingArgs = datadog.TrackTransactionArgs{
				Tags:            []string{apiDDTag, datadog.TagExternal},
				Client:          dataDogClient,
				AlertType:       datadog.Success,
				StartTime:       time.Now(),
				EventInfo:       []string{},
				MetricName:      datadog.MetricRequestDuration,
				CreateEvent:     statsd.NewEvent,
				CustomEventName: ddEventGetTrendingPackages,
			}
		)

		// Track the request with DataDog.
		defer datadog.TrackTransaction(trackingArgs)

		// Parse out the args.
		if args, err = extractGetTrendingPackagesRequestArgs(r); err != nil {
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
		if results, err = pkg.GetTrending(q, args.limit); err != nil {
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

// extractGetTrendingPackagesRequestArgs validates and extracts the necessary
// parameters for a get trending packages request.
func extractGetTrendingPackagesRequestArgs(
	r *http.Request,
) (getTrendingPackagesRequestArgs, error) {
	var (
		err      error
		args     getTrendingPackagesRequestArgs
		limitStr = r.URL.Query().Get(urlVarLimit)
	)

	if len(limitStr) == 0 {
		args.limit = maxTrendingPackagesLimit
		return args, nil
	}

	if args.limit, err = strconv.Atoi(limitStr); err != nil {
		return args, NewInvalidQueryStringParameterError(urlVarLimit, limitStr)
	}
	if args.limit > maxTrendingPackagesLimit {
		args.limit = maxTrendingPackagesLimit
	}

	return args, nil
}
