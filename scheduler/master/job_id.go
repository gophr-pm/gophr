package main

import (
	"github.com/Pallinder/go-randomdata"
)

// generateJobID generates a unique-ish string to identify a job with.
func generateJobID() string {
	return randomdata.SillyName()
}
