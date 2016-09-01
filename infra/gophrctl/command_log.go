package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

func logCommand(c *cli.Context) error {
	var (
		m          *module
		err        error
		env        environment
		exists     bool
		moduleName string
	)

	if env, err = readEnvironment(c); err != nil {
		goto exitWithError
	}

	moduleName = c.Args().First()
	if m, exists = modules[moduleName]; exists {
		if err = assertMinikubeRunning(); err != nil {
			goto exitWithError
		}
		if err = logModule(m, moduleName, env); err != nil {
			goto exitWithError
		}
	} else {
		err = newNoSuchModuleError(moduleName)
		goto exitWithErrorAndHelp
	}

	return nil
exitWithError:
	exit(exitCodeLogFailed, nil, "", err)
	return nil
exitWithErrorAndHelp:
	exit(exitCodeLogFailed, c, "log", err)
	return nil
}

func logModule(m *module, moduleName string, env environment) error {
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
		execK8SLogs(podNames[0], true)
		return nil
	}

	// TODO(skeswa): refine and standardize this error.
	exit(exitCodeLogFailed, nil, "", fmt.Errorf("Could not find any pods matching module \"%s\"", moduleName))
	return nil
}
