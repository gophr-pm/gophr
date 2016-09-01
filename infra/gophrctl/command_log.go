package main

import (
	"errors"

	"gopkg.in/urfave/cli.v1"
)

func logCommand(c *cli.Context) error {
	exit(exitCodeLogFailed, c, "log", errors.New("Command not yet implemented"))
	return nil
}
