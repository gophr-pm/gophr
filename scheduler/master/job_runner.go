package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	workerK8SServiceName = "scheduler-worker-svc"
)

type httpGetter func(url string) (resp *http.Response, err error)

func newJobRunner(j job, request httpGetter) func() {
	return func() {
		startTime := time.Now()
		startTimeStr := startTime.Format(time.RFC3339)

		log.Printf(`Started job "%s" at %s.\n`, j.name, startTimeStr)
		if resp, err := request(fmt.Sprintf(
			"http://%s%s",
			workerK8SServiceName,
			j.path)); err != nil {
			log.Printf(
				`Job "%s" started at %s failed to execute: %v.\n`,
				j.name,
				startTimeStr,
				err)
		} else if resp.StatusCode != 200 {
			log.Printf(
				`Job "%s" started at %s failed to execute: received status code %d.\n`,
				j.name,
				startTimeStr,
				resp.StatusCode)
		}
		log.Printf(
			`Job "%s" started at %s executed successfully after %s.\n`,
			j.name,
			startTimeStr,
			time.Since(startTime).String())
	}
}
