package gosearch

import (
	"net/http"
	"sync"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/db/model/package/awesome"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

// indexArgs is the arguments struct for index.
type indexArgs struct {
	q                            db.Queryable
	conf                         *config.Config
	ghSvc                        github.RequestService
	logger                       common.JobLogger
	packageInsertionFactoryCount int
}

// index is responsible for reading the list of packages from go-search, and
// inserting the ones gophr does not know about into the database.
func index(args indexArgs) {
	var (
		err                  error
		errs                 = make(chan error)
		newPackages          = make(chan packageSetEntry)
		insertPackagesWG     sync.WaitGroup
		goSearchPackages     *packageSet
		existingPackages     = make(chan pkg.Summary)
		packageInsertions    = make(chan pkg.InsertArgs)
		insertionFactoryWG   sync.WaitGroup
		goSearchPackageLimit = noPackageLimit
	)

	// In dev, the limit is reduced to go easy on the tiny database.
	if args.conf.IsDev {
		goSearchPackageLimit = devPackageLimit
	}

	args.logger.Info("Fetching go-search packages.")
	if goSearchPackages, err = fetchGoSearchPackages(
		http.Get,
		goSearchPackageLimit,
	); err != nil {
		args.logger.Errorf("Failed to fetch go-search packages: %v.", err)
		return
	}

	args.logger.Info("Started reading existing packages from the database.")
	go pkg.ReadAll(args.q, existingPackages, errs)

	args.logger.Info("Removing existing packages from the package set.")
	for existingPackage := range existingPackages {
		goSearchPackages.remove(existingPackage.Author, existingPackage.Repo)
	}

	if goSearchPackages.len() > 0 {
		args.logger.Infof(
			"Inserting %d new packages into the database.",
			goSearchPackages.len())

		// Pipe all of the remaining packages into a channel.
		go goSearchPackages.stream(newPackages)

		// Start up the insertion factories.
		insertionFactoryWG.Add(args.packageInsertionFactoryCount)
		for i := 0; i < args.packageInsertionFactoryCount; i++ {
			// Each one of these goroutines turns package set entries into
			// package insertions.
			go packageInsertionFactory(packageInsertionFactoryArgs{
				q:                 args.q,
				wg:                &insertionFactoryWG,
				errs:              errs,
				ghSvc:             args.ghSvc,
				isAwesome:         awesome.IncludesPackage,
				newPackages:       newPackages,
				packageInsertions: packageInsertions,
			})
		}

		// Start executing the resulting package insertions.
		insertPackagesWG.Add(1)
		go insertPackages(insertPackagesArgs{
			q:                 args.q,
			wg:                &insertPackagesWG,
			errs:              errs,
			packageInsertions: packageInsertions,
		})

		// Wait for the insertion factories to finish up.
		insertionFactoryWG.Wait()

		// Once the factories finish, close the insertions channel.
		close(packageInsertions)

		// Then wait for insert packages to exit.
		insertPackagesWG.Wait()
	} else {
		args.logger.Infof("No new packages found.")
	}
}
