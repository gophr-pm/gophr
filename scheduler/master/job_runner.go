package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	jobEndpointTemplate = "http://scheduler-worker-svc/%s?jobID=%s&startTime=%s"
)

// httpGetter executes an HTTP get.
type httpGetter func(url string) (resp *http.Response, err error)

// newJobRunner creates a new function that runs the specified job.
func newJobRunner(j job, request httpGetter) func() {
	return func() {
		var (
			jobID         = generateJobID()
			startTimeStr  = time.Now().Format(time.RFC3339)
			jobRequestURL = fmt.Sprintf(
				jobEndpointTemplate,
				j.path,
				jobID,
				startTimeStr)
		)

		log.Printf(
			`Scheduler fired off a job request to start "%s" with id "%s".`+"\n",
			j.name,
			jobID)
		if resp, err := request(jobRequestURL); err != nil {
			log.Printf(
				`Job "%s" with id "%s" started at %s failed to execute: %v.`+"\n",
				j.name,
				jobID,
				startTimeStr,
				err)
		} else if resp.StatusCode != 200 {
			// Turn the response body into a string.
			msgBytes, _ := ioutil.ReadAll(resp.Body)

			log.Printf(
				`Job "%s" with id "%s" started at %s failed to execute: received `+
					"status code %d: %s.\n",
				j.name,
				jobID,
				startTimeStr,
				resp.StatusCode,
				string(msgBytes[:]))
		}
	}
}
