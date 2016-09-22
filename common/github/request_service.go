package github

import (
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/models"
)

// GitHubGophrPackageOrgName is the  Github organization name for all versioned packages
const (
	GitHubGophrPackageOrgName = "gophr-packages"
)

// GitHubBaseAPIURL is the base Github API
const (
	GitHubBaseAPIURL = "https://api.github.com"
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
	FetchGitHubDataForPackageModel(models.PackageModel) (map[string]interface{}, error)
}

// requestService is the library responsible for managing all outbound
// requests to GitHub
type requestService struct {
	APIKeyChain *APIKeyChain
}

// NewRequestService initialies a new GitHubrequestService and APIKeyChain
func NewRequestService(args RequestServiceArgs) RequestService {
	svc := requestService{}
	svc.APIKeyChain = NewAPIKeyChain(args)

	return &svc
}

// RequestServiceArgs passes import Kubernetes configuration and secrets.
// Also can dictate whether request service is being used by indexer.
type RequestServiceArgs struct {
	Conf       *config.Config
	Session    *gocql.Session
	ForIndexer bool
}
