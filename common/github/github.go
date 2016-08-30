package github

import (
	"net/http"
	"time"
)

// GitHubGophrPackageOrgName is the  Github organization name for all versioned packages
const (
	GitHubGophrPackageOrgName = "gophr-packages"
	GitHubBaseAPIURL          = "https://api.github.com"
	githubRootTemplate        = "github.com/%s/%s"
	gitHubRemoteOrigin        = "git@github.com:gophr-packages/%s.git"
)

// TODO:(Shikkic) trim whatever isn't needed here
const (
	refsHead                                  = "HEAD"
	refsLineCap                               = "\n\x00"
	refsSpaceChar                             = ' '
	refsHeadPrefix                            = "refs/heads/"
	refsLineFormat                            = "%04x%s"
	refsHeadMaster                            = "refs/heads/master"
	refsMasterLineFormat                      = "%s refs/heads/master\n"
	refsSymRefAssignment                      = "symref="
	refsOldRefAssignment                      = "oldref="
	refsFetchURLTemplate                      = "https://%s.git/info/refs?service=git-upload-pack"
	refsAugmentedHeadLineFormat               = "%s HEAD\n"
	refsAugmentedSymrefHeadLineFormat         = "%s HEAD\x00symref=HEAD:%s\n"
	refsAugmentedHeadLineWithCapsFormat       = "%s HEAD\x00%s\n"
	refsAugmentedSymrefHeadLineWithCapsFormat = "%s HEAD\x00symref=HEAD:%s %s\n"
)

// TODO:(Shikkic) trim whatever isn't needed here
const (
	errorRefsFetchNoSuchRepo       = "Could not find a Github repository at %s"
	errorRefsFetchGithubError      = "Github responded with an error: %v"
	errorRefsFetchGithubParseError = "Cannot read refs from Github: %v"
	errorRefsFetchNetworkFailure   = "Could not reach Github at the moment; Please try again later"
	errorRefsParseSizeFormat       = "Could not parse refs line size: %s"
	errorRefsParseIncompleteRefs   = "Incomplete refs data received from GitHub"
)

var (
	commitsUntilParameter = "until"
	commitsAfterParameter = "after'"
	httpClient            = &http.Client{Timeout: 10 * time.Second}
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
