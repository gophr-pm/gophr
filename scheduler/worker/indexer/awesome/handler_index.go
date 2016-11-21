package awesome

import (
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

const (
	// The name of this job.
	jobName = "index-awesome-go"
	// ddEventName is the name of the custom datadog event for this handler.
	ddEventName = "scheduler.worker.indexer.awesome"
)

// IndexHandler exposes an endpoint that indexes all of the awesome go packages.
func IndexHandler(
	q db.BatchingQueryable,
	ddClient datadog.Client,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err          error
			errs         = make(chan error)
			logger       common.JobLogger
			jobParams    common.JobParams
			trackingArgs = datadog.TrackTransactionArgs{
				Tags:            []string{jobName, datadog.TagInternal},
				Client:          ddClient,
				AlertType:       datadog.Success,
				StartTime:       time.Now(),
				MetricName:      datadog.MetricJobDuration,
				CreateEvent:     statsd.NewEvent,
				CustomEventName: ddEventName,
			}
			errLogResults = make(chan common.ErrorLoggingResult)
		)

		// Read job params so we can build a logger.
		if jobParams, err = common.ReadJobParams(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Ensure that the transaction is tracked after the job finishes.
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, jobParams.String())
		defer datadog.TrackTransaction(&trackingArgs)

		// Build a logger for use in the sub-routines.
		logger = common.NewJobLogger(jobName, jobParams)

		// Log the runtime events of this job.
		logger.Start()
		defer logger.Finish()

		// Make sure those those errors get logged.
		go common.LogErrors(logger, errLogResults, errs)

		// Let the indexing begin!
		index(indexArgs{
			q:               q,
			errs:            errs,
			doHTTPGet:       http.Get,
			batchExecutor:   execBatch,
			packageFetcher:  fetchAwesomeGoList,
			persistPackages: persistAwesomePackages,
		})

		// There cannot be any more errors, so kill the errors channel.
		close(errs)

		// If there were errors, be sure to alter the tracking metadata.
		if errLogResult := <-errLogResults; len(errLogResult.Errors) > 0 {
			trackingArgs.AlertType = datadog.Error
			for _, err = range errLogResult.Errors {
				trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
			}
		}
	}
}
