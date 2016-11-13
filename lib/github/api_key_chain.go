package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/github/apikey"
)

const (
	devAPIKeysSecretFileName  = "github-api-keys.dev.json"
	prodAPIKeysSecretFileName = "github-api-keys.prod.json"
)

// APIKeyChain is responsible for managing GitHubAPIKeymodels
// and cycling through keys that hit their request limit
type apiKeyChain struct {
	q                db.BatchingQueryable
	lock             sync.Mutex
	keys             []*apiKey
	conf             *config.Config
	cursor           int
	forScheduledJobs bool
}

// newAPIKeyChain intializes and returns a new GitHubAPIKeyChain
// and instantiates all available keys in the db as apiKeyModels
func newAPIKeyChain(args RequestServiceArgs) (*apiKeyChain, error) {
	keyChain := apiKeyChain{
		q:                args.Queryable,
		conf:             args.Conf,
		forScheduledJobs: args.ForScheduledJobs,
	}

	// Go out an get keys for the first time.
	if err := keyChain.refreshKeys(); err != nil {
		return nil, fmt.Errorf(
			"Failed to perform initial key refresh to create new chain: %v.",
			err)
	}

	return &keyChain, nil
}

func (chain *apiKeyChain) refreshKeys() error {
	// Read keys from the database.
	keyStrings, err := apikey.GetAll(chain.q, chain.forScheduledJobs)
	if err != nil {
		return fmt.Errorf(
			"Failed to get Github API keys from the database for key chain: %v.",
			err)
	}

	// If there aren't any keys, go out and get some.
	if len(keyStrings) < 1 {
		keyStrings, err = readGithubKeysFromSecret(chain.conf, chain.q)
		if err != nil {
			return fmt.Errorf(
				"Failed to get Github API keys from secrets for new key chain: %v.",
				err)
		}
	}

	// Turn the key strings into fully-qualified keys.
	var keys []*apiKey
	for _, keyString := range keyStrings {
		// TODO(skeswa): make this parallel instead of blocking on each "newAPIKey".
		key, err := newAPIKey(keyString)
		if err != nil {
			return err
		}

		keys = append(keys, key)
	}

	// Change over the keys.
	chain.lock.Lock()
	chain.keys = keys
	chain.lock.Unlock()

	return nil
}

func (chain *apiKeyChain) acquireKey() *apiKey {
	chain.lock.Lock()
	keys := chain.keys
	cursor := chain.cursor
	chain.cursor++
	chain.lock.Unlock()

	if key := keys[cursor%len(keys)]; key.canBeUsed() {
		return key
	}

	// This key is spent. Gotta find another one.
	for i := 0; i < len(keys); i++ {
		if key := keys[(i+cursor)%len(keys)]; key.canBeUsed() {
			return key
		}
	}

	// There are no keys presently available. Wait for the original one.
	key := keys[cursor%len(keys)]
	key.waitUntilUseful()
	return key
}

func readGithubKeysFromSecret(
	conf *config.Config,
	q db.BatchingQueryable,
) ([]string, error) {
	log.Println(
		"There were no keys in the database. " +
			"Attempting to load from the github keys secret.")

	var filePath string
	if conf.IsDev {
		filePath = filepath.Join(conf.SecretsPath, devAPIKeysSecretFileName)
	} else {
		filePath = filepath.Join(conf.SecretsPath, prodAPIKeysSecretFileName)
	}

	var (
		err         error
		apiKeysJSON []byte
	)

	// Read the secret data.
	if apiKeysJSON, err = ioutil.ReadFile(filePath); err != nil {
		return nil, err
	}

	log.Println(
		"Loaded the data from the keys secret file successfully. " +
			"Now unmarshalling json.")

	// Create the struct for unmarshalling.
	type unmarshalledAPIKey struct {
		Key              string `json:"key"`
		ForScheduledJobs bool   `json:"forScheduledJobs"`
	}

	// Create the slice for unmarshalling.
	unmarshalledAPIKeys := []unmarshalledAPIKey{}
	if err = json.Unmarshal(apiKeysJSON, &unmarshalledAPIKeys); err != nil {
		return nil, err
	} else if len(unmarshalledAPIKeys) < 1 {
		return nil, fmt.Errorf("There were no keys in the secret!")
	}

	// Turn the unmarshalledAPIKeys into insertion tuples.
	var insertionTuples []apikey.InsertionTuple
	for _, uak := range unmarshalledAPIKeys {
		insertionTuples = append(insertionTuples, apikey.InsertionTuple{
			Key:              uak.Key,
			ForScheduledJobs: uak.ForScheduledJobs,
		})
	}

	log.Println(
		"Unmarshalled keys from the secret successfully. " +
			"Now inserting into the database.")

	// Execute the insert queries.
	if err = apikey.InsertAll(q, insertionTuples); err != nil {
		return nil, fmt.Errorf("Failed to read github keys from secret: %v.", err)
	}

	log.Println(
		"Inserted keys into the database successfully. " +
			"Returning the keys in string form.")

	var keyStrings []string
	for _, key := range unmarshalledAPIKeys {
		keyStrings = append(keyStrings, key.Key)
	}

	return keyStrings, nil
}
