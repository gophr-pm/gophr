package main

import (
	"net/http"
	"strconv"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/db/model/package"
)

const (
	maxNewPackagesLimit = 200
)

// getNewPackagesRequestArgs is the args struct for new packages requests.
type getNewPackagesRequestArgs struct {
	limit int
}

// GetNewPackagesHandler creates an HTTP request handler that responds to top
// packages get requests.
func GetNewPackagesHandler(
	q db.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			args    getNewPackagesRequestArgs
			json    []byte
			results pkg.Summaries
		)

		// Parse out the args.
		if args, err = extractGetNewPackagesRequestArgs(r); err != nil {
			errors.RespondWithError(w, err)
			return
		}
		// Get from the database.
		if results, err = pkg.GetNew(q, args.limit); err != nil {
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
