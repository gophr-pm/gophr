package awesome

import (
	"net/http"

	"github.com/gophr-pm/gophr/lib/db"
)

// IndexHandler exposes an endpoint that indexes all of the awesome go packages.
func IndexHandler(q db.Queryable) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(shikkic): copy in code here.
	}
}
