package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/common/dtos"
	"github.com/gophr-pm/gophr/common/errors"
	"github.com/gophr-pm/gophr/common/models"
)

const (
	queryStringSearchTextKey = "q"
)

// SearchHandler creates an HTTP request handler that responds to fuzzy package
// searches.
func SearchHandler(session *gocql.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var searchText string

		qs := r.URL.Query()
		if qs != nil {
			searchText = qs.Get(queryStringSearchTextKey)
		}

		if len(searchText) < 1 {
			errors.RespondWithError(w, NewInvalidQueryStringParameterError(
				queryStringSearchTextKey,
				searchText,
			))
			return
		}

		packageModels, err := models.FuzzySearchPackages(session, searchText)
		if err != nil {
			errors.RespondWithError(w, err)
			return
		}

		packageListDTO := dtos.NewPackageListDTOFromPackageModelList(packageModels)
		if len(packageListDTO) > 0 {
			json, err := packageListDTO.MarshalJSON()
			if err != nil {
				errors.RespondWithError(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(json)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(dtos.EmptyListJSON))
		}
	}
}
