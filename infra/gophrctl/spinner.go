package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

var (
	loadingSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
)

func startSpinner(message string) {
	loadingSpinner.Color("green")
	loadingSpinner.Suffix = " " + message
	loadingSpinner.Start()
}

func stopSpinner(operationSuccessful bool) {
	if operationSuccessful {
		loadingSpinner.FinalMSG = "done."
	} else {
		loadingSpinner.FinalMSG = "failed."
	}
	loadingSpinner.Stop()
	fmt.Println()
}
