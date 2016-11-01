package verdeps

import (
	"github.com/gophr-pm/gophr/lib/errors"
	"fmt"

	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
)

// depsProcessor is a function type that de-couples verdeps.VersionDeps from
// verdeps.processDeps.
type depsProcessor func(args processDepsArgs) error

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
	// processDeps version-locks all appropriate dependencies in a package
	// directory, while keeping in mind chronological accuracy. If unspecified,
	// verdeps.processDeps will be used.
	processDeps depsProcessor
	// GithubService is the service, with which, requests can be made of the
	// Github API.
	GithubService github.RequestService
}

// VersionDeps version locks all of the Github-based Go dependencies referenced
// in the source code of a package. It takes a variety of package metadata and
// the path to the source code, and changes its dependencies accordingly.
func VersionDeps(args VersionDepsArgs) error {
	if args.IO == nil {
		return errors.New("Invalid IO.")
	} else if len(args.SHA) < 1 {
		return errors.New("Invalid SHA.")
	} else if len(args.Path) < 1 {
		return errors.New("Invalid Path.")
	} else if len(args.Repo) < 1 {
		return errors.New("Invalid Repo.")
	} else if len(args.Author) < 1 {
		return errors.New("Invalid Author.")
	} else if args.GithubService == nil {
		return errors.New("Invalid GithubService.")
	}

	// Fallback to verdeps.processDeps if no override is supplied.
	if args.processDeps == nil {
		args.processDeps = processDeps
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

	return args.processDeps(processDepsArgs{
		io:                      args.IO,
		ghSvc:                   args.GithubService,
		fetchSHA:                fetchSHA,
		reviseDeps:              reviseDeps,
		packageSHA:              args.SHA,
		packagePath:             args.Path,
		packageRepo:             args.Repo,
		packageAuthor:           args.Author,
		readPackageDir:          readPackageDir,
		packageVersionDate:      commitDate,
		newSpecWaitingList:      newSpecWaitingList,
		newSyncedStringMap:      newSyncedStringMap,
		newSyncedWaitingListMap: newSyncedWaitingListMap,
	})
}
