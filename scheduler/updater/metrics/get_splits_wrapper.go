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

type getSplitsWrapperResult struct {
	err    error
	splits download.Splits
}

type getSplitsWrapperArgs struct {
	q         db.Queryable
	wg        *sync.WaitGroup
	repo      string
	author    string
	result    *getSplitsWrapperResult
	getSplits downloadSplitGetter
}

func getSplitsWrapper(args getSplitsWrapperArgs) {
	var result getSplitsWrapperResult

	result.splits, result.err = args.getSplits(args.q, args.author, args.repo)
	*args.result = result
	args.wg.Done()
}
