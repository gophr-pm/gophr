package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"gopkg.in/urfave/cli.v1"
)

const (
	kubectl          = "kubectl"
	k8sNamespace     = "gophr"
	k8sDevContext    = "minikube"
	k8sSecretsName   = "gophr-secrets"
	k8sNamespaceFlag = "--namespace=gophr"
)

// K8SPod is a kubernetes pod.
type K8SPod struct {
	Status   K8SPodStatus   `json:"status"`
	Metadata K8SPodMetadata `json:"metadata"`
}

// K8SPodList is a list of kubernetes pods.
type K8SPodList struct {
	Pods []K8SPod `json:"items"`
}

// K8SPodStatus is the status of a kubernetes pod.
type K8SPodStatus struct {
	Phase string `json:"phase"`
}

// K8SPodMetadata is the metadata of a kubernetes pod.
type K8SPodMetadata struct {
	Name string `json:"name"`
}

var (
	prodK8SImageURLRegex   = regexp.MustCompile(`gcr\.io/([a-zA-Z0-9\-_{}]+)/([a-zA-Z0-9\-:\.]+)`)
	persistentK8SFileRegex = regexp.MustCompile(`(?:service|claim|volume)s?\.[a-z]+(\.template)?\.yml$`)
)

func isPersistentK8SResource(k8sfile string) bool {
	return persistentK8SFileRegex.MatchString(k8sfile)
}

func readK8SProdContext(c *cli.Context) (string, error) {
	context := c.GlobalString(flagNameK8SProdContext)
	if len(context) < 1 {
		return context, errors.New("The kubernetes production context must be specified for this command to function.")
	}

	return context, nil
}

// deleteK8STemplateFiles deletes any generated kubernetes config files.
func deleteGeneratedK8SFiles(k8sfilePaths []string) error {
	for _, k8sfilePath := range k8sfilePaths {
		if isTemplateK8SFile(k8sfilePath) {
			if err := os.Remove(k8sfilePath); err != nil {
				return err
			}
		}
	}

	return nil
}

// getModuleK8SFilePaths gets the appropriate kubernetes config files for the
// module depending on the environment. Also compiles all templates.
func getModuleK8SFilePaths(c *cli.Context, m *module) ([]string, error) {
	var (
		err       error
		paths     []string
		realPath  string
		gophrRoot string
		realPaths []string
	)

	if env := readEnvironment(c); env == environmentProd {
		paths = m.prodK8SFiles
	} else {
		paths = m.devK8SFiles
	}

	// Figure out where everything is.
	if gophrRoot, err = readGophrRoot(c); err != nil {
		return nil, err
	}

	// Compile any templates that may exist.
	for _, path := range paths {
		realPath = filepath.Join(gophrRoot, path)

		// If its a template, then compile it first.
		if isTemplateK8SFile(realPath) {
			if realPath, err = compileK8STemplateFile(c, realPath); err != nil {
				return nil, err
			}
		}

		realPaths = append(realPaths, realPath)
	}

	return realPaths, nil
}

// Returns the old kubernetes context, whether the context needs to be switched
// back, and the error.
func switchK8SContext(newK8SContext string, switchingBack bool) (string, bool, error) {
	// First, get the current context.
	output, err := exec.Command(kubectl, "config", "current-context").CombinedOutput()
	if err != nil {
		return "", false, err
	}

	// Save the old context, so that we can return it later.
	oldK8SContext := strings.TrimSpace(string(output[:]))

	// If the k8s context is already switched, then return.
	if newK8SContext == oldK8SContext {
		return oldK8SContext, false, nil
	}

	// Only say something if the context is changing.
	if switchingBack {
		startSpinner(fmt.Sprintf("Switching back to the \"%s\" kubernetes context", newK8SContext))
	} else {
		startSpinner(fmt.Sprintf("Switching to the \"%s\" kubernetes context", newK8SContext))
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
	if oldK8SContext, mustSwitchK8SContextBack, err = switchK8SContext(k8sContext, false); err != nil {
		return err
	} else if mustSwitchK8SContextBack {
		// If a message was printed, add a new line for padding.
		fmt.Println()
	}

	// Make sure we have a namespace before doing anything else.
	if err = assertNamespaceInK8S(); err != nil {
		return err
	}

	// Check if the secrets are here. If not, scream & shout.
	if c.Command.FullName() != "secrets cycle" && !secretsExistInK8S() {
		return errors.New("The gophr secrets have not been installed yet. Use `gophrctl secrets cycle` to correct that.")
	}

	// Execute fn now that the context has been switched.
	if err = fn(); err != nil {
		// Before returning with an error, return the context back to where it was.
		if mustSwitchK8SContextBack {
			if _, _, switchErr := switchK8SContext(oldK8SContext, true); switchErr != nil {
				printError("Failed to reset the kubernetes context:", switchErr)
			}
		}

		return err
	}

	// Switch the context back, error out if there was a problem switching.
	if mustSwitchK8SContextBack {
		fmt.Println()
		if _, _, err := switchK8SContext(oldK8SContext, true); err != nil {
			return err
		}
	}

	return nil
}

func existsInK8S(k8sConfigFilePath string) bool {
	_, err := exec.Command(kubectl, k8sNamespaceFlag, "describe", "--show-events=false", "-f", k8sConfigFilePath).CombinedOutput()
	if err != nil {
		return false
	}

	return true
}

func secretsExistInK8S() bool {
	if _, err := exec.Command(kubectl, k8sNamespaceFlag, "describe", "secret", "gophr-secrets").CombinedOutput(); err != nil {
		return false
	}

	return true
}

func assertNamespaceInK8S() error {
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "get", "namespaces", "--output=jsonpath={.items..metadata.name}").CombinedOutput()
	if err != nil {
		return err
	}

	// Loop through all the namespaces, and look for our namespace.
	namespaces := strings.Split(strings.TrimSpace(string(output[:])), " ")
	for _, namespace := range namespaces {
		if k8sNamespace == strings.TrimSpace(namespace) {
			return nil
		}
	}

	// If we're here then the namespace does not exist. Time to create it.
	startSpinner("Creating namespace in kubernetes")
	_, err = exec.Command(kubectl, k8sNamespaceFlag, "create", "namespace", k8sNamespace).CombinedOutput()
	if err != nil {
		stopSpinner(false)
		return err
	}

	stopSpinner(true)
	return nil
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

func getPodsInK8S() (string, error) {
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "get", "pods").CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output[:]), nil
}

func filterRunningK8SPods(moduleName string) ([]string, error) {
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "get", "pods", "--selector=module="+moduleName, "--output=json").CombinedOutput()
	if err != nil {
		return nil, newExecError(output, err)
	}

	podList := K8SPodList{}
	if err = json.Unmarshal(output, &podList); err != nil {
		return nil, newExecError(output, errors.New("Could not read pod metadata"))
	}

	var runningPodNames []string
	for _, pod := range podList.Pods {
		if pod.Status.Phase == "Running" {
			runningPodNames = append(runningPodNames, pod.Metadata.Name)
		}
	}

	return runningPodNames, nil
}

func areK8SPodsDead(moduleName string) (bool, error) {
	output, err := exec.Command(kubectl, k8sNamespaceFlag, "get", "pods", "--selector=module="+moduleName, "--output=json").CombinedOutput()
	if err != nil {
		return false, newExecError(output, err)
	}

	podList := K8SPodList{}
	if err = json.Unmarshal(output, &podList); err != nil {
		return false, newExecError(output, errors.New("Could not read pod metadata"))
	}

	runningCount := 0
	notRunningCount := 0
	for _, pod := range podList.Pods {
		if pod.Status.Phase == "Running" {
			runningCount = runningCount + 1
		} else {
			notRunningCount = notRunningCount + 1
		}
	}

	return runningCount == 0 && notRunningCount > 0, nil
}

func waitForK8SPods(moduleName string, waitTilFinished bool) error {
	if waitTilFinished {
		startSpinner(fmt.Sprintf("Waiting for module \"%s\" to execute", moduleName))
	} else {
		startSpinner(fmt.Sprintf("Waiting for module \"%s\" to start", moduleName))
	}

	// TODO(skeswa): scale the checks for prod.
	// Make enough attempts to span 1 minute.
	for i := 0; i < 60; i++ {
		// Pause for a second before trying again (after first attempt).
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

		if waitTilFinished {
			// Check if there pods that have finished.
			if podsAreDead, err := areK8SPodsDead(moduleName); err != nil {
				stopSpinner(false)
				return err
			} else if podsAreDead {
				stopSpinner(true)
				return nil
			}
		} else {
			// Check if there are running pods.
			if pods, err := filterRunningK8SPods(moduleName); err != nil {
				stopSpinner(false)
				return err
			} else if len(pods) > 0 {
				stopSpinner(true)
				return nil
			}
		}
	}

	// Ran out of attempts!
	stopSpinner(false)
	return fmt.Errorf(
		"Timed out waiting for pods of module \"%s\" to come up.",
		moduleName)
}

func abbreviateK8SPath(k8sPath string) string {
	sep := string(os.PathSeparator)
	parts := strings.Split(k8sPath, sep)
	return strings.Join(parts[len(parts)-3:], sep)
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
	fileData, err := ioutil.ReadFile(k8sfilePath)
	if err != nil {
		return err
	}

	updatedFileData := prodK8SImageURLRegex.ReplaceAll(fileData, []byte(newImageURL))
	if err = ioutil.WriteFile(k8sfilePath, updatedFileData, 0644); err != nil {
		return err
	}

	return nil
}
