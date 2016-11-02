package main

import (
	"fmt"
	"net/http"
)

// StatusHandler creates an HTTP request handler that responds to status
// requests.
func StatusHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	}
}
