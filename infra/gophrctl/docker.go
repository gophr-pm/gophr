package main

import "os/exec"

type localDockerBuildArgs struct {
	latest         bool
	imageTag       string
	imageName      string
	dockerhub      bool
	contextPath    string
	dockerfilePath string
}

const (
	dockerHubUser = "gophr"
)

func localDockerBuild(args localDockerBuildArgs) error {
	imageIdentifier := args.imageName + ":" + args.imageTag
	startSpinner("Building " + imageIdentifier)

	cmdArgs := []string{
		"build",
		"-f",
		args.dockerfilePath,
		"--rm",
		"-t",
		imageIdentifier,
	}

	if args.latest {
		cmdArgs = append(cmdArgs, "-t", args.imageName+":latest")
	}
	if args.dockerhub {
		cmdArgs = append(cmdArgs, "-t", dockerHubUser+"/"+imageIdentifier)
	}

	cmdArgs = append(cmdArgs, ".")
	dockerCmd := exec.Command("docker", cmdArgs...)
	dockerCmd.Dir = args.contextPath

	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func localDockerPush(imageName, imageTag string) error {
	imageIdentifier := imageName + ":" + imageTag
	startSpinner("Pushing " + imageIdentifier)

	output, err := exec.Command("docker", "push", dockerHubUser+"/"+imageIdentifier).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}
