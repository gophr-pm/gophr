package main

import (
	"net/http"
	"strconv"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/model/package"
	"github.com/gorilla/mux"
)

const (
	dailyTimeSplit      = "daily"
	weeklyTimeSplit     = "weekly"
	monthlyTimeSplit    = "monthly"
	allTimeTimeSplit    = "alltime"
	maxTopPackagesLimit = 200
)

// getTopPackagesRequestArgs is the args struct for get top packages requests.
type getTopPackagesRequestArgs struct {
	limit     int
	timeSplit pkg.TimeSplit
}

// GetTopPackagesHandler creates an HTTP request handler that responds to top
// packages get requests.
func GetTopPackagesHandler(
	q db.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			args    getTopPackagesRequestArgs
			json    []byte
			results pkg.Summaries
		)

		// Parse out the args.
		if args, err = extractGetTopPackagesRequestArgs(r); err != nil {
			errors.RespondWithError(w, err)
			return
		}
		// Get from the database.
		if results, err = pkg.GetTopX(q, args.limit, args.timeSplit); err != nil {
			errors.RespondWithError(w, err)
			return
		}
		// Turn the result into JSON.
		if json, err = results.ToJSON(); err != nil {
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
