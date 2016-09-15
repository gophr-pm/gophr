package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/skeswa/gophr/common/errors"
)

const (
	queryStringRepoTextKey = "r"
	queryStringRefTextKey  = "hb"
)

var (
	repoName string
	ref      string
)

// ReadmeHandler creates an HTTP request handler that responds to fuzzy package
// searches.
func ReadmeHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

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

		url := fmt.Sprintf("http://depot-svc:3000/?p=%s;a=blob_plain;f=README.md;hb=%s", repoName, ref)
		resp, err := http.Get(url)
		if err != nil {
			errors.RespondWithError(w, err)
			return
		}
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
