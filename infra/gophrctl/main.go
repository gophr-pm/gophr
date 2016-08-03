package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "gophrctl"
	app.Usage = "Manages the gophr development and deployment environment."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env,e",
			Value: "dev",
			Usage: "gophr execution evironment",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "build",
			Usage: "Updates component images",
			Action: func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
	}
}
