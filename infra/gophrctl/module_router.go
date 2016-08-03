package main

import "gopkg.in/urfave/cli.v1"

var (
	routerModuleID   = "router"
	routerModuleDeps = []string{
		dbModuleID,
	}
)

type routerModule struct{}

func (m *routerModule) id() string {
	return routerModuleID
}

func (m *routerModule) deps() []string {
	return routerModuleDeps
}

func (m *routerModule) build(c *cli.Context, shallow bool) error {
	printInfo("Building", routerModuleID+"...")
	printSuccess("Built", routerModuleID, "successfully.")
	return nil
}

func (m *routerModule) start(c *cli.Context, shallow bool) error {
	return nil
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
