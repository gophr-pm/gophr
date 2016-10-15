package errors

import (
	"fmt"
	"net/http"
)

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
