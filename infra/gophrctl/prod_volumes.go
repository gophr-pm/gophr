package main

const (
	gophrVolumePrefix = "gophr-volume-"
	dbVolumeCapacity  = 120 // In gb.
)

var (
	prodVolumes = []gCloudVolume{
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
	}
)

func assertProdVolumes() error {
	return createGCloudVolumesIfNotExist(prodVolumes...)
}
