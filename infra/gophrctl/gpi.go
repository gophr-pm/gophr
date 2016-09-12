package main

import (
	"errors"

	"gopkg.in/urfave/cli.v1"
)

func readGPI(c *cli.Context) (string, error) {
	gpi := c.String(flagNameGPI)
	if len(gpi) < 1 {
		return gpi, errors.New("The google project id must be specified for this command to function.")
	}

	return gpi, nil
}
