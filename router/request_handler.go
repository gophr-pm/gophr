package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/common"
	"github.com/gophr-pm/gophr/common/config"
	"github.com/gophr-pm/gophr/common/errors"
	"github.com/gophr-pm/gophr/common/github"
	"github.com/gophr-pm/gophr/common/io"
	"github.com/gophr-pm/gophr/common/newrelic"
	"github.com/newrelic/go-agent"
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
	creds *config.Credentials,
	newRelicApp newrelic.Application,
) func(http.ResponseWriter, *http.Request) {
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
		// Log this transaction in new relic if in production.
		var nrTxn newrelic.Transaction
		if conf.IsDev {
			nrTxn = nr.CreateNewRelicTxn(newRelicApp, &w, r)
			defer nrTxn.End()
		}

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
			doHTTPHead:   github.DoHTTPHeadReq,
		}); err != nil {
			if nrTxn != nil {
				nr.ReportNewRelicError(nrTxn, err)
			}

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
			if nrTxn != nil {
				nr.ReportNewRelicError(nrTxn, err)
			}

			errors.RespondWithError(w, err)
			return
		}
	}
}
