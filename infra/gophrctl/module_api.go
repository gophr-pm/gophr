package main

import "gopkg.in/urfave/cli.v1"

var (
	apiModuleID   = "api"
	apiModuleDeps = []string{
		dbModuleID,
	}
	apiModuleVolumeMount    = "/go/src/github.com/skeswa/gophr"
	apiModuleDockerfilePath = "infra/images/dev/api/Dockerfile"
	apiModuleContainerPorts = []dockerPortMapping{
		{hostPort: 3000, containerPort: 3000},
	}
	apiModuleContainerLinks = []dockerLinkMapping{
		{moduleID: dbModuleID},
	}
	apiModuleContainerVolumes = []dockerVolumeMapping{
		{
			containerPath: apiModuleVolumeMount,
			hostPathGenerator: func(repoPath string) string {
				return repoPath
			},
		},
	}
)

type apiModule struct {
	baseModule // Extends base module.
}

func (m *apiModule) deps() []string {
	return apiModuleDeps
}

func (m *apiModule) dockerfile() string {
	return apiModuleDockerfilePath
}

func (m *apiModule) containerMetadata() ([]dockerPortMapping, []dockerLinkMapping, []dockerVolumeMapping) {
	return apiModuleContainerPorts, apiModuleContainerLinks, apiModuleContainerVolumes
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
