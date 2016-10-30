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
	addDailyBumpQuery(batch, day, author, repo)
	addAllTimeBumpQuery(batch, author, repo, sha)
	addAllTimeBumpQuery(batch, author, repo, "") // "" represents "all" versions.

	if err := batch.Execute(); err != nil {
		resultChan <- err
		return
	}

	resultChan <- nil
}

// addDailyBumpQuery adds a daily download total increment query to a batch.
func addDailyBumpQuery(
	b db.Batch,
	day time.Time,
	author string,
	repo string,
) {
	query.Update(dailyTableName).
		Increment(dailyColumnNameTotal, 1).
		Where(query.Column(dailyColumnNameDay).Equals(day)).
		And(query.Column(dailyColumnNameAuthor).Equals(author)).
		And(query.Column(dailyColumnNameRepo).Equals(repo)).
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
			Select(allTimeColumnNameTotal).
			From(allTimeTableName).
			Where(query.Column(allTimeColumnNameAuthor).Equals(author)).
			And(query.Column(allTimeColumnNameRepo).Equals(repo)).
			And(query.Column(allTimeColumnNameSHA).Equals("")).
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
