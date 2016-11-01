package awesome

import (
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
)

// IndexArgs is the args struct for indexing awesome-go packages.
type IndexArgs struct {
	Init            Init
	DoHTTPGet       httpGetter
	BatchExecutor   BatchExecutor
	PackageFetcher  PackageFetcher
	PersistPackages PersistPackages
}

// Init is responsible for setting up the app configuration and db
// connection.
type Init func() (*config.Config, db.Client)

// httpGetter executes an HTTP get to the specified URL and returns the
// corresponding response.
type httpGetter func(url string) ([]byte, error)

// BatchExecutor executes a batch cassandra query and returns errors via
// an error channel.
type BatchExecutor func(batch db.Batch, resultChan chan error)

// PackageFetcher is responsible for fetching packages found on awesome-go.
type PackageFetcher func(FetchAwesomeGoListArgs) ([]PackageTuple, error)

// PersistPackages is reponsible for grouping batch package queries.
type PersistPackages func(PersistAwesomePackagesArgs) error

// FetchAwesomeGoListArgs lol
type FetchAwesomeGoListArgs struct {
	doHTTPGet httpGetter
}

// PackageTuple represents packages found on awesome-go as a Tuple of
// author and repo.
type PackageTuple struct {
	author string
	repo   string
}

// PersistAwesomePackagesArgs is the args struct for PersistAwesomePackages.
type PersistAwesomePackagesArgs struct {
	Session         db.BatchingQueryable
	NewBatchCreator newBatch
	BatchExecutor   BatchExecutor
	PackageTuples   []PackageTuple
}

// newBatch returns a new cqlBatch query.
type newBatch func() db.Batch
