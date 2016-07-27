package common

import (
	"log"
	"time"
)

type GitHubAPIKeyChain struct {
	GitHubAPIKeys []GitHubAPIKeyModel
	CurrentKey    GitHubAPIKeyModel
}

func NewGitHubAPIKeyChain() *GitHubAPIKeyChain {
	log.Println("CREATING NEW KEYCHAIN")
	newGitHubAPIKeyChain := GitHubAPIKeyChain{}
	newGitHubAPIKeyChain.GitHubAPIKeys = make([]GitHubAPIKeyModel, 0)

	keys := []string{}

	// For each api key create a new key model and push to the the list of GitHubAPIKeys
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

		newGitHubAPIKeyChain.GitHubAPIKeys = append(newGitHubAPIKeyChain.GitHubAPIKeys, gitHubAPIKeyModel)
	}

	//TODO sort
	newGitHubAPIKeyChain.setCurrentKey()
	return &newGitHubAPIKeyChain
}

func (gitHubAPIKeyChain *GitHubAPIKeyChain) getAPIKeyModel() *GitHubAPIKeyModel {
	if gitHubAPIKeyChain.CurrentKey.RemainingUses > 0 {
		return &gitHubAPIKeyChain.CurrentKey
	}
	log.Println("KEY NEEDS TO BE SWAPED")
	newGitHubAPIKeys := make([]GitHubAPIKeyModel, 0)
	var swapModel GitHubAPIKeyModel
	for index, APIKeyModel := range gitHubAPIKeyChain.GitHubAPIKeys {
		if index == 0 {
			swapModel = APIKeyModel
		} else {
			newGitHubAPIKeys = append(newGitHubAPIKeys, APIKeyModel)
		}
	}
	newGitHubAPIKeys = append(newGitHubAPIKeys, swapModel)
	gitHubAPIKeyChain.GitHubAPIKeys = newGitHubAPIKeys
	gitHubAPIKeyChain.setCurrentKey()

	if gitHubAPIKeyChain.CurrentKey.RemainingUses == 0 {
		coolDown(gitHubAPIKeyChain.CurrentKey)
	}

	return &gitHubAPIKeyChain.CurrentKey
}

func (gitHubAPIKeyChain *GitHubAPIKeyChain) setCurrentKey() {
	if len(gitHubAPIKeyChain.GitHubAPIKeys) == 0 {
		gitHubAPIKeyChain.CurrentKey = GitHubAPIKeyModel{}
	}

	gitHubAPIKeyChain.CurrentKey = gitHubAPIKeyChain.GitHubAPIKeys[0]

}

func coolDown(apiKeyModel GitHubAPIKeyModel) {
	timeNow := time.Now()
	log.Println("Time now = ", timeNow)
	resetTime := apiKeyModel.RateLimitResetTime
	log.Println("Reset time = ", resetTime)

	sleepTime := resetTime.Sub(timeNow)
	log.Println("Sleep time is = ", sleepTime)
	time.Sleep(sleepTime)
}
