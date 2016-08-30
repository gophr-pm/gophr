package github

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/errors"
)

// APIKeyChain is responsible for managing GitHubAPIKeymodels
// and cycling through keys that hit their request limit
type APIKeyChain struct {
	GitHubAPIKeys []APIKeyModel
	CurrentKey    APIKeyModel
}

// NewAPIKeyChain intializes and returns a new GitHubAPIKeyChain
// and instantiates all available keys in the db as APIKeyModels
func NewAPIKeyChain() *APIKeyChain {
	log.Println("Creating new github api keychain")
	newGitHubAPIKeyChain := APIKeyChain{}

	_, session := common.Init()
	defer session.Close()

	gitHubAPIKeys, err := scanAllGitHubKey(session)
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

func scanAllGitHubKey(session *gocql.Session) ([]string, error) {
	var (
		err          error
		scanError    error
		closeError   error
		gitHubAPIKey string

		key string

		query = session.Query(`SELECT
			key
			FROM gophr.github_api_key`)
		iter          = query.Iter()
		gitHubAPIKeys = make([]string, 0)
	)

	for iter.Scan(&key) {
		gitHubAPIKey = key
		gitHubAPIKeys = append(gitHubAPIKeys, gitHubAPIKey)
	}

	if err = iter.Close(); err != nil {
		closeError = err
	}

	if scanError != nil || closeError != nil {
		return nil, errors.NewQueryScanError(scanError, closeError)
	}

	return gitHubAPIKeys, nil
}
