package awesome

import (
	"net/http"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

// The name of this job.
const jobName = "index-awesome-go"

// IndexHandler exposes an endpoint that indexes all of the awesome go packages.
func IndexHandler(
	q db.BatchingQueryable,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Let the indexing begin!
		index(indexArgs{
			q:               q,
			logger:          logger,
			doHTTPGet:       http.Get,
			batchExecutor:   execBatch,
			packageFetcher:  fetchAwesomeGoList,
			persistPackages: persistAwesomePackages,
		})
	}
}
