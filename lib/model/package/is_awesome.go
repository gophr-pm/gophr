package pkg

import (
	"fmt"

	"github.com/gophr-pm/gophr/lib/db/query"
)

// IsAwesome returns true if the package matching the specified author and repo
// has been recorded as an "awesome" package in the database.
func IsAwesome(q query.Queryable, author, repo string) (bool, error) {
	var (
		err   error
		count int
	)

	if err = query.SelectCount().
		From(awesomeTableName).
		Where(query.Column(awesomeColumnNameAuthor).Equals(author)).
		And(query.Column(awesomeColumnNameRepo).Equals(repo)).
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
