package main

import (
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

func readGophrRoot(c *cli.Context) (string, error) {
	return filepath.Abs(c.GlobalString(flagNameRepoPath))
}
