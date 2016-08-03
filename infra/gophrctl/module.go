package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

type module interface {
	id() string
	deps() []string

	build(*cli.Context, bool) error
	start(*cli.Context, bool) error
	stop(*cli.Context, bool) error
	log(*cli.Context, bool) error
	ssh(*cli.Context, bool) error
	test(*cli.Context, bool) error
	restart(*cli.Context, bool) error
}

var modules = map[string]module{
	allModuleID:     &allModule{},
	apiModuleID:     &apiModule{},
	dbModuleID:      &dbModule{},
	indexerModuleID: &indexerModule{},
	routerModuleID:  &routerModule{},
	webModuleID:     &webModule{},
}

func doModuleBuild(
	moduleID string,
	targetDev bool,
	recursive bool,
	workDir string,
	dockerfilePath string,
) {
	printInfo("Building", moduleID+".")

	// Perform the docker build.
	startSpinner("Running docker build...")
	err := doDockerBuild(workDir, dockerfilePath, dockerImageNameOf(moduleID), dockerDevImageTag)
	stopSpinner()

	// Report on results.
	if err != nil {
		printError("Failed to build", moduleID+":")
		print(err)

		// Only exit if recursive.
		if recursive {
			os.Exit(exitCodeBuildFailed)
		}
	} else {
		printSuccess("Built", moduleID, "successfully.")
	}
}
