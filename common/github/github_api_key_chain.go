package github

import (
	"encoding/json"
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

	// If there are no keys in the database, and this is in the dev environment,
	// then add the ones from the secret file (if it exists).
	if conf.IsDev && len(gitHubAPIKeys) < 1 && len(conf.SecretsPath) > 0 {
		log.Println("There were no keys in the database. Since this is the " +
			"dev environment, attempting to load from the github keys secret.")
		filePath := filepath.Join(conf.SecretsPath, devAPIKeysSecretFileName)

		// Fail silently all the way through.
		if apiKeysJSON, err := ioutil.ReadFile(filePath); err == nil {
			log.Println("Loaded the data from the keys secret file successfully. Now unmarshalling json.")
			// Create the struct for unmarshalling.
			type apiKey struct {
				Key                string `json:"key"`
				HasAdminPrivileges string `json:"hasAdminPrivileges"`
			}

			// Create the slice for unmarshalling.
			keys := []apiKey{}
			if err = json.Unmarshal(apiKeysJSON, &keys); err == nil && len(keys) > 0 {
				log.Println("Unmarshalled the keys successfuly. Now inserting into the database.")
				// Create an insert query.
				q := query.InsertInto(tableNameGithubAPIKey)
				for _, key := range keys {
					q.Value(columnNameGithubAPIKeyKey, key.Key)
				}

				// Execute said query.
				if err = q.Create(session).Exec(); err == nil {
					log.Println("Inserted keys into the database successfully. Returning the keys in string form.")
					// If the keys were inserted okay, then return them in string form.
					for _, key := range keys {
						gitHubAPIKeys = append(gitHubAPIKeys, key.Key)
					}
				}
			}
		}
	}

	return gitHubAPIKeys, nil
}
