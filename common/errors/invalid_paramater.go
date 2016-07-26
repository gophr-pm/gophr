package errors

import "fmt"

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
