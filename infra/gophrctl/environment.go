package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

type environment string

const (
	environmentDev  = environment("dev")
	environmentProd = environment("prod")
)

func readEnvironment(c *cli.Context) (environment, error) {
	switch c.GlobalString(flagNameEnv) {
	case "dev":
		return environmentDev, nil
	case "prod":
		return environmentProd, nil
	default:
		return environmentDev, fmt.Errorf("Invalid environment \"%s\" specified.", c.GlobalString(flagNameEnv))
	}
}
