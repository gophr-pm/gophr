package main

import (
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
	creds *config.Credentials) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == healthCheckRoute {
			w.Write(statusCheckResponse)
		} else {
			// Create a new package request.
			pr, err := newPackageRequest(r)
			if err != nil {
				errors.RespondWithError(w, err)
				return
			}

			// Use the package request to respond.
			if err = pr.respond(w, conf, creds, session); err != nil {
				errors.RespondWithError(w, err)
				return
			}
		}
	}
}
