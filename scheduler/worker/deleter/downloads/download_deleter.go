package downloads

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/db/model/package/download"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

// downloadDeleterArgs is the arguments struct for downloadDeleter.
type downloadDeleterArgs struct {
	q  db.Queryable
	wg *sync.WaitGroup
	// TODO(skeswa): synced errors go here.
	errs      chan error
	logger    common.JobLogger
	summaries chan pkg.Summary
}

// downloadDeleter is a worker for the Update function. It reads incoming
// packages from the summaries channel and updates each package's metrics. If
// any errors are encountered in the process, then they are put into the errors
// channel.
func downloadDeleter(args downloadDeleterArgs) {
	// Guarantee that the waitgroup is notified at the end.
	defer args.wg.Done()

	// For each package summary, attempt an update in the database.
	for summary := range args.summaries {
		args.logger.Infof(
			"Now deleting old downloads for package %s/%s\n",
			summary.Author,
			summary.Repo)

		err := download.DeleteOld(args.q, summary.Author, summary.Repo)
		if err != nil {
			args.errs <- err
			continue
		}

		args.logger.Infof(
			"Finished deleting old downloads for package %s/%s\n",
			summary.Author,
			summary.Repo)
	}
}
