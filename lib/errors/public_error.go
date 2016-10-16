package errors

import (
	"log"
	"net/http"
)

const (
	unexpectedErrorMessage    = "An unexpected internal server error occurred."
	nonPublicLogMessageFormat = "Could not repsond to request with non-public error: %v"
)

// PublicError is an error that has an outside-friendly error message, and a
// corresponding status code.
type PublicError interface {
	// PublicError returns an outside-friendly error message, and a
	// corresponding status code.
	PublicError() (int, string)
}

// ResultingError is an error that reuslted from other errors.
type ResultingError interface {
	// Causes returns the list of errors that preceded this one.
	Causes() []error
}

// RespondWithError responds, using the supplied response writer, with an error.
// If the error is a PublicError, then its message is sent in the response.
// Otherwise, a generic internal server error message is sent.
func RespondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case PublicError:
		statusCode, errorMessage := e.PublicError()
		w.WriteHeader(statusCode)
		w.Write([]byte(errorMessage))
	default:
		log.Printf(nonPublicLogMessageFormat, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(unexpectedErrorMessage))
	}
}
