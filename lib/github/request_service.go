package github

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/dtos"
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

// RequestService is the external interface of the internal requestService.
type RequestService interface {
	FetchCommitSHA(string, string, time.Time) (string, error)
	FetchCommitTimestamp(string, string, string) (time.Time, error)
	FetchGitHubDataForPackageModel(
		author string,
		repo string) (dtos.GithubRepo, error)
}

// requestService is the library responsible for managing all outbound
// requests to GitHub
type requestService struct {
	keyChain *apiKeyChain
}

// RequestServiceArgs passes import Kubernetes configuration and secrets.
// Also can dictate whether request service is being used by indexer.
type RequestServiceArgs struct {
	Conf             *config.Config
	Queryable        db.BatchingQueryable
	ForScheduledJobs bool
}

// NewRequestService initialies a new GitHubrequestService and APIKeyChain
func NewRequestService(args RequestServiceArgs) (RequestService, error) {
	keyChain, err := newAPIKeyChain(args)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not create new Github request service: %v.",
			err)
	}

	svc := &requestService{keyChain: keyChain}
	return svc, nil
}
