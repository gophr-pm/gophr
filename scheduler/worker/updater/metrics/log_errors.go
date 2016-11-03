package metrics

import "github.com/gophr-pm/gophr/scheduler/worker/common"

// logErrors logs everything that goes wrong.
func logErrors(logger common.JobLogger, errs chan error) {
	errorCount := 1

	// Log every error that comes through the channel til it closes.
	for err := range errs {
		logger.Errorf("Encountered error #%d: %v.", errorCount, err)
		// Bump the error count for each encountered error.
		errorCount++
	}
}
