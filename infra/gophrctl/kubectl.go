package main

import "os/exec"

const (
	k8sNamespace = "gophr"
)

func createInK8S(k8sConfigFilePath string) error {
	output, err := exec.Command("kubectl", "--namespace="+k8sNamespace, "create", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		return newExecError(output, err)
	}

	return nil
}

func deleteInK8S(k8sConfigFilePath string) error {
	output, err := exec.Command("kubectl", "--namespace="+k8sNamespace, "delete", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		return newExecError(output, err)
	}

	return nil
}
