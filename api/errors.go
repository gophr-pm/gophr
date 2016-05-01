package main

import (
	"fmt"
	"net/http"
)

/*********************** INVALID QUERY STRING PARAMETER ***********************/

// InvalidQueryStringParameterError is an error that occurs when a particular
// query string parameter defies expectations.
type InvalidQueryStringParameterError struct {
	ParameterName  string
	ParameterValue interface{}
}

// NewInvalidQueryStringParameterError creates a new
// InvalidQueryStringParameterError.
func NewInvalidQueryStringParameterError(
	parameterName string,
	parameterValue interface{},
) InvalidQueryStringParameterError {
	return InvalidQueryStringParameterError{
		ParameterName:  parameterName,
		ParameterValue: parameterValue,
	}
}

func (err InvalidQueryStringParameterError) Error() string {
	return fmt.Sprintf(
		`Invalid value "%v" specified for query string parameter "%s".`,
		err.ParameterValue,
		err.ParameterName,
	)
}

func (err InvalidQueryStringParameterError) String() string {
	return err.Error()
}

// PublicError is an error that has an outside-friendly error message, and a
// corresponding status code.
func (err InvalidQueryStringParameterError) PublicError() (int, string) {
	return http.StatusBadRequest, err.Error()
}

/*************************** INVALID URL PARAMETER ****************************/

// InvalidURLParameterError is an error that occurs when a particular URL
// parameter defies expectations.
type InvalidURLParameterError struct {
	ParameterName  string
	ParameterValue interface{}
}

// NewInvalidURLParameterError creates a new InvalidURLParameterError.
func NewInvalidURLParameterError(
	parameterName string,
	parameterValue interface{},
) InvalidURLParameterError {
	return InvalidURLParameterError{
		ParameterName:  parameterName,
		ParameterValue: parameterValue,
	}
}

func (err InvalidURLParameterError) Error() string {
	return fmt.Sprintf(
		`Invalid value "%v" specified for URL parameter "%s".`,
		err.ParameterValue,
		err.ParameterName,
	)
}

func (err InvalidURLParameterError) String() string {
	return err.Error()
}

// PublicError is an error that has an outside-friendly error message, and a
// corresponding status code.
func (err InvalidURLParameterError) PublicError() (int, string) {
	return http.StatusBadRequest, err.Error()
}
