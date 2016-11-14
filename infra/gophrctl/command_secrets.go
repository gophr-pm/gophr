package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

const (
	secretsDir = "./infra/k8s/secrets"
)

func secretsNewKeyCommand(c *cli.Context) error {
	printInfo("Creating a new key")
	keyFilePath := c.Args().First()
	if len(keyFilePath) < 1 {
		exit(exitCodeNewKeyFailed, nil, "", fmt.Errorf("Invalid key file path: \"%s\".", keyFilePath))
	}

	keyFilePath, err := filepath.Abs(keyFilePath)
	if err != nil {
		exit(exitCodeNewKeyFailed, nil, "", err)
		return nil
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		exit(exitCodeNewKeyFailed, nil, "", errors.New("Failed to generate the nonce."))
	}
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		exit(exitCodeNewKeyFailed, nil, "", errors.New("Failed to generate the key."))
	}
	if err := writeKeyFile(keyFilePath, key, nonce); err != nil {
		exit(exitCodeNewKeyFailed, nil, "", fmt.Errorf("Invalid key file path: \"%s\".", keyFilePath))
	}

	printSuccess(fmt.Sprintf("New keyfile written at \"%s\".", keyFilePath))
	return nil
}

func secretsRecordCommand(c *cli.Context) error {
	var (
		err            error
		gophrRoot      string
		keyFilePath    string
		secretFilePath string
		secretFileName string
	)

	printInfo("Recording a new secret")
	if gophrRoot, err = readGophrRoot(c); err != nil {
		exit(exitCodeRecordSecretFailed, nil, "", err)
	}

	keyFilePath = c.String(flagNameKeyPath)
	if len(keyFilePath) < 1 {
		exit(exitCodeRecordSecretFailed, nil, "", fmt.Errorf("Invalid key file path: \"%s\".", keyFilePath))
	}
	keyFilePath, err = filepath.Abs(keyFilePath)
	if err != nil {
		exit(exitCodeRecordSecretFailed, nil, "", err)
	}

	secretFilePath = c.Args().First()
	if len(secretFilePath) < 1 {
		exit(exitCodeRecordSecretFailed, nil, "", fmt.Errorf("Invalid secret file path: \"%s\".", secretFilePath))
	}
	secretFilePath, err = filepath.Abs(secretFilePath)
	if err != nil {
		exit(exitCodeRecordSecretFailed, nil, "", err)
	}

	encryptedSecret, err := encryptSecret(secretFilePath, keyFilePath)
	if err != nil {
		exit(exitCodeRecordSecretFailed, nil, "", err)
	}

	// Concat the the output path together.
	_, secretFileName = filepath.Split(secretFilePath)
	outputFilePath := filepath.Join(gophrRoot, secretsDir, secretFileName)

	// Write the decrypted secret to the tmp file.
	if err = ioutil.WriteFile(
		outputFilePath,
		encryptedSecret,
		0644); err != nil {
		exit(exitCodeRecordSecretFailed, nil, "", err)
	}

	printSuccess(fmt.Sprintf("New secret recorded at \"%s\".", outputFilePath))
	return nil
}

func secretsCycleCommand(c *cli.Context) error {
	if err := runInK8S(c, func() error {
		var (
			err                  error
			gophrRoot            string
			keyFilePath          string
			secretFilePath       string
			decryptedSecretPath  string
			decryptedSecretPaths []string
		)

		printInfo("Cycling all recorded secrets")
		if gophrRoot, err = readGophrRoot(c); err != nil {
			exit(exitCodeCycleSecretsFailed, nil, "", err)
		}

		keyFilePath = c.String(flagNameKeyPath)
		if len(keyFilePath) < 1 {
			exit(exitCodeCycleSecretsFailed, nil, "", fmt.Errorf("Invalid key file path: \"%s\".", keyFilePath))
		}
		keyFilePath, err = filepath.Abs(keyFilePath)
		if err != nil {
			exit(exitCodeCycleSecretsFailed, nil, "", err)
		}
		secretFiles, err := ioutil.ReadDir(filepath.Join(gophrRoot, secretsDir))
		if err != nil {
			return err
		}

		for _, secretFile := range secretFiles {
			secretFilePath = filepath.Join(gophrRoot, secretsDir, secretFile.Name())
			if decryptedSecretPath, err = generateDecryptedSecret(secretFilePath, keyFilePath); err != nil {
				return err
			}

			decryptedSecretPaths = append(decryptedSecretPaths, decryptedSecretPath)
		}

		if secretExistsInK8S() {
			if err = deleteSecretsInK8S(); err != nil {
				return err
			}
		}
		if err = createSecretsInK8S(decryptedSecretPaths); err != nil {
			return err
		}

		// Delete all of the generated secret files.
		startSpinner("Cleaning up generated files")
		for _, decryptedSecretPath := range decryptedSecretPaths {
			if err = os.Remove(decryptedSecretPath); err != nil {
				return err
			}
		}
		stopSpinner(true)

		printSuccess("Secrets cycled successfully")
		return nil
	}); err != nil {
		exit(exitCodeCycleSecretsFailed, nil, "", err)
	}

	return nil
}

func secretsRevealCommand(c *cli.Context) error {
	var (
		err            error
		keyFilePath    string
		secretFilePath string
	)

	keyFilePath = c.String(flagNameKeyPath)
	if len(keyFilePath) < 1 {
		exit(exitCodeRevealSecretFailed, nil, "", fmt.Errorf("Invalid key file path: \"%s\".", keyFilePath))
	}
	keyFilePath, err = filepath.Abs(keyFilePath)
	if err != nil {
		exit(exitCodeRevealSecretFailed, nil, "", err)
	}

	secretFilePath = c.Args().First()
	if len(secretFilePath) < 1 {
		exit(exitCodeRevealSecretFailed, nil, "", fmt.Errorf("Invalid secret file path: \"%s\".", secretFilePath))
	}
	secretFilePath, err = filepath.Abs(secretFilePath)
	if err != nil {
		exit(exitCodeRevealSecretFailed, nil, "", err)
	}

	decryptedSecret, err := decryptSecret(secretFilePath, keyFilePath)
	if err != nil {
		exit(exitCodeRevealSecretFailed, nil, "", err)
	}

	bytes.NewBuffer(decryptedSecret).WriteTo(os.Stdout)
	return nil
}
