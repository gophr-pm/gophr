package main

import (
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

var (
	webModuleID   = "web"
	webModuleDeps = []string{
		apiModuleID,
		webModuleID,
	}
	webModuleVolumePath     = "web/dist"
	webModuleVolumeMount    = "/usr/share/nginx/html"
	webModuleDockerfilePath = "infra/images/dev/ws/Dockerfile"
	webModuleContainerPorts = []dockerPortMapping{
		{hostPort: 80, containerPort: 80},
		{hostPort: 443, containerPort: 443},
	}
	webModuleContainerLinks = []dockerLinkMapping{
		{moduleID: apiModuleID},
		{moduleID: routerModuleID},
	}
	webModuleContainerVolumes = []dockerVolumeMapping{
		{
			containerPath: webModuleVolumeMount,
			hostPathGenerator: func(repoPath string) string {
				return filepath.Join(repoPath, webModuleVolumePath)
			},
		},
	}
)

type webModule struct {
	baseModule // Extends base module.
}

func (m *webModule) deps() []string {
	return webModuleDeps
}

func (m *webModule) dockerfile() string {
	return webModuleDockerfilePath
}

func (m *webModule) containerMetadata() ([]dockerPortMapping, []dockerLinkMapping, []dockerVolumeMapping) {
	return webModuleContainerPorts, webModuleContainerLinks, webModuleContainerVolumes
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
