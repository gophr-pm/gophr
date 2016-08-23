package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"gopkg.in/urfave/cli.v1"
)

var (
	loadingSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
)

func newFailedModulesError(action string, failedModuleIds []string) error {
	// Create a buffer to list the modules that failed in the resulting error.
	buffer := bytes.Buffer{}
	buffer.WriteString("The following modules failed to ")
	buffer.WriteString(action)
	buffer.WriteByte(':')

	// Write out the list like [a,b,c,d] => "a, b, c and d".
	for i, failedModuleID := range failedModuleIds {
		if i > 0 {
			buffer.WriteString(", ")

			if i == (len(failedModuleIds) - 1) {
				buffer.WriteString("and ")
			}
		}

		buffer.WriteString(failedModuleID)
	}

	// Puncutate this message because we're civilized.
	buffer.WriteByte('.')

	return errors.New(buffer.String())
}

func modulesToString(excludeAll bool) string {
	var (
		buffer        bytes.Buffer
		isFirstModule = true
	)

	for moduleID := range modules {
		// Skip if this is the "all" module and we're supposed to skip it.
		if excludeAll && moduleID == allModuleID {
			continue
		}

		if !isFirstModule {
			buffer.WriteString(", ")
		} else {
			isFirstModule = false
		}

		buffer.WriteString(moduleID)
	}

	return buffer.String()
}

func exit(
	code int,
	c *cli.Context,
	command string,
	args ...interface{},
) {
	printError(args...)

	if c != nil {
		fmt.Println()
		if len(command) > 0 {
			cli.ShowCommandHelp(c, command)
		} else {
			cli.ShowAppHelp(c)
		}
	}

	os.Exit(code)
}

func startSpinner(message string) {
	loadingSpinner.Color("green")
	loadingSpinner.Suffix = " " + message
	loadingSpinner.FinalMSG = "done."
	loadingSpinner.Start()
}

func stopSpinner() {
	loadingSpinner.Stop()
	fmt.Println()
}
