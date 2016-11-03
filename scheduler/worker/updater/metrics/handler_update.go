package metrics

import (
	"net/http"
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

// The name of this job.
const jobName = "update-metrics"

// UpdateHandler exposes an endpoint that reads every package from the database
// and updates the metrics of each.
func UpdateHandler(
	q db.Queryable,
	ghSvc github.RequestService,
	numWorkers int,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			wg        sync.WaitGroup
			errs      = make(chan error)
			summaries = make(chan pkg.Summary)
		)

		// Read job params so we can build a logger.
		jobParams, err := common.ReadJobParams(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Build a logger for use in the sub-routines.
		logger := common.NewJobLogger(jobName, jobParams)

		// Log the runtime events of this job.
		logger.Start()
		defer logger.Finish()

		// Start reading packages.
		go logErrors(logger, errs)
		go pkg.ReadAll(q, summaries, errs)

		// Create all of the update workers, then wait for them.
		wg.Add(numWorkers)
		logger.Infof("Spinning up %d workers.\n", numWorkers)
		for i := 0; i < numWorkers; i++ {
			go packageUpdater(packageUpdaterArgs{
				q:         q,
				wg:        &wg,
				errs:      errs,
				ghSvc:     ghSvc,
				logger:    logger,
				summaries: summaries,
			})
		}
		wg.Wait()

		// Close the errors channel since nothing else will ever go through.
		close(errs)
	}
}
