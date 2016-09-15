package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/skeswa/gophr/common/depot"
	"github.com/skeswa/gophr/common/errors"
)

const (
	queryStringRepoTextKey = "r"
	queryStringRefTextKey  = "hb"
)

// ReadmeHandler creates an HTTP request handler that responds to fuzzy package
// searches.
func ReadmeHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			repoName string
			ref      string
		)

		// Get URL Query Params
		qs := r.URL.Query()
		if qs != nil {
			// Get repo name
			repoName = qs.Get(queryStringRepoTextKey)
			ref = qs.Get(queryStringRefTextKey)
		}

		// Check if repo name exists
		if len(repoName) < 1 {
			errors.RespondWithError(w, NewInvalidQueryStringParameterError(
				queryStringRepoTextKey,
				repoName,
			))
			return
		}

		// Check if ref exists
		if len(ref) < 1 {
			errors.RespondWithError(w, NewInvalidQueryStringParameterError(
				queryStringRefTextKey,
				ref,
			))
			return
		}

		// Request the README from depot gitweb
		url := fmt.Sprintf("http://%s/?p=%s;a=blob_plain;f=README.md;hb=%s", depot.DepotInternalServiceAddress, repoName, ref)
		resp, err := http.Get(url)
		if err != nil {
			errors.RespondWithError(w, err)
			return
		}

		// If no README was found return 404
		if resp.StatusCode == 404 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte{})
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errors.RespondWithError(w, err)
			return
		}

		if len(body) > 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
		}
	}
}
