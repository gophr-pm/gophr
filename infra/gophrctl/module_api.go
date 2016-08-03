package main

import "gopkg.in/urfave/cli.v1"

var (
	apiModuleID   = "api"
	apiModuleDeps = []string{
		dbModuleID,
	}
)

type apiModule struct{}

func (m *apiModule) id() string {
	return apiModuleID
}

func (m *apiModule) deps() []string {
	return apiModuleDeps
}

func (m *apiModule) build(c *cli.Context, shallow bool) error {
	printInfo("Building", apiModuleID+"...")
	printSuccess("Built", apiModuleID, "successfully.")
	return nil
}

func (m *apiModule) start(c *cli.Context, shallow bool) error {
	return nil
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
