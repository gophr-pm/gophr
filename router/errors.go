package main

import (
	"fmt"
	"net/http"
)

// TODO(skeswa): centralize all the errors here, give them codes, add logging

func respondWithInvalidURL(resp http.ResponseWriter, url string) {
	resp.WriteHeader(http.StatusBadRequest)
	resp.Write([]byte(fmt.Sprintf(
		"Failed to process URL \"%s\". Please refer to the gophr docs for information on how to use it.",
		url,
	)))
}

func respondWithError(resp http.ResponseWriter, err error) {
	// TODO(skeswa): customize with custom formatting logic for different errors
	resp.WriteHeader(http.StatusInternalServerError)
	resp.Write([]byte(err.Error()))
}
