package main

import "gopkg.in/urfave/cli.v1"

type environment string

const (
	environmentDev  = environment("dev")
	environmentProd = environment("prod")
)

func readEnvironment(c *cli.Context) environment {
	if c.GlobalBool(flagNameProd) {
		return environmentProd
	}

	return environmentDev
}
