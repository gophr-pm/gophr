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
	// Create other useful closure variables here.
	var (
		rd defaultRefsDownloader
	)

	return func(w http.ResponseWriter, r *http.Request) {
		// Make sure that this isn't a simple health check before getting more
		// complicated.
		if r.URL.Path == healthCheckRoute {
			w.Write(statusCheckResponse)
			return
		}

		// First, create the necessary variables.
		var (
			pr  *packageRequest
			err error
		)

		// Create a new package request.
		if pr, err = newPackageRequest(newPackageRequestArgs{
			req:        r,
			downloader: rd,
		}); err != nil {
			errors.RespondWithError(w, err)
			return
		}

		// Use the package request to respond.
		if err = pr.respond(respondToPackageRequestArgs{
			res:     w,
			conf:    conf,
			creds:   creds,
			session: session,
		}); err != nil {
			errors.RespondWithError(w, err)
			return
		}
	}
}
