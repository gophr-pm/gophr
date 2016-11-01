package awesome

import (
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/model/package/awesome"
)

const (
	numPackagesPerBatch = 50
)

// PersistAwesomePackages batch inserts awesome packages to help reduce network traffic.
func PersistAwesomePackages(args PersistAwesomePackagesArgs) error {
	var (
		currentBatch = args.NewBatchCreator()
		resultChan   = make(chan error)
		numBatches   = 0
		resultCount  = 0
		queryErrors  []error
	)

	for i, pkg := range args.PackageTuples {
		awesome.AppendAddPackageQuery(currentBatch, pkg.author, pkg.repo)
		if last := i == len(args.PackageTuples)-1; i%numPackagesPerBatch == 0 && i > 0 || last {
			numBatches++
			go args.BatchExecutor(currentBatch, resultChan)
			if !last {
				currentBatch = args.NewBatchCreator()
			}
		}
	}

	for err := range resultChan {
		resultCount++
		queryErrors = append(queryErrors, err)
		if resultCount == numBatches {
			close(resultChan)
		}
	}

	for _, err := range queryErrors {
		if err != nil {
			return errors.ComposeErrors("Failed to persist awesome packages", queryErrors)
		}
	}

	return nil
}
