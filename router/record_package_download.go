package main

import (
	"log"

	"github.com/gophr-pm/gophr/common/models"
)

// recordPackageDownload is a helper function that records the download of a
// specific version of a package.
func recordPackageDownload(args packageDownloadRecorderArgs) {
	// TODO(skeswa): support "version" + "sha" in daily downloads with sha
	// remaining the primary identfier.
	if err := models.RecordDailyDownload(
		args.db,
		args.author,
		args.repo,
		args.sha); err != nil {
		// Instead of bubbling this error, just commit it to the logs. This is
		// necessary because this function is executed asynchronously.
		log.Printf(
			"[ERR] Failed to record download for package %s/%s@%s: %v\n",
			args.author,
			args.repo,
			args.sha,
			err)
	}
}
