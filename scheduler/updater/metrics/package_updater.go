package metrics

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
)

// packageUpdater is a worker for the Update function. It reads incoming
// packages from the summaries channel and updates each package's metrics. If
// any errors are encountered in the process, then they are put into the errors
// channel.
func packageUpdater(
	q db.Queryable,
	wg *sync.WaitGroup,
	// TODO(skeswa): synced errors go here.
	errs chan error,
	summaries chan pkg.Summary,
) {
	// Guarantee that the waitgroup is notified at the end.
	defer wg.Done()

	// For each package summary, attempt an update in the database.
	for summary := range summaries {
		metrics, err := getPackageMetrics(q, summary)
		if err != nil {
			errs <- err
			continue
		}

		err = pkg.UpdateMetrics(metrics)
		if err != nil {
			errs <- err
		}
	}
}
