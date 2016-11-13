package main

import (
	"net/http"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
	"github.com/gophr-pm/gophr/lib/newrelic"
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
	io io.IO,
	conf *config.Config,
	client db.Client,
	creds *config.Credentials,
	ghSvc github.RequestService,
	newRelicApp newrelic.Application,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log this transaction in new relic if in production.
		var nrTxn newrelic.Transaction
		if !conf.IsDev {
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
			req:           r,
			downloadRefs:  lib.FetchRefs,
			fetchFullSHA:  github.FetchFullSHAFromPartialSHA,
			DoHTTPHeadReq: github.DoHTTPHeadReq,
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
			db:                    client,
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
