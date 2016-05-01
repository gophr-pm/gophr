package main

import (
	"fmt"
	"net/http"
)

/*************************** INVALID REQUEST BODY *****************************/

// InvalidRequestBodyError is an error that occurs when a particular parameter
// defies expectations.
type InvalidRequestBodyError struct {
	ExpectedFormat     string
	RequestBodyContent string
	CausedBy           []error
}

// NewInvalidRequestBodyError creates a new InvalidRequestBodyError.
func NewInvalidRequestBodyError(
	expectedFormat string,
	requestBodyContent string,
	causes ...error,
) InvalidRequestBodyError {
	return InvalidRequestBodyError{
		ExpectedFormat:     expectedFormat,
		RequestBodyContent: requestBodyContent,
		CausedBy:           causes,
	}
}

func (err InvalidRequestBodyError) Error() string {
	return fmt.Sprintf(
		`Invalid request body format. Expected %s, but received:\n\n%v`,
		err.ExpectedFormat,
		err.RequestBodyContent,
	)
}

func (err InvalidRequestBodyError) String() string {
	return err.Error()
}

// Causes returns the error(s) that caused this error.
func (err InvalidRequestBodyError) Causes() []error {
	return err.CausedBy
}

// PublicError returns an outside-friendly error message, and a
// corresponding status code.
func (err InvalidRequestBodyError) PublicError() (int, string) {
	return http.StatusBadRequest, err.Error()
}
