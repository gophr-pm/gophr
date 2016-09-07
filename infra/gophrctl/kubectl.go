package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"gopkg.in/urfave/cli.v1"
)

const (
	kubectl          = "kubectl"
	k8sDevContext    = "minikube"
	k8sSecretsName   = "gophr-secrets"
	k8sNamespaceFlag = "--namespace=gophr"
)

var (
	prodK8SImageURLRegex = regexp.MustCompile(`gcr\.io/([a-zA-Z0-9\-]+)/([a-zA-Z0-9\-:\.]+)`)
)

func readK8SProdContext(c *cli.Context) (string, error) {
	context := c.GlobalString(flagNameK8SProdContext)
	if len(context) < 1 {
		return context, errors.New("The kubernetes production context must be specified for this command to function.")
	}

	return context, nil
}

// Returns the old kubernetes context, whether the context needs to be switched
// back, and the error.
func switchK8SContext(newK8SContext string) (string, bool, error) {
	startSpinner(fmt.Sprintf("Switching to the \"%s\" kubernetes context", newK8SContext))

	// First, get the current context.
	output, err := exec.Command(kubectl, "config", "current-context").CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return "", false, err
	}

	// Save the old context, so that we can return it later.
	oldK8SContext := strings.TrimSpace(string(output[:]))

	// If the k8s context is already switched, then return.
	if newK8SContext == oldK8SContext {
		return oldK8SContext, false, nil
	}

	// Switch to the new context.
	_, err = exec.Command(kubectl, "config", "use-context", newK8SContext).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return "", false, err
	}

	stopSpinner(true)
	return oldK8SContext, true, nil
}

func runInK8S(c *cli.Context, fn func() error) error {
	var (
		err                      error
		oldK8SContext            string
		mustSwitchK8SContextBack bool

		env        = readEnvironment(c)
		k8sContext = k8sDevContext
	)

	// If the environment is prod, change the kubernetes context accordingly.
	if env == environmentProd {
		// Read the production context before continuing.
		if k8sContext, err = readK8SProdContext(c); err != nil {
			return err
		}
	}

	// Switch the kubernetes context before continuing.
	if oldK8SContext, mustSwitchK8SContextBack, err = switchK8SContext(k8sContext); err != nil {
		return err
	}

	// Execute fn now that the context has been switched.
	if err = fn(); err != nil {
		// Before returning with an error, return the context back to where it was.
		if mustSwitchK8SContextBack {
			if _, _, switchErr := switchK8SContext(oldK8SContext); switchErr != nil {
				printError("Failed to reset the kubernetes context:", switchErr)
			}
		}

		return err
	}

	// Switch the context back, error out if there was a problem switching.
	if _, _, err := switchK8SContext(oldK8SContext); err != nil {
		return err
	}

	return nil
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

func updateProdK8SFileImage(newImageURL, k8sfilePath string) error {
	versionfileData, err := ioutil.ReadFile(k8sfilePath)
	if err != nil {
		return err
	}

	updatedVersionfileData := prodK8SImageURLRegex.ReplaceAll(versionfileData, []byte(newImageURL))
	if err = ioutil.WriteFile(k8sfilePath, updatedVersionfileData, 0644); err != nil {
		return err
	}

	return nil
}
