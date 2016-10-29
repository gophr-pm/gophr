package download

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// splitType is an enum representing a type of split.
type splitType int

const (
	dailySplit = splitType(iota)
	weeklySplit
	monthlySplit
	allTimeSplit
)

// Splits contains the downloads count for every time split.
type Splits struct {
	Daily   int
	Weekly  int
	Monthly int
	AllTime int
}

// countResult is the result of countHistoricalDownloads.
type countResult struct {
	err   error
	count int
	split splitType
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

		today = time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			0,
			0,
			0,
			0,
			time.UTC)
	)

	// Execute the first update query. Exit if it fails.
	go countHistoricalDownloads(q, author, repo, today, dailySplit, resultsChan)
	go countHistoricalDownloads(q, author, repo, today, weeklySplit, resultsChan)
	go countHistoricalDownloads(q, author, repo, today, monthlySplit, resultsChan)
	go countHistoricalDownloads(q, author, repo, today, allTimeSplit, resultsChan)

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
		return splits, concatGetSplitErrors(errs)
	}

	return splits, nil
}

// countHistoricalDownloads queries the database to count the the number of
// downloads of a package over a specific split.
func countHistoricalDownloads(
	q db.Queryable,
	author string,
	repo string,
	today time.Time,
	split splitType,
	resultsChan chan countResult,
) {
	var (
		err   error
		from  time.Time
		count int
	)

	// All-time queries are different since the data is in a different table.
	if split == allTimeSplit {
		if err = query.
			SelectSum(dailyColumnNameTotal).
			From(dailyTableName).
			Where(query.Column(dailyColumnNameAuthor).Equals(author)).
			And(query.Column(dailyColumnNameRepo).Equals(repo)).
			And(query.Column(dailyColumnNameSHA).Equals("")).
			Create(q).
			Scan(&count); err != nil {
			resultsChan <- countResult{err: err}
			return
		}

		// Publish the recently fetched count to the database.
		resultsChan <- countResult{count: count, split: split}
		return
	}

	// Change the date boundary depending on the split.
	switch split {
	case dailySplit:
		from = today
	case weeklySplit:
		from = today.AddDate(0, 0, -7)
	case monthlySplit:
		from = today.AddDate(0, -1, 0)
	}

	// Run the query using the recently calculated date boundary.
	if err = query.
		SelectSum(dailyColumnNameTotal).
		From(dailyTableName).
		Where(query.Column(dailyColumnNameAuthor).Equals(author)).
		And(query.Column(dailyColumnNameRepo).Equals(repo)).
		And(query.Column(dailyColumnNameSHA).Equals("")).
		And(query.Column(dailyColumnNameDay).IsGreaterThanOrEqualTo(from)).
		And(query.Column(dailyColumnNameDay).IsLessThanOrEqualTo(today)).
		Create(q).
		Scan(&count); err != nil {
		resultsChan <- countResult{err: err}
		return
	}

	// Publish the recently fetched count to the database.
	resultsChan <- countResult{count: count, split: split}
}

// TODO(skeswa): this shouldn't need to exist. Get rid of this as soon as
// @Shikkic merges in his changes.
func concatGetSplitErrors(errs []error) error {
	var buffer bytes.Buffer

	buffer.WriteString("Failed to read download splits from the database.")
	buffer.WriteString(" Bumped into ")
	buffer.WriteString(strconv.Itoa(len(errs)))
	buffer.WriteString(" problems: [ ")

	for i, err := range errs {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(err.Error())
	}

	buffer.WriteString(" ].")

	return errors.New(buffer.String())
}
