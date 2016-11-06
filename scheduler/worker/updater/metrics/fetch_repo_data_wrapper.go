package metrics

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/dtos"
)

type repoDataFetcher func(author string, repo string) (dtos.GithubRepo, error)

type fetchRepoDataWrapperResult struct {
	err      error
	repoData dtos.GithubRepo
}

type fetchRepoDataWrapperArgs struct {
	wg            *sync.WaitGroup
	repo          string
	author        string
	result        *fetchRepoDataWrapperResult
	fetchRepoData repoDataFetcher
}

// fetchRepoDataWrapper wraps the fetchRepoData function and formats the outputs
// for use by packageUpdater.
func fetchRepoDataWrapper(args fetchRepoDataWrapperArgs) {
	var (
		err      error
		repoData dtos.GithubRepo
	)

	// Guarantee that the waitgroup is notified at the end.
	defer args.wg.Done()

	if repoData, err = args.fetchRepoData(args.author, args.repo); err != nil {
		*args.result = fetchRepoDataWrapperResult{err: err}
		return
	}

	*args.result = fetchRepoDataWrapperResult{repoData: repoData}
}
