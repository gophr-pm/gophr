package main

import (
	"errors"
	"os/exec"
	"path/filepath"
)

const (
	dockerDevImageTag = "dev"
)

func executeDockerComposeCommand(
	repoDir string,
	module string,
	isDev bool,
	cmd string,
	params ...string,
) error {
	var (
		command         exec.Command
		commandArgs     []string
		composeFilePath string
	)

	// Locate the compose file.
	if isDev {
		composeFilePath = filepath.Join(repoDir, "infra/docker-compose.dev.yml")
	} else {
		composeFilePath = filepath.Join(repoDir, "infra/docker-compose.prod.yml")
	}

	// Hack together the arguments.

	// Create the command.
	if len(module) > 0 {
		// If the module is explicitly specified, make normal call to compose.
		command = exec.Command("docker-compose", "-f", composeFilePath, cmd, module)
	} else {
		// Exclude
	}
	command = exec.Command("docker-compose", "-f", composeFilePath, cmd, module)
	command.Dir = repoDir

	// Execute the command.
	output, err := command.CombinedOutput()
	if err != nil || !command.ProcessState.Success() {
		return errors.New(string(output))
	}

	return nil
}
