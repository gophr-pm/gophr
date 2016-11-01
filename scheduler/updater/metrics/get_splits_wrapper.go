package metrics

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package/download"
)

type downloadSplitGetter func(
	q db.Queryable,
	author string,
	repo string,
) (download.Splits, error)

// getSplitsWrapperArgs is the reuslts struct for getSplitsWrapper.
type getSplitsWrapperResult struct {
	err    error
	splits download.Splits
}

// getSplitsWrapperArgs is the arguments struct for getSplitsWrapper.
type getSplitsWrapperArgs struct {
	q         db.Queryable
	wg        *sync.WaitGroup
	repo      string
	author    string
	result    *getSplitsWrapperResult
	getSplits downloadSplitGetter
}

// getSplitsWrapper wraps the getSplits function and formats the outputs for use
// by packageUpdater.
func getSplitsWrapper(args getSplitsWrapperArgs) {
	var result getSplitsWrapperResult

	result.splits, result.err = args.getSplits(args.q, args.author, args.repo)
	*args.result = result
	args.wg.Done()
}
