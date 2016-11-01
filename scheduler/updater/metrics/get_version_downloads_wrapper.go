package metrics

import (
	"sync"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/db"
)

type refsFetcher func(author string, repo string) (lib.Refs, error)
type versionDownloadsGetter func(
	q db.Queryable,
	author string,
	repo string,
	shaVersions map[string]string,
) (map[string]int, error)

type getVersionDownloadsWrapperArgs struct {
	q                   db.Queryable
	wg                  *sync.WaitGroup
	repo                string
	author              string
	result              *getVersionDownloadsWrapperResult
	fetchRefs           refsFetcher
	getVersionDownloads versionDownloadsGetter
}

type getVersionDownloadsWrapperResult struct {
	err              error
	versionDownloads map[string]int
}

func getVersionDownloadsWrapper(args getVersionDownloadsWrapperArgs) {
	var (
		err              error
		refs             lib.Refs
		shaVersions      = make(map[string]string)
		versionDownloads map[string]int
	)

	// Guarantee that the waitgroup is notified at the end.
	defer args.wg.Done()

	// Fetch the refs for this package to get the candidates.
	// TODO(skeswa): fetchRefs can be refactored to simply get candidates.
	if refs, err = args.fetchRefs(args.author, args.repo); err != nil {
		*args.result = getVersionDownloadsWrapperResult{err: err}
		return
	}

	// Compile a map of version SHAs to version names.
	for _, candidate := range refs.Candidates {
		// candidate.String() is the stringified semver version.
		shaVersions[candidate.GitRefHash] = candidate.String()
	}

	// shaVersions -> versionDownloads.
	if versionDownloads, err = args.getVersionDownloads(
		args.q,
		args.author,
		args.repo,
		shaVersions); err != nil {
		*args.result = getVersionDownloadsWrapperResult{err: err}
		return
	}

	// Operation was successful, return.
	*args.result = getVersionDownloadsWrapperResult{
		versionDownloads: versionDownloads,
	}
}
