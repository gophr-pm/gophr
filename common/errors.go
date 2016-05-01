package common

import (
	"fmt"
	"net/http"
)

/***************************** INVALID PARAMETER ******************************/

// InvalidParameterError is an error that occurs when a particular parameter
// defies expectations.
type InvalidParameterError struct {
	ParameterName  string
	ParameterValue interface{}
}

// NewInvalidParameterError creates a new InvalidParameterError.
func NewInvalidParameterError(
	parameterName string,
	parameterValue interface{},
) InvalidParameterError {
	return InvalidParameterError{
		ParameterName:  parameterName,
		ParameterValue: parameterValue,
	}
}

func (err InvalidParameterError) Error() string {
	return fmt.Sprintf(
		`Invalid value "%v" specified for parameter "%s".`,
		err.ParameterValue,
		err.ParameterName,
	)
}

func (err InvalidParameterError) String() string {
	return err.Error()
}

/********************************* QUERY SCAN *********************************/

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

/****************************** NO SUCH PACKAGE *******************************/

// NoSuchPackageError is an error that occurs when a requested package doesn't
// exist.
type NoSuchPackageError struct {
	PackageAuthor string
	PackageRepo   string
	CausedBy      []error
}

// NewNoSuchPackageError creates a new NoSuchPackageError.
func NewNoSuchPackageError(
	packageAuthor string,
	packageRepo string,
	causes ...error,
) NoSuchPackageError {
	return NoSuchPackageError{
		PackageAuthor: packageAuthor,
		PackageRepo:   packageRepo,
		CausedBy:      causes,
	}
}

func (err NoSuchPackageError) Error() string {
	return fmt.Sprintf(
		`Package "%s/%s" does not exist.`,
		err.PackageAuthor,
		err.PackageRepo,
	)
}

// Causes returns the error(s) that caused this error.
func (err NoSuchPackageError) Causes() []error {
	return err.CausedBy
}

// PublicError returns an outside-friendly error message, and a
// corresponding status code.
func (err NoSuchPackageError) PublicError() (int, string) {
	return http.StatusNotFound, err.Error()
}
