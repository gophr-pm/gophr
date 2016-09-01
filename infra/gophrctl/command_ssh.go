package main

import (
	"errors"

	"gopkg.in/urfave/cli.v1"
)

func sshCommand(c *cli.Context) error {
	exit(exitCodeSSHFailed, c, "ssh", errors.New("Command not yet implemented"))
	return nil
}
