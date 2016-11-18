package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

const (
	templateK8SFileSuffix           = ".template.yml"
	templateVarGCEProjectID         = "{{GCE_PROJECT_ID}}"
	templateVarDepotVolumeServiceIP = "{{DEPOT_VOL_SVC_IP}}"
	templateVarDataDogAPIKey        = "{{DATADOG_API_KEY}}"
	ddAPIKeySecretFileName          = "datadog-api-key.json"
)

var (
	templateVarProviders = map[string]templateVarProvider{
		templateVarGCEProjectID: func(c *cli.Context) ([]byte, error) {
			var (
				err error
				gpi string
			)

			// Get the gce project id for the push.
			if gpi, err = readGPI(c); err != nil {
				return nil, err
			}

			return []byte(gpi), nil
		},
		templateVarDepotVolumeServiceIP: func(c *cli.Context) ([]byte, error) {
			// Get the cluster IP from kubectl. If its empty, then that means the
			// cluster IP has not been assigned.

			output, err := exec.Command(
				kubectl,
				k8sNamespaceFlag,
				"get",
				"services",
				"--selector=module=depot-vol",
				"-o",
				"jsonpath={range .items[*]}{@.spec.clusterIP}").CombinedOutput()
			if err != nil {
				return nil, err
			}

			// Clean off the whitespace before reading the ip.
			ip := bytes.TrimSpace(output[:])
			if len(ip) < 1 {
				return nil, errors.New("Could not read the cluster IP of the depot volume service.")
			}

			return ip, nil
		},
		templateVarDataDogAPIKey: func(c *cli.Context) ([]byte, error) {
			var (
				err            error
				gophrRoot      string
				keyFilePath    string
				secretFilePath string
			)

			type datadogKey struct {
				DatadogAPIKey string `json:"datadog_api_key"`
			}

			if gophrRoot, err = readGophrRoot(c); err != nil {
				exit(exitCodeCycleSecretsFailed, nil, "", err)
			}
			keyFilePath = c.String(flagNameKeyPath)
			if len(keyFilePath) < 1 {
				return nil, fmt.Errorf("Invalid key file path: \"%s\".", keyFilePath)
			}
			keyFilePath, err = filepath.Abs(keyFilePath)
			if err != nil {
				return nil, err
			}
			secretFilePath, err = filepath.Abs(filepath.Join(gophrRoot, secretsDir, ddAPIKeySecretFileName))
			log.Println(secretFilePath)
			if err != nil {
				return nil, err
			}
			decryptedSecret, err := decryptSecret(secretFilePath, keyFilePath)
			if err != nil {
				return nil, err
			}

			var secret = datadogKey{}
			if err = json.Unmarshal(decryptedSecret, &secret); err != nil {
				return nil, err
			} else if len(secret.DatadogAPIKey) < 1 {
				return nil, fmt.Errorf("There were no keys in the secret!")
			}

			return []byte(secret.DatadogAPIKey), nil
		},
	}
)

type templateVarProvider func(*cli.Context) ([]byte, error)

// isTemplateK8SFile checks if the k8s file is a template file.
func isTemplateK8SFile(k8sfilePath string) bool {
	return strings.HasSuffix(k8sfilePath, templateK8SFileSuffix)
}

// reviseK8STemplateFileData replaces all the template variables in a file with
// their corresponding values.
func reviseK8STemplateFileData(c *cli.Context, fileData []byte) ([]byte, error) {
	// Replace any and all of the template variable providers.
	for templateVarName, templateVarValueProvider := range templateVarProviders {
		// Check to see if the template var is represented.
		templateVarNameBytes := []byte(templateVarName)
		if bytes.Index(fileData, templateVarNameBytes) != -1 {
			// Get a value of the template var.
			templateVarValue, err := templateVarValueProvider(c)
			if err != nil {
				return nil, err
			}

			// Replace all instances of the template var.
			fileData = bytes.Replace(
				fileData,
				templateVarNameBytes,
				templateVarValue,
				-1)
		}
	}

	return fileData, nil
}

// compileK8STemplateFile turns a template file into a real file in tmp.
func compileK8STemplateFile(c *cli.Context, k8sfilePath string) (string, error) {
	// Read the template file.
	templateFileData, err := ioutil.ReadFile(k8sfilePath)
	if err != nil {
		return "", err
	}

	// De-templatify!
	fileData, err := reviseK8STemplateFileData(c, templateFileData)
	if err != nil {
		return "", err
	}

	// Create a tmp file for the decrypted secret.
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	// Preserve the k8s file name.
	_, fileName := filepath.Split(k8sfilePath)
	outputFilePath := filepath.Join(tmpDir, fileName)

	// Write the decrypted secret to the tmp file.
	err = ioutil.WriteFile(outputFilePath, fileData, 0644)
	if err != nil {
		return "", err
	}

	return outputFilePath, nil
}
