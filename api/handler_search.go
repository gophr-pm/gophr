package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common"
)

const (
	queryStringSearchTextKey = "q"
)

func SearchHandler(session *gocql.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var searchText string

		qs := r.URL.Query()
		if qs != nil {
			searchText = qs.Get(queryStringSearchTextKey)
		}

		if len(searchText) < 1 {
			common.RespondWithError(w, NewInvalidQueryStringParameterError(
				queryStringSearchTextKey,
				searchText,
			))
			return
		}

		packageModels, err := common.FuzzySearchPackages(session, searchText)
		if err != nil {
			common.RespondWithError(w, err)
			return
		}

		marshalledPackageModels, err := ffjson.Marshal(packageModels)
		if err != nil {
			common.RespondWithError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(marshalledPackageModels)
	}
}
