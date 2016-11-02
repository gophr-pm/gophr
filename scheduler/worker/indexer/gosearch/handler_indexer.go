package gosearch

import (
	"net/http"

	"github.com/gophr-pm/gophr/lib/db"
)

// IndexHandler exposes an endpoint that indexes all of the go packages known
// to http://go-search.org/.
func IndexHandler(q db.Queryable) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(shikkic): copy in code here.
	}
}
