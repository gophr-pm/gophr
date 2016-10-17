package github

import (
	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/model/package/awesome"
)

const (
	numPackagesPerBatch = 50
)

// persistAwesomePackages batch inserts awesome packages to help reduce network traffic.
func persistAwesomePackages(session query.BatchingQueryable, pkgs []awesomePackage) error {
	var (
		currentBatch = session.NewBatch(gocql.UnloggedBatch)
		resultChan   = make(chan error)
		numBatches   = 0
		resultCount  = 0
		queryErrors  []error
	)

	for i, pkg := range pkgs {
		awesome.AppendAddPackageQuery(currentBatch, pkg.author, pkg.repo)
		if last := i == len(pkgs)-1; i%numPackagesPerBatch == 0 && i > 0 || last {
			numBatches++
			go execBatch(session, currentBatch, resultChan)
			if !last {
				currentBatch = session.NewBatch(gocql.UnloggedBatch)
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

	if queryErrors != nil {
		return errors.ComposeErrors("Failed to persist awesome packages", queryErrors)
	}

	return nil
}

func execBatch(
	session query.BatchingQueryable,
	batch *gocql.Batch,
	resultChan chan error,
) {
	err := session.ExecuteBatch(batch)
	resultChan <- err
}
