package gosearch

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
)

// insertPackagesArgs is the arguments struct for insertPackages.
type insertPackagesArgs struct {
	q                 db.Queryable
	wg                *sync.WaitGroup
	errs              chan error
	packageInsertions chan pkg.InsertArgs
}

// insertPackages inserts all of the package insertions from the provided
// channel into the database.
func insertPackages(args insertPackagesArgs) {
	// Unconditionally make sure the waitgroup is notified when finished.
	defer args.wg.Done()

	// For all insertions, insert.
	for packageInsertion := range args.packageInsertions {
		if err := pkg.Insert(packageInsertion); err != nil {
			args.errs <- err
		}
	}
}
