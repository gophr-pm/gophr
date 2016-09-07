package main

import "os/exec"

type dockerBuildArgs struct {
	gpi            string
	latest         bool
	imageTag       string
	imageName      string
	contextPath    string
	dockerfilePath string
}

const (
	dockerHubUser = "gophr"
)

func dockerBuild(args dockerBuildArgs) error {
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
	if len(args.gpi) > 0 {
		cmdArgs = append(cmdArgs, "-t", "gcr.io/"+args.gpi+"/"+imageIdentifier)
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

func dockerPush(gpi, imageName, imageTag string) error {
	imageIdentifier := imageName + ":" + imageTag
	cloudImagePath := "gcr.io/" + gpi + "/" + imageIdentifier
	startSpinner("Pushing " + imageIdentifier + " to " + cloudImagePath)

	output, err := exec.Command("gcloud", "docker", "push", cloudImagePath).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}
