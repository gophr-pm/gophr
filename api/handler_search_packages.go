package main

import (
	"net/http"
	"strconv"

	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/model/package"
)

const (
	maxSearchQueryLength   = 50
	maxSearchPackagesLimit = 20
)

type searchPackagesRequestArgs struct {
	limit       int
	searchQuery string
}

// SearchPackagesHandler creates an HTTP request handler that responds to top
// packages get requests.
func SearchPackagesHandler(
	q query.Queryable,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			args    searchPackagesRequestArgs
			json    []byte
			results pkg.Summaries
		)

		// Parse out the args.
		if args, err = extractSearchPackagesRequestArgs(r); err != nil {
			errors.RespondWithError(w, err)
			return
		}
		// Get from the database.
		if results, err = pkg.Search(q, args.searchQuery, args.limit); err != nil {
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

	if args.limit, err = strconv.Atoi(limitStr); err != nil {
		return args, NewInvalidURLParameterError(urlVarLimit, limitStr)
	} else if args.limit > maxSearchPackagesLimit {
		args.limit = maxSearchPackagesLimit
	}

	if len(searchQuery) < 1 || len(searchQuery) > maxSearchQueryLength {
		return args, NewInvalidURLParameterError(urlVarSearchQuery, searchQuery)
	}

	args.searchQuery = searchQuery

	return args, nil
}
