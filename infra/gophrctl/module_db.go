package main

import "gopkg.in/urfave/cli.v1"

var (
	dbModuleID             = "db"
	dbModuleDeps           = []string{}
	dbModuleVolumeName     = "cassandra-db-volume"
	dbModuleVolumeMount    = "/var/lib/cassandra"
	dbModuleDockerfilePath = "infra/images/dev/db/Dockerfile"
	dbModuleContainerPorts = []dockerPortMapping{
		{hostPort: 7000, containerPort: 7000},
		{hostPort: 7001, containerPort: 7001},
		{hostPort: 9042, containerPort: 9042},
	}
	dbModuleContainerVolumes = []dockerVolumeMapping{
		{
			containerPath: dbModuleVolumeMount,
			hostPathGenerator: func(repoPath string) string {
				return dbModuleVolumeName
			},
		},
	}
)

type dbModule struct {
	baseModule // Extends base module.
}

func (m *dbModule) deps() []string {
	return dbModuleDeps
}

func (m *dbModule) dockerfile() string {
	return dbModuleDockerfilePath
}

func (m *dbModule) containerMetadata() ([]dockerPortMapping, []dockerLinkMapping, []dockerVolumeMapping) {
	return dbModuleContainerPorts, nil, dbModuleContainerVolumes
}

func (m *dbModule) stop(c *cli.Context) error {
	return nil
}

func (m *dbModule) log(c *cli.Context) error {
	return nil
}

func (m *dbModule) ssh(c *cli.Context) error {
	return nil
}

func (m *dbModule) test(c *cli.Context) error {
	return nil
}

func (m *dbModule) restart(c *cli.Context) error {
	return nil
}
