package main

import (
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

var (
	apiModuleID   = "api"
	apiModuleDeps = []string{
		dbModuleID,
	}
	apiDockerfilePath = "infra/images/dev/api/Dockerfile"
)

type apiModule struct{}

func (m *apiModule) id() string {
	return apiModuleID
}

func (m *apiModule) deps() []string {
	return apiModuleDeps
}

func (m *apiModule) build(c *cli.Context, shallow bool) error {
	printInfo("Building", apiModuleID+".")

	// Localize requisite variables.
	workDir := c.GlobalString(flagNameRepoPath)
	dockerfilePath := filepath.Join(workDir, apiDockerfilePath)

	// Perform the docker build.
	startSpinner("Running docker build...")
	err := doDockerBuild(workDir, dockerfilePath, dockerImageNameOf(apiModuleID), dockerDevImageTag)
	stopSpinner()

	// Report on results.
	if err != nil {
		printError("Failed to build", apiModuleID+":")
		print(err)

		// Only exit if not shallow.
		os.Exit(1)
	} else {
		printSuccess("Built", apiModuleID, "successfully.")
	}

	return nil
}

func (m *apiModule) start(c *cli.Context, shallow bool) error {
	return nil
}

func (m *apiModule) stop(c *cli.Context, shallow bool) error {
	return nil
}

func (m *apiModule) log(c *cli.Context, shallow bool) error {
	return nil
}

func (m *apiModule) ssh(c *cli.Context, shallow bool) error {
	return nil
}

func (m *apiModule) test(c *cli.Context, shallow bool) error {
	return nil
}

func (m *apiModule) restart(c *cli.Context, shallow bool) error {
	return nil
}
