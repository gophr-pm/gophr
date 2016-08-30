package verdeps

import (
	"errors"
	"time"

	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/models"
)

// VersionDepsArgs is the arguments struct for VersionDeps(...).
type VersionDepsArgs struct {
	// SHA is the sha of the package being versioned.
	SHA string
	// SHA is the path to the package source code to be versioned.
	Path string
	// Date is date that the version of the package with that matches SHA was created.
	Date time.Time
	// Model is the package model of the package.
	Model *models.PackageModel
	// GithubServcie is the service, with which, requests can be made of the Github API.
	GithubService *github.GitHubRequestService
}

// VersionDeps version locks all of the Github-based Go dependencies referenced
// in the source code of a package. It takes a variety of package metadata and
// the path to the source code, and changes its dependencies accordingly.
func VersionDeps(args VersionDepsArgs) error {
	if len(args.SHA) < 1 {
		return errors.New("Invalid SHA.")
	} else if len(args.Path) < 1 {
		return errors.New("Invalid Path.")
	} else if args.Model == nil {
		return errors.New("Invalid Model.")
	} else if args.Model.Repo == nil || len(*args.Model.Repo) < 1 {
		return errors.New("Invalid Model.Repo.")
	} else if args.Model.Author == nil || len(*args.Model.Author) < 1 {
		return errors.New("Invalid Model.Author.")
	} else if args.GithubService == nil {
		return errors.New("Invalid GithubService.")
	}

	return processDeps(processDepsArgs{
		ghSvc:              args.GithubService,
		packageSHA:         args.SHA,
		packagePath:        args.Path,
		packageRepo:        *args.Model.Repo,
		packageAuthor:      *args.Model.Author,
		packageVersionDate: args.Date,
	})
}
