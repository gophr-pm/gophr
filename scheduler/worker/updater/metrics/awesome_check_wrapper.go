package metrics

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
)

// awesomeChecker is a proxy for awesome.IncludesPackage.
type awesomeChecker func(
	q db.Queryable,
	author string,
	repo string,
) (bool, error)

// awesomeCheckWrapperArgs is the reuslts struct for awesomeCheckWrapper.
type awesomeCheckWrapperResult struct {
	err     error
	awesome bool
}

// awesomeCheckWrapperArgs is the arguments struct for awesomeCheckWrapper.
type awesomeCheckWrapperArgs struct {
	q         db.Queryable
	wg        *sync.WaitGroup
	repo      string
	author    string
	result    *awesomeCheckWrapperResult
	isAwesome awesomeChecker
}

// awesomeCheckWrapper wraps the getSplits function and formats the outputs for
// use by packageUpdater.
func awesomeCheckWrapper(args awesomeCheckWrapperArgs) {
	var result awesomeCheckWrapperResult

	result.awesome, result.err = args.isAwesome(
		args.q,
		args.author,
		args.repo)

	*args.result = result
	args.wg.Done()
}
