package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

func buildCommand(c *cli.Context) error {
	var (
		err     error
		env     environment
		module  string
		modules []string
	)

	if env, err = readEnvironment(c); err != nil {
		goto exit
	}

	if modules, err = getModules(c); err != nil {
		goto exit
	}

	if module, err = readModule(c, modules); err != nil {
		goto exit
	}

	if len(module) < 1 {
		module = "all modules"
	}

	printInfo("Building " + module + ".")
	startSpinner("Executing docker build...")
	if err = doDockerComposeBuild(c.GlobalString(flagNameRepoPath), module, env == environmentDev); err != nil {
		stopSpinner()
		goto exit
	}

	return nil

exit:
	printError("Build failed.")
	print(err)
	os.Exit(exitCodeBuildFailed)
}
