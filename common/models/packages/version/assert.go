package versions

import "github.com/gophr-pm/gophr/common/db/query"

// AssertExistence asserts that a package version exists.
func AssertExistence(
	q query.Queryable,
	author string,
	repo string,
	sha string,
	version string,
) error {
	var (
		err   error
		count int
	)

	// Create the query to check if this package version exists.
	if err = query.SelectCount().
		From(tableName).
		Where(query.Column(columnNameAuthor).Equals(author)).
		And(query.Column(columnNameRepo).Equals(repo)).
		And(query.Column(columnNameVersion).Equals(version)).
		Create(q).
		Scan(&count); err != nil {
		return err
	}

	// If this package version doesn't exist, then make it exist.
	if count < 1 {
		if err := query.InsertInto(tableName).
			Value(columnNameAuthor, author).
			Value(columnNameRepo, repo).
			Value(columnNameSHA, sha).
			Value(columnNameVersion, version).
			Create(q).
			Exec(); err != nil {
			return err
		}
	}

	return nil
}
