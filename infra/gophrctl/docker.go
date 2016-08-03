package main

import (
	"errors"
	"os/exec"
)

const (
	dockerDevImageTag = "dev"
)

func doDockerBuild(
	workDir string,
	dockerfilePath string,
	imageName string,
	imageTag string,
) error {
	// Create the command.
	command := exec.Command("docker", "build", "-f", dockerfilePath, "-t", (imageName + ":" + imageTag), workDir)
	command.Dir = workDir

	// Execute the command.
	output, err := command.CombinedOutput()
	if err != nil || !command.ProcessState.Success() {
		return errors.New(string(output))
	}

	return nil
}
