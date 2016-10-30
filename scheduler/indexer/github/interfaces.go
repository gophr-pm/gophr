package githubIndexer

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/model"
)

// Init is responsible for setting up the app configuration and db
// connection.
type Init func() (*config.Config, *gocql.Session)

// githubRequestService lol
type githubRequestService func(args github.RequestServiceArgs) github.RequestService

// packageRetriever lol
type packageRetriever func(session *gocql.Session) ([]*models.PackageModel, error)

// packageInserter lol
type packageInserter func(session *gocql.Session, packageModel *models.PackageModel) error

// packageDeleter lol
type packageDeleter func(session *gocql.Session, packageModel *models.PackageModel) error

// IndexArgs lol
type IndexArgs struct {
	Init                    Init
	PackageDeleter          packageDeleter
	PackageInserter         packageInserter
	PackageRetriever        packageRetriever
	RequestTimeBuffer       time.Duration
	NewGithubRequestService githubRequestService
}

// packageRepoTuple lol
type packageRepoTuple struct {
	pkg      *models.PackageModel
	repoData dtos.GithubRepoDTO
}
