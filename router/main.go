package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	healthCheckRoute       = "/status"
	wildcardHandlerPattern = "/"
)

var (
	statusCheckResponse = []byte("ok")
)

func main() {
	http.HandleFunc(wildcardHandlerPattern, handler)

	portStr := os.Getenv("PORT")
	var port int
	if len(portStr) == 0 {
		fmt.Println("Port left unspecified; setting port to 3000.")
		port = 3000
	} else if portNum, err := strconv.Atoi(portStr); err == nil {
		fmt.Printf("Port was specified as %d.\n", portNum)
		port = portNum
	} else {
		fmt.Println("Port was invalid; setting port to 3000.")
		port = 3000
	}

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	requestID := generaterequestID()

	log.Printf("[%s] New request received: %s\n", requestID, r.URL.Path)

	if r.URL.Path == healthCheckRoute {
		log.Printf("[%s] Handling request for \"%s\" as a health check\n", requestID, r.URL.Path)

		w.Write(statusCheckResponse)
	} else {
		log.Printf("[%s] Handling request for \"%s\" as a package request\n", requestID, r.URL.Path)

		err := RespondToPackageRequest(requestID, r, w)
		if err != nil {
			respondWithError(w, err)
		}
	}
}
