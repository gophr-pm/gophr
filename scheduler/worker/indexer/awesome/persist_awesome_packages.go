package awesome

import (
	"fmt"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package/awesome"
)

const (
	numPackagesPerBatch = 50
)

// persistAwesomePackagesArgs is the args struct for persistAwesomePackages.
type persistAwesomePackagesArgs struct {
	q             db.BatchingQueryable
	batchExecutor batchExecutor
	packageTuples []packageTuple
}

// persistAwesomePackages batch inserts awesome packages to help reduce network traffic.
func persistAwesomePackages(args persistAwesomePackagesArgs) error {
	var (
		currentBatch = args.q.NewLoggedBatch()
		resultChan   = make(chan error)
		numBatches   = 0
		resultCount  = 0
		queryErrors  []error
	)

	for i, pkg := range args.packageTuples {
		awesome.AppendAddPackageQuery(currentBatch, pkg.author, pkg.repo)
		if last := i == len(args.packageTuples)-1; i%numPackagesPerBatch == 0 && i > 0 || last {
			numBatches++
			go args.batchExecutor(currentBatch, resultChan)
			if !last {
				currentBatch = args.q.NewLoggedBatch()
			}
		}
	}

	for err := range resultChan {
		resultCount++

		if err != nil {
			queryErrors = append(queryErrors, err)
		}

		if resultCount == numBatches {
			close(resultChan)
		}
	}

	if len(queryErrors) > 0 {
		// TODO(shikkic): concat errors goes here.
		return fmt.Errorf("Failed to persist awesome packages: %v.", queryErrors)
	}

	return nil
}
