package main

import "os/exec"

type localDockerBuildArgs struct {
	latest         bool
	imageTag       string
	imageName      string
	contextPath    string
	dockerfilePath string
}

func localDockerBuild(args localDockerBuildArgs) error {
	imageIdentifier := args.imageName + ":" + args.imageTag
	startSpinner("Building " + imageIdentifier)

	var dockerCmd *exec.Cmd
	if args.latest {
		dockerCmd = exec.Command("docker", "build", "-f", args.dockerfilePath, "--rm", "-t", args.imageName+":latest", "-t", imageIdentifier, ".")
	} else {
		dockerCmd = exec.Command("docker", "build", "-f", args.dockerfilePath, "--rm", "-t", imageIdentifier, ".")
	}

	dockerCmd.Dir = args.contextPath

	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}
