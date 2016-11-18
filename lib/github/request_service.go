package github

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/datadog"
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

// RequestService is an abstraction that enables reliable communication with the
// Github REST API.
type RequestService interface {
	// FetchRepoData fetches the Github repository metadata for the specified
	// package.
	FetchRepoData(author string, repo string) (dtos.GithubRepo, error)
	// FetchCommitSHA fetches the commit SHA that is chronologically closest to a
	// given timestamp.
	FetchCommitSHA(
		author string,
		repo string,
		timeStamp time.Time) (string, error)
	// ExpandPartialSHA is responsible for fetching a full commit SHA from a short
	// SHA. This works by sending a HEAD request to the git archive endpoint with
	// a short SHA. The request returns a full SHA of the archive in the `Etag`
	// of the request header that is sent back.
	ExpandPartialSHA(args ExpandPartialSHAArgs) (string, error)
	// FetchCommitTimestamp fetches the timestamp of a commit from Github API.
	FetchCommitTimestamp(
		author string,
		repo string,
		sha string) (time.Time, error)
}

// requestServiceImpl is the implementation of the RequestService.
type requestServiceImpl struct {
	keyChain *apiKeyChain
	ddClient datadog.Client
}

// RequestServiceArgs passes import Kubernetes configuration and secrets.
// Also can dictate whether request service is being used by indexer.
type RequestServiceArgs struct {
	Conf             *config.Config
	DDClient         datadog.Client
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

	svc := &requestServiceImpl{keyChain: keyChain, ddClient: args.DDClient}
	return svc, nil
}
