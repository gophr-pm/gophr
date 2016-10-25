package main

import "net/http"

const (
	contentTypeJSON   = "application/json"
	contentTypeHeader = "Content-Type"
)

// respondWithJSON sends a response with status code 200.
func respondWithJSON(w http.ResponseWriter, json []byte) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.Write(json)
}
