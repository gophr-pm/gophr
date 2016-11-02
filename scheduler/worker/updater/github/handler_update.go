package github

import (
	"net/http"

	"github.com/gophr-pm/gophr/lib/db"
)

// UpdateHandler exposes an endpoint that reads every package from the database
// and updates the github metadata of each.
func UpdateHandler(q db.Queryable) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(shikkic): copy in code here.
	}
}
