package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"gopkg.in/urfave/cli.v1"
)

const (
	kubectl          = "kubectl"
	k8sSecretsName   = "gophr-secrets"
	k8sNamespaceFlag = "--namespace=gophr"
)

func readK8SProdContext(c *cli.Context) (string, error) {
	context := c.GlobalString(flagNameK8SProdContext)
	if len(context) < 1 {
		return context, errors.New("The kubernetes production context must be specified for this command to function.")
	}

	return context, nil
}

func switchK8SContext(newK8SContext string) (string, error) {
	// First, get the current context.
	output, err := exec.Command(kubectl, "config", "current-context").CombinedOutput()
	if err != nil {
		return "", err
	}

	// Save the old context, so that we can return it later.
	oldK8SContext := strings.TrimSpace(string(output[:]))

	// Switch to the new context.
	_, err = exec.Command(kubectl, "config", "use-context", newK8SContext).CombinedOutput()
	if err != nil {
		return "", err
	}

	return oldK8SContext, nil
}

func existsInK8S(k8sConfigFilePath string) bool {
	_, err := exec.Command(kubectl, k8sNamespaceFlag, "describe", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		return false
	}

	return true
}

func applyInK8S(k8sConfigFilePath string) error {
	startSpinner(fmt.Sprintf("Applying \"%s\" in kubernetes", abbreviateK8SPath(k8sConfigFilePath)))
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "apply", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func createInK8S(k8sConfigFilePath string) error {
	startSpinner(fmt.Sprintf("Creating \"%s\" in kubernetes", abbreviateK8SPath(k8sConfigFilePath)))
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "create", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func execK8SLogs(podName string, follow bool) {
	// Find kubectl - panic if it ain't here.
	binary, err := exec.LookPath(kubectl)
	if err != nil {
		panic(err)
	}

	if follow {
		syscall.Exec(binary, []string{kubectl, k8sNamespaceFlag, "logs", "-f", podName}, os.Environ())
	} else {
		syscall.Exec(binary, []string{kubectl, k8sNamespaceFlag, "logs", podName}, os.Environ())
	}
}

func execK8SBash(podName string) {
	// Find kubectl - panic if it ain't here.
	binary, err := exec.LookPath(kubectl)
	if err != nil {
		panic(err)
	}

	syscall.Exec(binary, []string{kubectl, k8sNamespaceFlag, "exec", podName, "-i", "-t", "/bin/bash"}, os.Environ())
}

func deleteInK8S(k8sConfigFilePath string) error {
	startSpinner(fmt.Sprintf("Deleting \"%s\" in kubernetes", abbreviateK8SPath(k8sConfigFilePath)))
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "delete", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func filterK8SPods(moduleName string) ([]string, error) {
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "get", "pods", "--selector=module="+moduleName, "--output=jsonpath={.items..metadata.name}").CombinedOutput()
	if err != nil {
		return nil, newExecError(output, err)
	}

	return strings.Split(strings.Trim(string(output[:]), "\t\n "), " "), nil
}

func abbreviateK8SPath(k8sPath string) string {
	sep := string(os.PathSeparator)
	parts := strings.Split(k8sPath, sep)
	return strings.Join(parts[len(parts)-2:], sep)
}

func secretExistsInK8S() bool {
	_, err := exec.Command(kubectl, k8sNamespaceFlag, "describe", "secret", k8sSecretsName).CombinedOutput()
	if err != nil {
		return false
	}

	return true
}

func createSecretsInK8S(secretFilePaths []string) error {
	startSpinner("Creating secrets in kubernetes")
	args := []string{
		k8sNamespaceFlag,
		"create",
		"secret",
		"generic",
		k8sSecretsName,
	}

	for _, secretFilePath := range secretFilePaths {
		args = append(args, "--from-file="+secretFilePath)
	}

	output, err := exec.Command(kubectl, args...).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func deleteSecretsInK8S() error {
	startSpinner("Deleting secrets in kubernetes")
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "delete", "secret", k8sSecretsName).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}
