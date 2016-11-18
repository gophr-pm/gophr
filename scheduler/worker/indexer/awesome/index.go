package awesome

import (
	"fmt"

	"github.com/gophr-pm/gophr/lib/db"
)

// batchExecutor executes a batch cassandra query and returns errors via
// an error channel.
type batchExecutor func(batch db.Batch, resultChan chan error)

// packageFetcher is responsible for fetching packages found on awesome-go.
type packageFetcher func(fetchAwesomeGoListArgs) ([]packageTuple, error)

// persistPackages is reponsible for grouping batch package queries.
type persistPackages func(persistAwesomePackagesArgs) error

// indexArgs is the args struct for indexing awesome-go packages.
type indexArgs struct {
	q               db.BatchingQueryable
	errs            chan error
	doHTTPGet       httpGetter
	batchExecutor   batchExecutor
	packageFetcher  packageFetcher
	persistPackages persistPackages
}

// index is responsible for finding all go awesome packages and persisting them
// in `awesome_packages` table for later look up.
func index(args indexArgs) {
	packageTuples, err := args.packageFetcher(fetchAwesomeGoListArgs{
		doHTTPGet: args.doHTTPGet,
	})
	if err != nil {
		args.errs <- fmt.Errorf("Failed to fetch awesome packages: %v.", err)
		return
	}

	if err = args.persistPackages(persistAwesomePackagesArgs{
		q:             args.q,
		packageTuples: packageTuples,
		batchExecutor: args.batchExecutor,
	}); err != nil {
		args.errs <- fmt.Errorf("Failed to persist packages: %v.", err)
		return
	}
}
