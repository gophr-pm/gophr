package main

import (
	"github.com/skeswa/gophr/common/depot"
	"github.com/skeswa/gophr/common/models"
)

// isPackageArchived returns true if the specified package has been archived
// before.
func isPackageArchived(args packageArchivalArgs) (bool, error) {
	// Translate the args struct into the call to the models package.
	archivedInDB, err := models.IsPackageArchived(
		args.db,
		args.author,
		args.repo,
		args.sha)

	if err != nil {
		return false, err
	}

	if archivedInDB {
		return true, nil
	}

	archivedInDepot, err := depot.CheckIfRefExists(args.author, args.repo, args.sha)
	if err != nil {
		return false, err
	}

	if !archivedInDepot {
		return false, nil
	}

	// Since we wouldn't have gotten this far if this were already recorded,
	// make sure that we record it now.
	go args.recordPackageArchival(packageArchivalArgs{
		db:     args.db,
		sha:    args.sha,
		repo:   args.repo,
		author: args.author,
	})

	return true, nil
}
