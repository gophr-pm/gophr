package models

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/common/db/query"
)

// Database string constants.
const (
	tableNameDailyDownloads         = "daily_downloads"
	columnNameDailyDownloadsDay     = "day"
	columnNameDailyDownloadsAuthor  = "author"
	columnNameDailyDownloadsRepo    = "repo"
	columnNameDailyDownloadsVersion = "version"
	columnNameDailyDownloadsTotal   = "total"
)

// RecordDailyDownload records a single download of specific package version.
func RecordDailyDownload(
	session *gocql.Session,
	author string,
	repo string,
	version string,
) error {
	// TODO(skeswa): look into batching.

	// Derive today's date.
	now := time.Now()
	// Normalize the day date by setting all the time fields to zero.
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Create the update query for the specific version.
	update := query.Update(tableNameDailyDownloads).
		Increment(columnNameDailyDownloadsTotal, 1).
		Where(query.Column(columnNameDailyDownloadsDay).Equals(today)).
		And(query.Column(columnNameDailyDownloadsAuthor).Equals(author)).
		And(query.Column(columnNameDailyDownloadsRepo).Equals(repo)).
		And(query.Column(columnNameDailyDownloadsVersion).Equals(version)).
		Create(session)

	// Execute the first update query. Exit if it fails.
	if err := update.Exec(); err != nil {
		return err
	}

	// Create the update query for the whole package count.
	update = query.Update(tableNameDailyDownloads).
		Increment(columnNameDailyDownloadsTotal, 1).
		Where(query.Column(columnNameDailyDownloadsDay).Equals(today)).
		And(query.Column(columnNameDailyDownloadsAuthor).Equals(author)).
		And(query.Column(columnNameDailyDownloadsRepo).Equals(repo)).
		And(query.Column(columnNameDailyDownloadsVersion).Equals("")).
		Create(session)

	// Execute the first update query. Exit if it fails.
	return update.Exec()
}
