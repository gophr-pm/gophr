package metrics

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
)

// Update reads every package from the database and updates the metrics of each.
func Update(q db.Queryable, numWorkers int) {
	var (
		wg        sync.WaitGroup
		errs      = make(chan error)
		summaries = make(chan pkg.Summary)
	)

	// Start reading packages.
	go logErrors(errs)
	go pkg.ReadAll(q, summaries, errs)

	// Create all of the update workers, then wait for them.
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go packageUpdater(q, &wg, errs, summaries)
	}
	wg.Wait()

	// Close the errors channel since nothing else will ever go through.
	close(errs)
}
