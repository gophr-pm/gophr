package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/errors"
)

const (
	healthCheckRoute       = "/status"
	wildcardHandlerPattern = "/"
)

var (
	statusCheckResponse = []byte("ok")
)

func main() {
	// Initialize the router.
	config, session := common.Init()

	// Close the session right off the bat.
	// TODO(skeswa): add download reporting to the router.
	session.Close()

	// Start serving.
	http.HandleFunc(wildcardHandlerPattern, handler)
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	context := NewRequestContext(nil)

	log.Printf("[%s] New request received: %s\n", context.RequestID, r.URL.Path)

	if r.URL.Path == healthCheckRoute {
		log.Printf(
			"[%s] Handling request for \"%s\" as a health check\n",
			context.RequestID,
			r.URL.Path,
		)

		w.Write(statusCheckResponse)
	} else {
		log.Printf(
			"[%s] Handling request for \"%s\" as a package request\n",
			context.RequestID,
			r.URL.Path,
		)

		err := RespondToPackageRequest(context, r, w)
		if err != nil {
			errors.RespondWithError(w, err)
		}
	}
}
