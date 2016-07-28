package main

import (
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/errors"
)

const (
	healthCheckRoute       = "/status"
	wildcardHandlerPattern = "/"
)

var (
	statusCheckResponse = []byte("ok")
)

// RequestHandler creates an HTTP request handler that responds to all incoming
// router requests.
func RequestHandler(
	conf *config.Config,
	session *gocql.Session,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		context := NewRequestContext(nil)

		log.Printf("[%s] New request received: %s\n", context.RequestID, r.URL.Path)

		if r.URL.Path == healthCheckRoute {
			log.Printf(
				"[%s] Handling request for \"%s\" as a health check\n",
				context.RequestID,
				r.URL.Path,
			)

			w.Write(statusCheckResponse)
		} else {
			log.Printf(
				"[%s] Handling request for \"%s\" as a package request\n",
				context.RequestID,
				r.URL.Path,
			)

			err := RespondToPackageRequest(session, context, r, w)
			if err != nil {
				errors.RespondWithError(w, err)
			}
		}
	}
}
