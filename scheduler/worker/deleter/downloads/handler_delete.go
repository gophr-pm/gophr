package downloads

import (
	"net/http"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

const (
	// The name of this job.
	jobName = "delete-downloads"
	// ddEventName is the name of the custom datadog event for this handler.
	ddEventName = "scheduler.worker.indexer.awesome"
)

// DeleteHandler exposes an endpoint that deletes all of the hourly downloads
// older than one month.
func DeleteHandler(
	q db.Queryable,
	ddClient datadog.Client,
	numWorkers int,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err          error
			errs         = make(chan error)
			logger       common.JobLogger
			deleterWG    sync.WaitGroup
			jobParams    common.JobParams
			summaries    = make(chan pkg.Summary)
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
		defer datadog.TrackTransaction(trackingArgs)

		// Build a logger for use in the sub-routines.
		logger = common.NewJobLogger(jobName, jobParams)

		// Log the runtime events of this job.
		logger.Start()
		defer logger.Finish()

		// Spin up the error logger.
		go common.LogErrors(logger, errLogResults, errs)

		// Start reading packages.
		logger.Info("Reading all packages from the database.")
		go pkg.ReadAll(q, summaries, errs)

		// Create all of the update workers, then wait for them.
		deleterWG.Add(numWorkers)
		logger.Infof("Spinning up %d download deleters.\n", numWorkers)
		for i := 0; i < numWorkers; i++ {
			go downloadDeleter(downloadDeleterArgs{
				q:         q,
				wg:        &deleterWG,
				errs:      errs,
				logger:    logger,
				summaries: summaries,
			})
		}
		deleterWG.Wait()

		// Close the errors channel since nothing else will ever go through.
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
