package verdeps

import (
	"log"
	"time"

	"github.com/skeswa/gophr/common/github"
)

type fetchSHAArgs struct {
	ghSvc              *github.RequestService
	outputChan         chan *importPathSHA
	importPath         string
	packageSHA         string
	packageRepo        string
	packageAuthor      string
	pendingSHARequests *syncedInt
	packageVersionDate time.Time
}

func fetchSHA(args fetchSHAArgs) {
	var (
		err    error
		sha    string
		repo   string
		author string
	)

	// Signal that the sha fetching is in progress.
	args.pendingSHARequests.increment()
	defer args.pendingSHARequests.decrement()

	// Parse out the author and the repo.
	author, repo, _ = parseImportPath(args.importPath)

	// If the dep is a sub-package. If it is, don't fetch the commit sha.
	if isSubPackage(author, args.packageAuthor, repo, args.packageRepo) {
		sha = args.packageSHA
	} else {
		// Fetch the most appropriate commit sha for this package given the time
		// constraint.
		if sha, err = args.ghSvc.FetchCommitSHA(author, repo, args.packageVersionDate); err != nil {
			// Don't enqueue errors in the chan since they arent fatal. Just log the
			// failures.
			log.Printf("Failed to fetch the commit sha for %s: %v.\n", args.importPath, err)
			return
		} else if len(sha) == 0 {
			// Don't enqueue errors in the error chan since they arent fatal. Just log
			// the failures.
			log.Printf("Failed to fetch the commit sha for %s: came back from Github empty.\n", args.importPath)
			return
		}
	}

	// Put a new mapping struct into the output chan.
	args.outputChan <- &importPathSHA{
		sha:        sha,
		importPath: args.importPath,
	}
}
