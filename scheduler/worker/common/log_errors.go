package common

import "sync"

// LogErrors logs everything that goes wrong.
func LogErrors(logger JobLogger, wg *sync.WaitGroup, errs chan error) {
	errorCount := 1

	// Make sure the waitgroup is notified on exit.
	defer wg.Done()

	// Log every error that comes through the channel til it closes.
	for err := range errs {
		logger.Errorf("Encountered error #%d: %v.", errorCount, err)
		// Bump the error count for each encountered error.
		errorCount++
	}
}
