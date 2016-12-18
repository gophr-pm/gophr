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
	"github.com/gorilla/mux"
)

const (
	dailyTimeSplit      = "daily"
	weeklyTimeSplit     = "weekly"
	monthlyTimeSplit    = "monthly"
	allTimeTimeSplit    = "alltime"
	maxTopPackagesLimit = 200
)

// ddEventName is the name of the custom datadog event for this handler.
const ddEventGetTopPackages = "api.get-top-packages"

// getTopPackagesRequestArgs is the args struct for get top packages requests.
type getTopPackagesRequestArgs struct {
	limit     int
	timeSplit pkg.TimeSplit
}

// String serializes the arguments of the get top packages handler into a
// representative string.
func (args getTopPackagesRequestArgs) String() string {
	return fmt.Sprintf(
		`{ limit: %d, timeSplit: "%v" }`,
		args.limit,
		args.timeSplit)
}

// GetTopPackagesHandler creates an HTTP request handler that responds to top
// packages get requests.
func GetTopPackagesHandler(
	q db.Client,
	dataDogClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err          error
			args         getTopPackagesRequestArgs
			json         []byte
			results      pkg.Summaries
			trackingArgs = datadog.TrackTransactionArgs{
				Tags:            []string{apiDDTag, datadog.TagExternal},
				Client:          dataDogClient,
				AlertType:       datadog.Success,
				StartTime:       time.Now(),
				MetricName:      datadog.MetricRequestDuration,
				CreateEvent:     statsd.NewEvent,
				CustomEventName: ddEventGetTopPackages,
			}
		)

		// Track the request with DataDog.
		defer datadog.TrackTransaction(&trackingArgs)

		// Parse out the args.
		if args, err = extractGetTopPackagesRequestArgs(r); err != nil {
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
		if results, err = pkg.GetTopX(q, args.limit, args.timeSplit); err != nil {
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

// extractGetTopPackagesRequestArgs validates and extracts the necessary
// parameters for a get top packages request.
func extractGetTopPackagesRequestArgs(
	r *http.Request,
) (getTopPackagesRequestArgs, error) {
	var (
		err  error
		vars = mux.Vars(r)
		args getTopPackagesRequestArgs
	)

	if args.limit, err = strconv.Atoi(vars[urlVarLimit]); err != nil {
		return args, NewInvalidURLParameterError(urlVarLimit, vars[urlVarAuthor])
	} else if args.limit > maxTopPackagesLimit {
		args.limit = maxTopPackagesLimit
	}

	switch vars[urlVarTimeSplit] {
	case dailyTimeSplit:
		args.timeSplit = pkg.Daily
	case weeklyTimeSplit:
		args.timeSplit = pkg.Weekly
	case monthlyTimeSplit:
		args.timeSplit = pkg.Monthly
	case allTimeTimeSplit:
		args.timeSplit = pkg.AllTime
	default:
		return args, NewInvalidURLParameterError(
			urlVarTimeSplit,
			vars[urlVarTimeSplit])
	}

	return args, nil
}
