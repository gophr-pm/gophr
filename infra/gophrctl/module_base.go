package main

import "gopkg.in/urfave/cli.v1"

type baseModule struct {
	moduleID string
}

func (m *baseModule) id() string {
	return m.moduleID
}

func (m *baseModule) build(c *cli.Context, shallow bool) error {
	return doModuleBuild(
		m.moduleID,
		c.GlobalString(flagNameEnv) == envTypeDev,
		!shallow,
		c.GlobalString(flagNameRepoPath),
	)
}

func (m *baseModule) start(c *cli.Context, shallow bool) error {
	return doModuleStart(
		m.moduleID,
		c.GlobalString(flagNameEnv) == envTypeDev,
		!shallow,
		c.GlobalString(flagNameRepoPath),
		!c.Bool(flagNameForeground),
	)
}
