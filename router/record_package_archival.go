package main

import (
	"log"

	"github.com/gophr-pm/gophr/common/models/packages/archives"
)

// recordPackageArchival is a helper function that records the download of a
// specific version of a package.
func recordPackageArchival(args packageArchivalRecorderArgs) {
	// Use the package archive model to record this in the database.
	if err := archives.Create(
		args.db,
		args.author,
		args.repo,
		args.sha); err != nil {
		// Instead of bubbling this error, just commit it to the logs. This is
		// necessary because this function is executed asynchronously.
		log.Printf(
			"[ERR] Failed to record archival for package %s/%s@%s: %v\n",
			args.author,
			args.repo,
			args.sha,
			err)
	}
}
