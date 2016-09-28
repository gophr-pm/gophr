package main

import (
	"net/http"

	"github.com/skeswa/gophr/common/config"
)

// CreateRepoHandler creates a new repository in the depot.
func CreateRepoHandler(conf *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get request metadata.
		vars, err := readURLVars(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		alreadyExisted, err := createNewRepo(
			conf.DepotPath,
			vars.author,
			vars.repo,
			vars.sha)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// If the repo already existed, then return a 304.
		if alreadyExisted {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		// Otherwise, the repo was created successfully.
		w.WriteHeader(http.StatusOK)
		return
	}
}
