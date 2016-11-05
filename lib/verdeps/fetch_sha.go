package verdeps

import (
	"errors"
	"time"

	"github.com/gophr-pm/gophr/lib/github"
)

type fetchSHAArgs struct {
	ghSvc              github.RequestService
	outputChan         chan *fetchSHAResult
	importPath         string
	packageSHA         string
	packageRepo        string
	packageAuthor      string
	packageVersionDate time.Time
}

func fetchSHA(args fetchSHAArgs) {
	var (
		err    error
		sha    string
		repo   string
		author string
	)

	// Parse out the author and the repo.
	author, repo, _ = ParseImportPath(args.importPath)

	// If the dep is a sub-package. If it is, don't fetch the commit sha.
	if isSubPackage(author, args.packageAuthor, repo, args.packageRepo) {
		sha = args.packageSHA
	} else {
		// Fetch the most appropriate commit sha for this package given the time
		// constraint.
		if sha, err = args.ghSvc.FetchCommitSHA(
			author,
			repo,
			args.packageVersionDate,
		); err != nil {
			args.outputChan <- newFetchSHAFailure(err)
			return
		} else if len(sha) == 0 {
			args.outputChan <- newFetchSHAFailure(
				errors.New("Commit SHA it came back empty"))
			return
		}
	}

	// Put a new mapping struct into the output chan.
	args.outputChan <- newFetchSHASuccess(args.importPath, sha)
}
