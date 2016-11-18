package common

// ErrorLoggingResult is the results struct of LogErrors.
type ErrorLoggingResult struct {
	Errors []error
}

// LogErrors logs everything that goes wrong.
func LogErrors(
	logger JobLogger,
	resultChan chan ErrorLoggingResult,
	errs chan error,
) {
	result := ErrorLoggingResult{}

	// Log every error that comes through the channel til it closes.
	for err := range errs {
		logger.Errorf("Encountered error #%d: %v.", len(result.Errors)+1, err)
		// Bump the error count for each encountered error.
		result.Errors = append(result.Errors, err)
	}

	// Exit when the errors channel stops.
	resultChan <- result
	close(resultChan)
}
