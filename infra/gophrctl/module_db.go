package main

import (
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

var (
	dbModuleID       = "db"
	dbModuleDeps     = []string{}
	dbDockerfilePath = "infra/images/dev/db/Dockerfile"
)

type dbModule struct{}

func (m *dbModule) id() string {
	return dbModuleID
}

func (m *dbModule) deps() []string {
	return dbModuleDeps
}

func (m *dbModule) build(c *cli.Context, shallow bool) error {
	// Create parameters.
	workDir := c.GlobalString(flagNameRepoPath)
	targetDev := c.GlobalString(flagNameEnv) == envTypeDev
	recursive := !shallow
	dockerfilePath := filepath.Join(workDir, dbDockerfilePath)

	// Perform the operation.
	doModuleBuild(dbModuleID, targetDev, recursive, workDir, dockerfilePath)

	return nil
}

func (m *dbModule) start(c *cli.Context, shallow bool) error {
	printInfo("Starting", dbModuleID)
	return nil
}

func (m *dbModule) stop(c *cli.Context, shallow bool) error {
	return nil
}

func (m *dbModule) log(c *cli.Context, shallow bool) error {
	return nil
}

func (m *dbModule) ssh(c *cli.Context, shallow bool) error {
	return nil
}

func (m *dbModule) test(c *cli.Context, shallow bool) error {
	return nil
}

func (m *dbModule) restart(c *cli.Context, shallow bool) error {
	return nil
}
