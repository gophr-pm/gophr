package main

import "github.com/skeswa/gophr/common/models"

// isPackageArchived returns true if the specified package has been archived
// before.
func isPackageArchived(args packageArchivalArgs) (bool, error) {
	// Translate the args struct into the call to the models package.
	return models.IsPackageArchived(
		args.db,
		args.author,
		args.repo,
		args.sha)
}
