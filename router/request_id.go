package main

import "github.com/nu7hatch/gouuid"

var (
	requestIDFallback = "???"
)

func generateRequestID() string {
	u4, err := uuid.NewV4()
	if err != nil {
		return requestIDFallback
	}

	return u4.String()
}
