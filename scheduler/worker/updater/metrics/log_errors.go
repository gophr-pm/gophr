package metrics

import (
	"log"
	"time"
)

// logErrors logs everything that goes wrong.
func logErrors(errs chan error) {
	// TODO(skeswa): standardize scheduler job logging.
	var (
		startTime  = time.Now().Format(time.RFC3339)
		errorCount = 1
	)

	// Log every error that comes through the channel til it closes.
	for err := range errs {
		log.Printf(
			`[updater:metrics:%s] Encountered error #%d: %v`,
			startTime,
			errorCount,
			err)

		// Bump the error count for each encountered error.
		errorCount++
	}
}
