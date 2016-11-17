package awesome

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
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
	logger          common.JobLogger
	doHTTPGet       httpGetter
	batchExecutor   batchExecutor
	packageFetcher  packageFetcher
	persistPackages persistPackages
	dataDogClient   datadog.Client
}

// index is responsible for finding all go awesome packages and persisting them
// in `awesome_packages` table for later look up.
func index(args indexArgs) {
	trackingArgs := datadog.TrackTransactionArgs{
		Tags: []string{
			"awesome-indexer",
			"external",
		},
		Client:          args.dataDogClient,
		StartTime:       time.Now(),
		EventInfo:       []string{},
		MetricName:      "request.duration",
		CreateEvent:     statsd.NewEvent,
		CustomEventName: "api.get.trending.packages",
		AlertType:       datadog.Success,
	}

	defer datadog.TrackTransaction(trackingArgs)

	args.logger.Info("Fetching awesome go list.")
	packageTuples, err := args.packageFetcher(fetchAwesomeGoListArgs{
		doHTTPGet: args.doHTTPGet,
	})
	if err != nil {
		trackingArgs.AlertType = datadog.Error
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
		args.logger.Errorf("Failed to fetch awesome packages: %v.", err)
		return
	}

	args.logger.Info("Persisting awesome go list.")
	if err = args.persistPackages(persistAwesomePackagesArgs{
		q:             args.q,
		packageTuples: packageTuples,
		batchExecutor: args.batchExecutor,
	}); err != nil {
		trackingArgs.AlertType = datadog.Error
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
		args.logger.Errorf("Failed to persist packages: %v.", err)
	}
}
