package common

import (
	"log"
	"time"
)

type GitHubAPIKeyChain struct {
	GitHubAPIKeys []GitHubAPIKeyModel
}

func NewGitHubAPIKeyChain() *GitHubAPIKeyChain {
	log.Println("CREATING NEW KEYCHAIN")
	newGitHubAPIKeyChain := GitHubAPIKeyChain{}
	newGitHubAPIKeyChain.GitHubAPIKeys = make([]GitHubAPIKeyModel, 0)

	keys := []string{"79aed53b9efda8b1b62512d53b0308dbc3e0511f"}

	// For each api key create a new key model and push to the the list of GitHubAPIKeys
	for _, key := range keys {
		log.Println("KEY =", key)
		gitHubAPIKeyModel := GitHubAPIKeyModel{
			string(key),
			5000,
			5000,
			time.Time{},
		}

		newGitHubAPIKeyChain.GitHubAPIKeys = append(newGitHubAPIKeyChain.GitHubAPIKeys, gitHubAPIKeyModel)
	}

	return &newGitHubAPIKeyChain
}

func (gitHubAPIKeyChain *GitHubAPIKeyChain) getAPIKeyModel() *GitHubAPIKeyModel {
	return &gitHubAPIKeyChain.GitHubAPIKeys[0]
}
