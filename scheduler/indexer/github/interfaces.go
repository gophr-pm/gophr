package githubIndexer

import (
	"time"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/model"
)

// Init is responsible for setting up the app configuration and db
// connection.
type Init func() (*config.Config, db.Client)

// githubRequestService lol
type githubRequestService func(args github.RequestServiceArgs) github.RequestService

// packageRetriever lol
type packageRetriever func(session db.Client) ([]*models.PackageModel, error)

// packageInserter lol
type packageInserter func(session db.Client, packageModel *models.PackageModel) error

// packageDeleter lol
type packageDeleter func(session db.Client, packageModel *models.PackageModel) error

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
	repoData dtos.GithubRepo
}
