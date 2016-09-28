package main

import (
	"net/http"

	"github.com/skeswa/gophr/common/config"
)

// DeleteRepoHandler creates a new repository in the depot.
func DeleteRepoHandler(conf *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get request metadata.
		vars, err := readURLVars(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = destroyRepo(
			conf.DepotPath,
			vars.author,
			vars.repo,
			vars.sha)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Otherwise, the repo was deleted successfully.
		w.WriteHeader(http.StatusOK)
		return
	}
}
