package main

import (
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
)

const (
	healthCheckRoute       = "/status"
	wildcardHandlerPattern = "/"
)

var (
	statusCheckResponse    = []byte("ok")
	ddEventPackageDownload = "router.package.download"
)

// RequestHandler creates an HTTP request handler that responds to all incoming
// router requests.
func RequestHandler(
	io io.IO,
	conf *config.Config,
	creds *config.Credentials,
	ghSvc github.RequestService,
	client db.Client,
	dataDogClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingArgs := datadog.TrackTransactionArgs{
			Tags: []string{
				"package-download",
				"external",
			},
			Client:    dataDogClient,
			StartTime: time.Now(),
			AlertType: datadog.Success,
			EventInfo: []string{
				r.URL.Path, r.UserAgent(),
			},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: ddEventPackageDownload,
		}

		defer datadog.TrackTransaction(&trackingArgs)

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
			ghSvc:         ghSvc,
			downloadRefs:  lib.FetchRefs,
			DoHTTPHeadReq: github.DoHTTPHeadReq,
		}); err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
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
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			errors.RespondWithError(w, err)
			return
		}
	}
}
