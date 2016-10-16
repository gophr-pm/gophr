package archives

import (
	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// Create records that an archive of a package version exists.
func Create(
	session *gocql.Session,
	author string,
	repo string,
	sha string,
) error {
	// Execute the first update query. Exit if it fails.
	if err := query.InsertInto(tableName).
		Value(columnNameAuthor, author).
		Value(columnNameRepo, repo).
		Value(columnNameSHA, sha).
		Create(session).
		Exec(); err != nil {
		return err
	}

	return nil
}
