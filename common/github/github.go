package github

import (
	"net/http"
	"time"
)

// GitHubGophrPackageOrgName is the  Github organization name for all versioned packages
// GitHubBaseAPIURL is the base Github API
const (
	GitHubGophrPackageOrgName = "gophr-packages"
	GitHubBaseAPIURL          = "https://api.github.com"
)

const (
	githubRootTemplate    = "github.com/%s/%s"
	gitHubRemoteOrigin    = "git@github.com:gophr-packages/%s.git"
	commitsUntilParameter = "until"
	commitsAfterParameter = "after'"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
)

// GitHubRequestService is the library responsible for managing all outbound
// requests to GitHub
// TODO:(Shikkic) Rename this to just RequestService
type GitHubRequestService struct {
	APIKeyChain *GitHubAPIKeyChain
}

// NewGitHubRequestService initialies a new GitHubRequestService and APIKeyChain
func NewGitHubRequestService() *GitHubRequestService {
	newGitHubRequestService := GitHubRequestService{}
	newGitHubRequestService.APIKeyChain = NewGitHubAPIKeyChain()

	return &newGitHubRequestService
}
