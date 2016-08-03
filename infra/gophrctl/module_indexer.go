package main

import (
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

var (
	indexerModuleID   = "indexer"
	indexerModuleDeps = []string{
		dbModuleID,
	}
	indexerDockerfilePath = "infra/images/dev/indexer/Dockerfile"
)

type indexerModule struct{}

func (m *indexerModule) id() string {
	return indexerModuleID
}

func (m *indexerModule) deps() []string {
	return indexerModuleDeps
}

func (m *indexerModule) build(c *cli.Context, shallow bool) error {
	// Create parameters.
	workDir := c.GlobalString(flagNameRepoPath)
	targetDev := c.GlobalString(flagNameEnv) == envTypeDev
	recursive := !shallow
	dockerfilePath := filepath.Join(workDir, indexerDockerfilePath)

	// Perform the operation.
	doModuleBuild(indexerModuleID, targetDev, recursive, workDir, dockerfilePath)

	return nil
}

func (m *indexerModule) start(c *cli.Context, shallow bool) error {
	return nil
}

func (m *indexerModule) stop(c *cli.Context, shallow bool) error {
	return nil
}

func (m *indexerModule) log(c *cli.Context, shallow bool) error {
	return nil
}

func (m *indexerModule) ssh(c *cli.Context, shallow bool) error {
	return nil
}

func (m *indexerModule) test(c *cli.Context, shallow bool) error {
	return nil
}

func (m *indexerModule) restart(c *cli.Context, shallow bool) error {
	return nil
}
