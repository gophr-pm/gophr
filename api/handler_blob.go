package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/depot"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gorilla/mux"
)

type blobRequestArgs struct {
	author string
	repo   string
	sha    string
	path   string
}

// BlobHandler creates an HTTP request handler that responds to filepath lookups.
func BlobHandler(datadogClient *statsd.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingArgs := datadog.TrackTranscationArgs{
			Tags: []string{
				"repo-blob",
				"external",
			},
			Client:          datadogClient,
			StartTime:       time.Now(),
			EventInfo:       []string{},
			MetricName:      "request.duration",
			CreateEvent:     statsd.NewEvent,
			CustomEventName: "repo.blob",
		}

		// Get request metadata.
		args, err := extractBlobRequestArgs(r)
		// Track request metadata.
		trackingArgs.EventInfo = append(
			trackingArgs.EventInfo,
			fmt.Sprintf("%v", args),
		)
		if err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			defer datadog.TrackTranscation(trackingArgs)
			errors.RespondWithError(w, err)
			return
		}

		// Request the filepath from depot gitweb.
		hashedRepoName := depot.BuildHashedRepoName(args.author, args.repo, args.sha)
		depotBlobURL := fmt.Sprintf(
			"http://%s/?p=%s.git;a=blob_plain;f=%s;hb=refs/heads/master",
			depot.DepotInternalServiceAddress,
			hashedRepoName,
			args.path)
		depotBlobResp, err := http.Get(depotBlobURL)
		if err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			defer datadog.TrackTranscation(trackingArgs)
			errors.RespondWithError(w, err)
			return
		}

		// If path was not found return 404.
		if depotBlobResp.StatusCode == 404 {
			trackingArgs.AlertType = datadog.Info
			trackingArgs.Tags = append(trackingArgs.Tags, "404")
			defer datadog.TrackTranscation(trackingArgs)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte{})
		}

		body, err := ioutil.ReadAll(depotBlobResp.Body)
		if err != nil {
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			defer datadog.TrackTranscation(trackingArgs)
			errors.RespondWithError(w, err)
			return
		}

		depotBlobResp.Body.Close()
		trackingArgs.AlertType = datadog.Success
		defer datadog.TrackTranscation(trackingArgs)

		if len(body) > 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
		}
	}
}

func extractBlobRequestArgs(r *http.Request) (blobRequestArgs, error) {
	vars := mux.Vars(r)
	args := blobRequestArgs{}

	args.author = vars[urlVarAuthor]
	if len(args.author) < 0 {
		return args, NewInvalidURLParameterError(
			urlVarAuthor,
			args.author)
	}

	args.repo = vars[urlVarRepo]
	if len(args.repo) < 0 {
		return args, NewInvalidURLParameterError(urlVarRepo, args.repo)
	}

	args.sha = vars[urlVarSHA]
	if len(args.sha) < 0 {
		return args, NewInvalidURLParameterError(urlVarSHA, args.sha)
	}

	args.path = vars[urlVarPath]
	if len(args.path) < 0 {
		return args, NewInvalidURLParameterError(urlVarPath, args.path)
	}

	return args, nil
}
