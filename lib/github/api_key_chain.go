package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/github/apikey"
)

const (
	keyRefreshDaemonInterval  = 1 * time.Hour
	devAPIKeysSecretFileName  = "github-api-keys.dev.json"
	prodAPIKeysSecretFileName = "github-api-keys.prod.json"
)

// APIKeyChain is responsible for managing GitHubAPIKeymodels
// and cycling through keys that hit their request limit
type apiKeyChain struct {
	q                db.BatchingQueryable
	lock             sync.RWMutex
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

	// Start refreshing keys in a fixed interval.
	go keyChain.runKeyRefreshDaemon()

	return &keyChain, nil
}

// runKeyRefreshDaemon calls refreshKeys every so often to keep the keys in
// sync with those in the database.
func (chain *apiKeyChain) runKeyRefreshDaemon() {
	for {
		log.Printf(
			"Refreshing Github API keys in at %s.\n",
			time.Now().Add(keyRefreshDaemonInterval).String())

		time.Sleep(keyRefreshDaemonInterval)
		if err := chain.refreshKeys(); err != nil {
			log.Printf("Failed to refresh Github API keys: %v\n.", err)
		} else {
			log.Println("Github API keys refreshed successfully.")
		}
	}
}

// refreshKeys reads keys from the database, and adds the new ones to the in
// memory collection of keys.
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

	// Only change keys if there are now more keys.
	chain.lock.Lock()
	if len(chain.keys) <= len(keyStrings) {
		chain.lock.Unlock()
		return nil
	}

	// Put the existing keys in a set for easy lookups.
	oldKeysSet := make(map[string]bool)
	for _, oldKey := range chain.keys {
		oldKeysSet[oldKey.token] = true
	}
	chain.lock.Unlock()

	// Turn the new key strings into new fully-qualified keys.
	var newKeys []*apiKey
	for _, keyString := range keyStrings {
		if !oldKeysSet[keyString] {
			// TODO(skeswa): make this parallel instead of blocking on each
			// "newAPIKey".
			newKey, err := newAPIKey(keyString)
			if err != nil {
				return err
			}

			newKeys = append(newKeys, newKey)
		}
	}

	// Append the old keys to the new keys.
	chain.lock.Lock()
	chain.keys = append(newKeys, chain.keys...)
	chain.lock.Unlock()

	return nil
}

// acquireKey employs a round-robin policy to find the next Github API key. If
// no usable keys are found, it blocks until a key is available.
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

// readGithubKeysFromSecret reads Github API keys from a secrets file.
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
	type marshalledAPIKey struct {
		Key              string `json:"key"`
		ForScheduledJobs bool   `json:"forScheduledJobs"`
	}

	// Create the slice for unmarshalling.
	marshalledAPIKeys := []marshalledAPIKey{}
	if err = json.Unmarshal(apiKeysJSON, &marshalledAPIKeys); err != nil {
		return nil, err
	} else if len(marshalledAPIKeys) < 1 {
		return nil, fmt.Errorf("There were no keys in the secret!")
	}

	// Turn the marshalledAPIKeys into insertion tuples.
	var insertionTuples []apikey.InsertionTuple
	for _, uak := range marshalledAPIKeys {
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
	for _, key := range marshalledAPIKeys {
		keyStrings = append(keyStrings, key.Key)
	}

	return keyStrings, nil
}
