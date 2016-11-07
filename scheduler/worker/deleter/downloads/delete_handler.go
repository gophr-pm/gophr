package downloads

import (
	"net/http"
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

// The name of this job.
const jobName = "delete-downloads"

// DeleteHandler exposes an endpoint that deletes all of the hourly downloads
// older than one month.
func DeleteHandler(
	q db.Queryable,
	numWorkers int,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err       error
			errs      = make(chan error)
			logger    common.JobLogger
			loggerWG  sync.WaitGroup
			deleterWG sync.WaitGroup
			jobParams common.JobParams
			summaries = make(chan pkg.Summary)
		)

		// Read job params so we can build a logger.
		if jobParams, err = common.ReadJobParams(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Build a logger for use in the sub-routines.
		logger = common.NewJobLogger(jobName, jobParams)

		// Log the runtime events of this job.
		logger.Start()
		defer logger.Finish()

		// Spin up the error logger.
		loggerWG.Add(1)
		go common.LogErrors(logger, &loggerWG, errs)

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

		// Wait for the logger to exit.
		loggerWG.Wait()
	}
}
