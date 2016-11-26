package main

import "gopkg.in/urfave/cli.v1"

func cmdCommand(c *cli.Context) error {
	if err := runInK8S(c, func() error {
		args := []string{k8sNamespaceFlag}
		args = append(args, c.Args()...)

		return execInBackground(kubectl, args...)
	}); err != nil {
		exit(exitCodeCMD, nil, "", err)
	}

	return nil
}
