package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	k8sNamespace = "gophr"
)

func existsInK8S(k8sConfigFilePath string) bool {
	output, err := exec.Command("kubectl", "--namespace="+k8sNamespace, "describe", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		fmt.Println("FAIL", string(output[:]), err)
		return false
	}

	return true
}

func applyInK8S(k8sConfigFilePath string) error {
	startSpinner(fmt.Sprintf("Applying \"%s\" in kubernetes...", abbreviateK8SPath(k8sConfigFilePath)))
	output, err := exec.Command("kubectl", "--namespace="+k8sNamespace, "apply", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func createInK8S(k8sConfigFilePath string) error {
	startSpinner(fmt.Sprintf("Creating \"%s\" in kubernetes...", abbreviateK8SPath(k8sConfigFilePath)))
	output, err := exec.Command("kubectl", "--namespace="+k8sNamespace, "create", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func deleteInK8S(k8sConfigFilePath string) error {
	startSpinner(fmt.Sprintf("Deleting \"%s\" in kubernetes...", abbreviateK8SPath(k8sConfigFilePath)))
	output, err := exec.Command("kubectl", "--namespace="+k8sNamespace, "delete", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return newExecError(output, err)
	}

	stopSpinner(true)
	return nil
}

func abbreviateK8SPath(k8sPath string) string {
	sep := string(os.PathSeparator)
	parts := strings.Split(k8sPath, sep)
	return strings.Join(parts[len(parts)-2:], sep)
}
