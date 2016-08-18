package main

import "gopkg.in/urfave/cli.v1"

var (
	indexerModuleID   = "indexer"
	indexerModuleDeps = []string{
		dbModuleID,
	}
	indexerModuleVolumeMount    = "/go/src/github.com/skeswa/gophr"
	indexerModuleDockerfilePath = "infra/images/dev/indexer/Dockerfile"
	indexerModuleContainerPorts = []dockerPortMapping{
		{hostPort: 3000, containerPort: 3000},
	}
	indexerModuleContainerLinks = []dockerLinkMapping{
		{moduleID: dbModuleID},
	}
	indexerModuleContainerVolumes = []dockerVolumeMapping{
		{
			containerPath: indexerModuleVolumeMount,
			hostPathGenerator: func(repoPath string) string {
				return repoPath
			},
		},
	}
)

type indexerModule struct {
	baseModule // Extends base module.
}

func (m *indexerModule) deps() []string {
	return indexerModuleDeps
}

func (m *indexerModule) dockerfile() string {
	return indexerModuleDockerfilePath
}

func (m *indexerModule) containerMetadata() ([]dockerPortMapping, []dockerLinkMapping, []dockerVolumeMapping) {
	return indexerModuleContainerPorts, indexerModuleContainerLinks, indexerModuleContainerVolumes
}

func (m *indexerModule) stop(c *cli.Context) error {
	return nil
}

func (m *indexerModule) log(c *cli.Context) error {
	return nil
}

func (m *indexerModule) ssh(c *cli.Context) error {
	return nil
}

func (m *indexerModule) test(c *cli.Context) error {
	return nil
}

func (m *indexerModule) restart(c *cli.Context) error {
	return nil
}
