package main

// isPackageArchived returns true if the specified package has been archived
// before.
func isPackageArchived(args packageArchivalCheckerArgs) (bool, error) {
	// Translate the args struct into the call to the models package.
	archivedInDB, err := args.isPackageArchivedInDB(
		args.db,
		args.author,
		args.repo,
		args.sha)
	if err != nil {
		return false, err
	} else if archivedInDB {
		// If archived in the db, it is *very* likely that the package is archived.
		return true, nil
	}

	// Check if this package version is in depot already.
	archivedInDepot, err := args.packageExistsInDepot(args.author, args.repo, args.sha)
	if err != nil {
		return false, err
	}

	// If this is true, it means that the package is not archived in either place.
	if !archivedInDepot {
		return false, nil
	}

	// Since we wouldn't have gotten this far if this were already recorded,
	// make sure that we record it now.
	go args.recordPackageArchival(packageArchivalRecorderArgs{
		db:     args.db,
		sha:    args.sha,
		repo:   args.repo,
		author: args.author,
	})

	return true, nil
}
