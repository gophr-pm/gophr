package main

import (
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

type module interface {
	id() string
	deps() []string
	dockerfile() string
	containerMetadata() ([]dockerPortMapping, []dockerLinkMapping, []dockerVolumeMapping)

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
	apiModuleID:     &apiModule{baseModule{apiModuleID}},
	dbModuleID:      &dbModule{baseModule{dbModuleID}},
	indexerModuleID: &indexerModule{baseModule{indexerModuleID}},
	routerModuleID:  &routerModule{baseModule{routerModuleID}},
	webModuleID:     &webModule{baseModule{webModuleID}},
}

func doModuleBuild(
	moduleID string,
	targetDev bool,
	exitOnError bool,
	workDir string,
) error {
	printInfo("Building", moduleID+".")

	// Perform the docker build.
	startSpinner("Executing docker build...")
	err := doDockerBuild(
		workDir,
		filepath.Join(workDir, modules[moduleID].dockerfile()),
		dockerImageNameOf(moduleID),
		dockerDevImageTag,
	)
	stopSpinner()

	// Report on results.
	if err != nil {
		printError("Failed to build", moduleID+":")
		print(err)

		// Only exit if necessary.
		if exitOnError {
			os.Exit(exitCodeBuildFailed)
		}
	} else {
		printSuccess("Built", moduleID, "successfully.")
	}

	return nil
}

func doModuleStart(
	moduleID string,
	targetDev bool,
	exitOnError bool,
	workDir string,
	backgrounded bool,
) error {
	printInfo("Starting", moduleID+".")

	// Localize container metadata.
	ports, links, volumes := modules[moduleID].containerMetadata()

	// Perform the docker build.
	if backgrounded {
		startSpinner("Executing docker run...")
	}
	err := doDockerRun(
		workDir,
		dockerImageNameOf(moduleID),
		dockerDevImageTag,
		dockerContainerNameOf(moduleID),
		backgrounded,
		ports,
		links,
		volumes,
	)
	if backgrounded {
		stopSpinner()
	}

	// Report on results.
	if err != nil {
		printError("Failed to start", moduleID+":")
		print(err)

		// Only exit if necessary.
		if exitOnError {
			os.Exit(exitCodeStartFailed)
		}
	} else {
		printSuccess("Started", moduleID, "successfully.")
	}

	return nil
}
