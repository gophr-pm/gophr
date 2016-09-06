package github

import (
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common/config"
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

// RequestService is the library responsible for managing all outbound
// requests to GitHub
type RequestService struct {
	APIKeyChain *APIKeyChain
}

// NewRequestService initialies a new GitHubRequestService and APIKeyChain
func NewRequestService(params RequestServiceParams) *RequestService {
	newRequestService := RequestService{}
	newRequestService.APIKeyChain = NewAPIKeyChain(params)

	return &newRequestService
}

// RequestServiceParams passes import Kubernetes configuration and secrets.
// Also can dictate whether request service is being used by indexer.
type RequestServiceParams struct {
	ForIndexer bool
	Conf       *config.Config
	Session    *gocql.Session
}
