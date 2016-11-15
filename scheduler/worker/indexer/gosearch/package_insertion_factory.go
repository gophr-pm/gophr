package gosearch

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/gophr-pm/gophr/lib/github"
)

// awesomeChecker is a proxy for awesome.IncludesPackage.
type awesomeChecker func(
	q db.Queryable,
	author string,
	repo string,
) (bool, error)

// packageInsertionFactoryArgs is the arguments struct for
// packageInsertionFactory.
type packageInsertionFactoryArgs struct {
	q                 db.Queryable
	wg                *sync.WaitGroup
	errs              chan error
	ghSvc             github.RequestService
	isAwesome         awesomeChecker
	newPackages       chan packageSetEntry
	packageInsertions chan pkg.InsertArgs
}

// packageInsertionFactory collects the information necessary to turn each new
// package in entries into a package insertion.
func packageInsertionFactory(args packageInsertionFactoryArgs) {
	// Unconditionally make sure the waitgroup is notified when finished.
	defer args.wg.Done()

	var (
		err      error
		awesome  bool
		repoData dtos.GithubRepo
	)

	for newPackage := range args.newPackages {
		// Fetch metadata for this project from github.
		if repoData, err = args.ghSvc.FetchRepoData(
			newPackage.author,
			newPackage.repo,
		); err != nil {
			args.errs <- err
			continue
		}

		// Check if this package is awesome.
		if awesome, err = args.isAwesome(
			args.q,
			newPackage.author,
			newPackage.repo,
		); err != nil {
			args.errs <- err
			continue
		}

		// Use the github metadata to create the package insertion.
		args.packageInsertions <- pkg.InsertArgs{
			Repo:        newPackage.repo,
			Stars:       repoData.Stars,
			Author:      newPackage.author,
			Awesome:     awesome,
			Queryable:   args.q,
			Description: repoData.Description,
		}
	}
}
