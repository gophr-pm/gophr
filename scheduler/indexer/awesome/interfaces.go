package awesome

import (
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
)

// PackageTuple represents packages found on awesome-go as a Tuple of
// author and repo.
type PackageTuple struct {
	author string
	repo   string
}

// IndexArgs is the args struct for indexing awesome-go packages.
type IndexArgs struct {
	Init            Init
	DoHTTPGet       httpGetter
	BatchExecutor   batchExecutor
	PackageFetcher  packageFetcher
	PersistPackages persistPackages
}

// Init is responsible for setting up the app configuration and db
// connection.
type Init func() (*config.Config, db.Client)

// batchExecutor executes a batch cassandra query and returns errors via
// an error channel.
type batchExecutor func(batch db.Batch, resultChan chan error)

// packageFetcher is responsible for fetching packages found on awesome-go.
type packageFetcher func(FetchAwesomeGoListArgs) ([]PackageTuple, error)

// persistPackages is reponsible for grouping batch package queries.
type persistPackages func(PersistAwesomePackagesArgs) error

// FetchAwesomeGoListArgs is the args struct for fetching awesome go packages
// from godoc.
type FetchAwesomeGoListArgs struct {
	doHTTPGet httpGetter
}

// httpGetter executes an HTTP get to the specified URL and returns the
// corresponding response.
type httpGetter func(url string) ([]byte, error)

// PersistAwesomePackagesArgs is the args struct for PersistAwesomePackages.
type PersistAwesomePackagesArgs struct {
	Session         db.BatchingQueryable
	NewBatchCreator newBatch
	BatchExecutor   batchExecutor
	PackageTuples   []PackageTuple
}

// newBatch returns a new cqlBatch query.
type newBatch func() db.Batch
