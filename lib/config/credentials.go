package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// Credentials holds all the authentication metadata needed by gophr backend
// modules.
type Credentials struct {
	GithubPush UserPass `json:"githubPushCreds"`
}

// UserPass is a tuple of user name and password.
type UserPass struct {
	User string `json:"user"`
	Pass string `json:"password"`
}

const credentialsFileName = "credentials.json"

// ReadCredentials reads credentials from the credentials secret.
func ReadCredentials(conf *Config) (*Credentials, error) {
	data, err := ioutil.ReadFile(filepath.Join(conf.SecretsPath, credentialsFileName))
	if err != nil {
		return nil, err
	}

	creds := Credentials{}
	if err = json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}
