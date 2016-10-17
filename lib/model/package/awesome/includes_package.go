package awesome

import (
	"fmt"

	"github.com/gophr-pm/gophr/lib/db/query"
)

// IncludesPackage returns true if the package matching the specified author
// and repo has been recorded as an "awesome" package in the database.
func IncludesPackage(q query.Queryable, author, repo string) (bool, error) {
	var (
		err   error
		count int
	)

	if err = query.SelectCount().
		From(tableName).
		Where(query.Column(columnNameAuthor).Equals(author)).
		And(query.Column(columnNameRepo).Equals(repo)).
		Limit(1).
		Create(q).
		Scan(&count); err != nil {
		return false, fmt.Errorf(
			"Failed to check if package %s/%s is awesome: %v",
			author,
			repo,
			err)
	}

	return count > 0, nil
}
