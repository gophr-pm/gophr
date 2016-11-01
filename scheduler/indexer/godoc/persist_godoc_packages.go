package godoc

import (
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/model/package/godoc"
)

const (
	numPackagesPerBatch = 50
)

// persistGodocPackages batch inserts awesome packages to help reduce network traffic.
func persistGodocPackages(session db.BatchingQueryable, pkgs []PackageMetadata) error {
	var (
		currentBatch = session.NewUnloggedBatch()
		resultChan   = make(chan error)
		numBatches   = 0
		resultCount  = 0
		queryErrors  []error
	)

	for i, pkg := range pkgs {
		godoc.AppendAddPackageQuery(currentBatch, pkg.author, pkg.repo, pkg.description)
		if last := i == len(pkgs)-1; i%numPackagesPerBatch == 0 && i > 0 || last {
			numBatches++
			go execBatch(currentBatch, resultChan)
			if !last {
				currentBatch = session.NewUnloggedBatch()
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
	batch db.Batch,
	resultChan chan error,
) {
	err := batch.Execute()
	resultChan <- err
}
