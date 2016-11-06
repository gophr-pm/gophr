package common

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	queryStringVarID        = "jobID"
	queryStringVarStartTime = "startTime"
)

// JobParams holds the parameters passed to a job.
type JobParams struct {
	ID        string
	StartTime time.Time
}

// ReadJobParams reads JobParams from the query string of an http request.
func ReadJobParams(r *http.Request) (JobParams, error) {
	var (
		id           = r.URL.Query().Get(queryStringVarID)
		err          error
		startTime    time.Time
		startTimeStr = r.URL.Query().Get(queryStringVarStartTime)
	)

	if len(id) < 1 {
		return JobParams{}, errors.New(
			`Empty jobID parameter provided in query string.`)
	}
	if startTime, err = time.Parse(time.RFC3339, startTimeStr); err != nil {
		return JobParams{}, fmt.Errorf(
			`Invalid startTime parameter provided in query string: %v.`,
			err)
	}

	return JobParams{
		ID:        id,
		StartTime: startTime,
	}, nil
}
