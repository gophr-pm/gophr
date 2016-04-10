package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRequestID(t *testing.T) {
	var requestID, lastRequestID string

	for i := 0; i < 30; i++ {
		requestID = generateRequestID()
		assert.NotEqual(t, lastRequestID, requestID, "Request ids should be very unique")
	}
}
