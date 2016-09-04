package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io/ioutil"
	"path/filepath"
)

func splitKeyFileData(data []byte) (key []byte, nonce []byte) {
	// Read the key (the key is an aes key so it must be 32 bytes for AES-256).
	return data[:32], data[32:]
}

func writeKeyFile(keyFilePath string, key []byte, nonce []byte) error {
	buffer := bytes.Buffer{}
	buffer.Write(key)
	buffer.Write(nonce)
	return ioutil.WriteFile(keyFilePath, buffer.Bytes(), 0644)
}

func encryptSecret(secretFilePath string, keyFilePath string) ([]byte, error) {
	// Read the keyfile and secret.
	keyFileData, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	secret, err := ioutil.ReadFile(secretFilePath)
	if err != nil {
		return nil, err
	}

	// Split up the data in the key file.
	key, nonce := splitKeyFileData(keyFileData)
	// Use the key to create a cipher.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Create the gcm agent.
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// Encrypt the secret.
	encryptedSecret := aesgcm.Seal(nil, nonce, secret, nil)

	// Return the path to the tmp file.
	return encryptedSecret, nil
}

func generateDecryptedSecret(secretFilePath string, keyFilePath string) (string, error) {
	// Read the keyfile and secret.
	keyFileData, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return "", err
	}
	secret, err := ioutil.ReadFile(secretFilePath)
	if err != nil {
		return "", err
	}

	// Split up the data in the key file.
	key, nonce := splitKeyFileData(keyFileData)
	// Use the key to create a cipher.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// Create the gcm agent.
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	// Decrypt the secret.
	decryptedSecret, err := aesgcm.Open(nil, nonce, secret, nil)
	if err != nil {
		return "", err
	}
	// Create a tmp file for the decrypted secret.
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	// Preserve the secret file name.
	_, secretFileName := filepath.Split(secretFilePath)
	outputFilePath := filepath.Join(tmpDir, secretFileName)
	// Write the decrypted secret to the tmp file.
	err = ioutil.WriteFile(outputFilePath, decryptedSecret, 0644)
	if err != nil {
		return "", err
	}

	// Return the path to the tmp file.
	return outputFilePath, nil
}
