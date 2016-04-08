package main

import "github.com/nu7hatch/gouuid"

var (
	uuidNamespace     = []byte("gophr.pm")
	requestIDFallback = "???"
)

func generateRequestId() string {
	u5, err := uuid.NewV5(uuid.NamespaceURL, uuidNamespace)
	if err != nil {
		return requestIDFallback
	} else {
		return u5.String()
	}
}
