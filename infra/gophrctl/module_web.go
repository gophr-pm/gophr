package main

import (
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

var (
	webModuleID   = "web"
	webModuleDeps = []string{
		apiModuleID,
		routerModuleID,
	}
	webDockerfilePath = "infra/images/dev/ws/Dockerfile"
)

type webModule struct{}

func (m *webModule) id() string {
	return webModuleID
}

func (m *webModule) deps() []string {
	return webModuleDeps
}

func (m *webModule) build(c *cli.Context, shallow bool) error {
	// Create parameters.
	workDir := c.GlobalString(flagNameRepoPath)
	targetDev := c.GlobalString(flagNameEnv) == envTypeDev
	recursive := !shallow
	dockerfilePath := filepath.Join(workDir, webDockerfilePath)

	// Perform the operation.
	doModuleBuild(webModuleID, targetDev, recursive, workDir, dockerfilePath)

	return nil
}

func (m *webModule) start(c *cli.Context, shallow bool) error {
	return nil
}

func (m *webModule) stop(c *cli.Context, shallow bool) error {
	return nil
}

func (m *webModule) log(c *cli.Context, shallow bool) error {
	return nil
}

func (m *webModule) ssh(c *cli.Context, shallow bool) error {
	return nil
}

func (m *webModule) test(c *cli.Context, shallow bool) error {
	return nil
}

func (m *webModule) restart(c *cli.Context, shallow bool) error {
	return nil
}
