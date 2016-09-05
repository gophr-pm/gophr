package main

import "bytes"

type module struct {
	name         string
	k8sfiles     []string
	dockerfile   string
	buildContext string
}

var modules = map[string]*module{
	"api": &module{
		name: "api",
		k8sfiles: []string{
			"./infra/k8s/api/service",
			"./infra/k8s/api/controller",
		},
		dockerfile:   "./infra/docker/api/Dockerfile",
		buildContext: ".",
	},
	"db": &module{
		name: "db",
		k8sfiles: []string{
			"./infra/k8s/db/service",
			"./infra/k8s/db/daemonset",
		},
		dockerfile:   "./infra/docker/db/Dockerfile",
		buildContext: ".",
	},
	"indexer": &module{
		name: "indexer",
		k8sfiles: []string{
			"./infra/k8s/indexer/controller",
		},
		dockerfile:   "./infra/docker/indexer/Dockerfile",
		buildContext: ".",
	},
	"migrator": &module{
		name: "migrator",
		k8sfiles: []string{
			"./infra/k8s/migrator/pod",
		},
		dockerfile:   "./infra/docker/migrator/Dockerfile",
		buildContext: ".",
	},
	"router": &module{
		name: "router",
		k8sfiles: []string{
			"./infra/k8s/router/service",
			"./infra/k8s/router/controller",
		},
		dockerfile:   "./infra/docker/router/Dockerfile",
		buildContext: ".",
	},
	"web": &module{
		name: "web",
		k8sfiles: []string{
			"./infra/k8s/web/service",
			"./infra/k8s/web/controller",
		},
		dockerfile:   "./infra/docker/web/Dockerfile",
		buildContext: ".",
	},
}

func modulesToString() string {
	var (
		buffer        bytes.Buffer
		isFirstModule = true
	)

	for moduleName := range modules {
		if !isFirstModule {
			buffer.WriteString(", ")
		} else {
			isFirstModule = false
		}

		buffer.WriteString(moduleName)
	}

	return buffer.String()
}
