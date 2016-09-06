package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func sshCommand(c *cli.Context) error {
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
			goto exitWithError
		}
		if err = sshModule(m, moduleName, env); err != nil {
			goto exitWithError
		}
	} else {
		err = newNoSuchModuleError(moduleName)
		goto exitWithErrorAndHelp
	}

	return nil
exitWithError:
	exit(exitCodeSSHFailed, nil, "", err)
	return nil
exitWithErrorAndHelp:
	exit(exitCodeSSHFailed, c, "ssh", err)
	return nil
}

func sshModule(m *module, moduleName string, env environment) error {
	var (
		err      error
		podNames []string
	)

	if podNames, err = filterK8SPods(moduleName); err != nil {
		return err
	}

	// Only use the first pod that comes up - if there even is one.
	if len(podNames) > 0 {
		// TODO(skeswa): refine this to include the follow flag and also allow the
		// user to choose which pod.
		execK8SBash(podNames[0])
		return nil
	}

	// TODO(skeswa): refine and standardize this error.
	exit(exitCodeSSHFailed, nil, "", fmt.Errorf("Could not find any pods matching module \"%s\"", moduleName))
	return nil
}
