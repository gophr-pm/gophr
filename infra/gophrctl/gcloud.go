package main

import (
	"os/exec"
	"strconv"
	"strings"
)

const (
	gcloud = "gcloud"
	// TODO(skeswa): this should be configurable.
	gcloudVolumeZone = "us-central1-a"
)

type gCloudVolume struct {
	name string
	gigs int
	ssd  bool
}

func createGCloudVolumesIfNotExist(vols ...gCloudVolume) error {
	startSpinner("Checking if production volumes exist")
	output, err := exec.Command(gcloud, "compute", "disks", "list").CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return err
	}

	var (
		existentVolStr  = string(output[:])
		nonExistentVols []gCloudVolume
	)

	// Figure out what volumes don't exist yet.
	for _, vol := range vols {
		if strings.Index(existentVolStr, vol.name) == -1 {
			nonExistentVols = append(nonExistentVols, vol)
		}
	}
	stopSpinner(true)

	// Create the volumes that don't exist yet.
	if len(nonExistentVols) > 0 {
		startSpinner("Creating non-existent volumes")

		for _, nonExistentVol := range nonExistentVols {
			// Turn the disk type into a string.
			var diskType string
			if nonExistentVol.ssd {
				diskType = "pd-ssd"
			} else {
				diskType = "pd-standard"
			}

			// List out the args for volume.
			createVolumeArgs := []string{
				"compute",
				"disks",
				"create",
				nonExistentVol.name,
				"--zone=" + gcloudVolumeZone,
				"--size=" + strconv.Itoa(nonExistentVol.gigs) + "GB",
				"--type=" + diskType,
			}

			// Execute the command.
			if output, err = exec.Command(gcloud, createVolumeArgs...).CombinedOutput(); err != nil {
				stopSpinner(false)
				return newExecError(output, err)
			}
		}

		stopSpinner(true)
	}

	return nil
}
