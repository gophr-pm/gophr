package download

import (
	"time"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/gophr-pm/gophr/lib/github"
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
