package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/errors"
	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/io"
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
	// Instantiate the IO module for use in package downloading and versioning.
	io := io.NewIO()

	// Instantiate the the github request service to pass into new
	// package requests.
	ghSvc := github.NewRequestService(github.RequestServiceArgs{
		Conf:       conf,
		Session:    session,
		ForIndexer: false,
	})

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
			req:          r,
			downloadRefs: common.FetchRefs,
			fetchFullSHA: github.FetchFullSHAFromPartialSHA,
		}); err != nil {
			errors.RespondWithError(w, err)
			return
		}

		// Use the package request to respond.
		if err = pr.respond(respondToPackageRequestArgs{
			io:                    io,
			db:                    session,
			res:                   w,
			conf:                  conf,
			creds:                 creds,
			ghSvc:                 ghSvc,
			versionPackage:        versionAndArchivePackage,
			isPackageArchived:     isPackageArchived,
			recordPackageDownload: recordPackageDownload,
			recordPackageArchival: recordPackageArchival,
		}); err != nil {
			errors.RespondWithError(w, err)
			return
		}
	}
}
