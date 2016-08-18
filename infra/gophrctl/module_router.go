package main

import "gopkg.in/urfave/cli.v1"

var (
	routerModuleID   = "router"
	routerModuleDeps = []string{
		dbModuleID,
	}
	routerModuleVolumeMount    = "/go/src/github.com/skeswa/gophr"
	routerModuleDockerfilePath = "infra/images/dev/router/Dockerfile"
	routerModuleContainerPorts = []dockerPortMapping{
		{hostPort: 3000, containerPort: 3000},
	}
	routerModuleContainerLinks = []dockerLinkMapping{
		{moduleID: dbModuleID},
	}
	routerModuleContainerVolumes = []dockerVolumeMapping{
		{
			containerPath: routerModuleVolumeMount,
			hostPathGenerator: func(repoPath string) string {
				return repoPath
			},
		},
	}
)

type routerModule struct {
	baseModule // Extends base module.
}

func (m *routerModule) deps() []string {
	return routerModuleDeps
}

func (m *routerModule) dockerfile() string {
	return routerModuleDockerfilePath
}

func (m *routerModule) containerMetadata() ([]dockerPortMapping, []dockerLinkMapping, []dockerVolumeMapping) {
	return routerModuleContainerPorts, routerModuleContainerLinks, routerModuleContainerVolumes
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
