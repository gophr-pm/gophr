package metrics

import "log"

func logErrors(jobID string, errs chan error) {
	// TODO(skeswa): standardize scheduler job logging.
	errorCount := 1
	for err := range errs {
		log.Printf(
			`[scheduler:job:%s][updater:metrics] Encountered error #%d: %v`,
			jobID,
			errorCount,
			err)

		// Bump the error count for each encountered error.
		errorCount++
	}
}
