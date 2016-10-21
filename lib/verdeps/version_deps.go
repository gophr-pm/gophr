package verdeps

import (
	"github.com/gophr-pm/gophr/lib/errors"
	"fmt"

	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
)

// VersionDepsArgs is the arguments struct for VersionDeps(...).
type VersionDepsArgs struct {
	// IO is the input/output interface.
	IO io.IO
	// SHA is the sha of the package being versioned.
	SHA string
	// Repo is the repo of the package being versioned.
	Repo string
	// SHA is the path to the package source code to be versioned.
	Path string
	// Author is the author of the package being versioned.
	Author string
	// GithubServcie is the service, with which, requests can be made of the Github API.
	GithubService github.RequestService
}

// VersionDeps version locks all of the Github-based Go dependencies referenced
// in the source code of a package. It takes a variety of package metadata and
// the path to the source code, and changes its dependencies accordingly.
func VersionDeps(args VersionDepsArgs) error {
	if len(args.SHA) < 1 {
		return errors.New("Invalid SHA.")
	} else if len(args.Path) < 1 {
		return errors.New("Invalid Path.")
	} else if len(args.Repo) < 1 {
		return errors.New("Invalid Model.Repo.")
	} else if len(args.Author) < 1 {
		return errors.New("Invalid Model.Author.")
	} else if args.GithubService == nil {
		return errors.New("Invalid GithubService.")
	}

	// Fetch the timestamp of the commit SHA.
	commitDate, err := args.GithubService.FetchCommitTimestamp(
		args.Author,
		args.Repo,
		args.SHA,
	)
	if err != nil {
		return fmt.Errorf("Could not fetch commit timestamp: %v.", err)
	}

	return processDeps(processDepsArgs{
		io:                 args.IO,
		ghSvc:              args.GithubService,
		packageSHA:         args.SHA,
		packagePath:        args.Path,
		packageRepo:        args.Repo,
		packageAuthor:      args.Author,
		packageVersionDate: commitDate,
	})
}
