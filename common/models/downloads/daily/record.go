package downloads

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/common/db/query"
	"github.com/gophr-pm/gophr/common/models/packages"
	"github.com/gophr-pm/gophr/common/models/packages/version"
)

// Record records a single download of specific package version.
func Record(
	q query.BatchingQueryable,
	author string,
	repo string,
	sha string,
	version string,
) error {
	var (
		// Normalize the day date by setting all the time fields to zero.
		now         = time.Now()
		today       = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		resultsChan = make(chan error)
	)

	// Execute the first update query. Exit if it fails.
	go bumpDownloads(q, today, author, repo, sha, version, resultsChan)
	go assertPackageExistence(q, author, repo, resultsChan)
	go assertPackageVersionExistence(q, author, repo, sha, version, resultsChan)

	var (
		i    = 0
		errs []error
	)

	// Read all (3) of the responses.
	for err := range resultsChan {
		// Record all of the errors.
		if err != nil {
			errs = append(errs, err)
		}

		i++
		// Exit if every db call has returned.
		if i >= 3 {
			close(resultsChan)
		}
	}

	if len(errs) > 0 {
		// TODO(skeswa): this is also done in verdeps. This should be extracted to
		// a helper function.
		var buffer bytes.Buffer
		buffer.WriteString("Failed to record a daily download due to ")
		buffer.WriteString(strconv.Itoa(len(errs)))
		buffer.WriteString(" error(s): [ ")
		for i, err := range errs {
			if i > 0 {
				buffer.WriteString(", ")
			}

			buffer.WriteString(err.Error())
		}
		buffer.WriteString(" ]")

		return errors.New(buffer.String())
	}

	return nil
}

func assertPackageExistence(
	q query.Queryable,
	author string,
	repo string,
	resultChan chan error,
) {
	if err := packages.AssertExistence(q, author, repo); err != nil {
		resultChan <- err
		return
	}

	resultChan <- nil
}

func assertPackageVersionExistence(q query.Queryable,
	author string,
	repo string,
	sha string,
	version string,
	resultChan chan error,
) {
	if err := versions.AssertExistence(
		q,
		author,
		repo,
		sha,
		version,
	); err != nil {
		resultChan <- err
		return
	}

	resultChan <- nil
}

func bumpDownloads(
	b query.Batchable,
	day time.Time,
	author string,
	repo string,
	sha string,
	version string,
	resultChan chan error,
) {
	batch := b.NewBatch(gocql.UnloggedBatch)
	// Create the update query for the specific version.
	addBumpQuery(batch, day, author, repo, sha, version)
	// Create the update query for the whole package count.
	addBumpQuery(batch, day, author, repo, "", "")

	if err := b.ExecuteBatch(batch); err != nil {
		resultChan <- err
		return
	}

	resultChan <- nil
}

func addBumpQuery(
	q query.VoidQueryable,
	day time.Time,
	author string,
	repo string,
	sha string,
	version string,
) {
	query.Update(tableName).
		Increment(columnNameTotal, 1).
		Where(query.Column(columnNameDay).Equals(day)).
		And(query.Column(columnNameAuthor).Equals(author)).
		And(query.Column(columnNameRepo).Equals(repo)).
		And(query.Column(columnNameVersion).Equals(version)).
		CreateVoid(q)
}
