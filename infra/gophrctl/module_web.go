package main

import "gopkg.in/urfave/cli.v1"

var (
	webModuleID   = "web"
	webModuleDeps = []string{
		apiModuleID,
		routerModuleID,
	}
)

type webModule struct{}

func (m *webModule) id() string {
	return webModuleID
}

func (m *webModule) deps() []string {
	return webModuleDeps
}

func (m *webModule) build(c *cli.Context, shallow bool) error {
	printInfo("Building", webModuleID+"...")
	printSuccess("Built", webModuleID, "successfully.")
	return nil
}

func (m *webModule) start(c *cli.Context, shallow bool) error {
	return nil
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
