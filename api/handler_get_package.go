package main

import (
	"net/http"

	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/model/package"
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
	q query.Queryable,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			args   getPackageRequestArgs
			json   []byte
			result pkg.Details
		)

		// Parse out the args.
		if args, err = extractGetPackageRequestArgs(r); err != nil {
			errors.RespondWithError(w, err)
			return
		}
		// Get from the database.
		if result, err = pkg.Get(q, args.author, args.repo); err != nil {
			errors.RespondWithError(w, err)
			return
		}
		// Turn the result into JSON.
		if json, err = result.ToJSON(); err != nil {
			errors.RespondWithError(w, err)
			return
		}

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
