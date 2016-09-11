package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/db/query"
	"github.com/skeswa/gophr/common/errors"
)

// Database string constants.
const (
	tableNameGithubAPIKey     = "github_api_key"
	devAPIKeysSecretFileName  = "github-api-keys.dev.json"
	prodAPIKeysSecretFileName = "github-api-keys.prod.json"
	columnNameGithubAPIKeyKey = "key"
)

// APIKeyChain is responsible for managing GitHubAPIKeymodels
// and cycling through keys that hit their request limit
type APIKeyChain struct {
	GitHubAPIKeys []APIKeyModel
	CurrentKey    APIKeyModel
}

// NewAPIKeyChain intializes and returns a new GitHubAPIKeyChain
// and instantiates all available keys in the db as APIKeyModels
func NewAPIKeyChain(conf *config.Config, session *gocql.Session) *APIKeyChain {
	log.Println("Creating new github api keychain")
	newGitHubAPIKeyChain := APIKeyChain{}

	gitHubAPIKeys, err := scanAllGitHubKey(conf, session)
	if err != nil {
		log.Println("Could not scan github keys, fatal error occurred")
		log.Fatal(err)
	}
	log.Printf("Found %d keys %s \n", len(gitHubAPIKeys), gitHubAPIKeys)

	newGitHubAPIKeyChain.GitHubAPIKeys = initializeGitHubAPIKeyModels(gitHubAPIKeys)
	newGitHubAPIKeyChain.setCurrentKey()
	// TODO (Shikkic): Optimize sort and choose algo for keys

	return &newGitHubAPIKeyChain
}

func initializeGitHubAPIKeyModels(keys []string) []APIKeyModel {
	var gitHubAPIKeyModels = make([]APIKeyModel, 0)

	for _, key := range keys {
		log.Println("KEY =", key)
		gitHubAPIKeyModel := APIKeyModel{
			string(key),
			5000,
			5000,
			time.Time{},
		}
		log.Println("Priming API KEY")
		gitHubAPIKeyModel.prime()
		gitHubAPIKeyModel.print()

		gitHubAPIKeyModels = append(gitHubAPIKeyModels, gitHubAPIKeyModel)
	}

	return gitHubAPIKeyModels
}

func (gitHubAPIKeyChain *APIKeyChain) getAPIKeyModel() *APIKeyModel {
	if gitHubAPIKeyChain.CurrentKey.RemainingUses > 0 {
		return &gitHubAPIKeyChain.CurrentKey
	}

	log.Println("Current key has hit maxium limit, needs to be swaped")
	gitHubAPIKeyChain.shuffleKeys()
	gitHubAPIKeyChain.setCurrentKey()
	gitHubAPIKeyChain.CurrentKey.prime()

	if gitHubAPIKeyChain.CurrentKey.RemainingUses <= 0 {
		setRequestTimout(gitHubAPIKeyChain.CurrentKey)
	}

	return &gitHubAPIKeyChain.CurrentKey
}

func (gitHubAPIKeyChain *APIKeyChain) shuffleKeys() {
	var newGitHubAPIKeys = make([]APIKeyModel, 0)
	var firstAPIModelInArray APIKeyModel

	for index, APIKeyModel := range gitHubAPIKeyChain.GitHubAPIKeys {
		if index == 0 {
			firstAPIModelInArray = APIKeyModel
		} else {
			newGitHubAPIKeys = append(newGitHubAPIKeys, APIKeyModel)
		}
	}
	newGitHubAPIKeys = append(newGitHubAPIKeys, firstAPIModelInArray)

	gitHubAPIKeyChain.GitHubAPIKeys = newGitHubAPIKeys
}

// TODO(skeswa): this breaks when there are no keys in the database. @Shikkic,
// investigate this.
func (gitHubAPIKeyChain *APIKeyChain) setCurrentKey() {
	if len(gitHubAPIKeyChain.GitHubAPIKeys) == 0 {
		gitHubAPIKeyChain.CurrentKey = APIKeyModel{}
	}

	gitHubAPIKeyChain.CurrentKey = gitHubAPIKeyChain.GitHubAPIKeys[0]
}

func setRequestTimout(apiKeyModel APIKeyModel) {
	timeNow := time.Now()
	log.Printf("The current time is %s. \n", timeNow)
	resetTime := apiKeyModel.RateLimitResetTime
	log.Printf("APIKey Reset time is %s. \n", resetTime)
	sleepTime := resetTime.Sub(timeNow)
	log.Printf("Indexer will sleep for %s. \n", sleepTime)
	time.Sleep(sleepTime)
}

func scanAllGitHubKey(conf *config.Config, session *gocql.Session) ([]string, error) {
	var (
		err           error
		gitHubAPIKey  string
		gitHubAPIKeys []string
	)

	iter := query.Select(columnNameGithubAPIKeyKey).
		From(tableNameGithubAPIKey).
		Create(session).
		Iter()

	for iter.Scan(&gitHubAPIKey) {
		gitHubAPIKeys = append(gitHubAPIKeys, gitHubAPIKey)
	}

	if err = iter.Close(); err != nil {
		return nil, errors.NewQueryScanError(nil, err)
	}

	// If there are no keys in the database, then add the ones from the secret
	// file (if it exists).
	if len(gitHubAPIKeys) < 1 && len(conf.SecretsPath) > 0 {
		gitHubAPIKeys, err = readGithubKeysFromSecret(conf, session)
		if err != nil {
			log.Printf("Failed to read keys from secret: %v.", err)
		}
	}

	return gitHubAPIKeys, nil
}

func readGithubKeysFromSecret(conf *config.Config, session *gocql.Session) ([]string, error) {
	log.Println("There were no keys in the database. Attempting to load from the github keys secret.")

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

	log.Println("Loaded the data from the keys secret file successfully. Now unmarshalling json.")

	// Create the struct for unmarshalling.
	type apiKey struct {
		Key                string `json:"key"`
		HasAdminPrivileges bool   `json:"hasAdminPrivileges"`
	}

	// Create the slice for unmarshalling.
	keys := []apiKey{}
	if err = json.Unmarshal(apiKeysJSON, &keys); err != nil {
		return nil, err
	} else if len(keys) < 1 {
		return nil, fmt.Errorf("There were no keys in the secret!")
	}

	log.Println("Unmarshalled keys from the secret successfully. Now inserting into the database.")

	// Execute the insert queries.
	for _, key := range keys {
		if err = query.InsertInto(tableNameGithubAPIKey).
			Value(columnNameGithubAPIKeyKey, key.Key).
			Create(session).
			Exec(); err != nil {
			return nil, err
		}
	}

	log.Println("Inserted keys into the database successfully. Returning the keys in string form.")

	var keyStrings []string
	for _, key := range keys {
		keyStrings = append(keyStrings, key.Key)
	}

	return keyStrings, nil
}
