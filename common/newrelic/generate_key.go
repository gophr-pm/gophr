package nr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gophr-pm/gophr/common/config"
)

// NewRelicKey is a strut for representing new relic api keys
// configuration.
type NewRelicKey struct {
	NewRelicKey string `json:"newRelicKey"`
}

const (
	prodAPIKeysSecretFileName = "newrelic-key.json"
)

// generateKey returns a NewRelicKey from a secret.
func generateKey(conf *config.Config) (string, error) {
	if !conf.IsDev {
		var filePath string
		filePath = filepath.Join(conf.SecretsPath, prodAPIKeysSecretFileName)

		var (
			err        error
			apiKeyJSON []byte
		)

		// Read the secret data.
		if apiKeyJSON, err = ioutil.ReadFile(filePath); err != nil {
			return "", err
		}

		key := NewRelicKey{}
		if err = json.Unmarshal(apiKeyJSON, &key); err != nil {
			return "", err
		} else if len(key.NewRelicKey) < 1 {
			return "", fmt.Errorf("There were no keys in the secret!")
		}

		return key.NewRelicKey, nil
	}

	return "", nil
}
