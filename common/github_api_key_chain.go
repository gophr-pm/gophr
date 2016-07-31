package common

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common/errors"
)

type GitHubAPIKeyChain struct {
	GitHubAPIKeys []GitHubAPIKeyModel
	CurrentKey    GitHubAPIKeyModel
}

func NewGitHubAPIKeyChain() *GitHubAPIKeyChain {
	log.Println("CREATING NEW KEYCHAIN")
	newGitHubAPIKeyChain := GitHubAPIKeyChain{}

	_, session := Init()
	defer session.Close()

	gitHubAPIKeys, err := scanAllGitHubKey(session)
	if err != nil {
		log.Println(err)
	}
	log.Println(gitHubAPIKeys)

	// TODO RENAME TO initializeGitHubAPIKeyModels
	newGitHubAPIKeyChain.GitHubAPIKeys = initializeGitHubAPIKeys(gitHubAPIKeys)

	//TODO sort
	newGitHubAPIKeyChain.setCurrentKey()
	return &newGitHubAPIKeyChain
}

// For each api key create a new key model and push to the the list of GitHubAPIKeys
func initializeGitHubAPIKeys(keys []string) []GitHubAPIKeyModel {
	var gitHubAPIKeyModels = make([]GitHubAPIKeyModel, 0)

	for _, key := range keys {
		log.Println("KEY =", key)
		gitHubAPIKeyModel := GitHubAPIKeyModel{
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

func (gitHubAPIKeyChain *GitHubAPIKeyChain) getAPIKeyModel() *GitHubAPIKeyModel {
	if gitHubAPIKeyChain.CurrentKey.RemainingUses > 0 {
		return &gitHubAPIKeyChain.CurrentKey
	}

	log.Println("KEY NEEDS TO BE SWAPED")
	gitHubAPIKeyChain.shuffleKeys()
	gitHubAPIKeyChain.setCurrentKey()
	gitHubAPIKeyChain.CurrentKey.prime()

	// If new key's remaining use is 0 set a time out
	if gitHubAPIKeyChain.CurrentKey.RemainingUses <= 0 {
		setRequestTimout(gitHubAPIKeyChain.CurrentKey)
	}

	return &gitHubAPIKeyChain.CurrentKey
}

func (gitHubAPIKeyChain *GitHubAPIKeyChain) shuffleKeys() {
	newGitHubAPIKeys := make([]GitHubAPIKeyModel, 0)
	var firstAPIModelInArray GitHubAPIKeyModel

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

func (gitHubAPIKeyChain *GitHubAPIKeyChain) setCurrentKey() {
	if len(gitHubAPIKeyChain.GitHubAPIKeys) == 0 {
		gitHubAPIKeyChain.CurrentKey = GitHubAPIKeyModel{}
	}

	gitHubAPIKeyChain.CurrentKey = gitHubAPIKeyChain.GitHubAPIKeys[0]
}

func setRequestTimout(apiKeyModel GitHubAPIKeyModel) {
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
