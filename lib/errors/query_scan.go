package errors

import "fmt"

// QueryScanError is an error that occurs when the result of a database query
// was not processed successfully.
type QueryScanError struct {
	ScanError, CloseError error
}

// NewQueryScanError creates a new QueryScanError.
func NewQueryScanError(scanError, closeError error) QueryScanError {
	return QueryScanError{
		ScanError:  scanError,
		CloseError: closeError,
	}
}

func (err QueryScanError) Error() string {
	if err.ScanError != nil && err.CloseError == nil {
		return fmt.Sprintf(
			`Failed to scan the results of the db query: %v.`,
			err.ScanError,
		)
	}

	if err.ScanError == nil && err.CloseError != nil {
		return fmt.Sprintf(
			`Failed to close the iterator of the db query: %v.`,
			err.CloseError,
		)
	}

	return fmt.Sprintf(
		`Failed to both scan the results of the db query and close the iterator. The scan error: %v. The close error: %v.`,
		err.ScanError,
		err.CloseError,
	)
}

func (err QueryScanError) String() string {
	return err.Error()
}
