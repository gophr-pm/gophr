package main

import "bytes"

const (
	gophrVolumePrefix   = "gophr-volume-"
	dbVolumeCapacity    = 120 // In gb.
	depotVolumeCapacity = 600 // In gb.
	depotVolumeName     = "gophr-volume-depot"
)

type module struct {
	name         string
	dockerfile   string
	prodVolumes  []gCloudVolume
	versionfile  string
	devK8SFiles  []string
	prodK8SFiles []string
	buildContext string
}

var modules = map[string]*module{
	"db": &module{
		name: "db",
		devK8SFiles: []string{
			"./infra/k8s/db/service.dev.yml",
			"./infra/k8s/db/controllers.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/db/service.prod.yml",
			"./infra/k8s/db/controllers.prod.yml",
		},
		dockerfile: "./infra/docker/db/Dockerfile",
		prodVolumes: []gCloudVolume{
			gCloudVolume{
				name: gophrVolumePrefix + "db-a",
				gigs: dbVolumeCapacity,
				ssd:  true,
			},
			gCloudVolume{
				name: gophrVolumePrefix + "db-b",
				gigs: dbVolumeCapacity,
				ssd:  true,
			},
			gCloudVolume{
				name: gophrVolumePrefix + "db-c",
				gigs: dbVolumeCapacity,
				ssd:  true,
			},
		},
		versionfile:  "./infra/docker/db/Versionfile.prod",
		buildContext: ".",
	},
	"migrator": &module{
		name: "migrator",
		devK8SFiles: []string{
			"./infra/k8s/migrator/pod.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/migrator/pod.prod.yml",
		},
		dockerfile:   "./infra/docker/migrator/Dockerfile",
		versionfile:  "./infra/docker/migrator/Versionfile.prod",
		buildContext: ".",
	},
	"indexer": &module{
		name: "indexer",
		devK8SFiles: []string{
			"./infra/k8s/indexer/controller.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/indexer/controller.prod.yml",
		},
		dockerfile:   "./infra/docker/indexer/Dockerfile",
		versionfile:  "./infra/docker/indexer/Versionfile.prod",
		buildContext: ".",
	},
	"depot-vol": &module{
		name: "depot-vol",
		devK8SFiles: []string{
			"./infra/k8s/depot/volume/volume.dev.yml",
			"./infra/k8s/depot/volume/claim.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/depot/volume/service.prod.yml",
			"./infra/k8s/depot/volume/controller.prod.yml",
			"./infra/k8s/depot/volume/volume.prod.yml",
			"./infra/k8s/depot/volume/claim.prod.yml",
		},
		dockerfile: "./infra/docker/depot/volume/Dockerfile",
		prodVolumes: []gCloudVolume{
			gCloudVolume{
				name: depotVolumeName,
				gigs: depotVolumeCapacity,
				ssd:  false,
			},
		},
		versionfile:  "./infra/docker/depot/internal/Versionfile.prod",
		buildContext: ".",
	},
	"depot-int": &module{
		name: "depot-int",
		devK8SFiles: []string{
			"./infra/k8s/depot/internal/service.dev.yml",
			"./infra/k8s/depot/internal/controller.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/depot/internal/service.prod.yml",
			"./infra/k8s/depot/internal/controller.prod.yml",
		},
		dockerfile: "./infra/docker/depot/internal/Dockerfile",
		prodVolumes: []gCloudVolume{
			gCloudVolume{
				name: depotVolumeName,
				gigs: depotVolumeCapacity,
				ssd:  false,
			},
		},
		versionfile:  "./infra/docker/depot/internal/Versionfile.prod",
		buildContext: ".",
	},
	"depot-ext": &module{
		name: "depot-ext",
		devK8SFiles: []string{
			"./infra/k8s/depot/external/service.dev.yml",
			"./infra/k8s/depot/external/controller.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/depot/external/service.prod.yml",
			"./infra/k8s/depot/external/controller.prod.yml",
		},
		dockerfile: "./infra/docker/depot/external/Dockerfile",
		prodVolumes: []gCloudVolume{
			gCloudVolume{
				name: depotVolumeName,
				gigs: depotVolumeCapacity,
				ssd:  false,
			},
		},
		versionfile:  "./infra/docker/depot/external/Versionfile.prod",
		buildContext: ".",
	},
	"api": &module{
		name: "api",
		devK8SFiles: []string{
			"./infra/k8s/api/service.dev.yml",
			"./infra/k8s/api/controller.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/api/service.prod.yml",
			"./infra/k8s/api/controller.prod.yml",
		},
		dockerfile:   "./infra/docker/api/Dockerfile",
		versionfile:  "./infra/docker/api/Versionfile.prod",
		buildContext: ".",
	},
	"router": &module{
		name: "router",
		devK8SFiles: []string{
			"./infra/k8s/router/service.dev.yml",
			"./infra/k8s/router/controller.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/router/service.prod.yml",
			"./infra/k8s/router/controller.prod.yml",
		},
		dockerfile: "./infra/docker/router/Dockerfile",
		prodVolumes: []gCloudVolume{
			gCloudVolume{
				name: gophrVolumePrefix + "depot",
				gigs: dbVolumeCapacity,
				ssd:  false,
			},
		},
		versionfile:  "./infra/docker/router/Versionfile.prod",
		buildContext: ".",
	},
	"web": &module{
		name: "web",
		devK8SFiles: []string{
			"./infra/k8s/web/service.dev.yml",
			"./infra/k8s/web/controller.dev.yml",
		},
		prodK8SFiles: []string{
			"./infra/k8s/web/service.prod.yml",
			"./infra/k8s/web/controller.prod.yml",
		},
		dockerfile:   "./infra/docker/web/Dockerfile",
		versionfile:  "./infra/docker/web/Versionfile.prod",
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
