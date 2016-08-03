package main

import (
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

var (
	routerModuleID   = "router"
	routerModuleDeps = []string{
		dbModuleID,
	}
	routerDockerfilePath = "infra/images/dev/router/Dockerfile"
)

type routerModule struct{}

func (m *routerModule) id() string {
	return routerModuleID
}

func (m *routerModule) deps() []string {
	return routerModuleDeps
}

func (m *routerModule) build(c *cli.Context, shallow bool) error {
	// Create parameters.
	workDir := c.GlobalString(flagNameRepoPath)
	targetDev := c.GlobalString(flagNameEnv) == envTypeDev
	recursive := !shallow
	dockerfilePath := filepath.Join(workDir, routerDockerfilePath)

	// Perform the operation.
	doModuleBuild(routerModuleID, targetDev, recursive, workDir, dockerfilePath)

	return nil
}

func (m *routerModule) start(c *cli.Context, shallow bool) error {
	return nil
}

func (m *routerModule) stop(c *cli.Context, shallow bool) error {
	return nil
}

func (m *routerModule) log(c *cli.Context, shallow bool) error {
	return nil
}

func (m *routerModule) ssh(c *cli.Context, shallow bool) error {
	return nil
}

func (m *routerModule) test(c *cli.Context, shallow bool) error {
	return nil
}

func (m *routerModule) restart(c *cli.Context, shallow bool) error {
	return nil
}
