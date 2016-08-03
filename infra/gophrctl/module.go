package main

import "gopkg.in/urfave/cli.v1"

type module interface {
	id() string
	deps() []string

	build(*cli.Context, bool) error
	start(*cli.Context, bool) error
	stop(*cli.Context, bool) error
	log(*cli.Context, bool) error
	ssh(*cli.Context, bool) error
	test(*cli.Context, bool) error
	restart(*cli.Context, bool) error
}

var modules = map[string]module{
	allModuleID:     &allModule{},
	apiModuleID:     &apiModule{},
	dbModuleID:      &dbModule{},
	indexerModuleID: &indexerModule{},
	routerModuleID:  &routerModule{},
	webModuleID:     &webModule{},
}
