package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func sshCommand(c *cli.Context) error {
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
			if env == environmentDev {
				if err = assertMinikubeRunning(); err != nil {
					return err
				}
			}

			if err = sshModule(m, moduleName, env); err != nil {
				return err
			}
		} else {
			err = newNoSuchModuleError(moduleName)
			return err
		}

		return nil
	}); err != nil {
		exit(exitCodeSSHFailed, nil, "", err)
	}

	return nil
}

func sshModule(m *module, moduleName string, env environment) error {
	var (
		err      error
		podNames []string
	)

	// Get the running module pods.
	if podNames, err = filterK8SPods(moduleName); err != nil {
		return err
	}

	if len(podNames) <= 0 {
		// TODO(skeswa): refine and standardize this error.
		return fmt.Errorf("Could not find any pods matching module \"%s\"", moduleName)
	} else if len(podNames) == 1 {
		// TODO(skeswa): switch this to a process fork instead of a process
		// replacement.
		execK8SBash(podNames[0])
		return nil
	} else {
		// There are multiple pods to choose from - offer some options.
		podIndex := promptChoice(promptChoiceArgs{
			prompt:             "In which pod should bash be executed?",
			choice:             "Pod",
			options:            podNames,
			defaultOptionIndex: 0,
		})
		// TODO(skeswa): switch this to a process fork instead of a process
		// replacement.
		execK8SBash(podNames[podIndex])
		return nil
	}
}
