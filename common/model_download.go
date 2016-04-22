package common

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

// Constants directly related to interacting with the download model in the
// cassandra database.
const (
	// TableNameDownloads is the name of the table containing the download model.
	TableNamePackageDownloads               = "package_downloads"
	ColumnNamePackageDownloadsDate          = "date"
	ColumnNamePackageDownloadsCount         = "count"
	ColumnNamePackageDownloadsPackageRepo   = "package_repo"
	ColumnNamePackageDownloadsPackageAuthor = "package_author"
)

var (
	cqlQueryIncrementPackageDownloadCount = fmt.Sprintf(
		`UPDATE %s SET %s = %s + 1 WHERE %s = ? AND %s = ? AND %s = ?`,
		TableNamePackageDownloads,
		ColumnNamePackageDownloadsCount,
		ColumnNamePackageDownloadsCount,
		ColumnNamePackageDownloadsPackageAuthor,
		ColumnNamePackageDownloadsPackageRepo,
		ColumnNamePackageDownloadsDate,
	)
)

func IncrementDownloadCount(
	session *gocql.Session,
	packageAuthor string,
	packageRepo string,
) error {
	var (
		currentTime = time.Now().UTC()
		currentDate = time.Date(
			currentTime.Year(),
			currentTime.Month(),
			currentTime.Day(),
			0,
			0,
			0,
			0,
			currentTime.Location(),
		)
	)

	return session.Query(
		cqlQueryIncrementPackageDownloadCount,
		packageAuthor,
		packageRepo,
		currentDate,
	).Exec()
}
