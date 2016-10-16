package archives

import "github.com/gophr-pm/gophr/lib/db/query"

// Exists returns true if a package version matching the parameters exists.
func Exists(
	q query.Queryable,
	author string,
	repo string,
	sha string,
) (bool, error) {
	var (
		err   error
		count int
	)

	if err = query.SelectCount().
		From(tableName).
		Where(query.Column(columnNameAuthor).Equals(author)).
		And(query.Column(columnNameRepo).Equals(repo)).
		And(query.Column(columnNameSHA).Equals(sha)).
		Limit(1).
		Create(q).
		Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}
