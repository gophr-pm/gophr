package main

import (
	"errors"

	"gopkg.in/urfave/cli.v1"
)

func stopCommand(c *cli.Context) error {
	exit(exitCodeStopFailed, c, "stop", errors.New("Command not yet implemented"))
	return nil
}
