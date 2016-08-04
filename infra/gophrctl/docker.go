package main

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
)

const (
	dockerDevImageTag = "dev"
)

type dockerPortMapping struct {
	hostPort      int
	containerPort int
}

func (dpm dockerPortMapping) String() string {
	return fmt.Sprintf("%d:%d", dpm.hostPort, dpm.containerPort)
}

type dockerLinkMapping struct {
	moduleID string
	hostName string
}

func (dlm dockerLinkMapping) String() string {
	return fmt.Sprintf("%s:%s", dockerContainerNameOf(dlm.moduleID), dlm.hostName)
}

type dockerVolumeMapping struct {
	hostPath      int
	containerPath int
}

func (dvm dockerVolumeMapping) String() string {
	return fmt.Sprintf("%d:%d", dvm.hostPath, dvm.containerPath)
}

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

func doDockerRun(
	workDir string,
	imageName string,
	imageTag string,
	containerName string,
	backgrounded bool,
	portMappings []dockerPortMapping,
	linkMappings []dockerLinkMapping,
	volumeMappings []dockerVolumeMapping,
) error {
	// Compile the args to pass along to exec.
	args := []string{"run"}
	if backgrounded {
		args = append(args, "-d")
	}
	args = append(args, "--name", containerName)
	for _, volumeMapping := range volumeMappings {
		args = append(args, "-v", volumeMapping.String())
	}
	for _, linkMapping := range linkMappings {
		args = append(args, "--link", linkMapping.String())
	}
	for _, portMapping := range portMappings {
		args = append(args, "-p", portMapping.String())
	}
	args = append(args, (imageName + ":" + imageTag))

	// Create the command.
	command := exec.Command("docker", args...)
	command.Dir = workDir

	// Execute the command.
	if backgrounded {
		output, err := command.CombinedOutput()
		if err != nil || !command.ProcessState.Success() {
			return errors.New(string(output))
		}
	} else {
		// Need to pipe output since its running in the foreground.
		stdout, err := command.StdoutPipe()
		if err != nil {
			return err
		}

		// Start the command after having set up the pipe.
		if err := command.Start(); err != nil {
			return err
		}

		// Read command's stdout line by line.
		stdoutScanner := bufio.NewScanner(stdout)
		for stdoutScanner.Scan() {
			fmt.Println(stdoutScanner.Text())
		}

		// In the event of an error, exit promptly.
		if err := stdoutScanner.Err(); err != nil {
			return err
		}
	}

	return nil
}

func dockerImageNameOf(moduleID string) string {
	return "gophr-" + moduleID
}

func dockerContainerNameOf(moduleID string) string {
	return "gophr-" + moduleID
}
