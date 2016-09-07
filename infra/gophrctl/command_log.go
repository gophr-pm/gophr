package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func logCommand(c *cli.Context) error {
	if err := runInK8S(c, func() error {
		var (
			m          *module
			err        error
			env        = readEnvironment(c)
			exists     bool
			moduleName string
		)

		moduleName = c.Args().First()
		if m, exists = modules[moduleName]; exists {
			if err = assertMinikubeRunning(); err != nil {
				return err
			}
			if err = logModule(c, m, moduleName, env); err != nil {
				return err
			}
		} else {
			err = newNoSuchModuleError(moduleName)
			return err
		}

		return nil
	}); err != nil {
		exit(exitCodeLogFailed, nil, "", err)
	}

	return nil
}

func logModule(c *cli.Context, m *module, moduleName string, env environment) error {
	var (
		err      error
		podNames []string
	)

	if podNames, err = filterK8SPods(moduleName); err != nil {
		return err
	}

	// Only use the first pod that comes up - if there even is one.
	if len(podNames) > 0 {
		// TODO(skeswa): switch this to a process fork instead of a process
		// replacement.
		// TODO(skeswa): refine this to include the follow flag and also allow the
		// user to choose which pod.
		execK8SLogs(podNames[0], true)
		return nil
	}

	// TODO(skeswa): refine and standardize this error.
	return fmt.Errorf("Could not find any pods matching module \"%s\"", moduleName)
}
