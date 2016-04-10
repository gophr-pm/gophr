package main

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
)

func generateRequestID() string {
	u4, err := uuid.NewV4()
	if err != nil {
		// NB: its virtually impossible to get test coverage on this line since this
		// never happens
		panic(fmt.Sprintf("Failed to generate a UUID: %v", err))
	}

	return u4.String()
}
