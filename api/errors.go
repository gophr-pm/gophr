package main

import (
	"fmt"
	"net/http"
)

/*********************** INVALID QUERY STRING PARAMETER ***********************/

type InvalidQueryStringParameterError struct {
	ParameterName  string
	ParameterValue interface{}
}

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

func (err InvalidQueryStringParameterError) PublicError() (int, string) {
	return http.StatusBadRequest, err.Error()
}
