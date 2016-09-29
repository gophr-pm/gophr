package models

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/common/db/query"
)

// Database string constants.
const (
	tableNamePackageArchives        = "package_archive_records"
	columnNamePackageArchivesAuthor = "author"
	columnNamePackageArchivesRepo   = "repo"
	columnNamePackageArchivesRef    = "ref"
)

// RecordPackageArchival records that an archive of a package version exists.
func RecordPackageArchival(
	session *gocql.Session,
	author string,
	repo string,
	ref string,
) error {
	// Create the update query for the specific ref.
	insert := query.InsertInto(tableNamePackageArchives).
		Value(columnNamePackageArchivesAuthor, author).
		Value(columnNamePackageArchivesRepo, repo).
		Value(columnNamePackageArchivesRef, ref).
		Create(session)

	// Execute the first update query. Exit if it fails.
	if err := insert.Exec(); err != nil {
		return err
	}

	return nil
}

// IsPackageArchived returns true if a package version matching the parameters
// exists.
func IsPackageArchived(
	session *gocql.Session,
	author string,
	repo string,
	ref string,
) (bool, error) {
	var (
		err   error
		count int
	)

	if err = query.Select(fmt.Sprintf("COUNT(*)")).
		From(tableNamePackageArchives).
		Where(query.Column(columnNamePackageArchivesAuthor).Equals(author)).
		And(query.Column(columnNamePackageArchivesRepo).Equals(repo)).
		And(query.Column(columnNamePackageArchivesRef).Equals(ref)).
		Create(session).
		Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}
