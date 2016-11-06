package gosearch

import (
	"net/http"
	"sync"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

// The name of this job.
const jobName = "index-go-search"

// IndexHandler exposes an endpoint that indexes all of the go packages known
// to http://go-search.org/.
func IndexHandler(
	q db.Queryable,
	conf *config.Config,
	ghSvc github.RequestService,
	numWorkers int,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err       error
			errs      = make(chan error)
			logger    common.JobLogger
			loggerWG  sync.WaitGroup
			jobParams common.JobParams
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

		// Make sure those those errors get logged.
		loggerWG.Add(1)
		go common.LogErrors(logger, &loggerWG, errs)

		// Let the indexing begin!
		index(indexArgs{
			q:      q,
			conf:   conf,
			ghSvc:  ghSvc,
			logger: logger,
			packageInsertionFactoryCount: numWorkers,
		})

		// There cannot be any more errors, so kill the errors channel.
		close(errs)

		// Wait for the error logging to finish.
		loggerWG.Wait()
	}
}
