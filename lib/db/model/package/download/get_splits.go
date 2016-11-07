package download

import (
	"time"

	"github.com/gophr-pm/gophr/lib/db"
)

// Splits contains the downloads count for every time split.
type Splits struct {
	Daily   int
	Weekly  int
	Monthly int
	AllTime int
}

// GetSplits counts and returns the download totals for a package over different
// lengths of time.
func GetSplits(
	q db.Queryable,
	author string,
	repo string,
) (Splits, error) {
	var (
		// Normalize the day date by setting all the time fields to zero.
		errs         []error
		now          = time.Now()
		splits       Splits
		resultsChan  = make(chan countResult)
		resultsCount = 0
		resultsTotal = 4 // is called countHistoricalDownloads 4 times.
	)

	// Execute the first update query. Exit if it fails.
	go countHistoricalDownloads(q, author, repo, now, dailySplit, resultsChan)
	go countHistoricalDownloads(q, author, repo, now, weeklySplit, resultsChan)
	go countHistoricalDownloads(q, author, repo, now, monthlySplit, resultsChan)
	go countHistoricalDownloads(q, author, repo, now, allTimeSplit, resultsChan)

	// Read all of the results, then exit when we run out.
	for result := range resultsChan {
		if result.err != nil {
			errs = append(errs, result.err)
		} else {
			switch result.split {
			case dailySplit:
				splits.Daily = result.count
			case weeklySplit:
				splits.Weekly = result.count
			case monthlySplit:
				splits.Monthly = result.count
			case allTimeSplit:
				splits.AllTime = result.count
			}
		}

		if resultsCount++; resultsCount == resultsTotal {
			close(resultsChan)
		}
	}

	// If there were any errors, return them composed together.
	if len(errs) > 0 {
		return splits, concatErrors(
			"Failed to read download splits from the database.",
			errs)
	}

	return splits, nil
}
