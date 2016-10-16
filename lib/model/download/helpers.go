package download

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/model/package"
)

// assertPackageExistence is a wrapper around pkg.AssertExistence that puts the
// return value in a result channel instead of via a function return.
func assertPackageExistence(
	q query.Queryable,
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
	b query.Batchable,
	day time.Time,
	author string,
	repo string,
	sha string,
	resultChan chan error,
) {
	// Counter batches must be unlogdged (as of Cassandra 2.1).
	batch := b.NewBatch(gocql.UnloggedBatch)
	// Create the update queries for the specific version.
	addDailyBumpQuery(batch, day, author, repo, sha)
	addAllTimeBumpQuery(batch, author, repo, sha)
	// Create the update queries for the whole package count.
	addDailyBumpQuery(batch, day, author, repo, "")
	addAllTimeBumpQuery(batch, author, repo, "")

	if err := b.ExecuteBatch(batch); err != nil {
		resultChan <- err
		return
	}

	resultChan <- nil
}

// addDailyBumpQuery adds a daily download total increment query to a batch.
func addDailyBumpQuery(
	q query.VoidQueryable,
	day time.Time,
	author string,
	repo string,
	sha string,
) {
	query.Update(dailyTableName).
		Increment(dailyColumnNameTotal, 1).
		Where(query.Column(dailyColumnNameDay).Equals(day)).
		And(query.Column(dailyColumnNameAuthor).Equals(author)).
		And(query.Column(dailyColumnNameRepo).Equals(repo)).
		And(query.Column(dailyColumnNameSHA).Equals(sha)).
		CreateVoid(q)
}

// addAllTimeBumpQuery adds an all-time download total increment query to a
// batch.
func addAllTimeBumpQuery(
	q query.VoidQueryable,
	author string,
	repo string,
	sha string,
) {
	query.Update(allTimeTableName).
		Increment(allTimeColumnNameTotal, 1).
		And(query.Column(allTimeColumnNameAuthor).Equals(author)).
		And(query.Column(allTimeColumnNameRepo).Equals(repo)).
		And(query.Column(allTimeColumnNameSHA).Equals(sha)).
		CreateVoid(q)
}
