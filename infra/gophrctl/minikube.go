package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	envVarDockerPrefix    = "DOCKER"
	minikubeStatusStopped = "Stopped"
)

var (
	minikubeDockerEnvVarRegex = regexp.MustCompile("export ([^=]+)=\"([^\"]+)\"")
)

func isMinikubeRunning() (bool, error) {
	output, err := exec.Command("minikube", "status").CombinedOutput()
	if err != nil {
		return false, newExecError(output, err)
	}

	return strings.Index(string(output[:]), minikubeStatusStopped) == -1, nil
}

func startMinikube() error {
	startSpinner("Starting minikube...")

	output, err := exec.Command("minikube", "start").CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func assertMinikubeRunning() error {
	started, err := isMinikubeRunning()
	if err != nil {
		return err
	}

	if !started {
		err = startMinikube()
		if err != nil {
			return err
		}
	}

	return nil
}

func getMinikubeDockerEnv() ([]string, error) {
	output, err := exec.Command("minikube", "docker-env").CombinedOutput()
	if err != nil {
		return nil, newExecError(output, err)
	}

	var env []string
	for _, envVar := range os.Environ() {
		if !strings.HasPrefix(envVar, envVarDockerPrefix) {
			env = append(env, envVar)
		}
	}

	submatches := minikubeDockerEnvVarRegex.FindAllSubmatch(output, -1)
	for _, submatch := range submatches {
		env = append(env, string(submatch[1][:])+"="+string(submatch[2][:]))
	}

	return env, nil
}

type buildInMinikubeArgs struct {
	imageTag       string
	imageName      string
	contextPath    string
	dockerfilePath string
}

func buildInMinikube(args buildInMinikubeArgs) error {
	imageIdentifier := args.imageName + ":" + args.imageTag
	startSpinner("Building " + imageIdentifier + "...")

	dockerEnv, err := getMinikubeDockerEnv()
	if err != nil {
		stopSpinner(false)
		return err
	}

	dockerCmd := exec.Command("docker", "build", "-f", args.dockerfilePath, "--rm", "-t", imageIdentifier, ".")
	dockerCmd.Dir = args.contextPath
	dockerCmd.Env = dockerEnv

	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}
