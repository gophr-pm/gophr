package awesome

import "log"

// Index is responsible for finding all go awesome packages
// and persisting them in `awsome_packages` table for later look up.
func Index(args IndexArgs) error {
	_, session := args.Init()
	defer session.Close()

	log.Println("Fetching awesome go list.")
	packageTuples, err := args.PackageFetcher(FetchAwesomeGoListArgs{
		doHTTPGet: args.DoHTTPGet,
	})
	if err != nil {
		return err
	}

	log.Println("Persisting awesome go list.")
	if err = args.PersistPackages(
		PersistAwesomePackagesArgs{
			Session:         session,
			PackageTuples:   packageTuples,
			NewBatchCreator: session.NewBatch,
			BatchExecutor:   args.BatchExecutor,
		},
	); err != nil {
		return err
	}

	return nil
}
