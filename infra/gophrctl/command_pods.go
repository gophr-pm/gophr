package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func podsCommand(c *cli.Context) error {
	if err := runInK8S(c, func() error {
		var (
			env = readEnvironment(c)
		)

		if env == environmentDev {
			if err := assertMinikubeRunning(); err != nil {
				return err
			}
		}

		output, err := getPodsInK8S()
		if err != nil {
			return err
		}
		fmt.Println(output)

		return nil
	}); err != nil {
		exit(exitCodePodsFailed, nil, "", err)
	}

	return nil
}
