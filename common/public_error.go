package common

import (
	"log"
	"net/http"
)

const (
	unexpectedErrorMessage    = "An unexpected internal server error occurred."
	nonPublicLogMessageFormat = "Could not repsond to request with non-public error: %v"
)

type PublicError interface {
	PublicError() (int, string)
}

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
