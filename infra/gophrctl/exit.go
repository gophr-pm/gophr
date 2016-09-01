package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
)

const (
	exitCodeStartFailed  = 100
	exitCodeBuildFailed  = 101
	exitCodeCycleFailed  = 102
	exitCodeLogFailed    = 103
	exitCodeSSHFailed    = 104
	exitCodeStopFailed   = 105
	exitCodeUpdateFailed = 106
)

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
