package main

import (
	"log"

	"github.com/gophr-pm/gophr/lib/db/model/package/download"
)

// recordPackageDownload is a helper function that records the download of a
// specific version of a package, but doesn't bubble an error.
func recordPackageDownload(args packageDownloadRecorderArgs) {
	if err := download.Record(
		args.db,
		args.author,
		args.repo,
		args.sha,
		args.ghSvc); err != nil {
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
