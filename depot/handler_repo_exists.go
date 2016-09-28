package main

import (
	"net/http"

	"github.com/skeswa/gophr/common/config"
)

// RepoExistsHandler returns a 200 if the repo exists, or 404 if it doesn't.
func RepoExistsHandler(conf *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get request metadata.
		vars, err := readURLVars(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		exists, err := repoExists(
			conf.DepotPath,
			vars.author,
			vars.repo,
			vars.sha)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// If the repo doesn't exist, then return a 404.
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Otherwise, the repo exists.
		w.WriteHeader(http.StatusOK)
		return
	}
}
