package archives

import (
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// Create records that an archive of a package version exists.
func Create(
	q db.Queryable,
	author string,
	repo string,
	sha string,
) error {
	// Execute the first update query. Exit if it fails.
	if err := query.InsertInto(tableName).
		Value(columnNameAuthor, author).
		Value(columnNameRepo, repo).
		Value(columnNameSHA, sha).
		Create(q).
		Exec(); err != nil {
		return err
	}

	return nil
}
