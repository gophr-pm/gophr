package download

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/gophr-pm/gophr/lib/github"
)

// splitType is an enum representing a type of split.
type splitType int

const (
	dailySplit = splitType(iota)
	weeklySplit
	monthlySplit
	allTimeSplit
)

// countResult is the result of countHistoricalDownloads or
// countAllTimeDownloads.
type countResult struct {
	err   error
	sha   string
	count int
	split splitType
}

const (
	// In the database, the empty string SHA value represents the total of all
	// downloads of that package with any SHA.
	anySHA = ""
)

// assertPackageExistence is a wrapper around pkg.AssertExistence that puts the
// return value in a result channel instead of via a function return.
func assertPackageExistence(
	q db.Queryable,
	author string,
	repo string,
	ghSvc github.RequestService,
	resultChan chan error,
) {
	if err := pkg.AssertExistence(q, author, repo, ghSvc); err != nil {
		resultChan <- err
		return
	}

	resultChan <- nil
}

// bumpDownloads batches together all of the counter bump queries necessary to
// record this download in the database.
func bumpDownloads(
	b db.Batchable,
	day time.Time,
	author string,
	repo string,
	sha string,
	resultChan chan error,
) {
	// Counter batches must be unlogdged (as of Cassandra 2.1).
	batch := b.NewUnloggedBatch()
	// Create and add the update queries.
	addHourlyBumpQuery(batch, day, author, repo)
	addAllTimeBumpQuery(batch, author, repo, sha)
	addAllTimeBumpQuery(batch, author, repo, anySHA)

	if err := batch.Execute(); err != nil {
		resultChan <- err
		return
	}

	resultChan <- nil
}

// addHourlyBumpQuery adds a daily download total increment query to a batch.
func addHourlyBumpQuery(
	b db.Batch,
	hour time.Time,
	author string,
	repo string,
) {
	query.Update(hourlyTableName).
		Increment(hourlyColumnNameTotal, 1).
		Where(query.Column(hourlyColumnNameHour).Equals(hour)).
		And(query.Column(hourlyColumnNameAuthor).Equals(author)).
		And(query.Column(hourlyColumnNameRepo).Equals(repo)).
		AppendTo(b)
}

// addAllTimeBumpQuery adds an all-time download total increment query to a
// batch.
func addAllTimeBumpQuery(
	b db.Batch,
	author string,
	repo string,
	sha string,
) {
	query.Update(allTimeTableName).
		Increment(allTimeColumnNameTotal, 1).
		And(query.Column(allTimeColumnNameAuthor).Equals(author)).
		And(query.Column(allTimeColumnNameRepo).Equals(repo)).
		And(query.Column(allTimeColumnNameSHA).Equals(sha)).
		AppendTo(b)
}

// countHistoricalDownloads queries the database to count the the number of
// downloads of a package over a specific split.
func countHistoricalDownloads(
	q db.Queryable,
	author string,
	repo string,
	now time.Time,
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
		countAllTimeDownloads(q, author, repo, anySHA, resultsChan)
		return
	}

	// Change the date boundary depending on the split.
	switch split {
	case dailySplit:
		from = now.AddDate(0, 0, -1) // months, days, hours.
	case weeklySplit:
		from = now.AddDate(0, 0, -7) // months, days, hours.
	case monthlySplit:
		from = now.AddDate(0, -1, 0) // months, days, hours.
	}

	// Run the query using the recently calculated date boundary.
	if err = query.
		SelectSum(hourlyColumnNameTotal).
		From(hourlyTableName).
		Where(query.Column(hourlyColumnNameAuthor).Equals(author)).
		And(query.Column(hourlyColumnNameRepo).Equals(repo)).
		And(query.Column(hourlyColumnNameHour).IsGreaterThanOrEqualTo(from)).
		And(query.Column(hourlyColumnNameHour).IsLessThanOrEqualTo(now)).
		Create(q).
		Scan(&count); err != nil {
		// Check to see if the daily download simply does not exist. If so, just
		// return 0.
		if db.IsErrNotFound(err) {
			resultsChan <- countResult{count: 0, split: split}
			return
		}

		// If some kind of other error,
		resultsChan <- countResult{err: err}
		return
	}

	// Publish the recently fetched count to the channel of results.
	resultsChan <- countResult{count: count, split: split}
}

// countAllTimeDownloads queries the database to count the the number of
// downloads of a package that have ever been recorded.
func countAllTimeDownloads(
	q db.Queryable,
	author string,
	repo string,
	sha string,
	resultsChan chan countResult,
) {
	var (
		err   error
		count int
	)

	if err = query.
		Select(allTimeColumnNameTotal).
		From(allTimeTableName).
		Where(query.Column(allTimeColumnNameAuthor).Equals(author)).
		And(query.Column(allTimeColumnNameRepo).Equals(repo)).
		And(query.Column(allTimeColumnNameSHA).Equals(sha)).
		Create(q).
		Scan(&count); err != nil {
		// Check to see if the daily download simply does not exist. If so, just
		// return 0.
		if db.IsErrNotFound(err) {
			resultsChan <- countResult{count: 0, sha: sha, split: allTimeSplit}
			return
		}

		// If some kind of other error,
		resultsChan <- countResult{err: err}
		return
	}

	// Publish the recently fetched count to the channel of results.
	resultsChan <- countResult{
		sha:   sha,
		count: count,
		split: allTimeSplit,
	}
	return
}

// TODO(skeswa): this shouldn't need to exist. Get rid of this as soon as
// @Shikkic merges in his changes.
func concatErrors(msg string, errs []error) error {
	var buffer bytes.Buffer

	buffer.WriteString(msg)
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
