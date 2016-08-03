package main

import "gopkg.in/urfave/cli.v1"

var (
	dbModuleID   = "db"
	dbModuleDeps = []string{}
)

type dbModule struct{}

func (m *dbModule) id() string {
	return dbModuleID
}

func (m *dbModule) deps() []string {
	return dbModuleDeps
}

func (m *dbModule) build(c *cli.Context, shallow bool) error {
	printInfo("Building", dbModuleID+"...")
	printSuccess("Built", dbModuleID, "successfully.")
	return nil
}

func (m *dbModule) start(c *cli.Context, shallow bool) error {
	printInfo("Starting", dbModuleID)
	return nil
}

func (m *dbModule) stop(c *cli.Context, shallow bool) error {
	return nil
}

func (m *dbModule) log(c *cli.Context, shallow bool) error {
	return nil
}

func (m *dbModule) ssh(c *cli.Context, shallow bool) error {
	return nil
}

func (m *dbModule) test(c *cli.Context, shallow bool) error {
	return nil
}

func (m *dbModule) restart(c *cli.Context, shallow bool) error {
	return nil
}
