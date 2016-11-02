package metrics

import (
	"log"
	"time"
)

// logErrors logs everything that goes wrong.
func logErrors(errs chan error) {
	// TODO(skeswa): standardize scheduler job logging.
	var (
		startTime    = time.Now()
		errorCount   = 1
		startTimeStr = startTime.Format(time.RFC3339)
	)

	log.Printf(`[updater:metrics:%s] Started.`, startTimeStr)

	// Log every error that comes through the channel til it closes.
	for err := range errs {
		log.Printf(
			`[updater:metrics:%s] Encountered error #%d: %v.`,
			startTimeStr,
			errorCount,
			err)

		// Bump the error count for each encountered error.
		errorCount++
	}

	log.Printf(
		`[updater:metrics:%s] Finished in %s.`,
		startTimeStr,
		time.Since(startTime).String())
}
